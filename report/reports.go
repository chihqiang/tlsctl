package report

import (
	"fmt"
	"wangzhiqiang/tlsctl/report/dingtalk"
	"wangzhiqiang/tlsctl/report/discord"
	"wangzhiqiang/tlsctl/report/feishu"
	"wangzhiqiang/tlsctl/report/slack"
	"wangzhiqiang/tlsctl/report/telegram"
	"wangzhiqiang/tlsctl/report/workwx"
)

var reports = map[string]IReport{
	"feishu":   &feishu.Report{},
	"dingtalk": &dingtalk.Report{},
	"discord":  &discord.Report{},
	"telegram": &telegram.Report{},
	"slack":    &slack.Report{},
	"workwx":   &workwx.Report{},
}

func Register(name string, deploy IReport) {
	reports[name] = deploy
}

func Get(name string) (IReport, error) {
	if deploy, ok := reports[name]; ok {
		if err := deploy.WithEnvConfig(); err != nil {
			return nil, err
		}
		return deploy, nil
	}
	return nil, fmt.Errorf("report `%s` not found", name)
}

func All() map[string]IReport {
	return reports
}
