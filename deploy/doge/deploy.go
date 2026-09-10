package doge

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
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

// Deploy 上传证书并绑定到多吉云 CDN 域名（按证书 SHA256 去重复用）。
func (d *Deploy) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.Config.AccessKey == "" || d.Config.SecretKey == "" {
		return fmt.Errorf("doge access_key and secret_key are required")
	}
	domain := d.Config.Domain
	if domain == "" {
		domain = certificate.Domain
	}
	certStr := string(certificate.Certificate)
	sum := sha256.Sum256([]byte(certStr))
	note := fmt.Sprintf("tlsctl-%s", hex.EncodeToString(sum[:]))

	certId, err := d.findCertId(note)
	if err != nil {
		return err
	}
	if certId == 0 {
		certId, err = d.uploadCert(certStr, string(certificate.PrivateKey), note)
		if err != nil {
			return err
		}
	}
	if _, err := d.call("/cdn/cert/bind.json", map[string]any{"id": certId, "domain": domain}, true); err != nil {
		return err
	}
	return nil
}

// findCertId 按 note 在已有证书中查找证书 ID，不存在返回 0。
func (d *Deploy) findCertId(note string) (float64, error) {
	res, err := d.call("/cdn/cert/list.json", map[string]any{}, true)
	if err != nil {
		return 0, err
	}
	data, ok := res["data"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("doge list cert: data not found")
	}
	certs, ok := data["certs"].([]any)
	if !ok {
		return 0, fmt.Errorf("doge list cert: certs not found")
	}
	for _, item := range certs {
		c, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if c["note"] == note {
			if id, ok := c["id"].(float64); ok {
				return id, nil
			}
		}
	}
	return 0, nil
}

func (d *Deploy) uploadCert(cert, key, note string) (float64, error) {
	res, err := d.call("/cdn/cert/upload.json", map[string]any{"cert": cert, "private": key, "note": note}, true)
	if err != nil {
		return 0, err
	}
	data, ok := res["data"].(map[string]any)
	if !ok {
		return 0, fmt.Errorf("doge upload cert: data not found")
	}
	id, ok := data["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("doge upload cert: id not found")
	}
	return id, nil
}

// call 调用多吉云开放 API，使用 HMAC-SHA1 签名。
func (d *Deploy) call(apiPath string, data map[string]any, jsonMode bool) (map[string]any, error) {
	body := ""
	mime := ""
	if jsonMode {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = string(b)
		mime = "application/json"
	} else {
		values := url.Values{}
		for k, v := range data {
			values.Set(k, v.(string))
		}
		body = values.Encode()
		mime = "application/x-www-form-urlencoded"
	}

	mac := hmac.New(sha1.New, []byte(d.Config.SecretKey))
	mac.Write([]byte(apiPath + "\n" + body))
	sign := hex.EncodeToString(mac.Sum(nil))

	req, err := http.NewRequest(http.MethodPost, "https://api.dogecloud.com"+apiPath, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mime)
	req.Header.Set("Authorization", "TOKEN "+d.Config.AccessKey+":"+sign)

	resp, err := dcommon.NewHTTPClient(false).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(r, &result); err != nil {
		return nil, err
	}
	if code, ok := result["code"].(float64); ok && code != 200 {
		return nil, fmt.Errorf("doge api error: %v", result["msg"])
	}
	return result, nil
}
