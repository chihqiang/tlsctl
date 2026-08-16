package ssh

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/logx"
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
	// 连接
	client, err := newSSHClient(d.Config.Host, d.Config.Port, d.Config.Username, d.Config.Password, d.Config.Key, d.Config.KeyPassphrase)
	if err != nil {
		return fmt.Errorf("failed to create ssh client: %w", err)
	}
	defer client.Close()
	logx.Info("SSH connected")
	// 执行前置命令
	if d.Config.PreCommand != "" {
		stdout, stderr, err := execSSHCommand(client, d.Config.PreCommand)
		if err != nil {
			return fmt.Errorf("failed to execute pre-command, stdout: %s, stderr: %s: %w", stdout, stderr, err)
		}
		logx.Info("SSH pre-command executed %s", stdout)
	}
	if err := writeFile(client, d.Config.UseSCP, d.Config.CertPath, certificate.Certificate); err != nil {
		return fmt.Errorf("failed to upload certificate file: %w", err)
	}
	logx.Info("certificate file uploaded to %s", d.Config.CertPath)
	if err := writeFile(client, d.Config.UseSCP, d.Config.KeyPath, certificate.PrivateKey); err != nil {
		return fmt.Errorf("failed to upload private key file: %w", err)
	}
	logx.Info("private key file uploaded to %s", d.Config.KeyPath)
	// 执行后置命令
	if d.Config.PostCommand != "" {
		stdout, stderr, err := execSSHCommand(client, d.Config.PostCommand)
		if err != nil {
			return fmt.Errorf("failed to execute post-command, stdout: %s, stderr: %s: %w", stdout, stderr, err)
		}
		logx.Info("SSH post-command executed %s", stdout)
	}
	return nil
}
