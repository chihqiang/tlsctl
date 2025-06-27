package live

import (
	"fmt"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	aliyunLive "github.com/alibabacloud-go/live-20161101/client"
	"github.com/alibabacloud-go/tea/tea"
)

func newClient(accessKeyId, accessKeySecret, region string) (*aliyunLive.Client, error) {
	// 接入点一览 https://api.aliyun.com/product/live
	var endpoint string
	switch region {
	case
		"cn-qingdao",
		"cn-beijing",
		"cn-shanghai",
		"cn-shenzhen",
		"ap-northeast-1",
		"ap-southeast-5",
		"me-central-1":
		endpoint = "live.aliyuncs.com"
	default:
		endpoint = fmt.Sprintf("live.%s.aliyuncs.com", region)
	}

	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String(endpoint),
	}

	client, err := aliyunLive.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
