package clb

import (
	tccommon "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"
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
	sslClient, err := tccommon.NewSSLClient(secretId, secretKey, region) // 注意虽然官方文档中地域无需指定，但实际需要部署到 CLB 时必传
	if err != nil {
		return nil, err
	}

	clbClient, err := tcclb.NewClient(common.NewCredential(secretId, secretKey), region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &Clients{
		SSL: sslClient,
		CLB: clbClient,
	}, nil
}
