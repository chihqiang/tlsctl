package deploy

import (
	"context"
	"slices"

	"github.com/chihqiang/tlsctl/pkg/fp"
	"github.com/go-acme/lego/v4/certificate"
)

func RunWithJSONFile(filePath, name string, certificate *certificate.Resource) error {
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
		// 保持原有顺序，仅去重并追加新项，避免 map 遍历导致顺序随机。
		seen := make(map[string]bool, len(existing.Deploys))
		res := make([]string, 0, len(existing.Deploys)+len(newItem.Deploys))
		for _, v := range existing.Deploys {
			if !seen[v] {
				res = append(res, v)
				seen[v] = true
			}
		}
		for _, v := range newItem.Deploys {
			if !seen[v] {
				res = append(res, v)
				seen[v] = true
			}
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
