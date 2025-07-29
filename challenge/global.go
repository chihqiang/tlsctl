package challenge

import (
	"errors"

	"github.com/chihqiang/tlsctl/challenge/dns"
	"github.com/chihqiang/tlsctl/challenge/httpport"
	"github.com/chihqiang/tlsctl/challenge/memcached"
	"github.com/chihqiang/tlsctl/challenge/s3"
	"github.com/chihqiang/tlsctl/challenge/tls"
	"github.com/chihqiang/tlsctl/challenge/webroot"

	"github.com/go-acme/lego/v4/lego"
)

type Config struct {
	DNS               string
	Webroot           string
	HTTPMemcachedHost []string
	S3Bucket          string
	HTTPPort          string
	TLSPort           string
	TLS               bool
	Delay             int
}

func SetConfigChallenge(client *lego.Client, cfg Config) error {
	var challenge IChallenge
	if cfg.DNS != "" {
		challenge = &dns.Challenge{DNS: cfg.DNS}
	} else if cfg.Webroot != "" {
		challenge = &webroot.Challenge{WebRoot: cfg.Webroot}
	} else if len(cfg.HTTPMemcachedHost) > 0 {
		challenge = &memcached.Challenge{Hosts: cfg.HTTPMemcachedHost}
	} else if cfg.S3Bucket != "" {
		challenge = &s3.Challenge{Bucket: cfg.S3Bucket}
	} else if cfg.HTTPPort != "" {
		challenge = &httpport.Challenge{HostPort: cfg.HTTPPort}
	} else if cfg.TLSPort != "" || cfg.TLS {
		challenge = &tls.Challenge{HostPort: cfg.TLSPort, TLS: cfg.TLS, Delay: cfg.Delay}
	}
	if challenge == nil {
		return errors.New("challenge is required")
	}
	return SetChallenge(client, challenge)
}

func SetChallenge(client *lego.Client, challenge IChallenge) error {
	return challenge.Set(client)
}
