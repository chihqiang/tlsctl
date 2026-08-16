package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func newClient(accessKeyId, accessKeySecret, region string) (*oss.Client, error) {
	// 接入点一览 https://help.aliyun.com/zh/oss/user-guide/regions-and-endpoints
	var endpoint string
	switch region {
	case "":
		endpoint = "oss.aliyuncs.com"
	case
		"cn-hzjbp",
		"cn-hzjbp-a",
		"cn-hzjbp-b":
		endpoint = "oss-cn-hzjbp-a-internal.aliyuncs.com"
	case
		"cn-shanghai-finance-1",
		"cn-shenzhen-finance-1",
		"cn-beijing-finance-1",
		"cn-north-2-gov-1":
		endpoint = fmt.Sprintf("oss-%s-internal.aliyuncs.com", region)
	default:
		endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", region)
	}

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}

	return client, nil
}
