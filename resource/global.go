package resource

import (
	"golang.org/x/net/idna"
	"strings"
)

func SanitizedDomain(domain string) (string, error) {
	safe, err := idna.ToASCII(strings.NewReplacer(":", "-", "*", "_").Replace(domain))
	if err != nil {
		return "", err
	}
	return safe, nil
}
