package local

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/caarlos0/env/v11"
	"github.com/chihqiang/logx"
	"github.com/chihqiang/tlsctl/resource"
	"github.com/go-acme/lego/v4/certificate"
)

// deployLocalPath 是本地部署未配置输出路径时的默认证书目录。
const deployLocalPath = "/etc/nginx/ssl/"

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
	// 未配置输出路径时，使用默认目录 /etc/nginx/ssl/
	if p.Config.CertPath == "" && p.Config.KeyPath == "" {
		sanitizedDomain, err := resource.SanitizedDomain(certificate.Domain)
		if err != nil {
			return fmt.Errorf("failed to sanitize domain: %w", err)
		}
		p.Config.CertPath = filepath.Join(deployLocalPath, sanitizedDomain+resource.PemExt)
		p.Config.KeyPath = filepath.Join(deployLocalPath, sanitizedDomain+resource.KeyExt)
	}
	// 执行前置命令
	if p.Config.PreCommand != "" {
		stdout, stderr, err := ExecCommand(p.Config.PreCommand, "")
		if err != nil {
			return fmt.Errorf("failed to execute pre-command, stdout: %s, stderr: %s: %w", stdout, stderr, err)
		}
		logx.Info("pre-command executed %s", stdout)
	}
	if err := CopyFile(p.Config.CertPath, certificate.Certificate); err != nil {
		return fmt.Errorf("failed to save certificate file: %w", err)
	}
	logx.Info("certificate file saved to %s", p.Config.CertPath)
	if err := CopyFile(p.Config.KeyPath, certificate.PrivateKey); err != nil {
		return fmt.Errorf("failed to save private key file: %w", err)
	}
	logx.Info("private key file saved to %s", p.Config.KeyPath)
	// 执行后置命令
	if p.Config.PostCommand != "" {
		stdout, stderr, err := ExecCommand(p.Config.PostCommand, "")
		if err != nil {
			return fmt.Errorf("failed to execute post-command, stdout: %s, stderr: %s: %w", stdout, stderr, err)
		}
		logx.Info("post-command executed %s", stdout)
	}
	return nil
}
