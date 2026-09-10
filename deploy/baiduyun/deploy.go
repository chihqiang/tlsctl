package baiduyun

import (
	"context"
	"fmt"
	"time"

	baiduyuncdn "github.com/baidubce/bce-sdk-go/services/cdn"
	"github.com/baidubce/bce-sdk-go/services/cdn/api"
	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
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

// Deploy 上传证书并绑定到百度云 CDN 域名。
func (d *Deploy) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.Config.AccessKey == "" || d.Config.SecretKey == "" {
		return fmt.Errorf("baidu access_key and secret_key are required")
	}
	domain := d.Config.Domain
	if domain == "" {
		domain = certificate.Domain
	}

	client, err := baiduyuncdn.NewClient(d.Config.AccessKey, d.Config.SecretKey, "https://cdn.baidubce.com")
	if err != nil {
		return fmt.Errorf("failed to create baidu cdn client: %w", err)
	}
	certName := fmt.Sprintf("%s_tlsctl_%d", domain, time.Now().UnixMilli())
	if _, err := client.PutCert(domain, &api.UserCertificate{
		CertName:    certName,
		ServerData:  string(certificate.Certificate),
		PrivateData: string(certificate.PrivateKey),
	}, "ON"); err != nil {
		return fmt.Errorf("failed to deploy baidu cdn cert: %w", err)
	}
	return nil
}
