package vod

import (
	"fmt"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	aliyunVod "github.com/alibabacloud-go/vod-20170321/v4/client"
)

func newClient(accessKeyId, accessKeySecret, region string) (*aliyunVod.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/vod
	endpoint := fmt.Sprintf("vod.%s.aliyuncs.com", region)

	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := aliyunVod.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
