package workwx

import (
	"fmt"
	"wangzhiqiang/tlsctl/pkg/httpx"

	"github.com/caarlos0/env/v11"
)

type Report struct {
	// 企业微信机器人 Webhook 地址。
	URL string `json:"url" yaml:"URL" xml:"URL" env:"workwx_url"`
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

// SendText
// https://developer.work.weixin.qq.com/document/path/91770
func (d *Report) SendText(msg string) error {
	data := map[string]any{
		"msgtype": map[string]any{
			"content": msg,
		},
		"text": msg,
	}
	resp := map[string]any{}
	err := httpx.PostJSON(d.URL, data, &resp, map[string]string{})
	if err != nil {
		return err
	}
	if resp["errcode"].(float64) != 0 {
		return fmt.Errorf("feishu send failed: %s", resp["errmsg"].(string))
	}
	return nil
}
