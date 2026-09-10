package cos

import (
	"context"
	"fmt"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	tccommon "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Deploy struct {
	Config *Config
}

func (d *Deploy) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.Config = cfg
	return nil
}
func (d *Deploy) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	client, err := tccommon.NewSSLClient(d.Config.SecretId, d.Config.SecretKey, d.Config.Region)
	if err != nil {
		return err
	}
	certId, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return err
	}
	// 证书部署到 COS 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcssl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(certId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cos")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s#%s#%s", d.Config.Region, d.Config.Bucket, d.Config.Domain)})
	resp, err := client.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'ssl.DeployCertificateInstance': %w", err)
	}

	return tccommon.WaitForDeploy(ctx, client, *resp.Response.DeployRecordId)
}
