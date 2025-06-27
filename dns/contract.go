package dns

import "github.com/go-acme/lego/v4/challenge"

type IDNSProvider interface {
	WithEnvConfig() error
	NewChallengeProvider() (challenge.Provider, error)
}
