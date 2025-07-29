package challenge

import "github.com/go-acme/lego/v4/lego"

type IChallenge interface {
	Set(client *lego.Client) error
}
