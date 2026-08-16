package register

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type EABRegister struct {
	TermsOfServiceAgreed bool
	Kid                  string
	HmacEncoded          string
}

func (r *EABRegister) Register(lego *lego.Client) (*registration.Resource, error) {
	return lego.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
		TermsOfServiceAgreed: r.TermsOfServiceAgreed,
		Kid:                  r.Kid,
		HmacEncoded:          r.HmacEncoded,
	})
}
