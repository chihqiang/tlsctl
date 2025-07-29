package westcn

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/westcn"
	"time"
)

type Config struct {
	Username           string `json:"username" yaml:"username" xml:"username" env:"WESTCN_USERNAME"`
	Password           string `json:"password" yaml:"password" xml:"password" env:"WESTCN_PASSWORD"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout,omitempty" xml:"propagationTimeout" env:"WESTCN_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" yaml:"ttl,omitempty" xml:"ttl,omitempty" env:"WESTCN_TTL"`
}

type Provider struct {
	Config *Config
}

func (d *Provider) WithEnvConfig() error {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	d.Config = &cfg
	return nil
}
func (p *Provider) NewChallengeProvider() (challenge.Provider, error) {
	config := p.Config
	providerConfig := westcn.NewDefaultConfig()
	providerConfig.Username = config.Username
	providerConfig.Password = config.Password
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = int(config.TTL)
	}

	return westcn.NewDNSProviderConfig(providerConfig)
}
