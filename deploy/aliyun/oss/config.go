package oss

type Config struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"access_key_id,omitempty" yaml:"AccessKeyId" xml:"AccessKeyId" env:"ALIYUN_ACCESS_KEY_ID"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"access_key_secret,omitempty" yaml:"AccessKeySecret" xml:"AccessKeySecret" env:"ALIYUN_ACCESS_KEY_SECRET"`
	// 阿里云地域。
	Region string `json:"region,omitempty" yaml:"Region" xml:"Region" env:"ALIYUN_REGION"`
	// 存储桶名。
	Bucket string `json:"bucket,omitempty" yaml:"Bucket" xml:"Bucket" env:"ALIYUN_BUCKET"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain,omitempty" yaml:"Domain" xml:"Domain" env:"ALIYUN_DOMAIN"`
}
