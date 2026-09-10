package common

import (
	tccdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Clients struct {
	SSL *tcssl.Client
	CDN *tccdn.Client
}

func NewSSLClient(secretId, secretKey, region string) (*tcssl.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	return tcssl.NewClient(credential, region, profile.NewClientProfile())
}

func NewClients(secretId, secretKey string) (*Clients, error) {
	sslClient, err := NewSSLClient(secretId, secretKey, "")
	if err != nil {
		return nil, err
	}

	credential := common.NewCredential(secretId, secretKey)
	cdnClient, err := tccdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &Clients{
		SSL: sslClient,
		CDN: cdnClient,
	}, nil
}
