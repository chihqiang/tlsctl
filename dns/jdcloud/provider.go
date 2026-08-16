package jdcloud

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/dns/jdcloud/lego"
	"github.com/go-acme/lego/v4/challenge"
)

type Config struct {
	AccessKeyId        string `json:"access_key_id" yaml:"AccessKeyId" xml:"AccessKeyId" env:"JDCLOUD_ACCESS_KEY_ID"`
	AccessKeySecret    string `json:"access_key_secret" yaml:"AccessKeySecret" xml:"AccessKeySecret" env:"JDCLOUD_ACCESS_KEY_SECRET"`
	RegionId           string `json:"region_id" yaml:"regionId" xml:"RegionId" env:"JDCLOUD_REGION_ID"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagation_timeout,omitempty" env:"JDCLOUD_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" yaml:"ttl,omitempty" xml:"TTL,omitempty" env:"JDCLOUD_TTL"`
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
	regionId := config.RegionId
	if regionId == "" {
		// 京东云的 SDK 要求必须传一个区域，实际上 DNS-01 流程里用不到，但不传会报错
		regionId = "cn-north-1"
	}
	providerConfig := lego.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.AccessKeySecret = config.AccessKeySecret
	providerConfig.RegionId = regionId
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = config.TTL
	}
	return lego.NewDNSProviderConfig(providerConfig)
}
