package waf

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id" yaml:"secretId" xml:"secretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 腾讯云地域。
	Region string `json:"region" yaml:"region" xml:"region" env:"TENCENTCLOUD_REGION"`
	// 防护域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"TENCENTCLOUD_DOMAIN"`
	// 防护域名 ID。
	DomainId string `json:"domain_id" yaml:"domainId" xml:"domainId" env:"TENCENTCLOUD_DOMAIN_ID"`
	// 防护域名所属实例 ID。
	InstanceId string `json:"instance_id" yaml:"instanceId" xml:"instanceId" env:"TENCENTCLOUD_INSTANCE_ID"`
}
