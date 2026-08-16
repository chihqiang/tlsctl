package clb

type Config struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId" yaml:"secretId" xml:"SecretId" env:"TENCENTCLOUD_SECRET_ID"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey" yaml:"secretKey" xml:"SecretKey" env:"TENCENTCLOUD_SECRET_KEY"`
	// 腾讯云地域。
	Region string `json:"region" yaml:"region" xml:"Region" env:"TENCENTCLOUD_REGION"`
	// 部署资源类型。 ssl-deploy loadbalancer listener ruledomain
	ResourceType string `json:"resourceType" yaml:"resourceType" xml:"ResourceType" env:"TENCENTCLOUD_RESOURCE_TYPE"`
	// 负载均衡器 ID。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY]、[RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_RULEDOMAIN] 时必填。
	LoadbalancerId string `json:"loadbalancerId,omitempty" yaml:"loadbalancerId,omitempty" xml:"LoadbalancerId,omitempty" env:"TENCENTCLOUD_LOADBALANCER_ID"`
	// 负载均衡监听 ID。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY]、[RESOURCE_TYPE_LOADBALANCER]、[RESOURCE_TYPE_LISTENER]、[RESOURCE_TYPE_RULEDOMAIN] 时必填。
	ListenerId string `json:"listenerId,omitempty" yaml:"listenerId,omitempty" xml:"ListenerId,omitempty" env:"TENCENTCLOUD_LISTENER_ID"`
	// SNI 域名或七层转发规则域名（支持泛域名）。
	// 部署资源类型为 [RESOURCE_TYPE_SSLDEPLOY] 时选填；部署资源类型为 [RESOURCE_TYPE_RULEDOMAIN] 时必填。
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty" xml:"Domain,omitempty" env:"TENCENTCLOUD_DOMAIN"`
}
