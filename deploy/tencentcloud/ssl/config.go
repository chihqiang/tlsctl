package ssl

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id,omitempty" yaml:"SecretId" xml:"SecretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key,omitempty" yaml:"SecretKey" xml:"SecretKey" env:"TENCENTCLOUD_SECRET_KEY"`
}
