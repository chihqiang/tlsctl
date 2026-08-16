package cmd

import (
	"context"
	"path"
	"runtime"

	"github.com/chihqiang/tlsctl/localhost"
	"github.com/chihqiang/tlsctl/pkg/fp"
	"github.com/chihqiang/tlsctl/pkg/log"
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
	cStorage := setupResourceCache(cmd)
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
		log.Error("error saving certificate: %v", err)
	}
	log.Debug("Certificate for %s has been saved successfully at %s",
		"localhost",
		cStorage.GetSanitizedDomainSavePath("localhost"),
	)
	localhostSSLInstallHelp(path.Join(hostSSL.Path, "rootCA.pem"))
	return nil
}

func localhostSSLInstallHelp(rootCAPath string) {
	switch runtime.GOOS {
	case "darwin":
		log.Info("🔐 [macOS] Trust commands for the root certificate:")
		log.Info("🛠 Install: sudo security add-trusted-cert -d -k /Library/Keychains/System.keychain '%s'", rootCAPath)
		log.Info("🔎 Check: security find-certificate -c 'tlsctl'")
		log.Info("🧹 Uninstall: sudo security delete-certificate -c 'tlsctl'")
	case "linux":
		var targetPath string
		if fp.PathExists("/usr/local/share/ca-certificates/") {
			targetPath = "/usr/local/share/ca-certificates/tlsctl.crt"
		} else if fp.PathExists("/etc/pki/ca-trust/source/anchors/") {
			targetPath = "/etc/pki/ca-trust/source/anchors/tlsctl.pem"
		} else {
			log.Info("⚠️ Cannot determine system trust directory. Please install manually: '%s'", rootCAPath)
			return
		}
		log.Info("🔐 [Linux] Trust commands for the root certificate:")
		log.Info("🛠 Install: sudo cp '%s' '%s' && sudo update-ca-certificates", rootCAPath, targetPath)
		log.Info("🔎 Check: ls '%s'", targetPath)
		log.Info("🧹 Uninstall: sudo rm '%s' && sudo update-ca-certificates", targetPath)
	case "windows":
		log.Info("🔐 [Windows] Trust commands for the root certificate (run in Administrator PowerShell):")
		log.Info("🛠 Install: certutil -addstore -f Root \"%s\"", rootCAPath)
		log.Info("🔎 Check: certutil -store Root | findstr tlsctl")
		log.Info("🧹 Uninstall: certutil -delstore Root tlsctl")
	default:
		log.Info("⚠️ Trust installation is not supported on this OS. Please install manually: '%s'", rootCAPath)
	}
}
