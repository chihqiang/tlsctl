package ecdn

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId" yaml:"secretId" xml:"SecretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey" yaml:"secretKey" xml:"SecretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"Domain" env:"TENCENTCLOUD_DOMAIN"`
}
