package btpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
)

// Deploy 部署宝塔面板自身 SSL 证书。
type Deploy struct {
	client
}

func (d *Deploy) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *Deploy) Deploy(_ context.Context, certificate *certificate.Resource) error {
	data := url.Values{}
	data.Set("cert_type", "1")
	data.Set("privateKey", string(certificate.PrivateKey))
	data.Set("certPem", string(certificate.Certificate))
	if _, err := d.request(data, "config?action=SetPanelSSL"); err != nil {
		return fmt.Errorf("failed to deploy btpanel ssl: %w", err)
	}
	return nil
}

// DeploySite 部署证书到指定网站（支持多个，逗号分隔）。
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
		return fmt.Errorf("btpanel site_name is required")
	}
	sslHash, err := d.uploadCert(string(certificate.PrivateKey), string(certificate.Certificate))
	if err != nil {
		return err
	}
	batchInfo := make([]map[string]string, 0)
	for _, siteName := range strings.Split(d.SiteName, ",") {
		siteName = strings.TrimSpace(siteName)
		if siteName == "" {
			continue
		}
		batchInfo = append(batchInfo, map[string]string{
			"siteName": siteName,
			"ssl_hash": sslHash,
		})
	}
	if len(batchInfo) == 0 {
		return fmt.Errorf("btpanel site_name is required")
	}
	batchs, err := json.Marshal(batchInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal batch info: %w", err)
	}
	data := url.Values{}
	data.Set("BatchInfo", string(batchs))
	if _, err := d.request(data, "ssl?action=SetBatchCertToSite"); err != nil {
		return fmt.Errorf("failed to deploy btpanel site ssl: %w", err)
	}
	return nil
}

// DeployDockerSite 部署证书到宝塔 Docker 面板网站。
type DeployDockerSite struct {
	client
}

func (d *DeployDockerSite) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeployDockerSite) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.SiteName == "" {
		return fmt.Errorf("btpanel site_name is required")
	}
	data := url.Values{}
	data.Set("key", string(certificate.PrivateKey))
	data.Set("csr", string(certificate.Certificate))
	data.Set("siteName", d.SiteName)
	if _, err := d.request(data, "mod/docker/com/set_ssl"); err != nil {
		return fmt.Errorf("failed to deploy btpanel docker site ssl: %w", err)
	}
	return nil
}

// DeploySingleSite 部署证书到旧版本宝塔单个站点。
type DeploySingleSite struct {
	client
}

func (d *DeploySingleSite) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeploySingleSite) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.SiteName == "" {
		return fmt.Errorf("btpanel site_name is required")
	}
	data := url.Values{}
	data.Set("key", string(certificate.PrivateKey))
	data.Set("csr", string(certificate.Certificate))
	data.Set("siteName", d.SiteName)
	if _, err := d.request(data, "site?action=SetSSL"); err != nil {
		return fmt.Errorf("failed to deploy btpanel single site ssl: %w", err)
	}
	return nil
}
