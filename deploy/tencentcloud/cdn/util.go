package cdn

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

func newClients(secretId, secretKey string) (*Clients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcssl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	cdnClient, err := tccdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &Clients{
		SSL: sslClient,
		CDN: cdnClient,
	}, nil
}
