package challenge

import (
	"errors"
	"fmt"
	"time"

	"github.com/chihqiang/tlsctl/challenge/dns"
	"github.com/chihqiang/tlsctl/challenge/http"
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
	HTTPProxyHeader   string
	TLSPort           string
	TLS               bool
	Delay             time.Duration
}

func SetConfigChallenge(client *lego.Client, cfg Config) error {
	var count int

	// DNS-01
	if cfg.DNS != "" {
		if err := (&dns.Challenge{DNS: cfg.DNS}).Set(client); err != nil {
			return fmt.Errorf("DNS-01 challenge: %w", err)
		}
		count++
	}

	// HTTP-01：webroot / memcached / s3 / httpport 中至多选一个
	var httpChallenges []IChallenge
	if cfg.Webroot != "" {
		httpChallenges = append(httpChallenges, &webroot.Challenge{WebRoot: cfg.Webroot})
	}
	if len(cfg.HTTPMemcachedHost) > 0 {
		httpChallenges = append(httpChallenges, &memcached.Challenge{Hosts: cfg.HTTPMemcachedHost})
	}
	if cfg.S3Bucket != "" {
		httpChallenges = append(httpChallenges, &s3.Challenge{Bucket: cfg.S3Bucket})
	}
	if cfg.HTTPPort != "" {
		httpChallenges = append(httpChallenges, &httpport.Challenge{HostPort: cfg.HTTPPort})
	}
	if cfg.HTTPProxyHeader != "" {
		httpChallenges = append(httpChallenges, &http.Challenge{HeaderName: cfg.HTTPProxyHeader})
	}
	if len(httpChallenges) > 1 {
		return errors.New("conflicting HTTP-01 challenge providers specified")
	}
	if len(httpChallenges) == 1 {
		if err := httpChallenges[0].Set(client); err != nil {
			return fmt.Errorf("HTTP-01 challenge: %w", err)
		}
		count++
	}

	// TLS-ALPN-01
	if cfg.TLSPort != "" || cfg.TLS {
		if err := (&tls.Challenge{HostPort: cfg.TLSPort, TLS: cfg.TLS, Delay: cfg.Delay}).Set(client); err != nil {
			return fmt.Errorf("TLS-ALPN-01 challenge: %w", err)
		}
		count++
	}

	if count == 0 {
		return errors.New("challenge is required")
	}
	return nil
}
