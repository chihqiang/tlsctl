package volcengine

import (
	"fmt"
	"regexp"

	volccdn "github.com/volcengine/volcengine-go-sdk/service/cdn"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"github.com/volcengine/volcengine-go-sdk/volcengine/credentials"
	"github.com/volcengine/volcengine-go-sdk/volcengine/session"
)

// client 封装火山引擎 SDK 公共逻辑。
type client struct {
	Config
}

// newSession 创建火山引擎 SDK 会话。
func (c *client) newSession() (*session.Session, error) {
	cfg := volcengine.NewConfig().
		WithRegion(c.Region).
		WithCredentials(credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""))
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create volcengine session: %w", err)
	}
	return sess, nil
}

var duplicateCertRe = regexp.MustCompile(`cert-[a-f0-9]{32}`)

// uploadCert 上传证书到火山引擎证书中心，若已存在相同证书则复用其 ID。
func (c *client) uploadCert(cdnClient *volccdn.CDN, certPem, keyPem string) (string, error) {
	input := &volccdn.AddCertificateInput{
		Certificate: volcengine.String(certPem),
		PrivateKey:  volcengine.String(keyPem),
		Repeatable:  volcengine.Bool(false),
		Source:      volcengine.String("volc_cert_center"),
	}
	output, err := cdnClient.AddCertificate(input)
	if err != nil {
		if output != nil && output.Metadata != nil && output.Metadata.Error != nil && output.Metadata.Error.Code == "InvalidParameter.Certificate.Duplicated" {
			if certId := duplicateCertRe.FindString(output.Metadata.Error.Message); certId != "" {
				return certId, nil
			}
		}
		return "", fmt.Errorf("failed to upload volcengine cert: %w", err)
	}
	if output.CertId == nil {
		return "", fmt.Errorf("volcengine upload cert: empty cert id")
	}
	return *output.CertId, nil
}

func (c *client) domain(domain string) string {
	if c.Domain != "" {
		return c.Domain
	}
	return domain
}
