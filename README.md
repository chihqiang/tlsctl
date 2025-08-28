# 🔐 tlsctl

A command-line tool designed for developers and operators. It supports the application, renewal, and deployment of SSL/TLS certificates, helping you easily manage the entire HTTPS process.

## ✨ Features

🌍 **Supports multiple ACME CAs**: Compatible with mainstream ACME certificate authorities such as Let's Encrypt, ZeroSSL, and Google Trust Services

🔒 **Automatically apply for DV certificates**: Supports both DNS-01 and HTTP-01 verification methods

⏱️**Smart Renewal Mechanism**: Automatically detects certificate validity and renews it to ensure continuous service availability

📁 **Flexible certificate storage**: Certificates can be saved locally or in a specified custom directory

☁️ **Support multiple DNS service providers**: Compatible with mainstream DNS platforms such as Alibaba Cloud, Western Digital, JD Cloud, Baidu Cloud, Tencent Cloud, Huawei Cloud, AWS, GoDaddy, Cloudflare, etc.

🚀 **Multiple deployment methods supported**: Supports local deployment, SSH deployment, Tencent Cloud, Alibaba Cloud and other automatic upload deployment methods

## 🚀 One-click installation

```bash
curl -sSL https://cnb.cool/zhiqiangwang/tlsctl/-/git/raw/main/install.sh | bash
```

## 🛠️ Source code reading & local building (for advanced users)

```bash
git clone https://github.com/chihqiang/tlsctl.git
cd tlsctl && make build
```

## 🚀 Quick Start

### Apply for a certificate (taking webroot as an example)

```bash
tlsctl create --domain="test.example.com" --http.webroot="/data/wwwroot/test.example.com"
```

### By locally deploying to the nginx directory

```bash
tlsctl deploy --domain="test.example.com" --deploy="local"
```

> The default save path is：`/etc/nginx/ssl/`

## ⏱️ Scheduled tasks

### Perform certificate renewal tasks

```bash
tlsctl scheduled:run
```

### Add as systemd startup service (recommended method)

```bash
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
```

#### Common systemctl commands

```bash
# Reload all systemd service configuration files (globally)
systemctl daemon-reload
# Start the service
systemctl start tlsctl-scheduled.service
# Stop the service
systemctl stop tlsctl-scheduled.service
# Restart the service
systemctl restart tlsctl-scheduled.service
# Reload the configuration (this takes effect only if the service supports reloading)
systemctl reload tlsctl-scheduled.service
# Set it to start automatically at boot
systemctl enable tlsctl-scheduled.service
# Disable it from starting at boot
systemctl disable tlsctl-scheduled.service
# Check if it is set to start at boot
systemctl is-enabled tlsctl-scheduled.service
# Check the current status (running/stopped/abnormal)
systemctl status tlsctl-scheduled.service
# View all logs
journalctl -u tlsctl-scheduled.service
# Track log output in real time
journalctl -fu tlsctl-scheduled.service
# View logs from the last 10 minutes
journalctl -u tlsctl-scheduled --since "10 minutes ago"
# Required after modifying the .service file
systemctl daemon-reload
# Reset a "failed" status (e.g., recover from a service startup failure)
systemctl reset-failed tlsctl-scheduled.service
```

## ⚙️ Local configuration `.env`

You can configure automatic post-deployment actions via `.env` files, such as automatically reloading nginx:

```bash
cat > ~/.tlsctl/.env << EOF
LOCAL_POST_COMMAND="nginx -s reload"
EOF
```

## 📦 Nginx Configuration Example

### Add the Webroot path when applying for a certificate:

```bash
tlsctl create --domain="test.example.com" --http.webroot="/var/www/html"

//Add in nginx configuration
location  /.well-known/ {
  alias /var/www/html/.well-known/;
}
```

### Add the certificate to the nginx configuration

```bash
listen 443 ssl;
ssl_certificate /etc/nginx/ssl/test.example.com.pem;
ssl_certificate_key /etc/nginx/ssl/test.example.com.key;
ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;
ssl_prefer_server_ciphers on;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
```

## 🔐 Using EAB (External Account Binding)

When using ZeroSSL, Google CA, etc., you may need to configure `kid` and `hmacEncoded`:

- ZeroSSL generation method：<https://zerossl.com/documentation/acme/generate-eab-credentials/>
- Google Public CA Tutorial：<https://cloud.google.com/certificate-manager/docs/public-ca-tutorial?hl=zh-cn>
