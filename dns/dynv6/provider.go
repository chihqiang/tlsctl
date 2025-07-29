package dynv6

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/dns/dynv6/lego"
	"github.com/go-acme/lego/v4/challenge"
)

type Config struct {
	HttpToken          string `json:"http_token" yaml:"HttpToken" xml:"HttpToken" env:"DYNV6_HTTP_TOKEN"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout" xml:"propagationTimeout" env:"DYNV6_PROPAGATION_TIMEOUT"`
	TTL                int    `json:"ttl,omitempty" yaml:"ttl" xml:"ttl" env:"DYNV6_TTL"`
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
	providerConfig := lego.NewDefaultConfig()
	providerConfig.HTTPToken = p.Config.HttpToken
	if p.Config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(p.Config.PropagationTimeout) * time.Second
	}
	if p.Config.TTL != 0 {
		providerConfig.TTL = int(p.Config.TTL)
	}
	return lego.NewDNSProviderConfig(providerConfig)
}
