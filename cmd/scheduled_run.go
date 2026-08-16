package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chihqiang/tlsctl/deploy"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/urfave/cli/v3"
)

func scheduledRunCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "scheduled:run",
		Usage:                  "Automatically renew and complete deployment through scheduled tasks",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "day",
				Usage:   "When the expiration date is less than a few days, it will be regenerated",
				Value:   1,
				Sources: cli.EnvVars(envDay),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			interval := cmd.Duration(flgInterval)
			log.Info("Scheduled loop started, will check every %s", interval)
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			// 初次启动时先执行一次
			runScheduled(cmd)
			// 支持优雅退出
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			for {
				select {
				case <-ticker.C:
					runScheduled(cmd)
				case s := <-sig:
					log.Info("Received signal: %v, exiting loop...", s)
					return nil
				}
			}
		},
	}
}

func runScheduled(cmd *cli.Command) {
	domainDeploys, err := deploy.JSONFileLoad(getDeployJson(cmd))
	if err != nil {
		log.Warn("No plan to execute %d", err)
		return
	}
	storage := setupResourceCache(cmd)
	for _, domainDeploy := range domainDeploys {
		var (
			renew    bool
			domain   = domainDeploy.Domain
			resource *certificate.Resource
		)
		resource, err = storage.ReadResource(domain)
		if err != nil {
			log.Warn("Failed to read cert for %s: %v", domain, err)
			resource, err = buildLegoSSL(cmd, domain)
			if err != nil {
				log.Warn("Invalid certificate for %s: %v", domain, err)
			}
			renew = true
		} else {
			cert, err := storage.ParseResourceFindCertificate(resource)
			if err != nil {
				log.Warn("Invalid certificate for %s: %v", domain, err)
				continue
			}
			daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
			log.Info("%s will expire in %d days", domain, daysLeft)
			if daysLeft < cmd.Int("day") {
				renew = true
				resource, err = buildLegoSSL(cmd, domain)
				if err != nil {
					log.Warn("Invalid certificate Day for %s: %v", domain, err)
				}
			}
		}
		if renew {
			log.Warn("Renewing and deploying certificate for  %s", domain)
			for _, deployName := range domainDeploy.Deploys {
				if err := deploy.RunWithJSONFile(getDeployJson(cmd), deployName, resource); err != nil {
					log.Warn("Deployment failed: domain=%s deploy=%s err=%v", domain, deployName, err)
				} else {
					log.Info("Deployment success: %s => %s", domain, deployName)
				}
			}
		}
	}
}
