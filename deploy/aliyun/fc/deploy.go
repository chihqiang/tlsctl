package fc

import (
	"context"
	"fmt"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/certificate"

	"time"

	alifc3 "github.com/alibabacloud-go/fc-20230330/v4/client"
	alifc2 "github.com/alibabacloud-go/fc-open-20210406/v2/client"
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
	clients, err := NewClients(d.Config.AccessKeyId, d.Config.AccessKeySecret, d.Config.Region)
	if err != nil {
		return err
	}
	switch d.Config.Version {
	case "3", "3.0":
		if err := d.deployToFC3(clients.FC3, d.Config.Domain, certificate); err != nil {
			return err
		}

	case "2", "2.0":
		if err := d.deployToFC2(clients.FC2, d.Config.Domain, certificate); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported service version '%s'", d.Config.Version)
	}
	return nil
}

func (d *Deploy) deployToFC3(FC3 *alifc3.Client, domain string, certificate *certificate.Resource) error {
	// 获取自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-getcustomdomain
	getCustomDomainResp, err := FC3.GetCustomDomain(tea.String(domain))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.GetCustomDomain': %w", err)
	}
	// 更新自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-3-0/developer-reference/api-fc-2023-03-30-updatecustomdomain
	updateCustomDomainReq := &alifc3.UpdateCustomDomainRequest{
		Body: &alifc3.UpdateCustomDomainInput{
			CertConfig: &alifc3.CertConfig{
				CertName:    tea.String(fmt.Sprintf("tlsctl-%d", time.Now().UnixMilli())),
				Certificate: tea.String(string(certificate.Certificate)),
				PrivateKey:  tea.String(string(certificate.PrivateKey)),
			},
			Protocol:  getCustomDomainResp.Body.Protocol,
			TlsConfig: getCustomDomainResp.Body.TlsConfig,
		},
	}
	if tea.StringValue(updateCustomDomainReq.Body.Protocol) == "HTTP" {
		updateCustomDomainReq.Body.Protocol = tea.String("HTTP,HTTPS")
	}
	_, err = FC3.UpdateCustomDomain(tea.String(d.Config.Domain), updateCustomDomainReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.UpdateCustomDomain': %w", err)
	}
	return nil
}
func (d *Deploy) deployToFC2(FC2 *alifc2.Client, domain string, certificate *certificate.Resource) error {
	// 获取自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-2-0/developer-reference/api-fc-open-2021-04-06-getcustomdomain
	getCustomDomainResp, err := FC2.GetCustomDomain(tea.String(domain))
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.GetCustomDomain': %w", err)
	}
	// 更新自定义域名
	// REF: https://help.aliyun.com/zh/functioncompute/fc-2-0/developer-reference/api-fc-open-2021-04-06-updatecustomdomain
	updateCustomDomainReq := &alifc2.UpdateCustomDomainRequest{
		CertConfig: &alifc2.CertConfig{
			CertName:    tea.String(fmt.Sprintf("tlsctl-%d", time.Now().UnixMilli())),
			Certificate: tea.String(string(certificate.Certificate)),
			PrivateKey:  tea.String(string(certificate.PrivateKey)),
		},
		Protocol:  getCustomDomainResp.Body.Protocol,
		TlsConfig: getCustomDomainResp.Body.TlsConfig,
	}
	if tea.StringValue(updateCustomDomainReq.Protocol) == "HTTP" {
		updateCustomDomainReq.Protocol = tea.String("HTTP,HTTPS")
	}
	_, err = FC2.UpdateCustomDomain(tea.String(d.Config.Domain), updateCustomDomainReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'fc.UpdateCustomDomain': %w", err)
	}
	return nil
}
