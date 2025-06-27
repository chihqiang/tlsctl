package dcdn

import (
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunDcdn "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	"github.com/alibabacloud-go/tea/tea"
)

func newClient(accessKeyId, accessKeySecret string) (*aliyunDcdn.Client, error) {
	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("dcdn.aliyuncs.com"),
	}

	client, err := aliyunDcdn.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
