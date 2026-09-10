package vod

import "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"

type Config struct {
	common.BaseConfig
	// 点播应用 ID。
	SubAppId int64 `json:"sub_app_id" yaml:"subAppId" xml:"subAppId" env:"TENCENTCLOUD_SUB_APP_ID"`
	// 点播加速域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"TENCENTCLOUD_DOMAIN"`
}
