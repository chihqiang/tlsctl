package resource

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chihqiang/tlsctl/pkg/certificates"
	"github.com/chihqiang/tlsctl/pkg/fp"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"software.sslmate.com/src/go-pkcs12"
)

const (
	baseCertificatesFolderName = "certificates"

	IssuerExt   = ".issuer.crt"
	CertExt     = ".crt"
	KeyExt      = ".key"
	PemExt      = ".pem"
	PfxExt      = ".pfx"
	ResourceExt = ".json"
)

type Cache struct {
	rootPath    string
	pfxFormat   string //"DES", "RC2", "SHA256"
	pfxPassword string
}

func NewCache(path, pfxFormat string) (s *Cache, err error) {
	switch pfxFormat {
	case "DES", "RC2", "SHA256":
	default:
		err = fmt.Errorf("invalid PFX format: %s", pfxFormat)
	}
	if err != nil {
		return nil, err
	}
	rootPath := filepath.Join(path, baseCertificatesFolderName)
	if err = fp.CreateNonExistingFolder(rootPath); err != nil {
		return nil, err
	}
	return &Cache{
		rootPath:  rootPath,
		pfxFormat: pfxFormat,
	}, nil
}

func (s *Cache) SaveResource(resource *certificate.Resource) error {
	domain := resource.Domain
	sanitizedDomain, err := SanitizedDomain(domain)
	if err != nil {
		return err
	}
	if err := s.createPath(sanitizedDomain); err != nil {
		return err
	}
	if err := s.writeFile(sanitizedDomain, CertExt, resource.Certificate); err != nil {
		return fmt.Errorf("Unable to save Certificate for domain %s\n\t%v", domain, err)
	}
	if err := s.writeFile(sanitizedDomain, IssuerExt, resource.IssuerCertificate); err != nil {
		return fmt.Errorf("Unable to save IssuerCertificate for domain  %s\n\t%v", domain, err)
	}
	if resource.PrivateKey != nil {
		if err := s.writeCertificateFiles(sanitizedDomain, resource); err != nil {
			return fmt.Errorf("Unable to save PrivateKey for domain %s\n\t%v", domain, err)
		}
	}
	jsonBytes, err := json.MarshalIndent(resource, "", "\t")
	if err != nil {
		return fmt.Errorf("Unable to marshal CertResource for domain %s\n\t%v", domain, err)
	}
	return s.writeFile(sanitizedDomain, ResourceExt, jsonBytes)
}
func (s *Cache) GetAllDomainResources() ([]*certificate.Resource, error) {
	var resources []*certificate.Resource
	jsons, _ := filepath.Glob(filepath.Join(s.rootPath, "/*/*"+ResourceExt))
	for _, jsonFile := range jsons {
		var tmpResource *certificate.Resource
		jsonBytes, err := os.ReadFile(jsonFile)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(jsonBytes, &tmpResource); err != nil {
			continue
		}
		resource, err := s.ReadResource(tmpResource.Domain)
		if err != nil {
			continue
		}
		resources = append(resources, resource)
	}
	return resources, nil
}
func (s *Cache) ReadResource(domain string) (*certificate.Resource, error) {
	sanitizedDomain, err := SanitizedDomain(domain)
	if err != nil {
		return nil, err
	}
	raw, err := s.readFile(sanitizedDomain, ResourceExt)
	if err != nil {
		return nil, fmt.Errorf("Error while loading the meta data for domain %s\n\t%v", domain, err)
	}
	var resource certificate.Resource
	if err = json.Unmarshal(raw, &resource); err != nil {
		return nil, fmt.Errorf("Error while marshaling the meta data for domain %s\n\t%v", domain, err)
	}
	resource.PrivateKey, err = s.readFile(sanitizedDomain, KeyExt)
	if err != nil {
		return nil, fmt.Errorf("load key for domain %s\n\t%v", domain, err)
	}
	resource.IssuerCertificate, err = s.readFile(sanitizedDomain, IssuerExt)
	if err != nil {
		return nil, fmt.Errorf("load IssuerCertificate for domain %s\n\t%v", domain, err)
	}
	resource.Certificate, err = s.readFile(sanitizedDomain, CertExt)
	if err != nil {
		return nil, fmt.Errorf("load Certificate for domain %s\n\t%v", domain, err)
	}
	return &resource, nil
}

func (s *Cache) ParseResourceFindCertificate(resource *certificate.Resource) (*x509.Certificate, error) {
	bundle, err := certcrypto.ParsePEMBundle(resource.Certificate)
	if err != nil {
		return nil, err
	}
	for _, x := range bundle {
		if x.Subject.CommonName == resource.Domain {
			return x, nil
		}
	}
	return nil, errors.New("Certificate for domain " + resource.Domain + " not found")
}

func (s *Cache) writeCertificateFiles(domain string, certRes *certificate.Resource) error {
	sanitizedDomain, err := SanitizedDomain(domain)
	if err != nil {
		return err
	}
	if err = s.writeFile(sanitizedDomain, KeyExt, certRes.PrivateKey); err != nil {
		return fmt.Errorf("unable to save key file: %w", err)
	}
	if err = s.writeFile(sanitizedDomain, PemExt, bytes.Join([][]byte{certRes.Certificate, certRes.PrivateKey}, nil)); err != nil {
		return fmt.Errorf("unable to save PEM file: %w", err)
	}
	if err = s.writePFXFile(sanitizedDomain, certRes); err != nil {
		return fmt.Errorf("unable to save PFX file: %w", err)
	}
	return nil
}
func (s *Cache) writePFXFile(domain string, certRes *certificate.Resource) error {
	certPemBlock, _ := pem.Decode(certRes.Certificate)
	if certPemBlock == nil {
		return fmt.Errorf("unable to parse Certificate for domain %s", domain)
	}
	cert, err := x509.ParseCertificate(certPemBlock.Bytes)
	if err != nil {
		return fmt.Errorf("unable to load Certificate for domain %s: %w", domain, err)
	}
	certChain, err := s.getCertificateChain(certRes)
	if err != nil {
		return fmt.Errorf("unable to get certificate chain for domain %s: %w", domain, err)
	}
	privateKey, err := certificates.ParseX509(certRes.PrivateKey)
	if err != nil {
		return fmt.Errorf("unable to parse private key for domain %s: %w", domain, err)
	}
	encoder, err := s.getPFXEncoder(s.pfxFormat)
	if err != nil {
		return fmt.Errorf("PFX encoder: %w", err)
	}
	pfxBytes, err := encoder.Encode(privateKey, cert, certChain, s.pfxPassword)
	if err != nil {
		return fmt.Errorf("unable to encode PFX data for domain %s: %w", domain, err)
	}
	return s.writeFile(domain, PfxExt, pfxBytes)
}
func (s *Cache) getPFXEncoder(pfxFormat string) (*pkcs12.Encoder, error) {
	var encoder *pkcs12.Encoder
	switch pfxFormat {
	case "SHA256":
		encoder = pkcs12.Modern2023
	case "DES":
		encoder = pkcs12.LegacyDES
	case "RC2":
		encoder = pkcs12.LegacyRC2
	default:
		return nil, fmt.Errorf("invalid PFX format: %s", pfxFormat)
	}
	return encoder, nil
}
func (s *Cache) getCertificateChain(certRes *certificate.Resource) ([]*x509.Certificate, error) {
	chainCertPemBlock, rest := pem.Decode(certRes.IssuerCertificate)
	if chainCertPemBlock == nil {
		return nil, errors.New("unable to parse Issuer Certificate")
	}
	var certChain []*x509.Certificate
	for chainCertPemBlock != nil {
		chainCert, err := x509.ParseCertificate(chainCertPemBlock.Bytes)
		if err != nil {
			return nil, fmt.Errorf("unable to parse Chain Certificate: %w", err)
		}

		certChain = append(certChain, chainCert)
		chainCertPemBlock, rest = pem.Decode(rest) // Try decoding the next pem block
	}
	return certChain, nil
}

func (s *Cache) createPath(sanitizedDomain string) error {
	return fp.CreateNonExistingFolder(filepath.Join(s.rootPath, sanitizedDomain))
}
func (s *Cache) readFile(sanitizedDomain, extension string) ([]byte, error) {
	return os.ReadFile(s.sanitizedDomainFileName(sanitizedDomain, extension))
}
func (s *Cache) writeFile(sanitizedDomain, extension string, data []byte) error {
	return os.WriteFile(s.sanitizedDomainFileName(sanitizedDomain, extension), data, 0o600)
}
func (s *Cache) sanitizedDomainFileName(sanitizedDomain, extension string) string {
	return filepath.Join(s.GetSanitizedDomainSavePath(sanitizedDomain), sanitizedDomain+extension)
}
func (s *Cache) GetSanitizedDomainSavePath(sanitizedDomain string) string {
	return filepath.Join(s.rootPath, sanitizedDomain)
}
