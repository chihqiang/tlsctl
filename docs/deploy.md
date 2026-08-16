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

### 边缘安全加速 `eo`

```bash
TENCENTCLOUD_ZONE_ID=zone-xxx
tlsctl deploy --domain="example.com" --deploy="eo"
```

### VOD 点播 `vod`

```bash
TENCENTCLOUD_SUB_APP_ID=123456
TENCENTCLOUD_DOMAIN=example.com
tlsctl deploy --domain="example.com" --deploy="vod"
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

## 查看已登记任务

部署成功后，域名会登记到 `~/.tlsctl/deploy.json`。查看与管理：

```bash
tlsctl scheduled:list
tlsctl scheduled:remove --domain="example.com"
```
