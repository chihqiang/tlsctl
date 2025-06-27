package telegram

import (
	"fmt"
	"wangzhiqiang/tlsctl/pkg/httpx"

	"github.com/caarlos0/env/v11"
)

type Report struct {
	// Telegram Bot API Token。
	Token string `json:"token" xml:"token" yaml:"token" env:"TELEGRAM_TOKEN"`
	// Telegram Chat ID。
	ChatId int64 `json:"chat_id" xml:"chat_id" yaml:"chat_id" env:"TELEGRAM_CHAT_ID"`
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

func (d *Report) SendText(msg string) error {
	data := map[string]any{
		"chat_id": d.ChatId,
		"text":    msg,
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", d.Token)
	resp := map[string]any{}
	err := httpx.PostJSON(url, data, &resp, map[string]string{})
	if err != nil {
		return err
	}
	return nil
}
