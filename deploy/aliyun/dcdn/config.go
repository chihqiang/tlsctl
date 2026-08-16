package dcdn

type Config struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"access_key_id,omitempty" yaml:"AccessKeyId" xml:"AccessKeyId" env:"ALIYUN_ACCESS_KEY_ID"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"access_key_secret,omitempty" yaml:"AccessKeySecret" xml:"AccessKeySecret" env:"ALIYUN_ACCESS_KEY_SECRET"`
	// 加速域名（支持泛域名）。
	Domain string `json:"domain,omitempty" yaml:"Domain" xml:"Domain" env:"ALIYUN_DOMAIN"`
}
