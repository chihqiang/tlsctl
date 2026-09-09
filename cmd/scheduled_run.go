package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chihqiang/cli"
	"github.com/chihqiang/logx"
	"github.com/chihqiang/tlsctl/deploy"
	"github.com/go-acme/lego/v4/certificate"
)

func scheduledRunCommand() *cli.Command {
	return &cli.Command{
		Name:  "scheduled:run",
		Usage: "Automatically renew and complete deployment through scheduled tasks",
		Flags: scheduledRunFlags(),
		Action: func(ctx context.Context, in *cli.Input, _ *cli.Output) error {
			interval := in.Duration(flgInterval)
			logx.Info("Scheduled loop started, will check every %s", interval)
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			// 初次启动时先执行一次
			runScheduled(in)
			// 支持优雅退出
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			for {
				select {
				case <-ticker.C:
					runScheduled(in)
				case s := <-sig:
					logx.Info("Received signal: %v, exiting loop...", s)
					return nil
				}
			}
		},
	}
}

func runScheduled(in *cli.Input) {
	domainDeploys, err := deploy.JSONFileLoad(getDeployJson(in))
	if err != nil {
		logx.Warn("No plan to execute: %v", err)
		return
	}
	storage, err := setupResourceCache(in)
	if err != nil {
		logx.Warn("Failed to setup certificate cache: %v", err)
		return
	}
	for _, domainDeploy := range domainDeploys {
		var (
			renew    bool
			domain   = domainDeploy.Domain
			resource *certificate.Resource
		)
		resource, err = storage.ReadResource(domain)
		if err != nil {
			logx.Warn("Failed to read cert for %s: %v", domain, err)
			resource, err = buildLegoSSL(in, []string{domain})
			if err != nil {
				logx.Warn("Failed to obtain certificate for %s: %v", domain, err)
				continue
			}
			renew = true
		} else {
			cert, err := storage.ParseResourceFindCertificate(resource)
			if err != nil {
				logx.Warn("Invalid certificate for %s: %v", domain, err)
				continue
			}
			daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
			logx.Info("%s will expire in %d days", domain, daysLeft)
			if daysLeft < in.Int("day") {
				renew = true
				resource, err = buildLegoSSL(in, []string{domain})
				if err != nil {
					logx.Warn("Failed to renew certificate for %s: %v", domain, err)
					continue
				}
			}
		}
		if renew {
			logx.Warn("Renewing and deploying certificate for %s", domain)
			for _, deployName := range domainDeploy.Deploys {
				if err := deploy.RunWithJSONFile(getDeployJson(in), deployName, resource); err != nil {
					logx.Warn("Deployment failed: domain=%s deploy=%s err=%v", domain, deployName, err)
				} else {
					logx.Info("Deployment success: %s => %s", domain, deployName)
				}
			}
		}
	}
}
