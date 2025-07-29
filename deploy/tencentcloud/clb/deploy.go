package clb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	tcclb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
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
	client, err := newClient(d.Config.SecretId, d.Config.SecretKey, d.Config.Region)
	if err != nil {
		return err
	}
	id, err := ssl.FastDeploy(ctx, d.Config.SecretId, d.Config.SecretKey, certificate)
	if err != nil {
		return fmt.Errorf("failed to deploy SSL certificate: %w", err)
	}
	// 根据部署资源类型决定部署方式
	switch d.Config.ResourceType {
	case "ssl-deploy":
		if err := d.deployViaSslService(ctx, client.SSL, id); err != nil {
			return err
		}
	case "loadbalancer":
		if err := d.deployToLoadbalancer(ctx, client.CLB, id); err != nil {
			return err
		}

	case "listener":
		if err := d.deployToListener(ctx, client.CLB, id); err != nil {
			return err
		}

	case "ruledomain":
		if err := d.deployToRuleDomain(ctx, client.CLB, id); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported resource type '%s'", d.Config.ResourceType)
	}
	return nil
}

func (d *Deploy) deployViaSslService(ctx context.Context, SSL *tcssl.Client, cloudCertId string) error {
	if d.Config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.Config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 证书部署到 CLB 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcssl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(cloudCertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("clb")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	if d.Config.Domain == "" {
		// 未指定 SNI，只需部署到监听器
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s", d.Config.LoadbalancerId, d.Config.ListenerId)})
	} else {
		// 指定 SNI，需部署到域名
		deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s|%s|%s", d.Config.LoadbalancerId, d.Config.ListenerId, d.Config.Domain)})
	}
	deployCertificateInstanceResp, err := SSL.DeployCertificateInstance(deployCertificateInstanceReq)
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
		describeHostDeployRecordDetailResp, err := SSL.DescribeHostDeployRecordDetail(describeHostDeployRecordDetailReq)
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

	return nil
}

func (d *Deploy) deployToLoadbalancer(ctx context.Context, CLB *tcclb.Client, cloudCertId string) error {
	if d.Config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}

	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	listenerIds := make([]string, 0)
	describeListenersReq := tcclb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(d.Config.LoadbalancerId)
	describeListenersResp, err := CLB.DescribeListeners(describeListenersReq)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'clb.DescribeListeners': %w", err)
	} else {
		if describeListenersResp.Response.Listeners != nil {
			for _, listener := range describeListenersResp.Response.Listeners {
				if listener.Protocol == nil || (*listener.Protocol != "HTTPS" && *listener.Protocol != "TCP_SSL" && *listener.Protocol != "QUIC") {
					continue
				}

				listenerIds = append(listenerIds, *listener.ListenerId)
			}
		}
	}

	// 遍历更新监听器证书
	if len(listenerIds) == 0 {
		log.Info("no clb listeners to deploy")
	} else {
		log.Info("found https/tcpssl/quic listeners to deploy")
		var errs []error

		for _, listenerId := range listenerIds {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if err := d.modifyListenerCertificate(ctx, CLB, d.Config.LoadbalancerId, listenerId, cloudCertId); err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}
	}

	return nil
}

func (d *Deploy) deployToListener(ctx context.Context, CLB *tcclb.Client, cloudCertId string) error {
	if d.Config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.Config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}

	// 更新监听器证书
	if err := d.modifyListenerCertificate(ctx, CLB, d.Config.LoadbalancerId, d.Config.ListenerId, cloudCertId); err != nil {
		return err
	}

	return nil
}

func (d *Deploy) deployToRuleDomain(ctx context.Context, CLB *tcclb.Client, cloudCertId string) error {
	if d.Config.LoadbalancerId == "" {
		return errors.New("config `loadbalancerId` is required")
	}
	if d.Config.ListenerId == "" {
		return errors.New("config `listenerId` is required")
	}
	if d.Config.Domain == "" {
		return errors.New("config `domain` is required")
	}

	// 修改负载均衡七层监听器转发规则的域名级别属性
	// REF: https://cloud.tencent.com/document/api/214/38092
	modifyDomainAttributesReq := tcclb.NewModifyDomainAttributesRequest()
	modifyDomainAttributesReq.LoadBalancerId = common.StringPtr(d.Config.LoadbalancerId)
	modifyDomainAttributesReq.ListenerId = common.StringPtr(d.Config.ListenerId)
	modifyDomainAttributesReq.Domain = common.StringPtr(d.Config.Domain)
	modifyDomainAttributesReq.Certificate = &tcclb.CertificateInput{
		SSLMode: common.StringPtr("UNIDIRECTIONAL"),
		CertId:  common.StringPtr(cloudCertId),
	}
	modifyDomainAttributesResp, err := CLB.ModifyDomainAttributes(modifyDomainAttributesReq)
	log.Info("sdk request 'clb.ModifyDomainAttributes'  request: %#v response: %#v", modifyDomainAttributesReq, modifyDomainAttributesResp)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'clb.ModifyDomainAttributes': %w", err)
	}

	return nil
}

func (d *Deploy) modifyListenerCertificate(ctx context.Context, CLB *tcclb.Client, cloudLoadbalancerId, cloudListenerId, cloudCertId string) error {
	// 查询监听器列表
	// REF: https://cloud.tencent.com/document/api/214/30686
	describeListenersReq := tcclb.NewDescribeListenersRequest()
	describeListenersReq.LoadBalancerId = common.StringPtr(cloudLoadbalancerId)
	describeListenersReq.ListenerIds = common.StringPtrs([]string{cloudListenerId})
	describeListenersResp, err := CLB.DescribeListeners(describeListenersReq)
	log.Info("sdk request 'clb.DescribeListeners' request: %#v response: %#v", describeListenersReq, describeListenersResp)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'clb.DescribeListeners': %w", err)
	} else if len(describeListenersResp.Response.Listeners) == 0 {
		return errors.New("listener not found")
	}

	// 修改监听器属性
	// REF: https://cloud.tencent.com/document/product/214/30681
	modifyListenerReq := tcclb.NewModifyListenerRequest()
	modifyListenerReq.LoadBalancerId = common.StringPtr(cloudLoadbalancerId)
	modifyListenerReq.ListenerId = common.StringPtr(cloudListenerId)
	modifyListenerReq.Certificate = &tcclb.CertificateInput{CertId: common.StringPtr(cloudCertId)}
	if describeListenersResp.Response.Listeners[0].Certificate != nil && describeListenersResp.Response.Listeners[0].Certificate.SSLMode != nil {
		modifyListenerReq.Certificate.SSLMode = describeListenersResp.Response.Listeners[0].Certificate.SSLMode
		modifyListenerReq.Certificate.CertCaId = describeListenersResp.Response.Listeners[0].Certificate.CertCaId
	} else {
		modifyListenerReq.Certificate.SSLMode = common.StringPtr("UNIDIRECTIONAL")
	}
	modifyListenerResp, err := CLB.ModifyListener(modifyListenerReq)
	log.Info("sdk request 'clb.ModifyListener' request: %#v response: %#v", modifyListenerReq, modifyListenerResp)
	if err != nil {
		return fmt.Errorf("failed to execute sdk request 'clb.ModifyListener': %w", err)
	}

	return nil
}
