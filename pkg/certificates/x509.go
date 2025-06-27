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
	default:
		return nil, fmt.Errorf("unsupported key type: %s", keyPemBlock.Type)
	}
}
