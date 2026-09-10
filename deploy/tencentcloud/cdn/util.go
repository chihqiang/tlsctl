package cdn

import (
	"fmt"

	tccdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func getDomainsByCertificateId(CDN *tccdn.Client, cloudCertId string) ([]string, error) {
	describeCertDomainsReq := tccdn.NewDescribeCertDomainsRequest()
	describeCertDomainsReq.CertId = common.StringPtr(cloudCertId)
	describeCertDomainsReq.Product = common.StringPtr("cdn")
	resp, err := CDN.DescribeCertDomains(describeCertDomainsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sdk request 'cdn.DescribeCertDomains': %w", err)
	}
	domains := make([]string, 0)
	if resp.Response.Domains != nil {
		for _, domain := range resp.Response.Domains {
			domains = append(domains, *domain)
		}
	}
	return domains, nil
}
