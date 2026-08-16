package certificates

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ParseX509(privateKey []byte) (signer crypto.Signer, err error) {
	keyPemBlock, _ := pem.Decode(privateKey)
	if keyPemBlock == nil {
		return nil, fmt.Errorf("unable to parse private key")
	}
	switch keyPemBlock.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(keyPemBlock.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(keyPemBlock.Bytes)
	case "PRIVATE KEY": // PKCS#8，如 localhost CA 私钥
		key, err := x509.ParsePKCS8PrivateKey(keyPemBlock.Bytes)
		if err != nil {
			return nil, err
		}
		signer, ok := key.(crypto.Signer)
		if !ok {
			return nil, fmt.Errorf("private key does not implement crypto.Signer")
		}
		return signer, nil
	default:
		return nil, fmt.Errorf("unsupported key type: %s", keyPemBlock.Type)
	}
}
