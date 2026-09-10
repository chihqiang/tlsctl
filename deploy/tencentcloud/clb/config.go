package clb

import "github.com/chihqiang/tlsctl/deploy/tencentcloud/common"

type Config struct {
	common.BaseConfig
	// 腾讯云地域。
	Region string `json:"region" yaml:"region" xml:"region" env:"TENCENTCLOUD_REGION"`
	// 部署资源类型。 ssl-deploy loadbalancer listener ruledomain
	ResourceType string `json:"resource_type" yaml:"resourceType" xml:"resourceType" env:"TENCENTCLOUD_RESOURCE_TYPE"`
	// 负载均衡器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY]、[RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_RULEDOMAIN] 时必填。
	LoadbalancerId string `json:"loadbalancer_id,omitempty" yaml:"loadbalancerId,omitempty" xml:"loadbalancerId,omitempty" env:"TENCENTCLOUD_LOADBALANCER_ID"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY]、[RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER]、[RESOURCE_TYPE_RULEDOMAIN] 时必填。
	ListenerId string `json:"listener_id,omitempty" yaml:"listenerId,omitempty" xml:"listenerId,omitempty" env:"TENCENTCLOUD_LISTENER_ID"`
	// SNI 域名或七层转发规则域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY] 时选填；部署资源类型为 [RESOURCE_TYPE_RULEDOMAIN] 时必填。
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty" xml:"Domain,omitempty" env:"TENCENTCLOUD_DOMAIN"`
}
