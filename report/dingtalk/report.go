package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
	"wangzhiqiang/tlsctl/pkg/httpx"

	"github.com/caarlos0/env/v11"
)

type Report struct {
	WebHookUrl string `json:"web_hook_url" xml:"WebHookUrl" yaml:"WebHookUrl" env:"DINGTALK_WEB_HOOK_URL"`
	Secret     string `json:"secret" xml:"Secret" yaml:"Secret" env:"DINGTALK_SECRET"`
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
func (d *Report) sign() (string, error) {
	timestamp := time.Now().UnixNano() / 1000000
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, d.Secret)
	hash := hmac.New(sha256.New, []byte(d.Secret))
	hash.Write([]byte(stringToSign))
	sum := hash.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(sum)
	webhookurl, _ := url.Parse(d.WebHookUrl)
	query := webhookurl.Query()
	query.Set("timestamp", fmt.Sprint(timestamp))
	query.Set("sign", signature)
	webhookurl.RawQuery = query.Encode()
	return webhookurl.String(), nil
}

// SendText
// https://open.dingtalk.com/document/robots/custom-robot-access
func (d *Report) SendText(msg string) error {
	data := map[string]any{
		"text": map[string]any{
			"content": msg,
		},
		"msgtype": "text",
	}
	signUrl, err := d.sign()
	var res map[string]interface{}
	err = httpx.PostJSON(signUrl, data, &res, map[string]string{})
	if err != nil {
		return err
	}
	if res["errcode"].(float64) != 0 {
		return fmt.Errorf("dingtalk send failed: %s", res["errmsg"].(string))
	}
	return nil
}
