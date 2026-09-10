package common

type BaseConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secret_id" yaml:"secretId" xml:"secretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secret_key" yaml:"secretKey" xml:"secretKey" env:"TENCENTCLOUD_SECRET_KEY"`
}