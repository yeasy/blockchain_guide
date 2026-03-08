### 安装与配置 Golang 环境

Go 语言环境安装十分简单，可以通过包管理器或自行下载方式进行。为了使用最新版本的 Go 环境，推荐通过官方网站下载安装包。

#### 1. 下载安装包

访问 Go 语言官方下载页面 [https://go.dev/dl/](https://go.dev/dl/)，根据你的操作系统（Linux, macOS, Windows）选择对应的安装包。

以 Linux 为例，下载最新的 `tar.gz` 包（例如 `go1.21.0.linux-amd64.tar.gz`）。

#### 2. 安装

**Linux / macOS (Tarball):**

删除旧的安装（如果存在）并将新包解压到 `/usr/local/go`：

```bash
$ sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
```

**macOS (PKG):**
直接双击下载的 `.pkg` 文件，按照向导完成安装。

**Windows (MSI):**
直接双击下载的 `.msi` 文件，按照向导完成安装。

#### 3. 配置环境变量

你需要将 Go 的二进制文件路径添加到系统的 `PATH` 环境变量中，以便在终端中直接运行 `go` 命令。

编辑你的 Shell 配置文件（如 `$HOME/.bashrc`, `$HOME/.zshrc` 或 `/etc/profile`），添加以下行：

```bash
export PATH=$PATH:/usr/local/go/bin
```

保存文件并使配置生效（例如 `source ~/.bashrc`）。

**验证安装：**

在终端运行以下命令，如果能看到版本号，说明安装成功：

```bash
$ go version
go version go1.21.0 linux/amd64
```

#### 4. 关于项目工作区 （Workspace）

在早期的 Go 版本中，我们需要配置 `GOPATH` 环境变量，并将所有项目代码都放在 `$GOPATH/src` 目录下。

**自 Go 1.11 引入 Go Modules 以来，这已不再是必须的。**

现在的最佳实践是使用 **Go Modules** 进行依赖管理。你可以在文件系统的**任意位置**创建 Go 项目。

创建一个新项目只需：

```bash
mkdir my-project
cd my-project
go mod init example.com/my-project
```

这不仅简化了配置，还避免了版本冲突问题。

有关更多安装详情，请参考官方文档：[https://go.dev/doc/install](https://go.dev/doc/install)。
