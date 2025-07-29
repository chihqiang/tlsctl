package local

import (
	"context"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/pkg/errors"
)

type Deploy struct {
	Config *Config
}

func (p *Deploy) WithEnvConfig() error {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}
	p.Config = &cfg
	return nil
}
func (p *Deploy) Deploy(ctx context.Context, certificate *certificate.Resource) error {
	// 执行前置命令
	if p.Config.PreCommand != "" {
		stdout, stderr, err := ExecCommand(p.Config.PreCommand, "")
		if err != nil {
			return errors.Wrapf(err, "failed to execute pre-command, stdout: %s, stderr: %s", stdout, stderr)
		}
		log.Info("pre-command executed %s", stdout)
	}
	if err := CopyFile(p.Config.CertPath, certificate.Certificate); err != nil {
		return errors.Wrap(err, "failed to save certificate file")
	}
	log.Info("certificate file saved to %s", p.Config.CertPath)
	if err := CopyFile(p.Config.KeyPath, certificate.PrivateKey); err != nil {
		return errors.Wrap(err, "failed to save private key file")
	}
	log.Info("private key file saved to %s", p.Config.KeyPath)
	// 执行后置命令
	if p.Config.PostCommand != "" {
		stdout, stderr, err := ExecCommand(p.Config.PostCommand, "")
		if err != nil {
			return errors.Wrapf(err, "failed to execute post-command, stdout: %s, stderr: %s", stdout, stderr)
		}
		log.Info("post-command executed %s", stdout)
	}
	return nil
}
