# ⏱️ 定时任务

tlsctl 内置自动续签：周期性地检查已登记域名的证书有效期，即将过期时自动重新申请并重新部署。

## 工作原理

1. 执行 `tlsctl deploy` 时，域名被登记到 `~/.tlsctl/deploy.json`；
2. 运行 `tlsctl scheduled:run` 后进入循环，按 `--interval`（默认 24h）周期检查；
3. 每次检查读取每个域名的证书，若剩余天数小于 `--day`（默认 1）则自动续签并部署到登记的所有方式。

## 手动执行一次检查

```bash
tlsctl scheduled:run
```

## 相关命令

```bash
# 查看已登记的定时任务
tlsctl scheduled:list

# 删除指定域名的定时任务
tlsctl scheduled:remove --domain="example.com"
```

## 常用参数

| 参数 | 环境变量 | 默认值 | 说明 |
| ---- | -------- | ------ | ---- |
| `--interval` | `TLSCTL_RENEW_INTERVAL` | `24h` | 检查周期，如 `12h`、`30m`、`1h30m` |
| `--day` | `TLSCTL_DAY` | `1` | 距过期少于多少天时续签 |

示例（每 12 小时检查，剩余 7 天内续签）：

```bash
tlsctl scheduled:run --interval="12h" --day="7"
```

## 作为 systemd 服务运行（推荐）

### 1. 创建服务文件

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

> `--http.webroot` 与你的实际验证方式保持一致；若用 DNS 验证请改为 `--dns="服务商名"`。
> 如需自定义路径或环境变量，可加 `EnvironmentFile=/etc/tlsctl.env` 等。

### 2. 启动并设置开机自启

```bash
systemctl daemon-reload
systemctl enable tlsctl-scheduled.service
systemctl start tlsctl-scheduled.service
```

## 常用 systemctl 命令

```bash
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
# 重置"失败"状态（如服务启动失败后恢复）
systemctl reset-failed tlsctl-scheduled.service
```

## 定时任务的健壮性

- 单次部署有 10 分钟超时，云端任务卡死不会导致进程永久挂起；
- 单个域名申请/续签失败只记录警告并继续处理下一个域名，不会中断整个循环；
- 收到 `SIGINT` / `SIGTERM` 信号时优雅退出。
