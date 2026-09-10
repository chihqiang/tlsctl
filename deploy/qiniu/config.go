package qiniu

type Config struct {
	// 七牛云 AccessKey。
	AccessKey string `json:"access_key" yaml:"accessKey" xml:"accessKey" env:"QINIU_ACCESS_KEY"`
	// 七牛云 AccessSecret。
	AccessSecret string `json:"access_secret" yaml:"accessSecret" xml:"accessSecret" env:"QINIU_ACCESS_SECRET"`
	// 加速域名，留空时使用证书主域名。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"QINIU_DOMAIN"`
}
