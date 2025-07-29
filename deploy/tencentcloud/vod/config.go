package vod

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id,omitempty" yaml:"SecretId" xml:"SecretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key,omitempty" yaml:"SecretKey" xml:"SecretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 点播应用 ID。
	SubAppId int64 `json:"sub_app_id,omitempty" yaml:"SubAppId" xml:"SubAppId" env:"TENCENTCLOUD_SUB_APP_ID"`
	// 点播加速域名（不支持泛域名）。
	Domain string `json:"domain,omitempty" yaml:"Domain" xml:"Domain" env:"TENCENTCLOUD_DOMAIN"`
}
