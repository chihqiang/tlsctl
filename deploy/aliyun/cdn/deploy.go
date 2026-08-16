package cdn

import (
	"context"
	"fmt"
	"strings"
	"time"

	aliyunCdn "github.com/alibabacloud-go/cdn-20180510/v5/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/logx"
	"github.com/go-acme/lego/v4/certificate"
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
	domain := strings.TrimPrefix(d.Config.Domain, "*")
	client, err := newClient(d.Config.AccessKeyId, d.Config.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("failed to create sdk client: %w", err)
	}
	setCdnDomainSSLCertificateReq := &aliyunCdn.SetCdnDomainSSLCertificateRequest{
		DomainName:  tea.String(domain),
		CertName:    tea.String(fmt.Sprintf("certimate-%d", time.Now().UnixMilli())),
		CertType:    tea.String("upload"),
		SSLProtocol: tea.String("on"),
		SSLPub:      tea.String(string(certificate.Certificate)),
		SSLPri:      tea.String(string(certificate.PrivateKey)),
	}
	setCdnDomainSSLCertificateResp, err := client.SetCdnDomainSSLCertificate(setCdnDomainSSLCertificateReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'cdn.SetCdnDomainSSLCertificate': %w", err)
	}
	logx.Info("CDN domain name certificate has been set up %+v", setCdnDomainSSLCertificateResp)
	return nil
}
