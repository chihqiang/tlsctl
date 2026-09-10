package huaweicloud

import (
	"context"
	"fmt"
	"time"

	"github.com/chihqiang/logx"
	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	cdn "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/region"
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

// Deploy 更新华为云 CDN 域名证书。
func (d *Deploy) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.Config.AccessKey == "" || d.Config.SecretKey == "" {
		return fmt.Errorf("huawei access_key and secret_key are required")
	}
	domain := d.Config.Domain
	if domain == "" {
		domain = certificate.Domain
	}

	auth, err := global.NewCredentialsBuilder().WithAk(d.Config.AccessKey).WithSk(d.Config.SecretKey).SafeBuild()
	if err != nil {
		return fmt.Errorf("failed to build huawei credentials: %w", err)
	}
	Region, err := region.SafeValueOf("cn-north-1")
	if err != nil {
		return fmt.Errorf("failed to get huawei cdn region: %w", err)
	}
	builder, err := cdn.CdnClientBuilder().WithRegion(Region).WithCredential(auth).SafeBuild()
	if err != nil {
		return fmt.Errorf("failed to build huawei cdn client: %w", err)
	}
	client := cdn.NewCdnClient(builder)

	certStr := string(certificate.Certificate)
	keyStr := string(certificate.PrivateKey)
	certName := fmt.Sprintf("tlsctl(%s)", time.Now().Format(time.RFC3339))
	request := &model.UpdateDomainMultiCertificatesRequest{}
	httpsbody := &model.UpdateDomainMultiCertificatesRequestBodyContent{
		DomainName:  domain,
		HttpsSwitch: 1,
		CertName:    &certName,
		Certificate: &certStr,
		PrivateKey:  &keyStr,
	}
	request.Body = &model.UpdateDomainMultiCertificatesRequestBody{
		Https: httpsbody,
	}
	if _, err := client.UpdateDomainMultiCertificates(request); err != nil {
		return fmt.Errorf("failed to update huawei domain certificates: %w", err)
	}
	logx.Info("huawei cdn domain %s certificate updated", domain)
	return nil
}
