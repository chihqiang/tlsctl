package discord

import (
	"fmt"
	"wangzhiqiang/tlsctl/pkg/httpx"

	"github.com/caarlos0/env/v11"
)

type Report struct {
	// Discord Bot API Token。
	Token string `json:"token" xml:"token" yaml:"token" env:"DISCORD_TOKEN"`
	// Discord Channel ID。
	ChannelId string `json:"channel_id" xml:"channel_id" yaml:"channel_id" env:"DISCORD_CHANNEL_ID"`
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
//
//	https://discord.com/developers/docs/resources/message#create-message
func (d *Report) SendText(msg string) error {
	url := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", d.ChannelId)
	data := map[string]any{
		"content": msg,
	}
	var v map[string]interface{}
	if err := httpx.PostJSON(url, data, &v, map[string]string{
		"Authorization": "Bot " + d.Token,
		"User-Agent":    "tlsctl",
	}); err != nil {
		return fmt.Errorf("discord api error: failed to send request: %w", err)
	}
	return nil
}
