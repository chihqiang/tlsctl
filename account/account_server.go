package account

import "github.com/go-acme/lego/v4/lego"

const (
	LEDirectoryProduction string = lego.LEDirectoryProduction
	LEDirectoryStaging    string = lego.LEDirectoryStaging
	ZeroServer            string = "https://acme.zerossl.com/v2/DV90"
	GTS                   string = "https://dv.acme-v02.api.pki.goog/directory"
)

var (
	AllowServers = []string{
		LEDirectoryProduction,
		LEDirectoryStaging,
		ZeroServer,
		GTS,
	}
)
