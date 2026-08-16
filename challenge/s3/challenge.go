package s3

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/http/s3"
)

type Challenge struct {
	Bucket string
}

func (w *Challenge) Set(client *lego.Client) error {
	ps, err := s3.NewHTTPProvider(w.Bucket)
	if err != nil {
		return err
	}
	return client.Challenge.SetHTTP01Provider(ps)
}
