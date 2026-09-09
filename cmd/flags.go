package cmd

import (
	"os"
	"path"
	"time"

	"github.com/chihqiang/cli"
	"github.com/go-acme/lego/v4/lego"
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
	flgHTTPProxyHeader   = "http.proxy-header"

	flgTLS      = "tls"
	flgTLSPort  = "tls.port"
	flgTLSDelay = "tls.delay"

	flgDeploy = "deploy"

	flgKeyType = "key-type"

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

func defaultStorePath() string {
	homePath, _ := os.UserHomeDir()
	return path.Join(homePath, ".tlsctl")
}

func serverFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    flgServer,
		Aliases: []string{"s"},
		Sources: cli.EnvVars(envServer),
		Usage:   "CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client.",
		Value:   lego.LEDirectoryProduction,
	}
}

func emailFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    flgEmail,
		Aliases: []string{"m"},
		Sources: cli.EnvVars(envEmail),
		Usage:   "Email used for registration and recovery contact. If not provided, a system-generated email will be used.",
		Value:   generatedEmail(),
	}
}

func pathFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    flgPath,
		Sources: cli.EnvVars(envPath),
		Usage:   "Directory to use for storing the data.",
		Value:   defaultStorePath(),
	}
}

func kidFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    flgKID,
		Sources: cli.EnvVars(envEABKID),
		Usage:   "Key identifier from External CA. Used for External Account Binding.",
	}
}

func hmacFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    flgHMAC,
		Sources: cli.EnvVars(envEABHMAC),
		Usage:   "MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.",
	}
}

func dnsFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  flgDNS,
		Usage: `Solve a DNS-01 challenge using the specified provider. Can be mixed with other types of challenges.`,
	}
}

func keyTypeFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  flgKeyType,
		Usage: "Type of the private key to generate. Supported: RSA2048, RSA3072, RSA4096, RSA8192, EC256, EC384.",
		Value: "RSA2048",
	}
}

func httpWebrootFlag() cli.Flag {
	return &cli.StringFlag{
		Name: flgHTTPWebroot,
		Usage: "Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file." +
			" This disables the built-in server and expects the given directory to be publicly served with access to .well-known/acme-challenge",
		Value: "/var/www/html",
	}
}

func httpMemcachedHostFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:  flgHTTPMemcachedHost,
		Usage: "Set the memcached host(s) to use for HTTP-01 based challenges. Challenges will be written to all specified hosts.",
	}
}

func httpS3BucketFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  flgHTTPS3Bucket,
		Usage: "Set the S3 bucket name to use for HTTP-01 based challenges. Challenges will be written to the S3 bucket.",
	}
}

func httpPortFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  flgHTTPPort,
		Usage: "Set the port and interface to use for HTTP-01 based challenges to listen on. Supported: interface:port or :port.",
		Value: ":80",
	}
}

func httpProxyHeaderFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  flgHTTPProxyHeader,
		Usage: "Set the proxy header to use for HTTP-01 based challenges. Used when the HTTP-01 challenge is behind a reverse proxy.",
	}
}

func tlsFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  flgTLS,
		Usage: "Use the TLS-ALPN-01 challenge to solve challenges. Can be mixed with other types of challenges.",
	}
}

func tlsPortFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  flgTLSPort,
		Usage: "Set the port and interface to use for TLS-ALPN-01 based challenges to listen on. Supported: interface:port or :port.",
		Value: ":443",
	}
}

func tlsDelayFlag() cli.Flag {
	return &cli.DurationFlag{
		Name:  flgTLSDelay,
		Usage: "Delay between the start of the TLS listener (use for TLSALPN-01 based challenges) and the validation of the challenge.",
		Value: 0,
	}
}

func domainFlag() cli.Flag {
	return &cli.StringSliceFlag{
		Name:    flgDomain,
		Aliases: []string{"d"},
		Usage:   "Add a domain to the certificate. Can be specified multiple times.",
	}
}

func deployFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  flgDeploy,
		Usage: "Set your publishing method. When it is local, it is deployed to the /etc/nginx/ssl/ directory by default.",
		Value: "local",
	}
}

func intervalFlag() cli.Flag {
	return &cli.DurationFlag{
		Name:    flgInterval,
		Usage:   "Check interval (e.g. 12h, 30m, 1h30m)",
		Value:   24 * time.Hour,
		Sources: cli.EnvVars(envRenewInterval),
	}
}

func dayFlag() cli.Flag {
	return &cli.IntFlag{
		Name:    "day",
		Usage:   "When the expiration date is less than a few days, it will be regenerated",
		Value:   1,
		Sources: cli.EnvVars(envDay),
	}
}

func accountFlags() []cli.Flag {
	return []cli.Flag{
		serverFlag(),
		emailFlag(),
		pathFlag(),
		kidFlag(),
		hmacFlag(),
	}
}

func challengeFlags() []cli.Flag {
	return []cli.Flag{
		dnsFlag(),
		keyTypeFlag(),
		httpWebrootFlag(),
		httpMemcachedHostFlag(),
		httpS3BucketFlag(),
		httpPortFlag(),
		httpProxyHeaderFlag(),
		tlsFlag(),
		tlsPortFlag(),
		tlsDelayFlag(),
	}
}

// createFlags 是 create 命令使用的标志（账号 + 挑战 + 域名）。
func createFlags() []cli.Flag {
	return append(append(accountFlags(), challengeFlags()...), domainFlag())
}

// scheduledRunFlags 是 scheduled:run 命令使用的标志：签发所需 + 周期 + 续期天数阈值。
func scheduledRunFlags() []cli.Flag {
	return append(append(append(accountFlags(), challengeFlags()...), intervalFlag()), dayFlag())
}

// deployFlags 是 deploy 命令使用的标志：存储路径 + 域名 + 部署方式。
func deployFlags() []cli.Flag {
	return []cli.Flag{
		pathFlag(),
		domainFlag(),
		deployFlag(),
	}
}

// pathFlags 只包含存储路径标志（list/clear/scheduled:list 等命令使用）。
func pathFlags() []cli.Flag {
	return []cli.Flag{pathFlag()}
}
