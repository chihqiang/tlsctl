package volcengine

import (
	"context"
	"fmt"

	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	"github.com/go-acme/lego/v4/certificate"
	volccdn "github.com/volcengine/volcengine-go-sdk/service/cdn"
	volcdcdn "github.com/volcengine/volcengine-go-sdk/service/dcdn"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

// DeployCdn 部署证书到火山引擎 CDN。
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
	if d.AccessKey == "" || d.SecretKey == "" {
		return fmt.Errorf("volcengine access_key and secret_key are required")
	}
	if d.Region == "" {
		return fmt.Errorf("volcengine region is required")
	}
	sess, err := d.newSession()
	if err != nil {
		return err
	}
	cdnClient := volccdn.New(sess)

	certId, err := d.uploadCert(cdnClient, string(certificate.Certificate), string(certificate.PrivateKey))
	if err != nil {
		return err
	}
	input := &volccdn.BatchDeployCertInput{
		CertId: volcengine.String(certId),
		Domain: volcengine.String(d.domain(certificate.Domain)),
	}
	res, err := cdnClient.BatchDeployCert(input)
	if err != nil {
		return fmt.Errorf("failed to deploy volcengine cdn cert: %w", err)
	}
	if len(res.DeployResult) > 0 && res.DeployResult[0].Status != nil && *res.DeployResult[0].Status != "success" {
		msg := ""
		if res.DeployResult[0].ErrorMsg != nil {
			msg = *res.DeployResult[0].ErrorMsg
		}
		return fmt.Errorf("volcengine cdn deploy failed: %s", msg)
	}
	return nil
}

// DeployDcdn 部署证书到火山引擎 DCDN。
type DeployDcdn struct {
	client
}

func (d *DeployDcdn) WithEnvConfig() error {
	cfg, err := dcommon.ParseConfig[Config]()
	if err != nil {
		return err
	}
	d.client.Config = *cfg
	return nil
}

func (d *DeployDcdn) Deploy(_ context.Context, certificate *certificate.Resource) error {
	if d.AccessKey == "" || d.SecretKey == "" {
		return fmt.Errorf("volcengine access_key and secret_key are required")
	}
	if d.Region == "" {
		return fmt.Errorf("volcengine region is required")
	}
	sess, err := d.newSession()
	if err != nil {
		return err
	}
	cdnClient := volccdn.New(sess)
	dcdnClient := volcdcdn.New(sess)

	certId, err := d.uploadCert(cdnClient, string(certificate.Certificate), string(certificate.PrivateKey))
	if err != nil {
		return err
	}
	input := &volcdcdn.CreateCertBindInput{
		CertId:      volcengine.String(certId),
		DomainNames: volcengine.StringSlice([]string{d.domain(certificate.Domain)}),
	}
	if _, err := dcdnClient.CreateCertBind(input); err != nil {
		return fmt.Errorf("failed to deploy volcengine dcdn cert: %w", err)
	}
	return nil
}
