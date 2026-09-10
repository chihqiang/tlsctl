package scf

import "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"

type Config struct {
	common.BaseConfig
	// 腾讯云地域。
	Region string `json:"region" yaml:"region" xml:"region" env:"TENCENTCLOUD_REGION"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"TENCENTCLOUD_DOMAIN"`
}
