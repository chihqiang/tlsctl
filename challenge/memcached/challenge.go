package memcached

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/http/memcached"
)

type Challenge struct {
	Hosts []string
}

func (w *Challenge) Set(client *lego.Client) error {
	ps, err := memcached.NewMemcachedProvider(w.Hosts)
	if err != nil {
		return err
	}
	return client.Challenge.SetHTTP01Provider(ps)
}
