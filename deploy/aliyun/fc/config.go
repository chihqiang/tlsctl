package fc

type Config struct {
	// 阿里云 AccessKeyId。
	AccessKeyId string `json:"access_key_id,omitempty" yaml:"AccessKeyId" xml:"AccessKeyId" env:"ALIYUN_ACCESS_KEY_ID"`
	// 阿里云 AccessKeySecret。
	AccessKeySecret string `json:"access_key_secret,omitempty" yaml:"AccessKeySecret" xml:"AccessKeySecret" env:"ALIYUN_ACCESS_KEY_SECRET"`
	// 阿里云资源组 ID。
	ResourceGroupId string `json:"resource_group_id,omitempty" yaml:"ResourceGroupId" xml:"ResourceGroupId" env:"ALIYUN_RESOURCE_GROUP_ID"`
	// 阿里云地域。
	Region string `json:"region,omitempty" yaml:"Region" xml:"Region" env:"ALIYUN_REGION"`
	// 服务版本。可取值 "2.0"、"3.0"。
	Version string `json:"version,omitempty" yaml:"Version" xml:"Version" env:"ALIYUN_VERSION"`
	// 自定义域名（支持泛域名）。
	Domain string `json:"domain,omitempty" yaml:"Domain" xml:"Domain" env:"ALIYUN_DOMAIN"`
}
