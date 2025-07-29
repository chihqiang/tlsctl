package webroot

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/http/webroot"
)

type Challenge struct {
	WebRoot string
}

func (w *Challenge) Set(client *lego.Client) error {
	provider, err := webroot.NewHTTPProvider(w.WebRoot)
	if err != nil {
		return err
	}
	return client.Challenge.SetHTTP01Provider(provider)
}
