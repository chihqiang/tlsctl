package dns

import (
	"time"

	"github.com/chihqiang/tlsctl/dns"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
)

type Challenge struct {
	DNS                         string
	Servers                     []string
	PropagationWait             int
	Timeout                     int
	AuthoritativeNssPropagation bool
	RecursiveNssPropagation     bool
}

func (w *Challenge) Set(client *lego.Client) error {
	provider, err := dns.Get(w.DNS)
	if err != nil {
		return err
	}
	var opts = make([]dns01.ChallengeOption, 0)
	if len(w.Servers) > 0 {
		opts = append(opts, dns01.AddRecursiveNameservers(w.Servers))
	}
	if w.Timeout > 0 {
		opts = append(opts, dns01.AddDNSTimeout(time.Duration(w.Timeout)*time.Second))
	}
	if w.PropagationWait > 0 {
		opts = append(opts, dns01.PropagationWait(time.Duration(w.PropagationWait), true))
	}
	if w.AuthoritativeNssPropagation {
		opts = append(opts, dns01.DisableAuthoritativeNssPropagationRequirement())
	}
	if w.RecursiveNssPropagation {
		opts = append(opts, dns01.RecursiveNSsPropagationRequirement())
	}
	return client.Challenge.SetDNS01Provider(provider, opts...)

}
