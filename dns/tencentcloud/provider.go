package tencentcloud

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/tencentcloud"
	"time"
)

type Config struct {
	SecretId           string `json:"secret_id" yaml:"secretId" xml:"secretId" env:"TENCENTCLOUD_SECRET_ID"`
	SecretKey          string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout" xml:"propagationTimeout" env:"TENCENTCLOUD_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" yaml:"ttl" xml:"ttl" env:"TENCENTCLOUD_TTL"`
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
	providerConfig := tencentcloud.NewDefaultConfig()
	providerConfig.SecretID = config.SecretId
	providerConfig.SecretKey = config.SecretKey
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = int(config.TTL)
	}

	return tencentcloud.NewDNSProviderConfig(providerConfig)

}
