package eo

import "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"

type Config struct {
	common.BaseConfig
	// 站点 ID。
	ZoneId string `json:"zone_id" yaml:"zoneId" xml:"zoneId" env:"TENCENTCLOUD_ZONE_ID"`
	// 防护域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"TENCENTCLOUD_DOMAIN"`
}
