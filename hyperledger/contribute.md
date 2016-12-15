## 开发和提交代码

### 安装环境

推荐在 Ubuntu 14.04+ 环境中开发代码，并安装如下工具。

* git：用来获取代码；
* golang 1.6+：安装成功后配置 $GOPATH 等环境变量。

### 获取代码

首先注册 Linux foundation ID，并登陆 [https://gerrit.hyperledger.org/](https://gerrit.hyperledger.org/)，添加个人 ssh pub key。

查看项目列表，找到对应项目，以 fabric 为例，获取 `Clone with commit-msg hook` 的方式。

典型的，执行如下命令获取代码，其中 `LF_ID` 替换为你的 Linux foundation id。

```sh
$ git clone ssh://LF_ID@gerrit.hyperledger.org:29418/fabric && scp -p -P 29418 LF_ID@gerrit.hyperledger.org:hooks/commit-msg fabric/.git/hooks/
```

如果没有添加个人 ssh pubkey，则可以通过 https 方式 clone，需要输入用户名和密码信息。

```sh
git clone http://LF_ID@gerrit.hyperledger.org/r/fabric && (cd fabric && curl -kLo `git rev-parse --git-dir`/hooks/commit-msg http://LF_ID@gerrit.hyperledger.org/r/tools/hooks/commit-msg; chmod +x `git rev-parse --git-dir`/hooks/commit-msg)
```

clone 下代码后，为了方便后面的编译测试，需要放到 `$GOPATH/src/github.com/hyperledger/` 路径下。

```sh
$ mkdir 
```

### 编译和测试

大部分编译和安装过程都可以通过 Makefile 来执行。

#### 安装 go tools
执行 

```sh
$ make gotools
```

#### 语法格式检查

执行

```sh
$ make linter
```

#### 编译 peer

执行 

```sh
$ make peer
```

*注意：有时候会因为获取安装包不稳定而报错，需要执行 `make clean`，然后再次执行。*

#### 生成 Docker 镜像
执行 

```sh
$ make images
```

#### 执行单元测试

执行 

```sh
$ make unit-test
```

如果要运行某个特定单元测试，则可以通过类似如下格式。

```sh
$ go test -v -run=TestGetFoo
```

#### 执行 BDD 测试
需先生成本地 Docker 镜像。

执行 

```sh
$ make behave
```

### 提交代码

仍然使用 Linux foundation ID 登录 [jira.hyperledger.org](http://jira.hyperledger.org)，查看有没有未分配的任务，如果对某个任务感兴趣，可以添加自己为 assignee，如对 FAB-XXX 任务。

本地创建新的分支 FAB-XXX。

```sh
$ git checkout -b FAB-XXX
```

实现任务代码，完成后，执行语法格式检查和测试等，确保所有检查和测试都通过。

提交代码到本地仓库。

```sh
$ git commit -a -s
```

会打开一个窗口需要填写 commit 信息，格式一般要求为：

```
Simple words to describe main change

This fixes #FAB-XXX.

A more detailed description can be here, with several
paragraphs and sentences...
```

之后使用 git review 命令推送到远端仓库。

```sh
$ git review
```

提交成功后，可以打开 [gerrit.hyperledger.org/r/](https://gerrit.hyperledger.org/r/)，查看自己最新提交的 patchset 信息，添加几位 reviewer。之后就是等待开发者团队的 review 结果，如果得到通过，则会被项目的 maintainer 们 merge 到主分支。否则还需要针对大家提出的建议进一步的修正。

修正过程跟提交代码过程类似，唯一不同是提交的时候使用

```sh
$ git commit -a --amend
```

表示这个提交是对旧提交的一次修订。