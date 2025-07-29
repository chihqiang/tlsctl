package waf

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcwaf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
)

func newClient(secretId, secretKey, region string) (*tcwaf.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcwaf.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return client, nil
}
