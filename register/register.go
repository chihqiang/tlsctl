package register

import (
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type Register struct {
	TermsOfServiceAgreed bool
	Kid                  string
	HmacEncoded          string
}

func (r *Register) Register(client *lego.Client) (*registration.Resource, error) {
	if r.Kid != "" && r.HmacEncoded != "" {
		return client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: r.TermsOfServiceAgreed,
			Kid:                  r.Kid,
			HmacEncoded:          r.HmacEncoded,
		})
	}
	return client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: r.TermsOfServiceAgreed})
}
