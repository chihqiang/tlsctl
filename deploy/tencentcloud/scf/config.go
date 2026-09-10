package scf

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id" yaml:"secretId" xml:"secretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 腾讯云地域。
	Region string `json:"region" yaml:"region" xml:"region" env:"TENCENTCLOUD_REGION"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"TENCENTCLOUD_DOMAIN"`
}
