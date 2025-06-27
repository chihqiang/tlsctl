package tls

import (
	"net"
	"time"

	"github.com/go-acme/lego/v4/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/lego"
)

type Challenge struct {
	HostPort string
	TLS      bool
	Delay    int
}

func (w *Challenge) Set(client *lego.Client) error {
	if w.TLS {
		provider := tlsalpn01.NewProviderServer("", "")
		return client.Challenge.SetTLSALPN01Provider(provider, tlsalpn01.SetDelay(time.Duration(w.Delay)))
	}
	host, port, err := net.SplitHostPort(w.HostPort)
	if err != nil {
		return err
	}
	provider := tlsalpn01.NewProviderServer(host, port)
	return client.Challenge.SetTLSALPN01Provider(provider)
}
