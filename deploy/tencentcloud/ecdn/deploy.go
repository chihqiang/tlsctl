package ecdn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
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
	clients, err := newClients(d.Config.SecretId, d.Config.SecretKey)
	if err != nil {
		return err
	}
	certId, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return err
	}
	// 获取待部署的 CDN 实例
	// 如果是泛域名，根据证书匹配 CDN 实例
	instanceIds := make([]string, 0)
	if strings.HasPrefix(d.Config.Domain, "*.") {
		domains, err := getDomainsByCertificateId(clients.CDN, certId)
		if err != nil {
			return err
		}

		instanceIds = domains
	} else {
		instanceIds = append(instanceIds, d.Config.Domain)
	}
	if len(instanceIds) == 0 {
		log.Info("no ecdn instances to deploy")
	} else {
		log.Info("found ecdn instances to deploy")

		// 证书部署到 CDN 实例
		// REF: https://cloud.tencent.com/document/product/400/91667
		deployCertificateInstanceReq := tcssl.NewDeployCertificateInstanceRequest()
		deployCertificateInstanceReq.CertificateId = common.StringPtr(certId)
		deployCertificateInstanceReq.ResourceType = common.StringPtr("cdn")
		deployCertificateInstanceReq.Status = common.Int64Ptr(1)
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs(instanceIds)
		deployCertificateInstanceResp, err := clients.SSL.DeployCertificateInstance(deployCertificateInstanceReq)
		log.Info("sdk request 'ssl.DeployCertificateInstance' request: %#v response: %#v", deployCertificateInstanceReq, deployCertificateInstanceResp)
		if err != nil {
			return fmt.Errorf("failed to execute sdk request 'ssl.DeployCertificateInstance': %w", err)
		}

		// 循环获取部署任务详情，等待任务状态变更
		// REF: https://cloud.tencent.com.cn/document/api/400/91658
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			describeHostDeployRecordDetailReq := tcssl.NewDescribeHostDeployRecordDetailRequest()
			describeHostDeployRecordDetailReq.DeployRecordId = common.StringPtr(fmt.Sprintf("%d", *deployCertificateInstanceResp.Response.DeployRecordId))
			describeHostDeployRecordDetailResp, err := clients.SSL.DescribeHostDeployRecordDetail(describeHostDeployRecordDetailReq)
			log.Info("sdk request 'ssl.DescribeHostDeployRecordDetail' request: %#v response:%#v", describeHostDeployRecordDetailReq, describeHostDeployRecordDetailResp)
			if err != nil {
				return fmt.Errorf("failed to execute sdk request 'ssl.DescribeHostDeployRecordDetail': %w", err)
			}
			var runningCount, succeededCount, failedCount, totalCount int64
			if describeHostDeployRecordDetailResp.Response.TotalCount == nil {
				return errors.New("unexpected deployment job status")
			} else {
				if describeHostDeployRecordDetailResp.Response.RunningTotalCount != nil {
					runningCount = *describeHostDeployRecordDetailResp.Response.RunningTotalCount
				}
				if describeHostDeployRecordDetailResp.Response.SuccessTotalCount != nil {
					succeededCount = *describeHostDeployRecordDetailResp.Response.SuccessTotalCount
				}
				if describeHostDeployRecordDetailResp.Response.FailedTotalCount != nil {
					failedCount = *describeHostDeployRecordDetailResp.Response.FailedTotalCount
				}
				if describeHostDeployRecordDetailResp.Response.TotalCount != nil {
					totalCount = *describeHostDeployRecordDetailResp.Response.TotalCount
				}

				if succeededCount+failedCount == totalCount {
					break
				}
			}
			log.Info("waiting for deployment job completion (running: %d, succeeded: %d, failed: %d, total: %d) ...", runningCount, succeededCount, failedCount, totalCount)
			time.Sleep(time.Second * 5)
		}
	}
	return nil
}
