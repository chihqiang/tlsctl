# 🔐 tlsctl

一个专为开发者与运维设计的命令行工具，支持 SSL/TLS 证书的申请、续签与部署，助你轻松管理 HTTPS 全流程。

## ✨ 功能特点

🌍 **支持多家 ACME CA**：兼容 Let's Encrypt、ZeroSSL、Google Trust Services 等主流 ACME 证书颁发机构

🔒 **自动申请 DV 证书**：支持 DNS-01 与 HTTP-01 两种验证方式

⏱️ **智能续签机制**：自动检测证书有效期并续签，保障服务持续可用

📁 **灵活的证书存储**：可将证书保存到本地或指定的自定义目录

☁️ **支持多种 DNS 服务商**：兼容阿里云、西部数码、京东云、百度云、腾讯云、华为云、AWS、GoDaddy、Cloudflare 等主流 DNS 平台

🚀 **多种部署方式支持**：支持本地部署、SSH 部署、腾讯云、阿里云等自动上传部署方式

## ⚡ 安装

### 方式一：`go install`（推荐，需 Go 1.24+）

```bash
go install github.com/chihqiang/tlsctl/cmd/tlsctl@latest
```

> 安装完成后，请确保 `$(go env GOPATH)/bin` 已加入 `PATH`。

### 方式二：下载源码 + `make` 编译

```bash
git clone https://github.com/chihqiang/tlsctl.git
cd tlsctl
make build
```

> 编译产物为当前目录下的 `tlsctl`，可 `sudo cp tlsctl /usr/local/bin/`。

更多安装细节见 [安装与升级](docs/install.md)。

## 🚀 快速开始

### 申请证书（以 webroot 为例）

```bash
tlsctl create --domain="test.example.com" --http.webroot="/data/wwwroot/test.example.com"
```

### 本地部署到 nginx 目录

```bash
tlsctl deploy --domain="test.example.com" --deploy="local"
```

> 默认保存路径为：`/etc/nginx/ssl/`

### 定时任务（自动续签）

```bash
tlsctl scheduled:run
```

## 📚 详细文档

| 文档 | 说明 |
| ---- | ---- |
| [命令详解](docs/commands.md) | 全部子命令与参数 |
| [环境变量配置](docs/env.md) | `.env` 与 DNS/部署环境变量 |
| [部署方式](docs/deploy.md) | local / SSH / 腾讯云 / 阿里云 等 |
| [定时任务](docs/scheduled.md) | 自动续签与 systemd 配置 |
| [Nginx 配置](docs/nginx.md) | 与 Nginx 配合示例 |
| [EAB 使用](docs/eab.md) | ZeroSSL / Google CA 配置 |
| [安装与升级](docs/install.md) | 安装方式与常见问题 |
