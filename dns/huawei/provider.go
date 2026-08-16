package huawei

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/huaweicloud"
	"time"
)

type Config struct {
	AccessKeyId        string `json:"access_key_id" yaml:"AccessKeyID" xml:"AccessKeyID" env:"HUAWEICLOUD_ACCESS_KEY_ID"`
	SecretAccessKey    string `json:"secret_access_key" yaml:"SecretAccessKey" xml:"HUAWEICLOUD_SECRET_ACCESS_KEY"`
	Region             string `json:"region" yaml:"Region" xml:"Region" env:"HUAWEICLOUD_REGION"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" env:"HUAWEICLOUD_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" env:"HUAWEICLOUD_TTL"`
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
	region := config.Region
	if region == "" {
		// 华为云的 SDK 要求必须传一个区域，实际上 DNS-01 流程里用不到，但不传会报错
		region = "cn-north-1"
	}

	providerConfig := huaweicloud.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.SecretAccessKey = config.SecretAccessKey
	providerConfig.Region = region
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = config.TTL
	}
	return huaweicloud.NewDNSProviderConfig(providerConfig)
}
