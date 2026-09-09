package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/chihqiang/logx"
	"github.com/chihqiang/tlsctl/account"
	"github.com/chihqiang/tlsctl/challenge"
	"github.com/chihqiang/tlsctl/register"
	"github.com/chihqiang/tlsctl/resource"

	"github.com/chihqiang/cli"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
)

func getEmail(in *cli.Input) (string, error) {
	email := in.String(flgEmail)
	if email == "" {
		email = generatedEmail()
	}
	return email, nil
}

// generatedEmail 生成一个确定性的默认邮箱：tlsctl-<hostname>@<os>.com。
// 同一台机器每次生成一致（ACME 账号可稳定复用），不同机器各不相同（避免撞车）。
func generatedEmail() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		host = "localhost"
	}
	var b strings.Builder
	for _, r := range strings.ToLower(host) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.' {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	return "tlsctl-" + b.String() + "@" + runtime.GOOS + ".com"
}

func getDomain(in *cli.Input) ([]string, error) {
	domains := in.StringSlice(flgDomain)
	if len(domains) == 0 {
		return nil, fmt.Errorf("you have to pass a domain to the program using --%s or -d", flgDomain)
	}
	return domains, nil
}

func getDeployJson(in *cli.Input) string {
	return path.Join(in.String(flgPath), "deploy.json")
}

func setupAccountCache(in *cli.Input) (*account.Cache, error) {
	email, err := getEmail(in)
	if err != nil {
		return nil, err
	}
	cache, err := account.NewCache(in.String(flgPath), email, in.String(flgServer))
	if err != nil {
		return nil, fmt.Errorf("creating accounts cache: %w", err)
	}
	return cache, nil
}

func setupResourceCache(in *cli.Input) (*resource.Cache, error) {
	cCache, err := resource.NewCache(in.String(flgPath), "RC2")
	if err != nil {
		return nil, fmt.Errorf("creating certificates cache: %w", err)
	}
	return cCache, nil
}

func setupClient(in *cli.Input, ac *account.Cache) (*lego.Client, error) {
	loadAccount, err := ac.LoadAccount()
	if err != nil {
		logx.Warn("account.cache.LoadAccount: %v", err)
	}
	if loadAccount == nil {
		loadAccount = &account.Account{Email: ac.GetEmail()}
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("ecdsa.GenerateKey: %w", err)
		}
		loadAccount.Key = privateKey
	}
	keyType, err := account.GetKeyType(in.String(flgKeyType))
	if err != nil {
		return nil, fmt.Errorf("invalid key type: %w", err)
	}
	config := lego.NewConfig(loadAccount)
	config.CADirURL = ac.GetServer()
	config.Certificate.KeyType = keyType
	config.UserAgent = "github.com/chihqiang/tlsctl@main"
	client, err := lego.NewClient(config)
	if err != nil {
		ac.Remove()
		return nil, fmt.Errorf("lego.NewClient: %w", err)
	}
	if loadAccount.Registration == nil {
		loadAccount.Registration, err = register.GetRegister(in.String(flgKID), in.String(flgHMAC)).Register(client)
		if err != nil {
			ac.Remove()
			return nil, fmt.Errorf("Register: %w", err)
		}
	}
	if err := ac.Save(loadAccount); err != nil {
		logx.Warn("Account.Save %v", err)
	}
	return client, nil
}

func buildLegoSSL(in *cli.Input, domains []string) (*certificate.Resource, error) {
	rCache, err := setupResourceCache(in)
	if err != nil {
		return nil, err
	}
	cCache, err := setupAccountCache(in)
	if err != nil {
		return nil, err
	}
	client, err := setupClient(in, cCache)
	if err != nil {
		return nil, err
	}
	if err := challenge.SetConfigChallenge(client, challenge.Config{
		DNS:               in.String(flgDNS),
		Webroot:           in.String(flgHTTPWebroot),
		HTTPMemcachedHost: in.StringSlice(flgHTTPMemcachedHost),
		S3Bucket:          in.String(flgHTTPS3Bucket),
		HTTPPort:          in.String(flgHTTPPort),
		HTTPProxyHeader:   in.String(flgHTTPProxyHeader),
		TLSPort:           in.String(flgTLSPort),
		TLS:               in.Bool(flgTLS),
		Delay:             in.Duration(flgTLSDelay),
	}); err != nil {
		return nil, fmt.Errorf("SetConfigChallenge: %w", err)
	}
	res, err := client.Certificate.Obtain(certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("Obtain: %w", err)
	}
	if err := rCache.SaveResource(res); err != nil {
		logx.Warn("SaveResource err: %v", err)
	}
	sanitizedDomain, err := resource.SanitizedDomain(res.Domain)
	if err != nil {
		return nil, fmt.Errorf("sanitize domain: %w", err)
	}
	logx.Debug("Certificate for %s has been saved successfully at %s",
		res.Domain,
		rCache.GetSanitizedDomainSavePath(sanitizedDomain),
	)
	return res, nil
}
