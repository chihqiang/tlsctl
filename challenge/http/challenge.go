package http

import (
	"net"

	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
)

type Challenge struct {
	HostPort   string
	HeaderName string
}

func (w *Challenge) Set(client *lego.Client) error {
	host, port := "", ""
	if w.HostPort != "" {
		var err error
		host, port, err = net.SplitHostPort(w.HostPort)
		if err != nil {
			return err
		}
	}
	srv := http01.NewProviderServer(host, port)
	if header := w.HeaderName; header != "" {
		srv.SetProxyHeader(header)
	}
	return client.Challenge.SetHTTP01Provider(srv)
}
