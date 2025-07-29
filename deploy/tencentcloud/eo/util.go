package eo

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	tcteo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
)

type clients struct {
	SSL *tcssl.Client
	TEO *tcteo.Client
}

func newClient(secretId, secretKey string) (*clients, error) {
	credential := common.NewCredential(secretId, secretKey)

	sslClient, err := tcssl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	teoClient, err := tcteo.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return &clients{
		SSL: sslClient,
		TEO: teoClient,
	}, nil
}
