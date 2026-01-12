# 安装客户端

本节介绍如何安装 Geth（Go 语言实现的以太坊客户端）。Geth 是以太坊最主流的执行层客户端之一。

*注：自 The Merge（2022 年）之后，运行完整的以太坊节点需要同时安装执行层客户端（如 Geth）和共识层客户端（如 Prysm、Lighthouse）。仅安装 Geth 无法独立同步网络。*

## 快速安装（推荐）

### macOS

使用 Homebrew 安装：

```bash
brew tap ethereum/ethereum
brew install ethereum
```

### Ubuntu/Debian

添加 PPA 源并安装：

```bash
sudo add-apt-repository -y ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install ethereum
```

### Windows

从 [Geth 官方下载页面](https://geth.ethereum.org/downloads/) 下载 Windows 版本安装包。

### 验证安装

```bash
geth version
# 输出示例：
# Geth
# Version: 1.13.x-stable
# ...
```

## 从源码编译

如需使用最新开发版本或进行定制，可从源码编译。

### 1. 安装 Go 语言环境

访问 [go.dev/dl](https://go.dev/dl/) 下载并安装 Go 1.21 或更高版本。

```bash
# 验证安装
go version
```

*注：现代 Go 使用 Go Modules 管理依赖，无需配置 GOPATH。*

### 2. 克隆并编译

```bash
git clone https://github.com/ethereum/go-ethereum.git
cd go-ethereum
make geth
```

编译完成后，可执行文件位于 `build/bin/geth`。

```bash
./build/bin/geth version
```

## 运行节点

由于 The Merge 后需要执行层+共识层配合，推荐使用官方的 [ethereum-docker](https://github.com/eth-educators/ethstaker-guides) 或 [eth-docker](https://eth-docker.net/) 等工具来简化部署。

基本的 Geth 启动命令：

```bash
# 连接主网（需配合共识层客户端）
geth --http --http.api eth,net,engine,admin

# 连接 Sepolia 测试网
geth --sepolia --http
```

更多配置和同步选项请参考 [Geth 官方文档](https://geth.ethereum.org/docs/)。
