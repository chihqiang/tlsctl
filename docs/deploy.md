# 🚀 部署方式详解

`tlsctl deploy --domain="example.com" --deploy="<方式>"` 会把证书部署到目标位置，并登记到定时任务（`~/.tlsctl/deploy.json`），供后续自动续签后重新部署。

运行 `tlsctl help:deploy` 可查看所有部署方式的完整字段。

## 本地部署 `local`

把证书写入本机目录，默认 `/etc/nginx/ssl/`。

```bash
cat > ~/.tlsctl/.env << EOF
LOCAL_CERT_PATH=/etc/nginx/ssl/example.com.pem
LOCAL_KEY_PATH=/etc/nginx/ssl/example.com.key
LOCAL_POST_COMMAND="nginx -s reload"
EOF

tlsctl deploy --domain="example.com" --deploy="local"
```

> 不配置 `LOCAL_CERT_PATH` / `LOCAL_KEY_PATH` 时，自动使用 `/etc/nginx/ssl/<domain>.pem` 与 `.key`。
> `LOCAL_PRE_COMMAND` / `LOCAL_POST_COMMAND` 分别在写入前后执行。

## SSH 部署 `ssh`

通过 SSH 把证书上传到远程服务器。

### 密码登录

```bash
cat > ~/.tlsctl/.env << EOF
SSH_HOST=192.168.1.100
SSH_PORT=22
SSH_USERNAME=root
SSH_PASSWORD=your_password
SSH_CERT_PATH=/etc/nginx/ssl/example.com.pem
SSH_KEY_PATH=/etc/nginx/ssl/example.com.key
SSH_POST_COMMAND="nginx -s reload"
EOF

tlsctl deploy --domain="example.com" --deploy="ssh"
```

### 私钥登录

```bash
cat > ~/.tlsctl/.env << EOF
SSH_HOST=192.168.1.100
SSH_USERNAME=deploy
SSH_KEY=-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----
SSH_KEY_PASSPHRASE=your_passphrase
SSH_CERT_PATH=/etc/nginx/ssl/example.com.pem
SSH_KEY_PATH=/etc/nginx/ssl/example.com.key
EOF

tlsctl deploy --domain="example.com" --deploy="ssh"
```

> 默认使用 SFTP 传输，设置 `SSH_USE_SCP=true` 可改用 SCP。

## 腾讯云部署

公共配置（写在 `.env`）：

```bash
TENCENTCLOUD_SECRET_ID=your_secret_id
TENCENTCLOUD_SECRET_KEY=your_secret_key
```

### CDN `tcdn`

```bash
TENCENTCLOUD_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="tcdn"
```

### 边缘 CDN `ecdn`

```bash
TENCENTCLOUD_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="ecdn"
```

### SSL 证书上传 `tssl`

把证书上传到腾讯云 SSL 证书中心，供其它云产品引用：

```bash
tlsctl deploy --domain="example.com" --deploy="tssl"
```

### COS 对象存储 `cos`

```bash
TENCENTCLOUD_REGION=ap-guangzhou
TENCENTCLOUD_BUCKET=my-bucket
TENCENTCLOUD_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="cos"
```

### SCF 云函数 `scf`

```bash
TENCENTCLOUD_REGION=ap-guangzhou
TENCENTCLOUD_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="scf"
```

### CLB 负载均衡 `clb`

```bash
TENCENTCLOUD_REGION=ap-guangzhou
TENCENTCLOUD_LOADBALANCER_ID=lb-xxxx
TENCENTCLOUD_LISTENER_ID=lbl-xxxx
tlsctl deploy --domain="example.com" --deploy="clb"
```

### WAF `waf`

```bash
TENCENTCLOUD_REGION=ap-guangzhou
TENCENTCLOUD_DOMAIN=example.com
TENCENTCLOUD_DOMAIN_ID=xxx
tlsctl deploy --domain="example.com" --deploy="waf"
```

> 多实例场景可配置 `TENCENTCLOUD_INSTANCE_ID` 指定所属实例。

### 边缘安全加速 `eo`

```bash
TENCENTCLOUD_ZONE_ID=zone-xxx
tlsctl deploy --domain="example.com" --deploy="eo"
```

### VOD 点播 `tvod`

```bash
TENCENTCLOUD_SUB_APP_ID=123456
TENCENTCLOUD_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="tvod"
```

## 阿里云部署

公共配置（写在 `.env`）：

```bash
ALIYUN_ACCESS_KEY_ID=your_access_key_id
ALIYUN_ACCESS_KEY_SECRET=your_access_key_secret
```

### CDN `cdn`

```bash
ALIYUN_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="cdn"
```

### DCDN `dcdn`

```bash
ALIYUN_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="dcdn"
```

### 直播 `live`

```bash
ALIYUN_REGION=cn-shanghai
ALIYUN_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="live"
```

### OSS `oss`

```bash
ALIYUN_REGION=oss-cn-hangzhou
ALIYUN_BUCKET=my-bucket
ALIYUN_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="oss"
```

### VOD 视频点播 `vod`

```bash
ALIYUN_REGION=cn-shanghai
ALIYUN_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="vod"
```

### FC 函数计算 `fc`

```bash
ALIYUN_REGION=cn-hangzhou
ALIYUN_RESOURCE_GROUP_ID=rg-xxx
ALIYUN_VERSION=3.0
ALIYUN_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="fc"
```

## 宝塔面板

### 面板 SSL `btpanel`

```bash
cat > ~/.tlsctl/.env << EOF
BTPANEL_URL=https://example.com:8888
BTPANEL_API_KEY=your_api_key
EOF

tlsctl deploy --domain="example.com" --deploy="btpanel"
```

### 网站证书 `btpanel-site`

上传证书到宝塔证书库并绑定到指定网站（多个网站用逗号分隔）：

```bash
BTPANEL_SITE_NAME=example.com,blog.example.com
tlsctl deploy --domain="example.com" --deploy="btpanel-site"
```

> 另有 `btpanel-dockersite`（Docker 面板网站）、`btpanel-singlesite`（旧版单个站点）。
> 面板使用自签名证书时配置 `BTPANEL_IGNORE_SSL=true`。

## 1Panel

### 面板 SSL `1panel`

```bash
ONEPANEL_URL=https://example.com:8090
ONEPANEL_API_KEY=your_api_key
tlsctl deploy --domain="example.com" --deploy="1panel"
```

### 网站证书 `1panel-site`

```bash
ONEPANEL_SITE_ID=1
tlsctl deploy --domain="example.com" --deploy="1panel-site"
```

> 1Panel v2 配置 `ONEPANEL_VERSION=v2`；自签名证书配置 `ONEPANEL_IGNORE_SSL=true`。

## 雷池 WAF

```bash
SAFELINE_URL=https://example.com:9443
SAFELINE_API_TOKEN=your_api_token
SAFELINE_SITE_NAME=example.com   # 仅 safeline-site 需要

tlsctl deploy --domain="example.com" --deploy="safeline-site"      # 部署到站点
tlsctl deploy --domain="example.com" --deploy="safeline-panel"     # 部署到面板
tlsctl deploy --domain="example.com" --deploy="safeline-portal"    # 部署到认证中心
```

> 自签名证书配置 `SAFELINE_IGNORE_SSL=true`。

## 七牛云

```bash
QINIU_ACCESS_KEY=your_access_key
QINIU_ACCESS_SECRET=your_access_secret
QINIU_DOMAIN=example.com        # 可选，默认用证书主域名

tlsctl deploy --domain="example.com" --deploy="qiniu-cdn"   # 七牛 CDN
tlsctl deploy --domain="example.com" --deploy="qiniu-oss"   # 绑定证书中心
```

## 百度云 CDN `baidu-cdn`

```bash
BAIDU_ACCESS_KEY=your_access_key
BAIDU_SECRET_KEY=your_secret_key
BAIDU_DOMAIN=example.com        # 可选
tlsctl deploy --domain="example.com" --deploy="baidu-cdn"
```

## 华为云 CDN `huaweicloud-cdn`

```bash
HUAWEI_ACCESS_KEY=your_access_key
HUAWEI_SECRET_KEY=your_secret_key
HUAWEI_DOMAIN=example.com       # 可选
tlsctl deploy --domain="example.com" --deploy="huaweicloud-cdn"
```

## 火山引擎

```bash
VOLC_ACCESS_KEY=your_access_key
VOLC_SECRET_KEY=your_secret_key
VOLC_REGION=cn-north-1
VOLC_DOMAIN=example.com         # 可选

tlsctl deploy --domain="example.com" --deploy="volcengine-cdn"    # CDN
tlsctl deploy --domain="example.com" --deploy="volcengine-dcdn"   # DCDN
```

## 多吉云 CDN `doge-cdn`

```bash
DOGE_ACCESS_KEY=your_access_key
DOGE_SECRET_KEY=your_secret_key
DOGE_DOMAIN=example.com         # 可选
tlsctl deploy --domain="example.com" --deploy="doge-cdn"
```

## LeCDN `lecdn`

```bash
LECDN_URL=https://example.com
LECDN_USERNAME=your_username
LECDN_PASSWORD=your_password
LECDN_SITE_ID=1
LECDN_DOMAIN=example.com        # 可选
tlsctl deploy --domain="example.com" --deploy="lecdn"
```

## 雨云 SSL 中心 `rainyun`

```bash
RAINYUN_API_KEY=your_api_key
RAINYUN_CERT_ID=your_cert_id
tlsctl deploy --domain="example.com" --deploy="rainyun"
```

## Webhook `webhook`

通过 HTTP 请求推送证书，请求体模板支持 `__domain__` / `__cert__` / `__key__` 占位符：

```bash
cat > ~/.tlsctl/.env << EOF
WEBHOOK_URL=https://example.com/api/deploy
WEBHOOK_METHOD=POST
WEBHOOK_DATA={"domain":"__domain__","cert":"__cert__","key":"__key__"}
WEBHOOK_HEADERS="Authorization: Bearer your_token
Content-Type: application/json"
EOF

tlsctl deploy --domain="example.com" --deploy="webhook"
```

> 支持 GET（Query 参数）与 POST/PUT/PATCH；自签名目标配置 `WEBHOOK_IGNORE_SSL=true`。

## 查看已登记任务

部署成功后，域名会登记到 `~/.tlsctl/deploy.json`。查看与管理：

```bash
tlsctl scheduled:list
tlsctl scheduled:remove --domain="example.com"
```
