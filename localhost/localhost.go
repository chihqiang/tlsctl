package localhost

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/mail"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/chihqiang/tlsctl/pkg/fp"
	"github.com/go-acme/lego/v4/certificate"
)

const (
	rootName     = "rootCA.pem"
	rootKeyName  = "rootCA-key.pem"
	organization = "tlsctl local development CA"
)

type SSL struct {
	Path   string
	caCert *x509.Certificate
	caKey  crypto.PrivateKey
}

func NewLocalHostSSL(path string) (*SSL, error) {
	err := fp.CreateNonExistingFolder(path)
	if err != nil {
		return nil, err
	}
	localHost := &SSL{}
	localHost.Path = path
	localHost.caCert = nil
	localHost.caKey = nil
	return localHost, nil
}

func (l *SSL) LoadCA() error {
	ca, key, err := l.readCA()
	if err != nil {
		return err
	}
	l.caCert = ca
	l.caKey = key
	return nil
}

func (l *SSL) BuildResource(hosts []string) (*certificate.Resource, error) {
	cert, key, err := l.buildCert(hosts)
	if err != nil {
		return nil, err
	}
	resource := &certificate.Resource{}
	resource.Domain = "localhost"
	resource.Certificate = cert
	resource.PrivateKey = key
	resource.IssuerCertificate = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: l.caCert.Raw})
	return resource, nil
}

func (l *SSL) buildCert(hosts []string) (certificate []byte, privateKey []byte, err error) {
	if l.caKey == nil {
		return nil, nil, errors.New("no CA key provided")
	}
	priv, err := l.generateKey(false)
	if err != nil {
		return nil, nil, err
	}
	pub := priv.(crypto.Signer).Public()
	// Certificates last for 2 years and 3 months, which is always less than
	// 825 days, the limit that macOS/iOS apply to all certificates,
	// including custom roots. See https://support.apple.com/en-us/HT210176.
	expiration := time.Now().AddDate(2, 3, 0)
	serialNumber, err := l.randomSerialNumber()
	if err != nil {
		return nil, nil, err
	}
	tpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         "localhost",
			Organization:       []string{organization},
			OrganizationalUnit: []string{l.getUserAndHostname()},
		},
		NotBefore: time.Now(), NotAfter: expiration,
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			tpl.IPAddresses = append(tpl.IPAddresses, ip)
		} else if email, err := mail.ParseAddress(h); err == nil && email.Address == h {
			tpl.EmailAddresses = append(tpl.EmailAddresses, h)
		} else if uriName, err := url.Parse(h); err == nil && uriName.Scheme != "" && uriName.Host != "" {
			tpl.URIs = append(tpl.URIs, uriName)
		} else {
			tpl.DNSNames = append(tpl.DNSNames, h)
		}
	}
	if len(tpl.IPAddresses) > 0 || len(tpl.DNSNames) > 0 || len(tpl.URIs) > 0 {
		tpl.ExtKeyUsage = append(tpl.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}
	if len(tpl.EmailAddresses) > 0 {
		tpl.ExtKeyUsage = append(tpl.ExtKeyUsage, x509.ExtKeyUsageEmailProtection)
	}
	cert, err := x509.CreateCertificate(rand.Reader, tpl, l.caCert, pub, l.caKey)
	if err != nil {
		return nil, nil, err
	}
	certificate = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
	rsaKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("private key is not RSA")
	}
	privDER := x509.MarshalPKCS1PrivateKey(rsaKey)
	privateKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER})
	return certificate, privateKey, nil
}

func (l *SSL) readCA() (*x509.Certificate, crypto.PrivateKey, error) {
	var (
		certPEMBlock []byte
		keyPEMBlock  []byte
		err          error
	)
	certPath := filepath.Join(l.Path, rootName)
	keyPath := filepath.Join(l.Path, rootKeyName)
	if !fp.PathExists(certPath) && !fp.PathExists(keyPath) {
		certPEMBlock, keyPEMBlock, err = l.buildNewCa()
		if err != nil {
			return nil, nil, err
		}
		err = os.WriteFile(certPath, certPEMBlock, 0644)
		if err != nil {
			return nil, nil, err
		}
		err = os.WriteFile(keyPath, keyPEMBlock, 0400)
		if err != nil {
			return nil, nil, err
		}
	} else {
		certPEMBlock, err = os.ReadFile(certPath)
		if err != nil {
			return nil, nil, err
		}
		keyPEMBlock, err = os.ReadFile(keyPath)
		if err != nil {
			return nil, nil, err
		}
	}
	certDERBlock, _ := pem.Decode(certPEMBlock)
	if certDERBlock == nil || certDERBlock.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("failed to read the CA certificate: unexpected content")
	}
	certificate, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil || keyDERBlock.Type != "PRIVATE KEY" {
		return nil, nil, fmt.Errorf("failed to read the CA private key: unexpected content")
	}
	key, err := x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return certificate, key, err
}

func (l *SSL) buildNewCa() (certificateMemory []byte, keyMemory []byte, err error) {
	priv, err := l.generateKey(true)
	if err != nil {
		return nil, nil, err
	}
	pub := priv.(crypto.Signer).Public()
	spkiASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}
	var spki struct {
		Algorithm        pkix.AlgorithmIdentifier
		SubjectPublicKey asn1.BitString
	}
	_, err = asn1.Unmarshal(spkiASN1, &spki)
	if err != nil {
		return nil, nil, err
	}
	serialNumber, err := l.randomSerialNumber()
	if err != nil {
		return nil, nil, err
	}
	skid := sha1.Sum(spki.SubjectPublicKey.Bytes)
	tpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{organization},
			OrganizationalUnit: []string{l.getUserAndHostname()},
			// The CommonName is required by iOS to show the certificate in the
			// "Certificate Trust Settings" menu.
			CommonName: "tlsctl " + l.getUserAndHostname(),
		},
		SubjectKeyId: skid[:],

		NotAfter:  time.Now().AddDate(10, 0, 0),
		NotBefore: time.Now(),

		KeyUsage: x509.KeyUsageCertSign,

		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, tpl, tpl, pub, priv)
	if err != nil {
		return nil, nil, err
	}
	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}
	certificateMemory = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
	keyMemory = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})
	return certificateMemory, keyMemory, nil
}
func (l *SSL) generateKey(rootCA bool) (crypto.PrivateKey, error) {
	if rootCA {
		return rsa.GenerateKey(rand.Reader, 3072)
	}
	return rsa.GenerateKey(rand.Reader, 2048)
}
func (l *SSL) randomSerialNumber() (*big.Int, error) {
	return rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
}
func (l *SSL) getUserAndHostname() string {
	var userAndHostname string
	u, err := user.Current()
	if err == nil {
		userAndHostname = u.Username + "@"
	}
	if h, err := os.Hostname(); err == nil {
		userAndHostname += h
	}
	if err == nil && u.Name != "" && u.Name != u.Username {
		userAndHostname += " (" + u.Name + ")"
	}
	return userAndHostname
}
