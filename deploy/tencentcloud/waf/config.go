package waf

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id,omitempty" yaml:"SecretId" xml:"SecretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key,omitempty" yaml:"SecretKey" xml:"SecretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 腾讯云地域。
	Region string `json:"region" yaml:"Region" xml:"Region" env:"TENCENTCLOUD_REGION"`
	// 防护域名（不支持泛域名）。
	Domain string `json:"domain" yaml:"Domain" xml:"Domain" env:"TENCENTCLOUD_DOMAIN"`
	// 防护域名 ID。
	DomainId string `json:"domainId" yaml:"DomainId" xml:"DomainId" env:"TENCENTCLOUD_DOMAIN_ID"`
	// 防护域名所属实例 ID。
	InstanceId string `json:"instanceId" yaml:"InstanceId" xml:"InstanceId" env:"TENCENTCLOUD_INSTANCE_ID"`
}
