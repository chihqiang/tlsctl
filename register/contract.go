package register

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type IRegister interface {
	Register(lego *lego.Client) (*registration.Resource, error)
}
