package webhook

type Config struct {
	// Webhook 请求 URL。
	Url string `json:"url" yaml:"url" xml:"url" env:"WEBHOOK_URL"`
	// Webhook 请求方法，默认 POST。
	Method string `json:"method" yaml:"method" xml:"method" env:"WEBHOOK_METHOD"`
	// 请求体模板，支持 __domain__ / __cert__ / __key__ 占位符。
	Data string `json:"data" yaml:"data" xml:"data" env:"WEBHOOK_DATA"`
	// 请求头，每行一个 "Key: Value"。
	Headers string `json:"headers" yaml:"headers" xml:"headers" env:"WEBHOOK_HEADERS"`
	// 是否跳过 TLS 证书校验。
	IgnoreSSL bool `json:"ignore_ssl" yaml:"ignoreSsl" xml:"ignoreSsl" env:"WEBHOOK_IGNORE_SSL"`
}
