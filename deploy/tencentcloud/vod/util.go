package vod

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcvod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
)

func newClient(secretId, secretKey string) (*tcvod.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcvod.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
