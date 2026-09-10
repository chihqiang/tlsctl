package huaweicloud

type Config struct {
	// 华为云 AccessKey。
	AccessKey string `json:"access_key" yaml:"accessKey" xml:"accessKey" env:"HUAWEI_ACCESS_KEY"`
	// 华为云 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"HUAWEI_SECRET_KEY"`
	// CDN 加速域名，留空时使用证书主域名。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"HUAWEI_DOMAIN"`
}
