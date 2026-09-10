package onepanel

type Config struct {
	// 1Panel API 地址，如 https://example.com:8090。
	Url string `json:"url" yaml:"url" xml:"url" env:"ONEPANEL_URL"`
	// 1Panel API 密钥。
	ApiKey string `json:"api_key" yaml:"apiKey" xml:"apiKey" env:"ONEPANEL_API_KEY"`
	// 1Panel 版本：v1 或 v2，默认 v1。
	Version string `json:"version" yaml:"version" xml:"version" env:"ONEPANEL_VERSION"`
	// 是否跳过 TLS 证书校验。
	IgnoreSSL bool `json:"ignore_ssl" yaml:"ignoreSsl" xml:"ignoreSsl" env:"ONEPANEL_IGNORE_SSL"`
	// 站点 ID（仅 1panel-site 使用）。
	SiteId string `json:"site_id" yaml:"siteId" xml:"siteId" env:"ONEPANEL_SITE_ID"`
}
