package ecdn

import (
	"fmt"
	tccdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
)

type Clients struct {
	SSL *tcssl.Client
	CDN *tccdn.Client
}

func newClients(secretId, secretKey string) (*Clients, error) {
	credential := common.NewCredential(secretId, secretKey)
	sslClient, err := tcssl.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	cdnClient, err := tccdn.NewClient(credential, "", profile.NewClientProfile())
	if err != nil {
		return nil, err
	}
	return &Clients{
		SSL: sslClient,
		CDN: cdnClient,
	}, nil
}

func getDomainsByCertificateId(CDN *tccdn.Client, cloudCertId string) ([]string, error) {
	// 获取证书中的可用域名
	// REF: https://cloud.tencent.com/document/product/228/42491
	describeCertDomainsReq := tccdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(cloudCertId)
	describeCertDomainsReq.Product = common.StringPtr("ecdn")
	describeCertDomainsResp, err := CDN.DescribeCertDomains(describeCertDomainsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.DescribeCertDomains': %w", err)
	}
	domains := make([]string, 0)
	if describeCertDomainsResp.Response.Domains != nil {
		for _, domain := range describeCertDomainsResp.Response.Domains {
			domains = append(domains, *domain)
		}
	}
	return domains, nil
}
