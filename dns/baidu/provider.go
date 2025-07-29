package baidu

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/dns/baidu/lego"
	"github.com/go-acme/lego/v4/challenge"
)

type Config struct {
	AccessKeyId        string `json:"access_key_id" yaml:"accessKeyId" xml:"accessKeyId" env:"BAIDUCLOUD_ACCESS_KEY_ID"`
	SecretAccessKey    string `json:"secret_access_key" yaml:"secretAccessKey" xml:"secretAccessKey" env:"BAIDUCLOUD_SECRET_ACCESS_KEY"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout,omitempty" env:"BAIDUCLOUD_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" yaml:"ttl,omitempty" env:"BAIDUCLOUD_TTL"`
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
	providerConfig := lego.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.SecretAccessKey = config.SecretAccessKey
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = config.TTL
	}
	return lego.NewDNSProviderConfig(providerConfig)
}
