package lecdn

import (
	"context"
	"fmt"

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

// Deploy 上传证书并绑定到 LeCDN 站点域名。
func (d *Deploy) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.Config.Url == "" || d.Config.Username == "" || d.Config.Password == "" {
		return fmt.Errorf("lecdn url, username and password are required")
	}
	if d.Config.SiteId == 0 {
		return fmt.Errorf("lecdn site_id is required")
	}
	domain := d.Config.Domain
	if domain == "" {
		domain = certificate.Domain
	}
	certPem := string(certificate.Certificate)

	c := &client{
		baseURL:  d.Config.Url,
		username: d.Config.Username,
		password: d.Config.Password,
		ignore:   d.Config.IgnoreSSL,
	}
	if err := c.login(); err != nil {
		return err
	}

	name := certName(certPem)
	certId, err := c.findCertId(name)
	if err != nil {
		return err
	}
	if certId == 0 {
		certId, err = c.uploadCert(certPem, string(certificate.PrivateKey), name)
		if err != nil {
			return err
		}
	}

	domainId, err := c.domainIdFromSite(d.Config.SiteId, domain)
	if err != nil {
		return err
	}
	if err := c.deployCert(d.Config.SiteId, domainId, certId); err != nil {
		return err
	}
	return nil
}
