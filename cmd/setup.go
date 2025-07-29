package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"path"

	"github.com/chihqiang/tlsctl/account"
	"github.com/chihqiang/tlsctl/challenge"
	"github.com/chihqiang/tlsctl/pkg/log"
	"github.com/chihqiang/tlsctl/register"
	"github.com/chihqiang/tlsctl/resource"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/urfave/cli/v3"
)

func getEmail(ctx *cli.Command) string {
	email := ctx.String(flgEmail)
	if email == "" {
		log.Error("You have to pass an account (email address) to the program using --%s or -m", flgEmail)
	}
	return email
}

func getDomain(ctx *cli.Command) string {
	domain := ctx.String(flgDomain)
	if domain == "" {
		log.Error("You have to pass an Domain to the program using --%s or -d", flgDomain)
	}
	return domain
}

func getDeployJson(ctx *cli.Command) string {
	return path.Join(ctx.String(flgPath), "deploy.json")
}

func setupAccountCache(ctx *cli.Command) *account.Cache {
	cache, err := account.NewCache(ctx.String(flgPath), getEmail(ctx), ctx.String(flgServer))
	if err != nil {
		log.Error("creating accounts cache: %v", err)
	}
	return cache
}

func setupResourceCache(ctx *cli.Command) *resource.Cache {
	cCache, err := resource.NewCache(ctx.String(flgPath), "RC2")
	if err != nil {
		log.Error("creating certificates cache: %v", err)
	}
	return cCache
}

func setupClient(cmd *cli.Command, ac *account.Cache) *lego.Client {
	loadAccount, err := ac.LoadAccount()
	if err != nil {
		log.Warn("account.cache.LoadAccount: %v", err)
	}
	if loadAccount == nil {
		loadAccount = &account.Account{Email: ac.GetEmail()}
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Error("ecdsa.GenerateKey err:", err)
		}
		loadAccount.Key = privateKey
	}
	config := lego.NewConfig(loadAccount)
	config.CADirURL = ac.GetServer()
	config.Certificate.KeyType = certcrypto.RSA2048
	config.UserAgent = "github.com/chihqiang/tlsctl@main"
	client, err := lego.NewClient(config)
	if err != nil {
		ac.Remove()
		log.Error("lego.NewClient err:%v", err)
	}
	if loadAccount.Registration == nil {
		loadAccount.Registration, err = register.GetRegister(cmd.String(flgKID), cmd.String(flgHMAC)).Register(client)
		if err != nil {
			ac.Remove()
			log.Error("Register err:%v", err)
		}
	}
	err = ac.Save(loadAccount)
	if err != nil {
		log.Warn("Account.Save %v", err)
	}
	return client
}

func buildLegoSSL(cmd *cli.Command, domain string) (*certificate.Resource, error) {
	rCache := setupResourceCache(cmd)
	cCache := setupAccountCache(cmd)
	client := setupClient(cmd, cCache)
	if err := challenge.SetConfigChallenge(client, challenge.Config{
		DNS:               cmd.String(flgDNS),
		Webroot:           cmd.String(flgHTTPWebroot),
		HTTPMemcachedHost: cmd.StringSlice(flgHTTPMemcachedHost),
		S3Bucket:          cmd.String(flgHTTPS3Bucket),
		HTTPPort:          cmd.String(flgHTTPPort),
		TLSPort:           cmd.String(flgTLSPort),
		TLS:               cmd.Bool(flgTLS),
		Delay:             cmd.Int(flgTLSDelay),
	}); err != nil {
		log.Error("SetConfigChallenge %v", err)
	}
	res, err := client.Certificate.Obtain(certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	})
	if err != nil {
		log.Error("Obtain err: %v", err)
	}
	if err := rCache.SaveResource(res); err != nil {
		log.Warn("SaveResource err: %v", err)
	}
	sanitizedDomain, _ := resource.SanitizedDomain(domain)
	log.Debug("Certificate for %s has been saved successfully at %s",
		domain,
		rCache.GetSanitizedDomainSavePath(sanitizedDomain),
	)
	return res, nil
}
