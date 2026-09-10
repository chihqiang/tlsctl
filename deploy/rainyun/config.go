package rainyun

type Config struct {
	// 雨云 API Key。
	ApiKey string `json:"api_key" yaml:"apiKey" xml:"apiKey" env:"RAINYUN_API_KEY"`
	// 雨云 SSL 证书 ID。
	CertId string `json:"cert_id" yaml:"certId" xml:"certId" env:"RAINYUN_CERT_ID"`
}
