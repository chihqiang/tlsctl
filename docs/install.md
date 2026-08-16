# 🛠️ 安装与升级

tlsctl 需要 Go 1.24+ 环境，提供两种安装方式。

## 方式一：`go install`（推荐）

```bash
go install github.com/chihqiang/tlsctl/cmd/tlsctl@latest
```

安装完成后，可执行文件位于 `$(go env GOPATH)/bin`，请确保该目录已加入 `PATH`：

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

可将上面这行追加到 `~/.zshrc` / `~/.bashrc`，避免每次打开终端都要手动设置。

### 升级

```bash
go install github.com/chihqiang/tlsctl/cmd/tlsctl@latest
```

重复执行 `go install` 即可升级到最新版。

## 方式二：下载源码 + `make` 编译

```bash
git clone https://github.com/chihqiang/tlsctl.git
cd tlsctl
make build
```

编译产物为当前目录下的 `tlsctl`，可复制到系统路径（可选）：

```bash
sudo cp tlsctl /usr/local/bin/
```

### 交叉编译（可选）

`Makefile` 提供了 Linux amd64 的交叉编译目标：

```bash
make build_linux_amd64
```

### 其它 make 目标

```bash
make build    # 编译当前平台
make check    # 格式化 + go mod tidy
make test     # 编译并运行测试
make clean    # 清理产物
```

## 验证安装

```bash
tlsctl --version
# 输出示例：tlsctl version a4d3ef7 darwin/arm64
```

## 常见问题

### 命令找不到

确认 `$(go env GOPATH)/bin` 已加入 `PATH`：

```bash
echo "$(go env GOPATH)/bin"
which tlsctl
```

### 数据存在哪里

默认数据目录为 `~/.tlsctl/`，包含证书、ACME 账号、`.env` 配置等。可通过 `--path` 或 `TLSCTL_PATH` 环境变量修改。
