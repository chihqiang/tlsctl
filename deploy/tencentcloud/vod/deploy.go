package vod

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcvod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
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
	// 上传新证书
	client, err := newClient(d.Config.SecretId, d.Config.SecretKey)
	if err != nil {
		return err
	}
	certId, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return err
	}
	// 设置点播域名 HTTPS 证书
	// REF: https://cloud.tencent.com/document/api/266/102015
	setVodDomainCertificateReq := tcvod.NewSetVodDomainCertificateRequest()
	setVodDomainCertificateReq.Domain = common.StringPtr(d.Config.Domain)
	setVodDomainCertificateReq.Operation = common.StringPtr("Set")
	setVodDomainCertificateReq.CertID = common.StringPtr(certId)
	if d.Config.SubAppId != 0 {
		setVodDomainCertificateReq.SubAppId = common.Uint64Ptr(uint64(d.Config.SubAppId))
	}
	setVodDomainCertificateResp, err := client.SetVodDomainCertificate(setVodDomainCertificateReq)
	log.Info("sdk request 'vod.SetVodDomainCertificate' request: %#v response:%#v", setVodDomainCertificateReq, setVodDomainCertificateResp)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'vod.SetVodDomainCertificate': %w", err)
	}

	return nil
}
