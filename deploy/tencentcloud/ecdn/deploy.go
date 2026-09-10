package ecdn

import (
	"context"
	"fmt"
	"strings"

	"github.com/chihqiang/logx"
	dcommon "github.com/chihqiang/tlsctl/deploy/common"
	tccommon "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
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

func (d *Deploy) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	clients, err := tccommon.NewClients(d.Config.SecretId, d.Config.SecretKey)
	if err != nil {
		return err
	}

	certId, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return err
	}

	instanceIds, err := d.resolveDomains(clients, certId)
	if err != nil {
		return err
	}

	if len(instanceIds) == 0 {
		logx.Info("no ecdn instances to deploy")
		return nil
	}

	logx.Info("found ecdn instances to deploy %s", instanceIds)

	deployCertificateInstanceReq := tcssl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(certId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cdn")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(instanceIds)
	resp, err := clients.SSL.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'ssl.DeployCertificateInstance': %w", err)
	}

	return tccommon.WaitForDeploy(ctx, clients.SSL, *resp.Response.DeployRecordId)
}

func (d *Deploy) resolveDomains(clients *tccommon.Clients, certId string) ([]string, error) {
	instanceIds := make([]string, 0)

	if strings.HasPrefix(d.Config.Domain, "*.") {
		domains, err := getDomainsByCertificateId(clients.CDN, certId)
		if err != nil {
			return nil, err
		}
		instanceIds = domains
	} else {
		instanceIds = append(instanceIds, d.Config.Domain)
	}

	return instanceIds, nil
}
