package cdn

import "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"

type Config struct {
	common.BaseConfig
	// 加速域名（支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"TENCENTCLOUD_DOMAIN"`
}
