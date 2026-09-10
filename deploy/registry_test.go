package deploy

import (
	"testing"

	"github.com/chihqiang/tlsctl/pkg/structs"
)

// TestRegistryCompleteness 确保每个注册的部署方式都能实例化，
// 且 help:deploy 的标签反射能展示其环境变量字段。
func TestRegistryCompleteness(t *testing.T) {
	instances := All()
	if len(instances) != len(deploys) {
		t.Fatalf("All() returned %d instances, want %d", len(instances), len(deploys))
	}
	tagMaps := structs.TagsMaps(instances)
	if len(tagMaps.Keys) != len(deploys) {
		t.Fatalf("TagsMaps found env tags for %d deploys, want %d", len(tagMaps.Keys), len(deploys))
	}
	for _, name := range tagMaps.Keys {
		if len(tagMaps.Maps[name]) == 0 {
			t.Errorf("deploy %q has no env tags, help:deploy will hide it", name)
		}
	}
}

// TestRegistryNoNilInstance 校验每个注册项的工厂都能返回非 nil 实例。
func TestRegistryNoNilInstance(t *testing.T) {
	for name, factory := range deploys {
		if factory() == nil {
			t.Errorf("factory for %q returned nil", name)
		}
	}
}
