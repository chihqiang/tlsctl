package cmd

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/chihqiang/logx"
	"github.com/chihqiang/tlsctl/localhost"
	"github.com/chihqiang/tlsctl/pkg/fp"
	"github.com/urfave/cli/v3"
)

func localhostCommand() *cli.Command {
	return &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "localhost",
		Usage:                  "Build local development ssl certificate",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name: "hosts",
				Value: []string{
					"localhost",
					"127.0.0.1",
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return buildLocalHostSSL(cmd)
		},
	}
}

func buildLocalHostSSL(cmd *cli.Command) error {
	cStorage, err := setupResourceCache(cmd)
	if err != nil {
		return err
	}
	hostSSL, err := localhost.NewLocalHostSSL(path.Join(cmd.String(flgPath), "certificates", "localhost"))
	if err != nil {
		return err
	}
	if err = hostSSL.LoadCA(); err != nil {
		return err
	}
	resource, err := hostSSL.BuildResource(cmd.StringSlice("hosts"))
	if err != nil {
		return err
	}
	if err := cStorage.SaveResource(resource); err != nil {
		return fmt.Errorf("error saving certificate: %w", err)
	}
	logx.Debug("Certificate for %s has been saved successfully at %s",
		"localhost",
		cStorage.GetSanitizedDomainSavePath("localhost"),
	)
	localhostSSLInstallHelp(path.Join(hostSSL.Path, "rootCA.pem"))
	return nil
}

func localhostSSLInstallHelp(rootCAPath string) {
	switch runtime.GOOS {
	case "darwin":
		logx.Info("🔐 [macOS] Trust commands for the root certificate:")
		logx.Info("🛠 Install: sudo security add-trusted-cert -d -k /Library/Keychains/System.keychain '%s'", rootCAPath)
		logx.Info("🔎 Check: security find-certificate -c 'tlsctl'")
		logx.Info("🧹 Uninstall: sudo security delete-certificate -c 'tlsctl'")
	case "linux":
		var targetPath string
		if fp.PathExists("/usr/local/share/ca-certificates/") {
			targetPath = "/usr/local/share/ca-certificates/tlsctl.crt"
		} else if fp.PathExists("/etc/pki/ca-trust/source/anchors/") {
			targetPath = "/etc/pki/ca-trust/source/anchors/tlsctl.pem"
		} else {
			logx.Info("⚠️ Cannot determine system trust directory. Please install manually: '%s'", rootCAPath)
			return
		}
		logx.Info("🔐 [Linux] Trust commands for the root certificate:")
		logx.Info("🛠 Install: sudo cp '%s' '%s' && sudo update-ca-certificates", rootCAPath, targetPath)
		logx.Info("🔎 Check: ls '%s'", targetPath)
		logx.Info("🧹 Uninstall: sudo rm '%s' && sudo update-ca-certificates", targetPath)
	case "windows":
		logx.Info("🔐 [Windows] Trust commands for the root certificate (run in Administrator PowerShell):")
		logx.Info("🛠 Install: certutil -addstore -f Root \"%s\"", rootCAPath)
		logx.Info("🔎 Check: certutil -store Root | findstr tlsctl")
		logx.Info("🧹 Uninstall: certutil -delstore Root tlsctl")
	default:
		logx.Info("⚠️ Trust installation is not supported on this OS. Please install manually: '%s'", rootCAPath)
	}
}
