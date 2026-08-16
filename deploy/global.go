package deploy

import (
	"fmt"

	"github.com/chihqiang/tlsctl/deploy/aliyun/cdn"
	"github.com/chihqiang/tlsctl/deploy/aliyun/dcdn"
	"github.com/chihqiang/tlsctl/deploy/aliyun/fc"
	"github.com/chihqiang/tlsctl/deploy/aliyun/live"
	"github.com/chihqiang/tlsctl/deploy/aliyun/oss"
	"github.com/chihqiang/tlsctl/deploy/aliyun/vod"
	"github.com/chihqiang/tlsctl/deploy/local"
	"github.com/chihqiang/tlsctl/deploy/ssh"
	tcdn "github.com/chihqiang/tlsctl/deploy/tencentcloud/cdn"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/clb"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/cos"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/ecdn"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/eo"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/scf"
	tssl "github.com/chihqiang/tlsctl/deploy/tencentcloud/ssl"
	tvod "github.com/chihqiang/tlsctl/deploy/tencentcloud/vod"
	"github.com/chihqiang/tlsctl/deploy/tencentcloud/waf"
)

var deploys = map[string]func() IDeploy{
	// 本地与 SSH
	"local": func() IDeploy { return &local.Deploy{} },
	"ssh":   func() IDeploy { return &ssh.Deploy{} },

	// 腾讯云相关部署
	"tcdn": func() IDeploy { return &tcdn.Deploy{} },
	"ecdn": func() IDeploy { return &ecdn.Deploy{} },
	"tssl": func() IDeploy { return &tssl.Deploy{} },
	"cos":  func() IDeploy { return &cos.Deploy{} },
	"scf":  func() IDeploy { return &scf.Deploy{} },
	"tvod": func() IDeploy { return &tvod.Deploy{} },
	"clb":  func() IDeploy { return &clb.Deploy{} },
	"waf":  func() IDeploy { return &waf.Deploy{} },
	"eo":   func() IDeploy { return &eo.Deploy{} },

	// 阿里云相关部署
	"cdn":  func() IDeploy { return &cdn.Deploy{} },
	"dcdn": func() IDeploy { return &dcdn.Deploy{} },
	"live": func() IDeploy { return &live.Deploy{} },
	"oss":  func() IDeploy { return &oss.Deploy{} },
	"vod":  func() IDeploy { return &vod.Deploy{} },
	"fc":   func() IDeploy { return &fc.Deploy{} },
}

func Register(name string, deploy func() IDeploy) {
	deploys[name] = deploy
}

func Get(name string) (IDeploy, error) {
	if deployFun, ok := deploys[name]; ok {
		deploy := deployFun()
		if err := deploy.WithEnvConfig(); err != nil {
			return nil, err
		}
		return deploy, nil
	}
	return nil, fmt.Errorf("deploy `%s` not found", name)
}

func All() map[string]IDeploy {
	ds := deploys
	var m = make(map[string]IDeploy)
	for s, f := range ds {
		m[s] = f()
	}
	return m
}
