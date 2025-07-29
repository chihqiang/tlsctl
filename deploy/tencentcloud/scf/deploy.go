package scf

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v11"
	ssl "github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcscf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
)

type Deploy struct {
	Config *Config
}

func (d *Deploy) WithEnvConfig() error {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	d.Config = &cfg
	return nil
}
func (d *Deploy) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	client, err := newClient(d.Config.SecretId, d.Config.SecretKey, d.Config.Region)
	if err != nil {
		return fmt.Errorf("failed to create sdk client: %w", err)
	}
	getCustomDomainReq := tcscf.NewGetCustomDomainRequest()
	getCustomDomainReq.Domain = common.StringPtr(d.Config.Domain)
	getCustomDomainResp, err := client.GetCustomDomain(getCustomDomainReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'scf.GetCustomDomain': %w", err)
	}
	certId, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return fmt.Errorf("failed to deploy certificate: %w", err)
	}
	// 更新云函数自定义域名
	// REF: https://cloud.tencent.com/document/product/583/111922
	updateCustomDomainReq := tcscf.NewUpdateCustomDomainRequest()
	updateCustomDomainReq.Domain = common.StringPtr(d.Config.Domain)
	updateCustomDomainReq.CertConfig = &tcscf.CertConf{
		CertificateId: common.StringPtr(certId),
	}
	updateCustomDomainReq.Protocol = getCustomDomainResp.Response.Protocol
	updateCustomDomainResp, err := client.UpdateCustomDomain(updateCustomDomainReq)
	log.Info("sdk request 'scf.UpdateCustomDomain' request: %#v  response:%#v", updateCustomDomainReq, updateCustomDomainResp)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'scf.UpdateCustomDomain': %w", err)
	}
	return nil
}
