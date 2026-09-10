package ssl

import (
	"context"

	tccommon "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"
	"github.com/go-acme/lego/v4/certificate"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

func newClient(secretId, secretKey string) (*tcssl.Client, error) {
	return tccommon.NewSSLClient(secretId, secretKey, "")
}

func FastDeploy(ctx context.Context, secretId, secretKey string, certificate *certificate.Resource) (certId string, err error) {
	sslDeploy := &Deploy{Config: &Config{
		BaseConfig: tccommon.BaseConfig{SecretId: secretId, SecretKey: secretKey},
	}}
	if err := sslDeploy.Deploy(ctx, certificate); err != nil {
		return "", err
	}
	return sslDeploy.GetCertId(), nil
}
