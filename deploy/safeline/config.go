package safeline

type Config struct {
	// 雷池 WAF API 地址，如 https://example.com:9443。
	Url string `json:"url" yaml:"url" xml:"url" env:"SAFELINE_URL"`
	// 雷池 WAF API Token。
	ApiToken string `json:"api_token" yaml:"apiToken" xml:"apiToken" env:"SAFELINE_API_TOKEN"`
	// 是否跳过 TLS 证书校验。
	IgnoreSSL bool `json:"ignore_ssl" yaml:"ignoreSsl" xml:"ignoreSsl" env:"SAFELINE_IGNORE_SSL"`
	// 站点名称（仅 safeline-site 使用）。
	SiteName string `json:"site_name" yaml:"siteName" xml:"siteName" env:"SAFELINE_SITE_NAME"`
}
