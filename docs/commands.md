# ⌨️ 命令详解

所有子命令。运行 `tlsctl --help` 可查看全部全局参数。

## 全局参数

| 参数 | 别名 | 默认值 | 说明 |
| ---- | ---- | ------ | ---- |
| `--server` | `-s` | Let's Encrypt | ACME CA 服务器地址 |
| `--email` | `-m` | 系统生成 | 账号邮箱，不传则生成 `tlsctl-<hostname>@<os>.com` |
| `--path` | | `~/.tlsctl` | 数据存储目录 |
| `--kid` | | | EAB Key Identifier |
| `--hmac` | | | EAB HMAC Key |
| `--dns` | | | DNS-01 验证的服务商名 |
| `--key-type` | | `RSA2048` | 私钥类型（RSA2048/3072/4096/8192、EC256/384） |
| `--domain` | `-d` | | 申请/操作的域名，可多次指定 |
| `--deploy` | | `local` | 部署方式 |

## `create` — 申请证书

```bash
tlsctl create --domain="example.com" --http.webroot="/var/www/html"
```

### 验证方式

**HTTP-01**（需 80 端口或 webroot 可被公网访问）：

```bash
# webroot 验证（推荐搭配 nginx）
tlsctl create --domain="example.com" --http.webroot="/var/www/html"

# 内置 HTTP 服务（80 端口）
tlsctl create --domain="example.com" --http.port=":80"

# 反代场景下指定代理头
tlsctl create --domain="example.com" --http.proxy-header="X-Forwarded-Proto"

# 通过 S3 bucket 验证
tlsctl create --domain="example.com" --http.s3-bucket="my-bucket"

# 通过 memcached 验证
tlsctl create --domain="example.com" --http.memcached-host="127.0.0.1:11211"
```

**DNS-01**（无需开放端口，推荐通配符证书）：

```bash
tlsctl create --domain="example.com" --dns="aliyun"
```

**TLS-ALPN-01**（需 443 端口）：

```bash
tlsctl create --domain="example.com" --tls --tls.port=":443"
```

### 多域名（SAN 证书）

```bash
tlsctl create --domain="example.com" --domain="www.example.com" --dns="tencentcloud"
```

### 指定 CA 与 EAB

```bash
# 使用 ZeroSSL（需 EAB）
tlsctl create --domain="example.com" \
  --server="https://acme.zerossl.com/v2/DV90" \
  --kid="你的 kid" --hmac="你的 hmac"

# 使用 Google Public CA（需 EAB）
tlsctl create --domain="example.com" \
  --server="https://dv.acme-v02.api.pki.goog/directory" \
  --kid="你的 kid" --hmac="你的 hmac"
```

### 指定邮箱与私钥类型

```bash
tlsctl create --domain="example.com" --email="admin@example.com" --key-type="EC256"
```

> 不传 `--email` 时使用系统生成的邮箱（同一台机器稳定一致，可复用 ACME 账号）。

## `deploy` — 部署证书

```bash
tlsctl deploy --domain="example.com" --deploy="local"
```

部署成功后会把该域名登记到定时任务（`~/.tlsctl/deploy.json`）。详见 [部署方式](deploy.md)。

## `list` — 列出证书

```bash
tlsctl list
```

以表格形式显示本机所有证书的域名、有效期与保存路径。

## `clear` — 清理数据

```bash
tlsctl clear
```

删除 `--path` 下的全部数据（证书、ACME 账号、临时文件）。

> ⚠️ 安全保护：拒绝删除 `/`、`~`、`/etc`、`/usr` 等系统目录，防止误操作。

## `localhost` — 本地开发证书

```bash
tlsctl localhost
# 指定 hosts
tlsctl localhost --hosts="localhost" --hosts="127.0.0.1"
```

生成自签 CA 与 localhost 证书，并打印在系统上信任根证书的命令。

## `scheduled:*` — 定时任务

```bash
tlsctl scheduled:run      # 启动定时检查循环（每 24h，可配 --interval）
tlsctl scheduled:list     # 列出已登记的定时任务
tlsctl scheduled:remove --domain="example.com"   # 删除定时任务
```

详见 [定时任务](scheduled.md)。

## `help:*` — 帮助

```bash
tlsctl help:deploy   # 查看各部署方式的环境变量与字段
tlsctl help:dns      # 查看各 DNS 服务商的环境变量与字段
```
