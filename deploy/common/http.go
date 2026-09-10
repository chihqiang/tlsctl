package common

import (
	"crypto/tls"
	"net/http"
)

// NewHTTPClient 创建一个 HTTP 客户端，ignoreSSL 为 true 时跳过 TLS 证书校验。
// 用于宝塔/1Panel/雷池 WAF 等自建面板接口（通常使用自签名证书）。
func NewHTTPClient(ignoreSSL bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: ignoreSSL},
		DisableKeepAlives: true,
	}
	return &http.Client{Transport: tr}
}
