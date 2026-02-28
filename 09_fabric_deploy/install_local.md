# 安装 Fabric

本节主要介绍如何安装 Hyperledger Fabric 的相关组件（二进制文件和 Docker 镜像）。

Hyperledger Fabric 官方提供了一个方便的安装脚本，可以一键下载所需的 Docker 镜像、二进制文件（peer, orderer, configtxgen 等）以及示例代码（fabric-samples）。

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
*   **Docker**: 18.03 或更高版本。
*   **Docker Compose**: 推荐安装 Docker Desktop（包含 Docker Compose）或独立的 Docker Compose 插件。

## 使用官方脚本安装 (推荐)

最简单的方法是使用官方提供的 `install-fabric.sh` 脚本。

### 步骤 1：创建工作目录

```bash
mkdir -p $HOME/hyperledger
cd $HOME/hyperledger
```

### 步骤 2：运行安装脚本

以下命令将下载最新的生产版本（LTS）Fabric 镜像和二进制文件：

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s
```

如果你需要指定版本（例如 Fabric v2.5.9, CA v1.5.9），可以使用：

```bash
curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.5.9 1.5.9
```

该脚本会执行以下操作：
1.  克隆 `fabric-samples` 仓库。
2.  下载 Fabric 二进制文件（`peer`, `orderer`, `configtxgen` 等）到 `fabric-samples/bin` 目录。
3.  下载指定版本的 Hyperledger Fabric Docker 镜像。

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

如果你是开发者，或者需要最新的特性（如 Fabric 3.0 Beta），可以选择从源码编译。

### 1. 安装 Go 语言环境

确保安装了 Go 1.20+ 版本（参考附录 Go 语言安装章节）。

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
