package deploy

import (
	"context"
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/chihqiang/tlsctl/pkg/fp"
	"github.com/chihqiang/tlsctl/resource"
	"github.com/go-acme/lego/v4/certificate"
)

const (
	deployLocalPath = "/etc/nginx/ssl/"
)

func RunWithJSONFile(filePath, name string, certificate *certificate.Resource) error {
	sanitizedDomain, _ := resource.SanitizedDomain(certificate.Domain)
	switch name {
	case "local":
		if os.Getenv("LOCAL_CERT_PATH") == "" && os.Getenv("LOCAL_KEY_PATH") == "" {
			keyPath := path.Join(deployLocalPath, fmt.Sprintf("%s%s", sanitizedDomain, resource.PemExt))
			_ = os.Setenv("LOCAL_CERT_PATH", keyPath)
			cerPath := path.Join(deployLocalPath, fmt.Sprintf("%s%s", sanitizedDomain, resource.KeyExt))
			_ = os.Setenv("LOCAL_KEY_PATH", cerPath)
		}
	}
	envDeploy, err := Get(name)
	if err != nil {
		return err
	}
	if err := envDeploy.Deploy(context.Background(), certificate); err != nil {
		return err
	}
	return JSONFileSet(filePath, DomainDeploys{
		Domain:  certificate.Domain,
		Deploys: []string{name},
	})
}

type DomainDeploys struct {
	Domain  string   `json:"domain"`
	Deploys []string `json:"deploy"`
}

var deployJSONFile = &fp.JSONFile[DomainDeploys]{
	IsEqual: func(a, b DomainDeploys) bool {
		return a.Domain == b.Domain
	},
	Merge: func(existing *DomainDeploys, newItem DomainDeploys) {
		m := make(map[string]bool)
		for _, v := range existing.Deploys {
			m[v] = true
		}
		for _, v := range newItem.Deploys {
			m[v] = true
		}
		res := make([]string, 0, len(m))
		for k := range m {
			res = append(res, k)
		}
		existing.Deploys = res
	},
}

func JSONFileSet(file string, item DomainDeploys) error {
	return deployJSONFile.Save(file, item)
}
func JSONFileRemove(file string, domains []string) error {
	return deployJSONFile.Remove(file, func(t DomainDeploys) bool {
		return slices.Contains(domains, t.Domain)
	})
}
func JSONFileLoad(file string) ([]DomainDeploys, error) {
	return deployJSONFile.Load(file)
}
