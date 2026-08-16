package eo

import (
	"context"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcteo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
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
	client, err := newClient(d.Config.SecretId, d.Config.SecretKey)
	if err != nil {
		return err
	}
	certId, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return err
	}
	// 配置域名证书
	// REF: https://cloud.tencent.com/document/product/1552/80764
	modifyHostsCertificateReq := tcteo.NewModifyHostsCertificateRequest()
	modifyHostsCertificateReq.ZoneId = common.StringPtr(d.Config.ZoneId)
	modifyHostsCertificateReq.Mode = common.StringPtr("tlsctl")
	modifyHostsCertificateReq.Hosts = common.StringPtrs([]string{d.Config.Domain})
	modifyHostsCertificateReq.ServerCertInfo = []*tcteo.ServerCertInfo{{CertId: common.StringPtr(certId)}}
	modifyHostsCertificateResp, err := client.TEO.ModifyHostsCertificate(modifyHostsCertificateReq)
	if err != nil {
		return err
	}
	log.Info("sdk request 'teo.ModifyHostsCertificate request: %#v  response:%#v", modifyHostsCertificateReq, modifyHostsCertificateResp)
	return nil
}
