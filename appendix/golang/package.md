### 依赖管理：Go Modules

在 Go 1.11 之前，Go 语言缺乏官方统一的依赖管理机制，社区涌现了如 `dep`, `govendor`, `glide` 等第三方工具。自 Go 1.11 引入并在 Go 1.13 默认开启 **Go Modules** 以来，它已成为 Go 语言管理依赖的事实标准。

#### 什么是 Go Modules

Module（模块）是相关 Go 包（Package）的集合。Modules 是源代码交换和版本控制的单元。`go` 命令直接支持使用 Modules，包括记录和解析对其他模块的依赖性。

一个 Module 由根目录下的 `go.mod` 文件定义。

#### 核心命令

Go Modules 的日常使用非常简单，主要依赖 `go` 命令行工具。

**1. 初始化模块**

在项目根目录下运行：

```bash
$ go mod init <module-name>
# 例如：go mod init github.com/myuser/myproject
```

这会创建一个 `go.mod` 文件，其中包含模块名称和 Go 版本。

**2. 添加/更新依赖**

当你在代码中 `import` 一个包并运行 `go build` 或 `go test` 时，Go 会自动查找并下载该依赖的最新稳定版本，并将其添加到 `go.mod` 和 `go.sum` 文件中。

你也可以手动添加或升级依赖：

```bash
$ go get example.com/pkg
```

这会下载最新版本。如果需要特定版本：

```bash
$ go get example.com/pkg@v1.2.3
```

**3. 整理依赖**

在开发过程中，`go.mod` 可能会包含一些不再使用的依赖（例如删除了某些 import）。使用 `tidy` 命令可以自动清理：

```bash
$ go mod tidy
```

这个命令会添加缺失的模块并移除未使用的模块，确保 `go.mod` 与源代码不仅一致。

#### 配置文件

*   **go.mod**: 定义模块路径、Go 版本以及项目直接依赖的模块版本。
*   **go.sum**: 记录所有依赖模块（包括间接依赖）的加密哈希值，用于确保下载的模块未被篡改，保证构建的一致性和安全性。这两个文件都应提交到版本控制系统（Git）中。

#### 代理配置 (GOPROXY)

在中国大陆地区访问 Go 官方源可能较慢。可以通过设置 `GOPROXY` 环境变量来使用国内的代理服务加速下载。

**七牛云代理：**

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

**阿里云代理：**

```bash
go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

设置后，`go get` 将通过代理服务器高速下载依赖。

### 总结

现在，你在进行 Go 语言开发时，**不再需要**关注 `vendor` 目录（除非有特殊离线编译需求），也**不再需要**学习 `dep` 或 `glide` 等旧工具。只需掌握 `go mod` 系列命令即可轻松管理项目依赖。
