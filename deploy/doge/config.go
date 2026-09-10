package doge

type Config struct {
	// 多吉云 AccessKey。
	AccessKey string `json:"access_key" yaml:"accessKey" xml:"accessKey" env:"DOGE_ACCESS_KEY"`
	// 多吉云 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"DOGE_SECRET_KEY"`
	// CDN 加速域名，留空时使用证书主域名。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"DOGE_DOMAIN"`
}
