## 安装 Fabric

本节主要介绍如何安装 Hyperledger Fabric 的相关组件（二进制文件和 Docker 镜像）。

Hyperledger Fabric 官方提供了 `install-fabric.sh` 脚本，可以按需下载 Docker 镜像、二进制文件（`peer`、`orderer`、`configtxgen`、`osnadmin` 等）以及示例代码（`fabric-samples`）。

## 准备工作

### 1. 操作系统

Linux (Ubuntu/CentOS), macOS 或 Windows (通过 WSL2)。

### 2. 安装 Git 和 cURL

确保系统已安装 `git` 和 `curl`。
```bash
sudo apt-get install git curl
```

### 3. 安装 Docker

Fabric 运行依赖 Docker 容器环境。
*   **Docker**: 24.0 或更高版本。（推荐使用 24.0+ 以支持最新特性和安全补丁）
*   **Docker Compose**: 推荐安装 Docker Desktop（包含 Docker Compose）或独立的 Docker Compose 插件。

## 使用官方脚本安装 (推荐)

最简单的方法是使用官方提供的 `install-fabric.sh` 脚本。

### 步骤 1：创建工作目录

```bash
mkdir -p $HOME/hyperledger
cd $HOME/hyperledger
```

### 步骤 2：运行安装脚本

以下命令会先从官方仓库下载脚本，再执行安装。生产环境或受控环境中应先审查脚本内容，避免使用短链接和直接 pipe-to-bash 的方式执行远程脚本。

```bash
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
chmod +x install-fabric.sh
./install-fabric.sh docker binary samples
```

官方脚本默认使用当前 Fabric LTS 和 Fabric CA 版本。需要固定版本时显式指定，例如 Fabric v2.5.15 和 Fabric CA v1.5.17：

```bash
./install-fabric.sh --fabric-version 2.5.15 --ca-version 1.5.17 docker binary samples
```

该脚本会执行以下操作：
1.  克隆 `fabric-samples` 仓库。
2.  下载 Fabric 二进制文件（`peer`, `orderer`, `configtxgen` 等）到 `fabric-samples/bin` 目录。
3.  下载指定版本的 Hyperledger Fabric Docker 镜像。
4.  将 Fabric CA 客户端和服务端二进制文件下载到 `fabric-samples/bin`。

### 步骤 3：配置环境变量

为了方便使用下载的二进制文件，将其路径添加到环境变量中：

```bash
export PATH=$HOME/hyperledger/fabric-samples/bin:$PATH
```

验证安装是否成功：

```bash
peer version
```

## 从源码编译安装 (高级)

如果你是开发者，或者需要跟踪正式发布版本之外的源码变更，可以选择从源码编译。Fabric 3.x 已是正式发布系列，不应再按 Beta 版本处理。

### 1. 安装 Go 语言环境

确保安装目标 Fabric 版本发布说明中测试过的 Go 版本。例如 Fabric v2.5.15 和 v3.1.4 均使用 Go 1.26.0 进行测试。

### 2. 获取源码

```bash
git clone https://github.com/hyperledger/fabric.git
cd fabric
```

### 3. 编译

使用 `make` 命令编译组件。编译后的二进制文件位于 `build/bin` 目录。

```bash

# 编译所有组件
make all

# 或者单独编译
make peer
make orderer
make configtxgen
```

### 4. 安装

将编译好的二进制文件移动到系统路径或添加到 PATH 中。
