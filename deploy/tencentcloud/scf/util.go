package scf

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcscf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
)

func newClient(secretId, secretKey, region string) (*tcscf.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcscf.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
