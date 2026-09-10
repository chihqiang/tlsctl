package deploy

import (
	"fmt"

	"github.com/chihqiang/tlsctl/deploy/aliyun/cdn"
	"github.com/chihqiang/tlsctl/deploy/aliyun/dcdn"
	"github.com/chihqiang/tlsctl/deploy/aliyun/fc"
	"github.com/chihqiang/tlsctl/deploy/aliyun/live"
	"github.com/chihqiang/tlsctl/deploy/aliyun/oss"
	"github.com/chihqiang/tlsctl/deploy/aliyun/vod"
	"github.com/chihqiang/tlsctl/deploy/baiduyun"
	"github.com/chihqiang/tlsctl/deploy/btpanel"
	"github.com/chihqiang/tlsctl/deploy/doge"
	"github.com/chihqiang/tlsctl/deploy/huaweicloud"
	"github.com/chihqiang/tlsctl/deploy/lecdn"
	"github.com/chihqiang/tlsctl/deploy/local"
	"github.com/chihqiang/tlsctl/deploy/onepanel"
	"github.com/chihqiang/tlsctl/deploy/qiniu"
	"github.com/chihqiang/tlsctl/deploy/rainyun"
	"github.com/chihqiang/tlsctl/deploy/safeline"
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
	"github.com/chihqiang/tlsctl/deploy/volcengine"
	"github.com/chihqiang/tlsctl/deploy/webhook"
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

	// 宝塔面板
	"btpanel":            func() IDeploy { return &btpanel.Deploy{} },
	"btpanel-site":       func() IDeploy { return &btpanel.DeploySite{} },
	"btpanel-dockersite": func() IDeploy { return &btpanel.DeployDockerSite{} },
	"btpanel-singlesite": func() IDeploy { return &btpanel.DeploySingleSite{} },

	// 1Panel 面板
	"1panel":      func() IDeploy { return &onepanel.Deploy{} },
	"1panel-site": func() IDeploy { return &onepanel.DeploySite{} },

	// 雷池 WAF
	"safeline-panel":  func() IDeploy { return &safeline.DeployPanel{} },
	"safeline-site":   func() IDeploy { return &safeline.DeploySite{} },
	"safeline-portal": func() IDeploy { return &safeline.DeployPortal{} },

	// 七牛云
	"qiniu-cdn": func() IDeploy { return &qiniu.DeployCdn{} },
	"qiniu-oss": func() IDeploy { return &qiniu.DeployOss{} },

	// 百度云
	"baidu-cdn": func() IDeploy { return &baiduyun.Deploy{} },

	// 华为云
	"huaweicloud-cdn": func() IDeploy { return &huaweicloud.Deploy{} },

	// 火山引擎
	"volcengine-cdn":  func() IDeploy { return &volcengine.DeployCdn{} },
	"volcengine-dcdn": func() IDeploy { return &volcengine.DeployDcdn{} },

	// 多吉云
	"doge-cdn": func() IDeploy { return &doge.Deploy{} },

	// LeCDN
	"lecdn": func() IDeploy { return &lecdn.Deploy{} },

	// 雨云
	"rainyun": func() IDeploy { return &rainyun.Deploy{} },

	// Webhook
	"webhook": func() IDeploy { return &webhook.Deploy{} },
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
	var m = make(map[string]IDeploy, len(deploys))
	for s, f := range deploys {
		m[s] = f()
	}
	return m
}
