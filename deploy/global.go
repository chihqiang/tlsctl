package deploy

import (
	"wangzhiqiang/tlsctl/deploy/aliyun/cdn"
	"wangzhiqiang/tlsctl/deploy/aliyun/dcdn"
	"wangzhiqiang/tlsctl/deploy/aliyun/fc"
	"wangzhiqiang/tlsctl/deploy/aliyun/live"
	"wangzhiqiang/tlsctl/deploy/aliyun/oss"
	"wangzhiqiang/tlsctl/deploy/aliyun/vod"
	"wangzhiqiang/tlsctl/deploy/local"
	"wangzhiqiang/tlsctl/deploy/ssh"
	tcdn "wangzhiqiang/tlsctl/deploy/tencentcloud/cdn"
	"wangzhiqiang/tlsctl/deploy/tencentcloud/clb"
	"wangzhiqiang/tlsctl/deploy/tencentcloud/cos"
	"wangzhiqiang/tlsctl/deploy/tencentcloud/ecdn"
	"wangzhiqiang/tlsctl/deploy/tencentcloud/scf"
	tssl "wangzhiqiang/tlsctl/deploy/tencentcloud/ssl"
	tvod "wangzhiqiang/tlsctl/deploy/tencentcloud/vod"
	"fmt"
)

var deploys = map[string]IDeploy{
	"local": &local.Deploy{},
	"ssh":   &ssh.Deploy{},
	//腾讯相关的部署
	"tcdn": &tcdn.Deploy{},
	"ecdn": &ecdn.Deploy{},
	"tssl": &tssl.Deploy{},
	"cos":  &cos.Deploy{},
	"scf":  &scf.Deploy{},
	"tvod": &tvod.Deploy{},
	"clb":  &clb.Deploy{},
	//阿里云相关部署
	"cdn":  &cdn.Deploy{},
	"dcdn": &dcdn.Deploy{},
	"live": &live.Deploy{},
	"oss":  &oss.Deploy{},
	"vod":  &vod.Deploy{},
	"fc":   &fc.Deploy{},
}

func Register(name string, deploy IDeploy) {
	deploys[name] = deploy
}

func Get(name string) (IDeploy, error) {
	if deploy, ok := deploys[name]; ok {
		if err := deploy.WithEnvConfig(); err != nil {
			return nil, err
		}
		return deploy, nil
	}
	return nil, fmt.Errorf("deploy `%s` not found", name)
}

func All() map[string]IDeploy {
	return deploys
}
