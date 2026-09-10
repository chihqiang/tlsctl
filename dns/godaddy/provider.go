package godaddy

import (
	"time"

	"github.com/chihqiang/tlsctl/pkg/envconfig"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/godaddy"
)

type Config struct {
	ApiKey             string `json:"api_key" yaml:"apiKey" xml:"apiKey" env:"GODADDY_API_KEY"`
	ApiSecret          string `json:"api_secret" yaml:"apiSecret" xml:"apiSecret" env:"GODADDY_API_SECRET"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout" xml:"propagationTimeout" env:"GODADDY_PROPAGATION_TIMEOUT"`
	TTL                int    `json:"ttl,omitempty" yaml:"ttl" xml:"ttl" env:"GODADDY_TTL"`
}

type Provider struct {
	Config *Config
}

func (p *Provider) WithEnvConfig() error {
	cfg, err := envconfig.ParseConfig[Config]()
	if err != nil {
		return err
	}
	p.Config = cfg
	return nil
}
func (p *Provider) NewChallengeProvider() (challenge.Provider, error) {
	config := p.Config
	providerConfig := godaddy.NewDefaultConfig()
	providerConfig.APIKey = config.ApiKey
	providerConfig.APISecret = config.ApiSecret
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = config.TTL
	}

	return godaddy.NewDNSProviderConfig(providerConfig)
}
