# 📦 Nginx 配置示例

## 申请证书时添加 Webroot 路径

申请前先在 nginx 配置中把 `/.well-known/` 指向 webroot 目录：

```bash
tlsctl create --domain="example.com" --http.webroot="/var/www/html"
```

在 nginx 配置中添加：

```nginx
location /.well-known/ {
  alias /var/www/html/.well-known/;
}
```

> webroot 目录必须能被公网访问（80 端口放行）。

## HTTPS 站点配置

把证书部署到 nginx 后，站点配置如下：

```yml
listen 443 ssl;
ssl_certificate /etc/nginx/ssl/example.com.pem;
ssl_certificate_key /etc/nginx/ssl/example.com.key;
ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
ssl_protocols TLSv1.1 TLSv1.2 TLSv1.3;
ssl_prefer_server_ciphers on;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
```

## HTTP 跳转 HTTPS

```nginx
server {
    listen 80;
    server_name example.com www.example.com;
    return 301 https://$host$request_uri;
}
```

## 部署后自动重载

在 `~/.tlsctl/.env` 中配置，部署完成后自动重载 nginx：

```bash
cat > ~/.tlsctl/.env << EOF
LOCAL_POST_COMMAND="nginx -s reload"
EOF
```

配合定时任务（`scheduled:run`），续签后证书会重新部署并自动重载 nginx，实现全自动 HTTPS 续期。
