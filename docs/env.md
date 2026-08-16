# 🌐 环境变量配置

tlsctl 支持通过 `.env` 文件或系统环境变量配置，优先级：命令行参数 > `.env` 文件 > 系统环境变量。

`.env` 文件默认路径为 `~/.tlsctl/.env`，程序启动时会自动加载。

## 全局环境变量

| 变量 | 对应参数 | 说明 |
| ---- | -------- | ---- |
| `TLSCTL_SERVER` | `--server` | ACME CA 服务器地址 |
| `TLSCTL_EMAIL` | `--email` | 账号邮箱（不传则系统生成） |
| `TLSCTL_PATH` | `--path` | 数据存储目录（默认 `~/.tlsctl`） |
| `TLSCTL_EAB_KID` | `--kid` | EAB Key Identifier |
| `TLSCTL_EAB_HMAC` | `--hmac` | EAB HMAC Key |
| `TLSCTL_DAY` | `--day` | 距过期多少天内续签（默认 1） |
| `TLSCTL_RENEW_INTERVAL` | `--interval` | 定时检查间隔（默认 24h） |

## DNS 服务商环境变量（DNS-01 验证）

| 服务商 | 必填环境变量 | 可选变量 |
| ------ | ------------ | -------- |
| 阿里云 | `ALICLOUD_ACCESS_KEY`、`ALICLOUD_SECRET_KEY` | `ALICLOUD_PROPAGATION_TIMEOUT`、`ALICLOUD_TTL` |
| 腾讯云 | `TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY` | `TENCENTCLOUD_PROPAGATION_TIMEOUT`、`TENCENTCLOUD_TTL` |
| AWS | `AWS_ACCESS_KEY_ID`、`AWS_SECRET_ACCESS_KEY`、`AWS_REGION`、`AWS_HOSTED_ZONE_ID` | `AWS_PROPAGATION_TIMEOUT`、`AWS_TTL` |
| 百度云 | `BAIDUCLOUD_ACCESS_KEY_ID`、`BAIDUCLOUD_SECRET_ACCESS_KEY` | `BAIDUCLOUD_PROPAGATION_TIMEOUT`、`BAIDUCLOUD_TTL` |
| 华为云 | `HUAWEICLOUD_ACCESS_KEY_ID`、`HUAWEICLOUD_SECRET_ACCESS_KEY`、`HUAWEICLOUD_REGION` | `HUAWEICLOUD_PROPAGATION_TIMEOUT`、`HUAWEICLOUD_TTL` |
| 京东云 | `JDCLOUD_ACCESS_KEY_ID`、`JDCLOUD_ACCESS_KEY_SECRET`、`JDCLOUD_REGION_ID` | `JDCLOUD_PROPAGATION_TIMEOUT`、`JDCLOUD_TTL` |
| GoDaddy | `GODADDY_API_KEY`、`GODADDY_API_SECRET` | `GODADDY_PROPAGATION_TIMEOUT`、`GODADDY_TTL` |
| Cloudflare | `CLOUDFLARE_API_TOKEN` | `CLOUDFLARE_PROPAGATION_TIMEOUT`、`CLOUDFLARE_TTL` |
| 西部数码 | `WESTCN_USERNAME`、`WESTCN_PASSWORD` | `WESTCN_PROPAGATION_TIMEOUT`、`WESTCN_TTL` |
| dynv6 | `DYNV6_HTTP_TOKEN` | `DYNV6_PROPAGATION_TIMEOUT`、`DYNV6_TTL` |

> 完整字段可用 `tlsctl help:dns` 查看。

### DNS 配置示例

```bash
cat > ~/.tlsctl/.env << EOF
# 阿里云 DNS
ALICLOUD_ACCESS_KEY=your_access_key
ALICLOUD_SECRET_KEY=your_secret_key
EOF
```

申请通配符证书：

```bash
tlsctl create --domain="*.example.com" --dns="aliyun"
```

## 部署环境变量

### 本地部署 `local`

| 变量 | 说明 |
| ---- | ---- |
| `LOCAL_CERT_PATH` | 证书输出路径（默认 `/etc/nginx/ssl/<domain>.pem`） |
| `LOCAL_KEY_PATH` | 私钥输出路径（默认 `/etc/nginx/ssl/<domain>.key`） |
| `LOCAL_PRE_COMMAND` | 部署前执行的命令 |
| `LOCAL_POST_COMMAND` | 部署后执行的命令（如重载 nginx） |

### SSH 部署 `ssh`

| 变量 | 说明 |
| ---- | ---- |
| `SSH_HOST` | 主机（默认 localhost） |
| `SSH_PORT` | 端口（默认 22） |
| `SSH_USERNAME` | 用户名 |
| `SSH_PASSWORD` | 密码 |
| `SSH_KEY` | 私钥内容 |
| `SSH_KEY_PASSPHRASE` | 私钥口令 |
| `SSH_USE_SCP` | 是否使用 SCP（默认 SFTP） |
| `SSH_CERT_PATH` | 远端证书路径 |
| `SSH_KEY_PATH` | 远端私钥路径 |
| `SSH_PRE_COMMAND` / `SSH_POST_COMMAND` | 远端前置/后置命令 |

### 腾讯云部署

公共变量：`TENCENTCLOUD_SECRET_ID`、`TENCENTCLOUD_SECRET_KEY`，各方式另有专属变量：

| 部署方式 | 专属变量 |
| -------- | -------- |
| `tcdn` | `TENCENTCLOUD_DOMAIN` |
| `ecdn` | `TENCENTCLOUD_DOMAIN` |
| `tssl` | 无（自动上传证书） |
| `cos` | `TENCENTCLOUD_REGION`、`TENCENTCLOUD_BUCKET`、`TENCENTCLOUD_DOMAIN` |
| `scf` | `TENCENTCLOUD_REGION`、`TENCENTCLOUD_DOMAIN` |
| `clb` | `TENCENTCLOUD_REGION`、`TENCENTCLOUD_LOADBALANCER_ID`、`TENCENTCLOUD_LISTENER_ID` |
| `waf` | `TENCENTCLOUD_REGION`、`TENCENTCLOUD_DOMAIN`、`TENCENTCLOUD_DOMAIN_ID` |
| `eo` | `TENCENTCLOUD_ZONE_ID` |
| `vod` | `TENCENTCLOUD_SUB_APP_ID` |

### 阿里云部署

公共变量：`ALIYUN_ACCESS_KEY_ID`、`ALIYUN_ACCESS_KEY_SECRET`，各方式另有专属变量：

| 部署方式 | 专属变量 |
| -------- | -------- |
| `cdn` | `ALIYUN_DOMAIN` |
| `dcdn` | `ALIYUN_DOMAIN` |
| `live` | `ALIYUN_REGION`、`ALIYUN_DOMAIN` |
| `oss` | `ALIYUN_REGION`、`ALIYUN_BUCKET`、`ALIYUN_DOMAIN` |
| `vod` | `ALIYUN_REGION`、`ALIYUN_DOMAIN` |
| `fc` | `ALIYUN_REGION`、`ALIYUN_RESOURCE_GROUP_ID`、`ALIYUN_VERSION` |

> 完整字段可用 `tlsctl help:deploy` 查看。
