package cmd

import (
	"os"
	"path"
	"runtime"
	"time"

	"github.com/go-acme/lego/v4/lego"
	"github.com/urfave/cli/v3"
)

const (
	flgDomain = "domain"

	flgServer = "server"
	flgEmail  = "email"
	flgPath   = "path"
	flgKID    = "kid"
	flgHMAC   = "hmac"
	flgDNS    = "dns"

	flgHTTPWebroot       = "http.webroot"
	flgHTTPMemcachedHost = "http.memcached-host"
	flgHTTPS3Bucket      = "http.s3-bucket"
	flgHTTPPort          = "http.port"

	flgTLS      = "tls"
	flgTLSPort  = "tls.port"
	flgTLSDelay = "tls.delay"

	flgDeploy = "deploy"

	flgInterval = "interval"
)

const (
	envEmail         = "TLSCTL_EMAIL"
	envServer        = "TLSCTL_SERVER"
	envPath          = "TLSCTL_PATH"
	envEABHMAC       = "TLSCTL_EAB_HMAC"
	envEABKID        = "TLSCTL_EAB_KID"
	envDay           = "TLSCTL_DAY"
	envRenewInterval = "TLSCTL_RENEW_INTERVAL"
)

func CreateFlags() []cli.Flag {
	homePath, _ := os.UserHomeDir()
	return []cli.Flag{
		&cli.StringFlag{
			Name:    flgServer,
			Aliases: []string{"s"},
			Sources: cli.EnvVars(envServer),
			Usage:   "CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client.",
			Value:   lego.LEDirectoryProduction,
		},
		&cli.StringFlag{
			Name:    flgEmail,
			Aliases: []string{"m"},
			Sources: cli.EnvVars(envEmail),
			Usage:   "Email used for registration and recovery contact.",
			Value:   "tlsctl@" + runtime.GOOS + ".com",
		},
		&cli.StringFlag{
			Name:    flgPath,
			Sources: cli.EnvVars(envPath),
			Usage:   "Directory to use for storing the data.",
			Value:   path.Join(homePath, ".tlsctl"),
		},
		&cli.StringFlag{
			Name:    flgKID,
			Sources: cli.EnvVars(envEABKID),
			Usage:   "Key identifier from External CA. Used for External Account Binding.",
		},
		&cli.StringFlag{
			Name:    flgHMAC,
			Sources: cli.EnvVars(envEABHMAC),
			Usage:   "MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.",
		},
		&cli.StringFlag{
			Name:  flgDNS,
			Usage: `Solve a DNS-01 challenge using the specified provider. Can be mixed with other types of challenges.`,
		},
		&cli.StringFlag{
			Name: flgHTTPWebroot,
			Usage: "Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file." +
				" This disables the built-in server and expects the given directory to be publicly served with access to .well-known/acme-challenge",
			Value: "/var/www/html",
		},
		&cli.StringSliceFlag{
			Name:  flgHTTPMemcachedHost,
			Usage: "Set the memcached host(s) to use for HTTP-01 based challenges. Challenges will be written to all specified hosts.",
		},
		&cli.StringFlag{
			Name:  flgHTTPS3Bucket,
			Usage: "Set the S3 bucket name to use for HTTP-01 based challenges. Challenges will be written to the S3 bucket.",
		},
		&cli.StringFlag{
			Name:  flgHTTPPort,
			Usage: "Set the port and interface to use for HTTP-01 based challenges to listen on. Supported: interface:port or :port.",
			Value: ":80",
		},
		&cli.BoolFlag{
			Name:  flgTLS,
			Usage: "Use the TLS-ALPN-01 challenge to solve challenges. Can be mixed with other types of challenges.",
		},
		&cli.StringFlag{
			Name:  flgTLSPort,
			Usage: "Set the port and interface to use for TLS-ALPN-01 based challenges to listen on. Supported: interface:port or :port.",
			Value: ":443",
		},
		&cli.DurationFlag{
			Name:  flgTLSDelay,
			Usage: "Delay between the start of the TLS listener (use for TLSALPN-01 based challenges) and the validation of the challenge.",
			Value: 0,
		},
		&cli.StringFlag{
			Name:    flgDomain,
			Aliases: []string{"d"},
			Usage:   "Add a domain to the process. Can be specified multiple times.",
		},
		&cli.StringFlag{
			Name:  flgDeploy,
			Usage: "Set your publishing method. When it is local, it is deployed to the /etc/nginx/ssl/ directory by default.",
			Value: "local",
		},
		&cli.DurationFlag{
			Name:    flgInterval,
			Usage:   "Check interval (e.g. 12h, 30m, 1h30m)",
			Value:   24 * time.Hour,
			Sources: cli.EnvVars(envRenewInterval),
		},
	}
}
