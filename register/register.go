package register

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type Register struct {
}

func (r *Register) Register(lego *lego.Client) (*registration.Resource, error) {
	return lego.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
}
