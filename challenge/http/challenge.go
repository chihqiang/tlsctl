package http

import (
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
)

type Challenge struct {
	HeaderName string
}

func (w *Challenge) Set(client *lego.Client) error {
	srv := http01.NewProviderServer("", "")
	if header := w.HeaderName; header != "" {
		srv.SetProxyHeader(header)
	}
	return client.Challenge.SetHTTP01Provider(srv)
}
