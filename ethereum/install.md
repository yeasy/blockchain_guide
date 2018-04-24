## 安装客户端

本节将介绍如何安装 Geth，即 Go 语言实现的以太坊客户端。这里以 Ubuntu 16.04 操作系统为例，介绍从 PPA 仓库和从源码编译这两种方式来进行安装。

### 从 PPA 直接安装

首先安装必要的工具包。

```sh
$ apt-get install software-properties-common
```

之后用以下命令添加以太坊的源。

```sh
$ add-apt-repository -y ppa:ethereum/ethereum
$ apt-get update
```

最后安装 go-ethereum。

```sh
$ apt-get install ethereum
```

安装成功后，则可以开始使用命令行客户端 Geth。可用 `geth --help` 查看各命令和选项，例如，用以下命令可查看 Geth 版本为 1.6.1-stable。

```sh
$ geth version

Geth
Version: 1.6.1-stable
Git Commit: 021c3c281629baf2eae967dc2f0a7532ddfdc1fb
Architecture: amd64
Protocol Versions: [63 62]
Network Id: 1
Go Version: go1.8.1
Operating System: linux
GOPATH=
GOROOT=/usr/lib/go-1.8
```

### 从源码编译

也可以选择从源码进行编译安装。

#### 安装 Go 语言环境

Go 语言环境可以自行访问 [golang.org](https://golang.org) 网站下载二进制压缩包安装。注意不推荐通过包管理器安装版本，往往比较旧。

如下载 Go 1.8 版本，可以采用如下命令。

```bash
$ curl -O https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz
```

下载完成后，解压目录，并移动到合适的位置（推荐为 /usr/local 下）。

```bash
$ tar -xvf go1.8.linux-amd64.tar.gz
$ sudo mv go /usr/local
```

安装完成后记得配置 GOPATH 环境变量。

```bash
$ export GOPATH=YOUR_LOCAL_GO_PATH/Go
$ export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
```

此时，可以通过 `go version` 命令验证安装 是否成功。

```bash
$ go version

go version go1.8 linux/amd64
```

#### 下载和编译 Geth

用以下命令安装 C 的编译器。

```sh
$ apt-get install -y build-essential
```

下载选定的 go-ethereum 源码版本，如最新的社区版本：

```bash
$ git clone https://github.com/ethereum/go-ethereum
```

编译安装 Geth。

```bash
$ cd go-ethereum
$ make geth
```

安装成功后，可用 `build/bin/geth --help` 查看各命令和选项。例如，用以下命令可查看 Geth 版本为 1.6.3-unstable。

```bash
$ build/bin/geth version
Geth
Version: 1.6.3-unstable
Git Commit: 067dc2cbf5121541aea8c6089ac42ce07582ead1
Architecture: amd64
Protocol Versions: [63 62]
Network Id: 1
Go Version: go1.8
Operating System: linux
GOPATH=/usr/local/gopath/
GOROOT=/usr/local/go
```

