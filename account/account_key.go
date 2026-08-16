package account

import (
	"errors"
	"github.com/go-acme/lego/v4/certcrypto"
	"strings"
)

// GetKeyType the type from which private keys should be generated.
func GetKeyType(keyType string) (certcrypto.KeyType, error) {
	switch strings.ToUpper(keyType) {
	case "RSA2048":
		return certcrypto.RSA2048, nil
	case "RSA3072":
		return certcrypto.RSA3072, nil
	case "RSA4096":
		return certcrypto.RSA4096, nil
	case "RSA8192":
		return certcrypto.RSA8192, nil
	case "EC256":
		return certcrypto.EC256, nil
	case "EC384":
		return certcrypto.EC384, nil
	}
	return "", errors.New("unknown key type")
}
