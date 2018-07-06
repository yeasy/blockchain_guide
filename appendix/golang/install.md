### 安装与配置 Golang 环境

Golang 环境安装十分简单，可以通过包管理器或自行下载方式进行，为了使用最新版本的 Golang 环境，推荐大家通过下载环境包方式进行安装。

首先，从 [https://golang.org/dl/](https://golang.org/dl/) 页面查看最新的软件包，并根据自己的平台进行下载，例如 Linux 环境下，目前最新的环境包为 [https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz](https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz)。

下载后，直接进行环境包的解压，存放到默认的 `/usr/local/go` 目录（否则需要配置 $GOROOT 环境变量指向自定义位置）下。

```bash
$ sudo tar -C /usr/local -xzf go1.8.linux-amd64.tar.gz
```

此时，查看 `/usr/local/go` 路径下，可能看到如下几个子目录。

* api：Go API 检查器的辅助文件，记录了各个版本的 API 特性。
* bin：Go 语言相关的工具的二进制命令。
* doc：存放文档。
* lib：一些第三方库。
* misc：编辑器和开发环境的支持插件。
* pkg：存放不同平台的标准库的归档文件（.a 文件）。
* src：所有实现的源码。
* test：存放测试文件。

安装完毕后，可以添加 Golang 工具命令所在路径到系统路径，方便后面使用。并创建 `$GOPATH` 环境变量，指向某个本地创建好的目录（如 $HOME/Go），作为后面 Golang 项目的存放目录。

添加如下环境变量到用户启动配置（如 `$HOME/.bashrc`）中。

```sh
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/Go
export GOROOT=/usr/local/go
```

其它更多平台下安装，可以参考 [https://golang.org/doc/install](https://golang.org/doc/install)。

