package vod

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id" yaml:"secretId" xml:"secretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 点播应用 ID。
	SubAppId int64 `json:"sub_app_id" yaml:"subAppId" xml:"subAppId" env:"TENCENTCLOUD_SUB_APP_ID"`
	// 点播加速域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"TENCENTCLOUD_DOMAIN"`
}
