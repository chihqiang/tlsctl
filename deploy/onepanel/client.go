package onepanel

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
)

// client 封装 1Panel API 调用与签名逻辑。
type client struct {
	Config
}

// request 向 1Panel API 发起 JSON 请求，并携带签名头。
func (c *client) request(data map[string]any, method, requestUrl string) (map[string]any, error) {
	if c.Url == "" {
		return nil, fmt.Errorf("onepanel url is required")
	}
	if c.ApiKey == "" {
		return nil, fmt.Errorf("onepanel api_key is required")
	}
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	version := c.Version
	switch version {
	case "", "1", "v1":
		requestUrl = "api/v1/" + requestUrl
	case "2", "v2":
		data["ssl"] = "Enable"
		if strings.Split(requestUrl, "/")[0] == "settings" {
			requestUrl = "core/" + requestUrl
		}
		requestUrl = "api/v2/" + requestUrl
	default:
		return nil, fmt.Errorf("unsupported 1panel version: %s", version)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal onepanel request: %w", err)
	}
	parsed, err := url.Parse(c.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse onepanel url: %w", err)
	}
	base := fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host)

	req, err := http.NewRequest(method, base+requestUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to build onepanel request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("1Panel-Timestamp", timestamp)
	req.Header.Set("1Panel-Token", token(timestamp, c.ApiKey))

	resp, err := dcommon.NewHTTPClient(c.IgnoreSSL).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request onepanel: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read onepanel response: %w", err)
	}

	var res map[string]any
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to parse onepanel response: %w", err)
	}
	code, ok := res["code"].(float64)
	if !ok {
		return nil, fmt.Errorf("onepanel request failed: no code in response")
	}
	if code != 200 {
		msg, _ := res["message"].(string)
		return nil, fmt.Errorf("onepanel request failed: %s", msg)
	}
	return res, nil
}

// token 生成 1Panel 签名：md5("1panel" + api_key + timestamp)。
func token(timestamp, apiKey string) string {
	sum := md5.Sum([]byte("1panel" + apiKey + timestamp))
	return hex.EncodeToString(sum[:])
}
