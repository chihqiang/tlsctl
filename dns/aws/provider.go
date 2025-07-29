package aws

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns/route53"
	"time"
)

type Config struct {
	AccessKeyId        string `json:"access_key_id" yaml:"accessKeyId" xml:"accessKeyId" env:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey    string `json:"secret_access_key" yaml:"secretAccessKey" xml:"secretAccessKey" env:"AWS_SECRET_ACCESS_KEY"`
	Region             string `json:"region" yaml:"region" xml:"region" env:"AWS_REGION"`
	HostedZoneId       string `json:"hosted_zone_id" yaml:"hostedZoneId" xml:"hostedZoneId" env:"AWS_HOSTED_ZONE_ID"`
	PropagationTimeout int32  `json:"propagation_timeout,omitempty" yaml:"propagationTimeout" env:"AWS_PROPAGATION_TIMEOUT"`
	TTL                int32  `json:"ttl,omitempty" yaml:"ttl" env:"AWS_TTL"`
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

	providerConfig := route53.NewDefaultConfig()
	providerConfig.AccessKeyID = config.AccessKeyId
	providerConfig.SecretAccessKey = config.SecretAccessKey
	providerConfig.Region = config.Region
	providerConfig.HostedZoneID = config.HostedZoneId
	if config.PropagationTimeout != 0 {
		providerConfig.PropagationTimeout = time.Duration(config.PropagationTimeout) * time.Second
	}
	if config.TTL != 0 {
		providerConfig.TTL = int(config.TTL)
	}
	provider, err := route53.NewDNSProviderConfig(providerConfig)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
