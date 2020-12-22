## 本地编译组件

本地编译生成 Fabric 网络的各个组件，可以形成更直观的认识。Fabric 采用 Go 语言实现，推荐使用 Golang 1.10+ 版本进行编译。

下面将讲解如何编译生成 fabric-peer、fabric-orderer 和 fabric-ca 等组件的二进制文件，以及如何安装一些配置和开发辅助工具。如果用户在多服务器环境下进行部署，需要注意将文件复制到对应的服务器上。

### 环境配置

#### 操作系统

常见的 Linux 操作系统发行版（包括 Ubuntu、Redhat、CentOS 等）和 macOS 等都可以支持 Fabric。

内核推荐 3.10+ 版本，支持 64 位环境。下面将默认以 Ubuntu 18.04 操作系统为例进行讲解。

*注：运行 Fabric 节点需要的资源并不苛刻，作为实验，Fabric 节点甚至可以在树莓派（Raspberry Pi）上正常运行。但生产环境中往往需要较高的 CPU 和内存资源。*

#### 安装 Go 语言环境

可以访问 [golang.org](https://golang.org) 网站下载压缩包进行安装。不推荐使用自带的包管理器安装，版本往往比较旧。

如下载最新的 Go 1.13.4 稳定版本，可以采用如下命令：

```bash
$ curl -O https://dl.google.com/go/go1.13.4.linux-amd64.tar.gz
```

下载完成后，解压目录，并移动到合适的位置（如 /usr/local）：

```bash
$ tar -xvf go1.13.4.linux-amd64.tar.gz
$ sudo mv go /usr/local
```

配置 GOPATH 环境变量，同时可以加入 .bash_profile 文件中以长期生效：

```bash
export GOPATH=YOUR_LOCAL_GO_PATH/Go
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

此时，可以通过 `go version` 命令验证安装是否成功：

```bash
$ go version

go version go1.13.4 linux/amd64
```

#### 安装依赖包

编译 Fabric 代码依赖一些开发库，可以通过如下命令安装：

```bash
$ sudo apt-get update \
    && sudo apt-get install -y libsnappy-dev zlib1g-dev libbz2-dev libyaml-dev libltdl-dev libtool
```

#### 安装 Docker 环境

Fabric 目前采用 Docker 容器作为链码执行环境，因此即使在本地运行，Peer 节点上也需要安装 Docker 环境，推荐使用 1.18 或者更新的版本。

Linux 操作系统下可以通过如下命令来安装 Docker 最新版本：

```bash
$ curl -fsSL https://get.docker.com/ | sh
```

macOS 可以访问 https://docs.docker.com/docker-for-mac/install 下载 `Docker for Mac` 安装包自行安装：

### 获取代码

目前，Fabric 官方仓库托管在 Github 仓库（github.com/hyperledger/fabric）中供下载使用。

如果使用 1.13 之前版本的 Go 环境，需要将 Fabric 项目放到 $GOPATH 路径下。如下命令所示，创建 `$GOPATH/src/github.com/hyperledger` 目录结构并切换到该路径：

```bash
$ mkdir -p $GOPATH/src/github.com/hyperledger
$ cd $GOPATH/src/github.com/hyperledger
```

获取 Peer 和 Orderer 组件编译所需要的代码，两者目前在同一个 fabric 仓库中：

```bash
$ git clone https://github.com/hyperledger/fabric.git
```

为节约下载时间，读者可以指定 `--single-branch -b master --depth 1` 命令选项来指定只获取 master 分支最新代码：

```bash
$ git clone --single-branch -b master --depth 1 https://github.com/hyperledger/fabric.git
```

Fabric CA 组件在独立的 fabric-ca 仓库中，可以通过如下命令获取：

```bash
$ git clone https://github.com/hyperledger/fabric-ca.git
```

读者也可以直接访问 https://github.com/hyperledger/fabric/releases 和 https://github.com/hyperledger/fabric-ca/releases 来下载特定的 fabric 和 fabric-ca 发行版。

最后，检查确认 fabric 和 fabric-ca 两个仓库下载成功：

```bash
$ ls $GOPATH/src/github.com/hyperledger
fabric fabric-ca
```

### 编译安装 Peer 组件

配置版本号和编译参数：

```bash
$ PROJECT_VERSION=2.0.0
$ LD_FLAGS="-X github.com/hyperledger/fabric/common/metadata.Version=${PROJECT_VERSION} \
             -X github.com/hyperledger/fabric/common/metadata.BaseDockerLabel=org.hyperledger.fabric \
             -X github.com/hyperledger/fabric/common/metadata.DockerNamespace=hyperledger \
             -X github.com/hyperledger/fabric/common/metadata.BaseDockerNamespace=hyperledger"
```

通过如下命令编译并安装 fabric 的 peer 组件到 $GOPATH/bin 下：

```bash
$ CGO_CFLAGS=" " go install -tags "" -ldflags "$LD_FLAGS" \
    github.com/hyperledger/fabric/cmd/peer
```

当然，用户也可直接使用源码中的 Makefile 来进行编译，相关命令如下：

```bash
$ make peer
```

这种情况下编译生成的 peer 组件会默认放在 build/bin 路径下。

### 编译安装 Orderer 组件
通过如下命令编译并安装 fabric orderer 组件到 $GOPATH/bin 下：

```bash
$ CGO_CFLAGS=" " go install -tags "" -ldflags "$LD_FLAGS" \
    github.com/hyperledger/fabric/cmd/orderer
```

同样的，也可使用 Makefile 来编译安装 orderer 组件到 build/bin 路径下：

```bash
$ make orderer
```

### 编译安装 Fabric CA 组件
采用如下命令编译并安装 fabric-ca 相关组件到 $GOPATH/bin 下：

```bash
$ go install -ldflags "-X github.com/hyperledger/fabric-ca/lib/metadata.Version=$PROJECT_VERSION -linkmode external -extldflags '-static -lpthread'" \
    github.com/hyperledger/fabric-ca/cmd/...
```

### 编译安装配置辅助工具

Fabric 中还提供了一系列配置辅助工具，包括 cryptogen（本地生成组织结构和身份文件）、configtxgen（生成配置区块和配置交易）、configtxlator（解析转换配置信息）、discover（拓扑探测）、idemixgen（Idemix 证书生成）等，可以通过如下命令来快速编译和安装：

```bash

# 编译安装 cryptogen，等价于执行 make cryptogen
$ CGO_CFLAGS=" " \
    go install -tags "" -ldflags "$LD_FLAGS" \
    github.com/hyperledger/fabric/cmd/cryptogen

# 编译安装 configtxgen，等价于执行 make configtxgen
$ CGO_CFLAGS=" " \
    go install -tags "" -ldflags "$LD_FLAGS" \
    github.com/hyperledger/fabric/cmd/configtxgen

# 编译安装 configtxlator，等价于执行 make configtxlator
$ CGO_CFLAGS=" " \
    go install -tags "" -ldflags "$LD_FLAGS" \
    github.com/hyperledger/fabric/cmd/configtxlator

# 编译安装 discover，等价于执行 make discover
$ CGO_CFLAGS=" " \
    go install -tags "" -ldflags "$LD_FLAGS" \
    github.com/hyperledger/fabric/cmd/discover

# 编译安装 idemixgen，等价于执行 make idemixgen
$ CGO_CFLAGS=" " \
    go install -tags "" -ldflags "$LD_FLAGS" \
    github.com/hyperledger/fabric/cmd/idemixgen
```

另外，fabric 项目还提供了不少常见的编译命令，可以参考 Makefile 文件，例如编译所有的二进制文件可以使用如下命令：

```bash
$ make native
```

### 安装 Protobuf 支持和 Go 语言相关工具

Fabric 代码由 Go 语言构建，开发者可以选择安装如下的 Go 语言相关工具，方便开发和调试：

```bash
$ go get github.com/golang/protobuf/protoc-gen-go \
    && go get github.com/maxbrunsfeld/counterfeiter/v6 \
    && go get github.com/axw/gocov/... \
    && go get github.com/AlekSi/gocov-xml \
    && go get golang.org/x/tools/cmd/goimports \
    && go get golang.org/x/lint/golint \
    && go get github.com/estesp/manifest-tool \
    && go get github.com/client9/misspell/cmd/misspell \
    && go get github.com/estesp/manifest-tool \
    && go get github.com/onsi/ginkgo/ginkgo
```

### 示例配置

sampleconfig 目录下包括了一些示例配置文件，可以作为参考，包括：

* configtx.yaml：示例配置区块生成文件夹；
* orderer.yaml：示例 Orderer 节点配置文件；
* core.yaml：示例 Peer 节点配置文件；
* msp/config.yaml：示例组织身份配置文件；
* msp/：示例证书和秘钥文件。

可以将它们复制到默认的配置目录（/etc/hyperledger/fabric）下进行使用。
