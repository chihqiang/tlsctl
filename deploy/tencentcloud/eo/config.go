package eo

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id,omitempty" yaml:"SecretId" xml:"SecretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key,omitempty" yaml:"SecretKey" xml:"SecretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 站点 ID。
	ZoneId string `json:"zoneId" yaml:"ZoneId" xml:"ZoneId" env:"TENCENTCLOUD_ZONE_ID"`
	// 防护域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"Domain" xml:"Domain" env:"TENCENTCLOUD_DOMAIN"`
}
