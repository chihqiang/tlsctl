package slack

import (
	"wangzhiqiang/tlsctl/pkg/httpx"

	"github.com/caarlos0/env/v11"
)

type Report struct {
	// Slack Bot API Token。
	Token string `json:"token" xml:"token" yaml:"token" env:"slack_token"`
	//  Slack Channel ID。
	ChannelId int64 `json:"channel_id" xml:"channelId" yaml:"ChannelId" env:"SLACK_CHANNEL_ID"`
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
// https://docs.slack.dev/messaging/sending-and-scheduling-messages#publishing
func (d *Report) SendText(msg string) error {
	data := map[string]any{
		"token":   d.Token,
		"channel": d.ChannelId,
		"text":    msg,
	}
	url := "https://slack.com/api/chat.postMessage"
	resp := map[string]any{}
	err := httpx.PostJSON(url, data, &resp, map[string]string{})
	if err != nil {
		return err
	}
	return nil
}
