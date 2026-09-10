package safeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
)

// client 封装雷池 WAF 开放 API 调用。
type client struct {
	Config
}

// request 向雷池 WAF 开放 API 发起请求，携带 X-SLCE-API-TOKEN。
func (c *client) request(data map[string]any, method, requestUrl string) (map[string]any, error) {
	if c.Url == "" {
		return nil, fmt.Errorf("safeline url is required")
	}
	if c.ApiToken == "" {
		return nil, fmt.Errorf("safeline api_token is required")
	}
	parsed, err := url.Parse(c.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse safeline url: %w", err)
	}
	base := fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal safeline request: %w", err)
	}
	req, err := http.NewRequest(method, base+requestUrl, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to build safeline request: %w", err)
	}
	req.Header.Set("X-SLCE-API-TOKEN", c.ApiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := dcommon.NewHTTPClient(c.IgnoreSSL).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request safeline: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read safeline response: %w", err)
	}

	var res map[string]any
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to parse safeline response: %w", err)
	}
	if msg, _ := res["msg"].(string); msg != "" {
		return nil, fmt.Errorf("safeline request failed: %s", msg)
	}
	return res, nil
}

// uploadCert 上传证书到雷池 WAF。certId 为 0 时新建证书，否则覆盖更新。
func (c *client) uploadCert(certId int64, keyPem, certPem string) (int64, error) {
	data := map[string]any{
		"type": 2,
		"manual": map[string]any{
			"crt": certPem,
			"key": keyPem,
		},
	}
	if certId != 0 {
		data["id"] = certId
	}
	response, err := c.request(data, http.MethodPost, "api/open/cert")
	if err != nil {
		return 0, err
	}
	id, ok := response["data"].(float64)
	if !ok {
		return 0, fmt.Errorf("safeline upload cert: invalid response")
	}
	return int64(id), nil
}

// siteCertId 按站点名称匹配出证书 ID，站点未启用证书时返回 0。
func (c *client) siteCertId(siteName string) (int64, error) {
	requestUrl := fmt.Sprintf("api/open/site?page=1&page_size=100&site=%s", siteName)
	response, err := c.request(map[string]any{}, http.MethodGet, requestUrl)
	if err != nil {
		return 0, err
	}
	data, ok := response["data"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("safeline site list: invalid response")
	}
	items, ok := data["data"].([]any)
	if !ok {
		return 0, fmt.Errorf("safeline site list: data.data not found")
	}
	for _, item := range items {
		site, ok := item.(map[string]any)
		if !ok {
			continue
		}
		comment, _ := site["comment"].(string)
		if comment == siteName {
			if certId, ok := site["cert_id"].(float64); ok {
				return int64(certId), nil
			}
			return 0, nil
		}
	}
	return 0, fmt.Errorf("safeline site %s not found", siteName)
}
