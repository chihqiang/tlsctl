package lecdn

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
)

type client struct {
	baseURL  string
	username string
	password string
	token    string
	ignore   bool
}

// login 登录 LeCDN 并缓存 token。
func (c *client) login() error {
	res, err := c.request("/prod-api/login", http.MethodPost, "", map[string]any{
		"username": c.username,
		"password": c.password,
	})
	if err != nil {
		return fmt.Errorf("lecdn login failed: %w", err)
	}
	data, ok := res["data"].(map[string]any)
	if !ok {
		return fmt.Errorf("lecdn login: invalid response")
	}
	token, ok := data["token"].(string)
	if !ok {
		return fmt.Errorf("lecdn login: token not found")
	}
	c.token = token
	return nil
}

// uploadCert 上传证书，返回证书 ID。
func (c *client) uploadCert(cert, key, name string) (int, error) {
	params := map[string]any{
		"ssl_key":      base64.StdEncoding.EncodeToString([]byte(key)),
		"ssl_pem":      base64.StdEncoding.EncodeToString([]byte(cert)),
		"type":         "upload",
		"auto_renewal": false,
		"name":         name,
	}
	res, err := c.request("/prod-api/certificate", http.MethodPost, c.token, params)
	if err != nil {
		return 0, err
	}
	data, ok := res["data"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("lecdn upload cert: invalid response")
	}
	id, ok := data["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("lecdn upload cert: id not found")
	}
	return int(id), nil
}

// findCertId 按名称查找已有证书 ID，不存在返回 0。
func (c *client) findCertId(name string) (int, error) {
	res, err := c.request("/prod-api/certificate?current_page=1&total=9999&page_size=9999", http.MethodGet, c.token, nil)
	if err != nil {
		return 0, err
	}
	data, ok := res["data"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("lecdn list cert: invalid response")
	}
	items, ok := data["items"].([]any)
	if !ok {
		return 0, fmt.Errorf("lecdn list cert: items not found")
	}
	for _, item := range items {
		cert, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if certName, ok := cert["name"].(string); ok && certName == name {
			if id, ok := cert["id"].(float64); ok {
				return int(id), nil
			}
		}
	}
	return 0, nil
}

// domainIdFromSite 获取站点内域名对应的域名 ID。
func (c *client) domainIdFromSite(siteId int, domain string) (int, error) {
	res, err := c.request(fmt.Sprintf("/prod-api/site/%d/domain_name", siteId), http.MethodGet, c.token, nil)
	if err != nil {
		return 0, err
	}
	items, ok := res["data"].([]any)
	if !ok {
		return 0, fmt.Errorf("lecdn site domains: invalid response")
	}
	for _, item := range items {
		domainData, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if name, _ := domainData["domain_name"].(string); name == domain {
			if id, ok := domainData["id"].(float64); ok {
				return int(id), nil
			}
		}
	}
	return 0, fmt.Errorf("lecdn domain %s not found in site %d", domain, siteId)
}

// deployCert 绑定证书到域名的 HTTPS 配置。
func (c *client) deployCert(siteId, domainId, certId int) error {
	params := map[string]any{"certificate_enable": true, "certificate_id": certId, "domain_name_id": domainId}
	if _, err := c.request(fmt.Sprintf("/prod-api/site/%d/domain_name/certificate", siteId), http.MethodPut, c.token, params); err != nil {
		return err
	}
	return nil
}

// request 发送请求到 LeCDN API。
func (c *client) request(path, method, token string, params map[string]any) (map[string]any, error) {
	var res map[string]any
	jsonData, err := json.Marshal(params)
	if err != nil {
		return res, err
	}
	req, err := http.NewRequest(method, c.baseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return res, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("cookie", "LeCDN-Client="+token)
	}

	resp, err := dcommon.NewHTTPClient(c.ignore).Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("lecdn request failed: status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return res, fmt.Errorf("lecdn response parse failed: %v, body: %s", err, string(body))
	}
	return res, nil
}

// certName 生成稳定证书名，用于去重复用。
func certName(certPem string) string {
	sum := sha256.Sum256([]byte(certPem))
	return "tlsctl-" + hex.EncodeToString(sum[:])
}
