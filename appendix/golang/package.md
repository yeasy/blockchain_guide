### 依赖管理

#### govendor 工具

长期以来，Go 语言对外部依赖都没有很好的管理方式，只能从 `$GOPATH` 下查找依赖。这就造成不同用户在安装同一个项目时可能从外部获取到不同的依赖库版本，同时当无法联网时，无法编译依赖缺失的项目。

Golang 自 1.5 版本开始重视第三方依赖的管理，将项目依赖的外部包统一放到 vendor 目录下（类比 Nodejs 的 node_modules 目录），并通过 vendor.json 文件来记录依赖包的版本，方便用户使用相对稳定的依赖。

Daniel Theophanes 等人开发了 govendor 工具，方便对第三方依赖进行管理。

govendor 的安装十分简单，可以通过 go get 命令：

```bash
$ go get -u -v github.com/kardianos/govendor
```

对于 govendor 来说，主要存在三种位置的包：项目自身的包组织为本地（local）包；传统的存放在 $GOPATH 下的依赖包为外部（external）依赖包；被 govendor 管理的放在 vendor 目录下的依赖包则为 vendor 包。

具体来看，这些包可能的类型如下：

状态 | 缩写状态 | 含义
--- | ------- | ---
+local | l | 本地包，即项目自身的包组织
+external | e | 外部包，即被 $GOPATH 管理，但不在 vendor 目录下
+vendor | v | 已被 govendor 管理，即在 vendor 目录下
+std | s | 标准库中的包
+unused | u | 未使用的包，即包在 vendor 目录下，但项目并没有用到
+missing | m | 代码引用了依赖包，但该包并没有找到
+program | p | 主程序包，意味着可以编译为执行文件
+outside | | 外部包和缺失的包
+all     | | 所有的包

常见的命令如下，格式为 `govendor COMMAND`。

通过指定包类型，可以过滤仅对指定包进行操作。

命令 | 功能
-- | ---
`init` | 初始化 vendor 目录
`list` | 列出所有的依赖包
`add` | 添加包到 vendor 目录，如 govendor add +external 添加所有外部包
`add PKG_PATH` | 添加指定的依赖包到 vendor 目录
`update` | 从 $GOPATH 更新依赖包到 vendor 目录
`remove` | 从 vendor 管理中删除依赖
`status` | 列出所有缺失、过期和修改过的包
`fetch` | 添加或更新包到本地 vendor 目录
`sync` | 本地存在 vendor.json 时候拉去依赖包，匹配所记录的版本
`get` | 类似 `go get` 目录，拉取依赖包到 vendor 目录

#### dep 工具

为了方便管理依赖，Go 团队 2016 年 4 月开始开发了 dep 工具，试图进一步简化在 Go 项目中对第三方依赖的管理。该工具目前已经被试验性支持，相信很快会成为官方支持的工具。

dep 目前需要 Go 1.7+ 版本，兼容其他依赖管理工具如 glide、godep、vndr、govend、gb、gvt、govendor、glock 等。

类似于 govendor 工具，dep 将依赖都放在本地的 vendor 目录下，通过 Gopkg.toml 和 Gopkg.lock 文件来追踪依赖的状态。

* Gopkg.toml 文件：手动编写或通过 dep init 命令生成。描述了项目对第三库的依赖规则，例如允许的版本范围等。用户可以通过编辑该文件表达预期的依赖控制目标。
* Gopkg.lock 文件：通过 dep init 或 dep ensure 命令自动生成。根据项目代码和 Gopkg.toml 文件，计算出一个符合要求的具体的依赖关系并锁定，其中包括每个第三方库的具体版本。vendor 目录下的依赖库需要匹配这些版本。

安装可以通过 go get 命令：

```bash
$ go get -v -u github.com/golang/dep/cmd/dep
```

dep 使用保持简洁的原则，包括四个子命令。

* init：对一个新的 Go 项目，初始化依赖管理，生成配置文件和 vendor 目录等；
* status：查看当前项目依赖的状态，包括依赖包名称、限制范围、指定版本等。可以通过 -old 参数来只显示过期的依赖；
* ensure：更新依赖，确保满足指定的版本条件。如果本地缺乏某个依赖，会自动安装；
* version：显示 dep 工具的版本信息。

其中，ensure 命令最为常用，支持的子命令参数主要包括：

* -add：添加新的依赖，如 dep ensure -add github.com/pkg/foo@^1.0.0；
* -dry-run：模拟执行，打印参考改动但不实施；
* -no-vendor：根据计算结果更新 Gopkg.lock 文件，但不更新 vendor 中依赖包；
* -update：更新 Gopkg.lock 中的依赖到 Gopkg.toml 中允许的最新版本，默认同时更新 vendor 包中内容；
* -v：输出调试信息方面了解执行过程；
* -vendor-only：按照 Gopkg.lock 中条件更新 vendor 包中内容。

#### go module

Go 自 1.11 版本开始引入模块（module），在 1.13 版本中开始正式支持，以取代传统基于 $GOPATH 的方案。模块作为若干个包（package）的集合，带有语义化版本号，统一管理所有依赖。所有依赖模块缓存在 $GOPATH/pkg 目录下的 `mod` 和 `sum` 子目录中，未来计划迁移到 `$GOCACHE` 目录下。另外，不同项目的相同依赖模块全局只会保存一份，极大节约了存储空间。

模块需要两个配置文件，go.mod 和 go.sum。

前者管理项目中模块的依赖信息，可以通过 go mod init <module name> 命令生成；后者记录当前项目直接或间接依赖的所有模块的路径、版本、校验值等。

项目模块可以通过 go mod 子命令来显式操作，也会在编译、测试等命令中被隐式更新。

go mod 支持的子命令包括：

* download：下载依赖模块到本地的缓存；
* edit：编辑 go.mod 文件；
* graph：查看当前的依赖结构；
* init：初始化，创建 go.mod 文件；
* tidy：整理依赖模块：添加新模块，并删除未使用模块；
* vendor：将依赖模块复制到本地的 vendor 目录，方便兼容原来的 vendor 方式（该命令未来会遗弃）；
* verify：校验当前依赖是否正确，未被篡改；
* why：解释为何需要某个依赖包。

基本使用过程为：

* 使用 `go mod init <module name>` 来初始化本地的 go.mod 文件；
* 使用 `go get -u <package name>@<version>` 来获取某个依赖包（不添加版本号会默认获取当前最新），同时自动更新 go.mod 文件。更新全部模块可以使用 `go get -u all`；
* 编译时使用 `go build -mod=readonly` 可以避免在编译过程中修改 go.mod；
* 如果要使用本地的 vendor 目录进行编译，可以使用 `go build -mod=vendor`；
* 如果要检查可更新的依赖，可以使用 `go list -m -u all`。如果要执行更新，可以使用 `go get -u`；
* 此外，执行 `go` 相关命令（build、get、list、test 等）时，也会自动下载依赖并更新 go.mod 文件。

module 的开启可以通过 GO111MODULE=[auto|on|off] 等环境变量来控制。例如始终使用 go module，可以使用如下命令

```bash
$ go env -w GO111MODULE=on # 记录到 os.UserConfigDir 指定路径（默认为 $HOME/.config/go/）下的 env 文件中
```

此外，1.13 版本起，Go 还支持 GOPROXY 环境变量来指定拉取包的代理服务，GOPRIVATE 指定私有仓库地址。
