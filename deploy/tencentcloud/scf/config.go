package scf

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id,omitempty" yaml:"SecretId" xml:"SecretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key,omitempty" yaml:"SecretKey" xml:"SecretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 腾讯云地域。
	Region string `json:"region,omitempty" yaml:"Region" xml:"Region" env:"TENCENTCLOUD_REGION"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain,omitempty" yaml:"Domain" xml:"Domain" env:"TENCENTCLOUD_DOMAIN"`
}
