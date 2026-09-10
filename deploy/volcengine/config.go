package volcengine

type Config struct {
	// 火山引擎 AccessKey。
	AccessKey string `json:"access_key" yaml:"accessKey" xml:"accessKey" env:"VOLC_ACCESS_KEY"`
	// 火山引擎 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"VOLC_SECRET_KEY"`
	// 地域，如 cn-north-1。
	Region string `json:"region" yaml:"region" xml:"region" env:"VOLC_REGION"`
	// CDN/DCDN 加速域名，留空时使用证书主域名。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"VOLC_DOMAIN"`
}
