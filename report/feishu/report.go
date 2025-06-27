package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
	"wangzhiqiang/tlsctl/pkg/httpx"

	"github.com/caarlos0/env/v11"
)

type Report struct {
	WebHookUrl string `json:"web_hook_url" xml:"WebHookUrl" yaml:"WebHookUrl" env:"FEISHU_WEB_HOOK_URL"`
	Secret     string `json:"secret" xml:"Secret" yaml:"Secret" env:"FEISHU_SECRET"`
}

func (d *Report) WithEnvConfig() error {
	var cfg Report
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	d = &cfg
	return nil
}
func (d *Report) sign(params *map[string]any) error {
	timestamp := time.Now().Unix()
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + d.Secret
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return fmt.Errorf("failed to generate signature: %v", err)
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	(*params)["timestamp"] = timestamp
	(*params)["sign"] = signature
	return nil
}
func (d *Report) SendText(msg string) error {
	data := map[string]any{
		"timestamp": 0,
		"content": map[string]any{
			"text": msg,
		},
		"msg_type": "text",
		"sign":     "",
	}
	if err := d.sign(&data); err != nil {
		return err
	}
	var res map[string]interface{}
	err := httpx.PostJSON(d.WebHookUrl, data, &res, map[string]string{})
	if err != nil {
		return err
	}
	if res["code"].(float64) != 0 {
		return fmt.Errorf("feishu send failed: %s", res["msg"].(string))
	}
	return nil
}
