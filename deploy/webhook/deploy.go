package webhook

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chihqiang/logx"
	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-resty/resty/v2"
)

type Deploy struct {
	Config *Config
}

func (d *Deploy) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.Config = cfg
	return nil
}

// placeholderRe 匹配 __domain__ / __cert__ / __key__ 占位符。
var placeholderRe = regexp.MustCompile(`__([a-zA-Z0-9_]+)__`)

// replacePlaceholders 用证书信息替换请求体模板中的占位符，
// 并对替换值做 JSON 转义，保证可作为 JSON 字符串安全嵌入。
func replacePlaceholders(data string, vars map[string]string) string {
	return placeholderRe.ReplaceAllStringFunc(data, func(match string) string {
		key := placeholderRe.FindStringSubmatch(match)[1]
		val, ok := vars[key]
		if !ok {
			return match
		}
		quoted := strconv.Quote(val)
		return quoted[1 : len(quoted)-1]
	})
}

func (d *Deploy) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.Config.Url == "" {
		return fmt.Errorf("webhook url is required")
	}
	method := strings.ToUpper(d.Config.Method)
	if method == "" {
		method = http.MethodPost
	}

	data := replacePlaceholders(d.Config.Data, map[string]string{
		"domain": certificate.Domain,
		"cert":   string(certificate.Certificate),
		"key":    string(certificate.PrivateKey),
	})

	client := resty.New().SetTimeout(30 * time.Second)
	if d.Config.IgnoreSSL {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec // 用户显式配置跳过校验
	}
	req := client.R()

	headers, err := parseHeaders(d.Config.Headers)
	if err != nil {
		return err
	}
	req.Header = headers

	switch method {
	case http.MethodGet:
		params, err := stringMap(data)
		if err != nil {
			return err
		}
		req.SetQueryParams(params)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		switch contentType(req.Header.Get("Content-Type")) {
		case "application/x-www-form-urlencoded":
			params, err := stringMap(data)
			if err != nil {
				return err
			}
			req.SetFormData(params)
		case "multipart/form-data":
			params, err := stringMap(data)
			if err != nil {
				return err
			}
			req.SetMultipartFormData(params)
		default:
			var body any
			if err := json.Unmarshal([]byte(data), &body); err != nil {
				return fmt.Errorf("failed to parse webhook data as json: %w", err)
			}
			req.SetBody(body)
		}
	default:
		return fmt.Errorf("unsupported webhook method: %s", method)
	}

	resp, err := req.Execute(method, d.Config.Url)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %w", err)
	}
	if resp.IsError() {
		return fmt.Errorf("webhook returned error status %d: %s", resp.StatusCode(), resp.String())
	}
	logx.Info("webhook delivered to %s", d.Config.Url)
	return nil
}

func contentType(v string) string {
	if i := strings.IndexByte(v, ';'); i >= 0 {
		return strings.TrimSpace(v[:i])
	}
	return strings.TrimSpace(v)
}

func stringMap(data string) (map[string]string, error) {
	var m map[string]string
	if data == "" {
		return m, nil
	}
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil, fmt.Errorf("failed to parse webhook data as string map: %w", err)
	}
	return m, nil
}

// parseHeaders 解析每行一个 "Key: Value" 的请求头配置。
func parseHeaders(headerStr string) (http.Header, error) {
	headers := make(http.Header)
	for i, line := range strings.Split(headerStr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header at line %d: %s", i+1, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			return nil, fmt.Errorf("header key/value cannot be empty at line %d", i+1)
		}
		headers.Add(http.CanonicalHeaderKey(key), value)
	}
	return headers, nil
}
