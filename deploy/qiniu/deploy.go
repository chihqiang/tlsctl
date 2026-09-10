package qiniu

import (
	"context"
	"fmt"
	"net/http"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/qiniu/go-sdk/v7/auth"
	qiniucli "github.com/qiniu/go-sdk/v7/client"
)

type commonResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type sslCertResponse struct {
	CertID string `json:"certID"`
}

// client 封装七牛云开放 API 调用。
type client struct {
	Config
}

// request 调用七牛云开放 API。
func (c *client) request(path string, m map[string]any, method string, response any) error {
	if c.AccessKey == "" || c.AccessSecret == "" {
		return fmt.Errorf("qiniu access_key and access_secret are required")
	}
	credentials := auth.New(c.AccessKey, c.AccessSecret)
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	uri := fmt.Sprintf("https://api.qiniu.com/%s", path)
	if err := qiniucli.DefaultClient.CredentialedCallWithJson(context.Background(), credentials, auth.TokenQBox, response, method, uri, header, m); err != nil {
		return fmt.Errorf("qiniu api %s failed: %w", path, err)
	}
	return nil
}

// uploadCert 上传证书到七牛云，返回证书 ID。
func (c *client) uploadCert(certificate *certificate.Resource) (string, error) {
	m := map[string]any{
		"pri": string(certificate.PrivateKey),
		"ca":  string(certificate.Certificate),
	}
	var response sslCertResponse
	if err := c.request("sslcert", m, http.MethodPost, &response); err != nil {
		return "", err
	}
	if response.CertID == "" {
		return "", fmt.Errorf("qiniu upload cert: empty cert id")
	}
	return response.CertID, nil
}

func (c *client) domain(certificate *certificate.Resource) string {
	if c.Domain != "" {
		return c.Domain
	}
	return certificate.Domain
}

// DeployCdn 部署到七牛云 CDN。
type DeployCdn struct {
	client
}

func (d *DeployCdn) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeployCdn) Deploy(_ context.Context, certificate *certificate.Resource) error {
	certId, err := d.uploadCert(certificate)
	if err != nil {
		return err
	}
	var response commonResponse
	if err := d.request(fmt.Sprintf("domain/%s/sslize", d.domain(certificate)), map[string]any{"certid": certId}, http.MethodPut, &response); err != nil {
		return err
	}
	return nil
}

// DeployOss 部署到七牛云证书中心并绑定域名。
type DeployOss struct {
	client
}

func (d *DeployOss) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeployOss) Deploy(_ context.Context, certificate *certificate.Resource) error {
	certId, err := d.uploadCert(certificate)
	if err != nil {
		return err
	}
	var response commonResponse
	m := map[string]any{
		"certid": certId,
		"domain": d.domain(certificate),
	}
	if err := d.request("cert/bind", m, http.MethodPost, &response); err != nil {
		return err
	}
	return nil
}
