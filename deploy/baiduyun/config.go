package baiduyun

type Config struct {
	// 百度智能云 AccessKey。
	AccessKey string `json:"access_key" yaml:"accessKey" xml:"accessKey" env:"BAIDU_ACCESS_KEY"`
	// 百度智能云 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"BAIDU_SECRET_KEY"`
	// 加速域名，留空时使用证书主域名。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"BAIDU_DOMAIN"`
}
