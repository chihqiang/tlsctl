package cdn

import (
	aliyunCdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	aliyunOpen "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

func newClient(accessKeyId, accessKeySecret string) (*aliyunCdn.Client, error) {
	config := &aliyunOpen.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("cdn.aliyuncs.com"),
	}
	client, err := aliyunCdn.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
