package safeline

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chihqiang/logx"
	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
)

// DeployPanel 部署雷池 WAF 面板自身 SSL 证书。
type DeployPanel struct {
	client
}

func (d *DeployPanel) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeployPanel) Deploy(_ context.Context, certificate *certificate.Resource) error {
	certId, err := d.uploadCert(0, string(certificate.PrivateKey), string(certificate.Certificate))
	if err != nil {
		return err
	}
	if _, err := d.request(map[string]any{"cert_id": certId}, http.MethodPut, "api/open/system"); err != nil {
		return fmt.Errorf("failed to deploy safeline panel ssl: %w", err)
	}
	return nil
}

// DeploySite 部署证书到雷池 WAF 站点（自动新建或更新站点证书）。
type DeploySite struct {
	client
}

func (d *DeploySite) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeploySite) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.SiteName == "" {
		return fmt.Errorf("safeline site_name is required")
	}
	certId, err := d.siteCertId(d.SiteName)
	if err != nil {
		return err
	}
	if certId == 0 {
		certId, err = d.uploadCert(0, string(certificate.PrivateKey), string(certificate.Certificate))
		if err != nil {
			return fmt.Errorf("safeline site %s upload cert failed: %w", d.SiteName, err)
		}
		logx.Info("safeline site %s uploaded cert id %d, bind manually if needed", d.SiteName, certId)
		return nil
	}
	if _, err := d.uploadCert(certId, string(certificate.PrivateKey), string(certificate.Certificate)); err != nil {
		return fmt.Errorf("safeline site %s update cert failed: %w", d.SiteName, err)
	}
	return nil
}

// DeployPortal 部署证书到雷池 WAF 认证中心。
type DeployPortal struct {
	client
}

func (d *DeployPortal) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeployPortal) Deploy(_ context.Context, certificate *certificate.Resource) error {
	response, err := d.request(map[string]any{}, http.MethodGet, "api/open/portal")
	if err != nil {
		return err
	}
	data, ok := response["data"].(map[string]any)
	if !ok {
		return fmt.Errorf("safeline portal config: invalid response")
	}
	var portalCertId int64
	if v, ok := data["cert_id"].(float64); ok {
		portalCertId = int64(v)
	}
	if portalCertId == 0 {
		certId, err := d.uploadCert(0, string(certificate.PrivateKey), string(certificate.Certificate))
		if err != nil {
			return fmt.Errorf("safeline portal upload cert failed: %w", err)
		}
		logx.Info("safeline portal uploaded cert id %d, bind manually if needed", certId)
		return nil
	}
	if _, err := d.uploadCert(portalCertId, string(certificate.PrivateKey), string(certificate.Certificate)); err != nil {
		return fmt.Errorf("safeline portal update cert failed: %w", err)
	}
	return nil
}
