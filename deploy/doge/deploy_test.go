package doge

import (
	"context"
	"strings"
	"testing"
)

// TestDeployRequiresKeys 校验缺少密钥时在发起网络请求前即报错。
func TestDeployRequiresKeys(t *testing.T) {
	d := &Deploy{Config: &Config{}}
	err := d.Deploy(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "required") {
		t.Fatalf("Deploy() = %v, want required error", err)
	}
}
