package dns

import (
	"wangzhiqiang/tlsctl/dns/aliyun"
	"wangzhiqiang/tlsctl/dns/aws"
	"wangzhiqiang/tlsctl/dns/baidu"
	"wangzhiqiang/tlsctl/dns/cloudflare"
	"wangzhiqiang/tlsctl/dns/dynv6"
	"wangzhiqiang/tlsctl/dns/godaddy"
	"wangzhiqiang/tlsctl/dns/huawei"
	"wangzhiqiang/tlsctl/dns/jdcloud"
	"wangzhiqiang/tlsctl/dns/tencentcloud"
	"wangzhiqiang/tlsctl/dns/westcn"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns"
)

var providers = map[string]IDNSProvider{
	"aliyun":       &aliyun.Provider{},
	"aws":          &aws.Provider{},
	"godaddy":      &godaddy.Provider{},
	"huawei":       &huawei.Provider{},
	"tencentcloud": &tencentcloud.Provider{},
	"westcn":       &westcn.Provider{},
	"jdcloud":      &jdcloud.Provider{},
	"baidu":        &baidu.Provider{},
	"cloudflare":   &cloudflare.Provider{},
	"dynv6":        &dynv6.Provider{},
}

func Register(name string, provider IDNSProvider) {
	providers[name] = provider
}
func Get(name string) (challenge.Provider, error) {
	if provider, ok := providers[name]; ok {
		if err := provider.WithEnvConfig(); err != nil {
			return nil, err
		}
		return provider.NewChallengeProvider()
	}
	return dns.NewDNSChallengeProviderByName(name)
}
func All() map[string]IDNSProvider {
	return providers
}
