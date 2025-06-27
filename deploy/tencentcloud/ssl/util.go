package ssl

import (
	"context"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func newClient(secretId, secretKey string) (*tcssl.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcssl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func FastDeploy(ctx context.Context, secretId, secretKey string, certificate *certificate.Resource) (certId string, err error) {
	sslDeploy := &Deploy{Config: &Config{SecretId: secretId, SecretKey: secretKey}}
	if err := sslDeploy.Deploy(ctx, certificate); err != nil {
		return "", err
	}
	return sslDeploy.GetCertId(), nil
}
