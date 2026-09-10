package btpanel

type Config struct {
	// 宝塔面板 API 地址，如 https://example.com:8888。
	Url string `json:"url" yaml:"url" xml:"url" env:"BTPANEL_URL"`
	// 宝塔面板 API 密钥。
	ApiKey string `json:"api_key" yaml:"apiKey" xml:"apiKey" env:"BTPANEL_API_KEY"`
	// 是否跳过 TLS 证书校验。
	IgnoreSSL bool `json:"ignore_ssl" yaml:"ignoreSsl" xml:"ignoreSsl" env:"BTPANEL_IGNORE_SSL"`
	// 目标网站名称，多个用英文逗号分隔（仅 btpanel-site 使用）。
	SiteName string `json:"site_name" yaml:"siteName" xml:"siteName" env:"BTPANEL_SITE_NAME"`
}
