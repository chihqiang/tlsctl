package btpanel

import (
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

// client 封装宝塔面板 API 调用与公共逻辑。
type client struct {
	Config
}

// request 向宝塔面板 API 发起请求。
// 根据宝塔开放接口规范，使用 api_key 生成带时间戳的签名。
func (c *client) request(data url.Values, requestUrl string) (map[string]any, error) {
	if c.Url == "" {
		return nil, fmt.Errorf("btpanel url is required")
	}
	if c.ApiKey == "" {
		return nil, fmt.Errorf("btpanel api_key is required")
	}
	timestamp := time.Now().Unix()
	data.Set("request_time", fmt.Sprintf("%d", timestamp))
	data.Set("request_token", sign(timestamp, c.ApiKey))

	parsed, err := url.Parse(c.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse btpanel url: %w", err)
	}
	base := fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host)

	req, err := http.NewRequest(http.MethodPost, base+requestUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to build btpanel request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := dcommon.NewHTTPClient(c.IgnoreSSL).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request btpanel: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read btpanel response: %w", err)
	}

	var res map[string]any
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to parse btpanel response: %w, body: %s", err, string(body))
	}
	if status, ok := res["status"].(bool); ok && !status {
		msg, _ := res["msg"].(string)
		return nil, fmt.Errorf("btpanel request failed: %s", msg)
	}
	return res, nil
}

// uploadCert 上传证书到宝塔证书库，返回 ssl_hash 用于后续绑定。
func (c *client) uploadCert(keyPem, certPem string) (string, error) {
	data := url.Values{}
	data.Set("key", keyPem)
	data.Set("csr", certPem)
	response, err := c.request(data, "ssl/cert/save_cert")
	if err != nil {
		return "", err
	}
	sslHash, ok := response["ssl_hash"].(string)
	if !ok {
		return "", fmt.Errorf("failed to upload btpanel cert: ssl_hash not found")
	}
	return sslHash, nil
}

// sign 生成宝塔接口签名：request_token = md5(timestamp + md5(api_key))。
func sign(timestamp int64, apiKey string) string {
	keyMd5 := md5.Sum([]byte(apiKey))
	keyMd5Hex := strings.ToLower(hex.EncodeToString(keyMd5[:]))
	signMd5 := md5.Sum([]byte(fmt.Sprintf("%d", timestamp) + keyMd5Hex))
	return strings.ToLower(hex.EncodeToString(signMd5[:]))
}
