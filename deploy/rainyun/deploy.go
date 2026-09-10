package rainyun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
)

const apiBase = "https://api.v2.rainyun.com"

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

// Deploy 更新雨云证书中心的证书内容。
func (d *Deploy) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.Config.ApiKey == "" {
		return fmt.Errorf("rainyun api_key is required")
	}
	if d.Config.CertId == "" {
		return fmt.Errorf("rainyun cert_id is required")
	}
	body, err := json.Marshal(map[string]any{
		"cert":   string(certificate.Certificate),
		"key":    string(certificate.PrivateKey),
		"domain": certificate.Domain,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, apiBase+"/product/sslcenter/"+d.Config.CertId, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("X-Api-Key", d.Config.ApiKey)

	resp, err := dcommon.NewHTTPClient(false).Do(req)
	if err != nil {
		return fmt.Errorf("failed to request rainyun: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 50*1024*1024))
	if err != nil {
		return fmt.Errorf("failed to read rainyun response: %w", err)
	}

	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return fmt.Errorf("failed to parse rainyun response: %w, body: %s", err, string(respBody))
	}
	msg := res.Message
	if msg == "" {
		msg = res.Error
	}
	if res.Code != 200 {
		return fmt.Errorf("rainyun api error: %s", strip(msg))
	}
	return nil
}

func strip(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown error"
	}
	return s
}
