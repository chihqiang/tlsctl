# 🔐 使用 EAB（External Account Binding）

部分 ACME CA（如 ZeroSSL、Google Public CA）要求使用 EAB（External Account Binding）才能注册账号。EAB 包含两个凭据：`kid` 和 `hmac`。

## 获取 EAB 凭据

- **ZeroSSL**：[生成 EAB 凭据](https://zerossl.com/documentation/acme/generate-eab-credentials/)
- **Google Public CA**：[教程](https://cloud.google.com/certificate-manager/docs/public-ca-tutorial?hl=zh-cn)

## 使用方式

### 方式一：命令行参数

```bash
# ZeroSSL
tlsctl create --domain="example.com" \
  --server="https://acme.zerossl.com/v2/DV90" \
  --kid="你的 kid" --hmac="你的 hmac"
```

### 方式二：`.env` 文件（推荐，避免命令历史泄露）

```bash
cat > ~/.tlsctl/.env << EOF
TLSCTL_SERVER=https://acme.zerossl.com/v2/DV90
TLSCTL_EAB_KID=你的 kid
TLSCTL_EAB_HMAC=你的 hmac
EOF

tlsctl create --domain="example.com"
```

### 方式三：环境变量

```bash
export TLSCTL_EAB_KID="你的 kid"
export TLSCTL_EAB_HMAC="你的 hmac"
tlsctl create --domain="example.com" --server="https://acme.zerossl.com/v2/DV90"
```

## 常用 EAB 服务商

| 服务商 | 服务器地址 | 说明 |
| ------ | ---------- | ---- |
| ZeroSSL | `https://acme.zerossl.com/v2/DV90` | 需要 EAB |
| Google Public CA | `https://dv.acme-v02.api.pki.goog/directory` | 需要 EAB |

> Let's Encrypt（默认 `--server`）不需要 EAB，直接使用即可。
