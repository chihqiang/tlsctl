package onepanel

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
)

// Deploy 部署 1Panel 面板自身 SSL 证书。
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
	data := map[string]any{
		"cert":    string(certificate.Certificate),
		"key":     string(certificate.PrivateKey),
		"ssl":     "enable",
		"sslType": "import-paste",
	}
	if _, err := d.request(data, http.MethodPost, "settings/ssl/update"); err != nil {
		return fmt.Errorf("failed to deploy onepanel ssl: %w", err)
	}
	return nil
}

// DeploySite 部署证书到指定网站。
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
	if d.SiteId == "" {
		return fmt.Errorf("onepanel site_id is required")
	}
	siteId := d.SiteId

	siteData, err := d.request(map[string]any{}, http.MethodGet, fmt.Sprintf("websites/%s/https", siteId))
	if err != nil {
		return fmt.Errorf("failed to get onepanel site config: %w", err)
	}
	data, ok := siteData["data"].(map[string]any)
	if !ok {
		return fmt.Errorf("failed to get onepanel site config: data not found")
	}

	SSLProtocol, ok := data["SSLProtocol"].([]any)
	if !ok || len(SSLProtocol) == 0 {
		SSLProtocol = []any{"TLSv1.3", "TLSv1.2", "TLSv1.1", "TLSv1"}
	}
	algorithm, ok := data["algorithm"].(string)
	if !ok {
		return fmt.Errorf("failed to get onepanel site config: algorithm not found")
	}
	hsts, ok := data["hsts"].(bool)
	if !ok {
		return fmt.Errorf("failed to get onepanel site config: hsts not found")
	}
	httpConfig, ok := data["httpConfig"].(string)
	if !ok {
		return fmt.Errorf("failed to get onepanel site config: httpConfig not found")
	}
	if httpConfig == "" {
		httpConfig = "HTTPToHTTPS"
	}
	websiteId, err := strconv.Atoi(siteId)
	if err != nil {
		return fmt.Errorf("invalid onepanel site_id: %w", err)
	}

	deployData := map[string]any{
		"SSLProtocol": SSLProtocol,
		"algorithm":   algorithm,
		"certificate": string(certificate.Certificate),
		"privateKey":  string(certificate.PrivateKey),
		"enable":      true,
		"hsts":        hsts,
		"httpConfig":  httpConfig,
		"importType":  "paste",
		"type":        "manual",
		"websiteId":   websiteId,
	}
	if _, err := d.request(deployData, http.MethodPost, fmt.Sprintf("websites/%s/https", siteId)); err != nil {
		return fmt.Errorf("failed to deploy onepanel site ssl: %w", err)
	}
	return nil
}
