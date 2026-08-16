package ssl

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Deploy struct {
	Config *Config
	certId string
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
	// REF: https://cloud.tencent.com/document/product/400/41665
	uploadCertificateReq := tcssl.NewUploadCertificateRequest()
	uploadCertificateReq.CertificatePublicKey = common.StringPtr(string(certificate.Certificate))
	uploadCertificateReq.CertificatePrivateKey = common.StringPtr(string(certificate.PrivateKey))
	uploadCertificateReq.Repeatable = common.BoolPtr(false)
	uploadCertificateResp, err := client.UploadCertificate(uploadCertificateReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'ssl.UploadCertificate': %w", err)
	}
	d.certId = *uploadCertificateResp.Response.CertificateId
	return nil
}

func (d *Deploy) GetCertId() string {
	return d.certId
}
