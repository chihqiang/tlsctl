package live

import (
	"context"
	"fmt"
	"strings"
	"time"

	aliyunLive "github.com/alibabacloud-go/live-20161101/client"
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
	client, err := newClient(d.Config.AccessKeyId, d.Config.AccessKeySecret, d.Config.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create sdk client")
	}
	// "*.example.com" → ".example.com"，适配阿里云 Live 要求的泛域名格式
	domain := strings.TrimPrefix(d.Config.Domain, "*")
	// 设置域名证书
	// REF: https://help.aliyun.com/zh/live/developer-reference/api-live-2016-11-01-setlivedomaincertificate
	setLiveDomainSSLCertificateReq := &aliyunLive.SetLiveDomainCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(string(certificate.Certificate)),
		SSLPri:      tea.String(string(certificate.PrivateKey)),
	}
	setLiveDomainSSLCertificateResp, err := client.SetLiveDomainCertificate(setLiveDomainSSLCertificateReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute sdk request 'live.SetLiveDomainCertificate'")
	}
	log.Info("Domain name certificate has been set up %+v", setLiveDomainSSLCertificateResp)
	return nil
}
