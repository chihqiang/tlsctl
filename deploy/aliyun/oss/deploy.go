package oss

import (
	"context"
	"errors"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/caarlos0/env/v11"
	"github.com/go-acme/lego/v4/certificate"
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
	client, err := newClient(d.Config.AccessKeyId, d.Config.AccessKeySecret, d.Config.Region)
	if err != nil {
		return fmt.Errorf("failed to create sdk client: %w", err)
	}
	if d.Config.Bucket == "" {
		return errors.New("config `bucket` is required")
	}
	if d.Config.Domain == "" {
		return errors.New("config `domain` is required")
	}
	return client.PutBucketCnameWithCertificate(d.Config.Bucket, oss.PutBucketCname{
		Cname: d.Config.Domain,
		CertificateConfiguration: &oss.CertificateConfiguration{
			Certificate: string(certificate.Certificate),
			PrivateKey:  string(certificate.PrivateKey),
			Force:       true,
		},
	})
}
