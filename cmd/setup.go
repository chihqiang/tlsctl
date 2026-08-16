package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"path"

	"github.com/chihqiang/logx"
	"github.com/chihqiang/tlsctl/account"
	"github.com/chihqiang/tlsctl/challenge"
	"github.com/chihqiang/tlsctl/register"
	"github.com/chihqiang/tlsctl/resource"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/urfave/cli/v3"
)

func getEmail(ctx *cli.Command) (string, error) {
	email := ctx.String(flgEmail)
	if email == "" {
		return "", fmt.Errorf("you have to pass an account (email address) to the program using --%s or -m", flgEmail)
	}
	return email, nil
}

func getDomain(ctx *cli.Command) (string, error) {
	domain := ctx.String(flgDomain)
	if domain == "" {
		return "", fmt.Errorf("you have to pass a domain to the program using --%s or -d", flgDomain)
	}
	return domain, nil
}

func getDeployJson(ctx *cli.Command) string {
	return path.Join(ctx.String(flgPath), "deploy.json")
}

func setupAccountCache(ctx *cli.Command) (*account.Cache, error) {
	email, err := getEmail(ctx)
	if err != nil {
		return nil, err
	}
	cache, err := account.NewCache(ctx.String(flgPath), email, ctx.String(flgServer))
	if err != nil {
		return nil, fmt.Errorf("creating accounts cache: %w", err)
	}
	return cache, nil
}

func setupResourceCache(ctx *cli.Command) (*resource.Cache, error) {
	cCache, err := resource.NewCache(ctx.String(flgPath), "RC2")
	if err != nil {
		return nil, fmt.Errorf("creating certificates cache: %w", err)
	}
	return cCache, nil
}

func setupClient(cmd *cli.Command, ac *account.Cache) (*lego.Client, error) {
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
	keyType, err := account.GetKeyType(cmd.String(flgKeyType))
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
		loadAccount.Registration, err = register.GetRegister(cmd.String(flgKID), cmd.String(flgHMAC)).Register(client)
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

func buildLegoSSL(cmd *cli.Command, domain string) (*certificate.Resource, error) {
	rCache, err := setupResourceCache(cmd)
	if err != nil {
		return nil, err
	}
	cCache, err := setupAccountCache(cmd)
	if err != nil {
		return nil, err
	}
	client, err := setupClient(cmd, cCache)
	if err != nil {
		return nil, err
	}
	if err := challenge.SetConfigChallenge(client, challenge.Config{
		DNS:               cmd.String(flgDNS),
		Webroot:           cmd.String(flgHTTPWebroot),
		HTTPMemcachedHost: cmd.StringSlice(flgHTTPMemcachedHost),
		S3Bucket:          cmd.String(flgHTTPS3Bucket),
		HTTPPort:          cmd.String(flgHTTPPort),
		TLSPort:           cmd.String(flgTLSPort),
		TLS:               cmd.Bool(flgTLS),
		Delay:             cmd.Duration(flgTLSDelay),
	}); err != nil {
		return nil, fmt.Errorf("SetConfigChallenge: %w", err)
	}
	res, err := client.Certificate.Obtain(certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("Obtain: %w", err)
	}
	if err := rCache.SaveResource(res); err != nil {
		logx.Warn("SaveResource err: %v", err)
	}
	sanitizedDomain, err := resource.SanitizedDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("sanitize domain: %w", err)
	}
	logx.Debug("Certificate for %s has been saved successfully at %s",
		domain,
		rCache.GetSanitizedDomainSavePath(sanitizedDomain),
	)
	return res, nil
}
