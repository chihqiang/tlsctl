package aliyun

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/alidns"
	"time"
)

type Config struct {
	AccessKeyId        string `json:"access_key_id" yaml:"accessKeyId" xml:"accessKeyId" env:"ALICLOUD_ACCESS_KEY"`
	AccessKeySecret    string `json:"access_key_secret" yaml:"accessKeySecret" xml:"accessKeySecret" env:"ALICLOUD_SECRET_KEY"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout" env:"ALICLOUD_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" yaml:"ttl" xml:"ttl" env:"ALICLOUD_TTL"`
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
	providerConfig := alidns.NewDefaultConfig()
	providerConfig.APIKey = config.AccessKeyId
	providerConfig.SecretKey = config.AccessKeySecret
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = int(config.TTL)
	}
	provider, err := alidns.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
