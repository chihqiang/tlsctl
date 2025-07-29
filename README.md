### 🔐 tlsctl

一个专为开发者与运维设计的命令行工具，支持 SSL/TLS 证书的申请、续签与部署，助你轻松管理 HTTPS 全流程。

## ✨ 功能特点

🌍 **支持多家 ACME CA**：兼容 Let's Encrypt、ZeroSSL、Google Trust Services 等主流 ACME 证书颁发机构

🔒 **自动申请 DV 证书**：支持 DNS-01 与 HTTP-01 两种验证方式

⏱️ **智能续签机制**：自动检测证书有效期并续签，保障服务持续可用

📁 **灵活的证书存储**：可将证书保存到本地或指定的自定义目录

☁️ **支持多种 DNS 服务商**：兼容阿里云、西部数码、京东云、百度云、腾讯云、华为云、AWS、GoDaddy、Cloudflare 等主流 DNS 平台

🚀 **多种部署方式支持**：支持本地部署、SSH 部署、腾讯云、阿里云等自动上传部署方式

### 🛠️ 源码编译安装（适合折腾的你）

```
git clone https://github.com/chihqiang/tlsctl.git
cd tlsctl && make build
```

## 🚀 快速开始

#### 申请证书(以webroot为例)

~~~
tlsctl create --domain="test.example.com" --http.webroot="/data/wwwroot/test.example.com"
~~~

#### 通过本地部署到nginx目录

~~~
tlsctl deploy --domain="test.example.com" --deploy="local"
~~~

> 默认保存路径为：`/etc/nginx/ssl/`

## ⏱️ 定时任务

#### 执行证书续签任务

~~~
tlsctl scheduled:run
~~~

#### 添加为 systemd 启动服务（推荐方式）

~~~
cat > /etc/systemd/system/tlsctl-scheduled.service << EOF
[Unit]
Description=Start tlsctl Scheduled Task
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/tlsctl scheduled:run --http.webroot="/var/www/html"
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF
~~~

#### 常用 systemctl 命令：

~~~
# 重新加载所有 systemd 服务配置文件（全局作用）
systemctl daemon-reload
# 启动服务
systemctl start tlsctl-scheduled.service
# 停止服务
systemctl stop tlsctl-scheduled.service
# 重启服务
systemctl restart tlsctl-scheduled.service
# 重载配置（服务支持 reload 才生效）
systemctl reload tlsctl-scheduled.service
# 设置为开机自启
systemctl enable tlsctl-scheduled.service
# 禁用开机启动
systemctl disable tlsctl-scheduled.service
# 查询是否已设置开机启动
systemctl is-enabled tlsctl-scheduled.service
# 查看当前状态（运行中 / 停止 / 异常）
systemctl status tlsctl-scheduled.service
# 查看全部日志
journalctl -u tlsctl-scheduled.service
# 实时跟踪日志输出
journalctl -fu tlsctl-scheduled.service
# 查看最近 10 分钟内的日志
journalctl -u tlsctl-scheduled --since "10 minutes ago"
# 修改了 .service 文件后必须执行
systemctl daemon-reload
# 重置“失败”状态（如服务启动失败后恢复）
systemctl reset-failed tlsctl-scheduled.service
~~~

## ⚙️ 本地配置 `.env`

你可以通过 `.env` 文件配置自动部署后动作，例如自动重载 nginx：

~~~
cat > ~/.tlsctl/.env << EOF
LOCAL_POST_COMMAND="nginx -s reload"
EOF
~~~

## 📦 Nginx 配置示例

#### 申请证书时添加 Webroot 路径：

~~~
tlsctl create --domain="test.example.com" --http.webroot="/var/www/html"

//Add in nginx configuration
location  /.well-known/ {
  alias /var/www/html/.well-known/;
}
~~~

#### 添加到 nginx 配置中：

```
listen 443 ssl;
ssl_certificate /etc/nginx/ssl/test.example.com.pem;
ssl_certificate_key /etc/nginx/ssl/test.example.com.key;
ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;
ssl_prefer_server_ciphers on;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
```

## 🔐 使用 EAB（External Account Binding）

当使用 ZeroSSL、Google CA 等服务商时，你可能需要配置 `kid` 和 `hmacEncoded`：

- ZeroSSL 生成方式：https://zerossl.com/documentation/acme/generate-eab-credentials/
- Google Public CA 教程：https://cloud.google.com/certificate-manager/docs/public-ca-tutorial?hl=zh-cn
