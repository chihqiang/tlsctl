package lecdn

type Config struct {
	// LeCDN 面板地址，如 https://example.com。
	Url string `json:"url" yaml:"url" xml:"url" env:"LECDN_URL"`
	// LeCDN 登录用户名。
	Username string `json:"username" yaml:"username" xml:"username" env:"LECDN_USERNAME"`
	// LeCDN 登录密码。
	Password string `json:"password" yaml:"password" xml:"password" env:"LECDN_PASSWORD"`
	// 是否跳过 TLS 证书校验。
	IgnoreSSL bool `json:"ignore_ssl" yaml:"ignoreSsl" xml:"ignoreSsl" env:"LECDN_IGNORE_SSL"`
	// 站点 ID。
	SiteId int `json:"site_id" yaml:"siteId" xml:"siteId" env:"LECDN_SITE_ID"`
	// 加速域名，留空时使用证书主域名。
	Domain string `json:"domain" yaml:"domain" xml:"domain" env:"LECDN_DOMAIN"`
}
