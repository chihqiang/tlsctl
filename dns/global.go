package dns

import (
	"github.com/chihqiang/tlsctl/dns/aliyun"
	"github.com/chihqiang/tlsctl/dns/aws"
	"github.com/chihqiang/tlsctl/dns/baidu"
	"github.com/chihqiang/tlsctl/dns/cloudflare"
	"github.com/chihqiang/tlsctl/dns/dynv6"
	"github.com/chihqiang/tlsctl/dns/godaddy"
	"github.com/chihqiang/tlsctl/dns/huawei"
	"github.com/chihqiang/tlsctl/dns/jdcloud"
	"github.com/chihqiang/tlsctl/dns/tencentcloud"
	"github.com/chihqiang/tlsctl/dns/westcn"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns"
)

var providers = map[string]func() IDNSProvider{
	"aliyun":       func() IDNSProvider { return &aliyun.Provider{} },
	"aws":          func() IDNSProvider { return &aws.Provider{} },
	"godaddy":      func() IDNSProvider { return &godaddy.Provider{} },
	"huawei":       func() IDNSProvider { return &huawei.Provider{} },
	"tencentcloud": func() IDNSProvider { return &tencentcloud.Provider{} },
	"westcn":       func() IDNSProvider { return &westcn.Provider{} },
	"jdcloud":      func() IDNSProvider { return &jdcloud.Provider{} },
	"baidu":        func() IDNSProvider { return &baidu.Provider{} },
	"cloudflare":   func() IDNSProvider { return &cloudflare.Provider{} },
	"dynv6":        func() IDNSProvider { return &dynv6.Provider{} },
}

func Register(name string, provider func() IDNSProvider) {
	providers[name] = provider
}
func Get(name string) (challenge.Provider, error) {
	if providerFun, ok := providers[name]; ok {
		provider := providerFun()
		if err := provider.WithEnvConfig(); err != nil {
			return nil, err
		}
		return provider.NewChallengeProvider()
	}
	return dns.NewDNSChallengeProviderByName(name)
}
func All() map[string]IDNSProvider {
	ds := providers
	var m = make(map[string]IDNSProvider)
	for s, f := range ds {
		m[s] = f()
	}
	return m
}
