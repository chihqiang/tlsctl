package clb

import (
	tcclb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Clients struct {
	SSL *tcssl.Client
	CLB *tcclb.Client
}

func newClient(secretId, secretKey, region string) (*Clients, error) {
	credential := common.NewCredential(secretId, secretKey)

	// 注意虽然官方文档中地域无需指定，但实际需要部署到 CLB 时必传
	sslClient, err := tcssl.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	clbClient, err := tcclb.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &Clients{
		SSL: sslClient,
		CLB: clbClient,
	}, nil
}
