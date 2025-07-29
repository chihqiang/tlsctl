package cloudflare

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"time"
)

type Config struct {
	ApiToken           string `json:"api_token" yaml:"ApiToken" xml:"ApiToken" env:"CLOUDFLARE_API_TOKEN"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout" xml:"propagationTimeout" env:"CLOUDFLARE_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" yaml:"ttl" xml:"ttl" env:"CLOUDFLARE_TTL"`
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
	providerConfig := cloudflare.NewDefaultConfig()
	providerConfig.AuthToken = config.ApiToken
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = int(config.TTL)
	}
	return cloudflare.NewDNSProviderConfig(providerConfig)
}
