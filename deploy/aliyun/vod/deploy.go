package vod

import (
	"context"
	"fmt"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	aliyunVod "github.com/alibabacloud-go/vod-20170321/v4/client"
	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
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
	client, err := newClient(d.Config.AccessKeyId, d.Config.AccessKeySecret, d.Config.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create sdk client")
	}
	// 设置域名证书
	// REF: https://help.aliyun.com/zh/vod/developer-reference/api-vod-2017-03-21-setvoddomainsslcertificate
	setVodDomainSSLCertificateReq := &aliyunVod.SetVodDomainSSLCertificateRequest{
		DomainName:  tea.String(d.Config.Domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(string(certificate.Certificate)),
		SSLPri:      tea.String(string(certificate.PrivateKey)),
	}
	setVodDomainSSLCertificateResp, err := client.SetVodDomainSSLCertificate(setVodDomainSSLCertificateReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'vod.SetVodDomainSSLCertificate': %w", err)
	}
	log.Info("已设置域名证书  %+v", setVodDomainSSLCertificateResp)
	return nil
}
