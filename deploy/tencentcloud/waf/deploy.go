package waf

import (
	"context"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcwaf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
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
		return err
	}
	certId, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return err
	}
	describeDomainDetailsSaasReq := tcwaf.NewDescribeDomainDetailsSaasRequest()
	describeDomainDetailsSaasReq.Domain = common.StringPtr(d.Config.Domain)
	describeDomainDetailsSaasReq.DomainId = common.StringPtr(d.Config.DomainId)
	describeDomainDetailsSaasReq.InstanceId = common.StringPtr(d.Config.InstanceId)
	describeDomainDetailsSaasResp, err := client.DescribeDomainDetailsSaas(describeDomainDetailsSaasReq)
	log.Debug("sdk request 'waf.DescribeDomainDetailsSaas request: %#v response: %#v", describeDomainDetailsSaasReq, describeDomainDetailsSaasResp)
	if err != nil {
		return err
	}
	// 编辑 SaaS 型 WAF 域名
	// REF: https://cloud.tencent.com/document/api/627/94309
	modifySpartaProtectionReq := tcwaf.NewModifySpartaProtectionRequest()
	modifySpartaProtectionReq.Domain = common.StringPtr(d.Config.Domain)
	modifySpartaProtectionReq.DomainId = common.StringPtr(d.Config.DomainId)
	modifySpartaProtectionReq.InstanceID = common.StringPtr(d.Config.InstanceId)
	modifySpartaProtectionReq.CertType = common.Int64Ptr(2)
	modifySpartaProtectionReq.SSLId = common.StringPtr(certId)
	modifySpartaProtectionResp, err := client.ModifySpartaProtection(modifySpartaProtectionReq)
	if err != nil {
		return err
	}
	log.Debug("sdk request 'waf.ModifySpartaProtection' request: %#v response: %#v", modifySpartaProtectionReq, modifySpartaProtectionResp)
	return nil
}
