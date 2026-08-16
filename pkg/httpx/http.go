package httpx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var (
	defaultClient = &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableKeepAlives: true,
		},
		Timeout: 10 * time.Second,
	}
)

// PostJSON 发送 JSON POST 请求
func PostJSON(targetURL string, data any, v any, headers map[string]string) error {
	// 序列化 JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 构建请求
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	// 设置默认 Header
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := defaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 199 && resp.StatusCode < 300 {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return fmt.Errorf("postJson:%s failed with status code %d", targetURL, resp.StatusCode)
}
