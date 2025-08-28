package dcdn

import (
	"context"
	"fmt"
	"strings"
	"time"

	aliyunDcdn "github.com/alibabacloud-go/dcdn-20180115/v3/client"
	"github.com/alibabacloud-go/tea/tea"
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
	client, err := newClient(d.Config.AccessKeyId, d.Config.AccessKeySecret)
	if err != nil {
		return err
	}
	// "*.example.com" → ".example.com"，适配阿里云 DCDN 要求的泛域名格式
	domain := strings.TrimPrefix(d.Config.Domain, "*")
	// 配置域名证书
	// REF: https://help.aliyun.com/zh/edge-security-acceleration/dcdn/developer-reference/api-dcdn-2018-01-15-setdcdndomainsslcertificate
	setDcdnDomainSSLCertificateReq := &aliyunDcdn.SetDcdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(string(certificate.Certificate)),
		SSLPri:      tea.String(string(certificate.PrivateKey)),
	}
	setDcdnDomainSSLCertificateResp, err := client.SetDcdnDomainSSLCertificate(setDcdnDomainSSLCertificateReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'dcdn.SetDcdnDomainSSLCertificate'")
	}
	log.Info("CDN domain name certificate has been set up %+v", setDcdnDomainSSLCertificateResp)
	return nil
}
