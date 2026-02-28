# 配置容器环境

Hyperledger Fabric 深度依赖容器技术。为了运行 Fabric 网络，你需要准备好 Docker 和 Docker Compose 环境。

## 1. 安装 Docker

### Linux

我们推荐使用 Docker 官方提供的安装脚本：

```bash
curl -fsSL https://get.docker.com | sh
```

安装完成后，建议将当前用户加入 `docker` 用户组，以免每次运行命令都需要 `sudo`：

```bash
sudo usermod -aG docker $USER
# 注销并重新登录后生效
```

### macOS / Windows

直接下载并安装 **Docker Desktop**：https://www.docker.com/products/docker-desktop

## 2. 安装 Docker Compose

Docker Compose 用于定义和运行多容器应用。

*   **Docker Desktop (Mac/Windows)**: 已经内置了 Docker Compose，无需额外安装。
*   **Linux**: 如果你通过官方脚本安装了最新的 Docker Engine，现在它通常包含 `docker compose` 插件（注意命令中间没有连字符）。

验证安装：

```bash
docker compose version
# 或者旧版本
docker-compose --version
```

## 3. Fabric 核心镜像

Fabric 的网络节点大多运行在 Docker 容器中。主要涉及以下核心镜像：

| 镜像名称 | 描述 |
| :--- | :--- |
| `hyperledger/fabric-peer` | Peer 节点，负责维护账本和执行链码。 |
| `hyperledger/fabric-orderer` | Orderer 节点，负责交易排序和打包区块（基于 Raft 共识）。 |
| `hyperledger/fabric-tools` | 包含 `peer`, `osnadmin`, `configtxgen` 等工具的命令行环境，常用于测试。 |
| `hyperledger/fabric-ccenv` | Go/Java 链码的编译环境。 |
| `hyperledger/fabric-baseos` | 链码运行时的基础操作系统环境。 |
| `hyperledger/fabric-ca` | 证书授权服务（CA），用于身份管理。 |
| `couchdb` | (可选) 用作 Peer 节点的状态数据库，支持富查询。 |

**状态数据库选择建议**：
*   **LevelDB (首选)**：Peer 节点的默认状态数据库（内嵌在 `fabric-peer` 镜像中）。它性能更高、运维更简单。除非你有明确的富查询需求，否则应优先选择 LevelDB。
*   **CouchDB**：支持基于 JSON 数据的复杂查询（Rich Queries）。如果你的应用需要对链上数据进行多维度的检索和统计，才建议配置 CouchDB。

**注意**：早期版本的 Kafka 和 Zookeeper 镜像已在 Fabric 2.x 中被弃用，因为 Raft (etcdraft) 已成为默认的排序服务共识机制。

## 4. 获取镜像

如果你使用了上一节提到的 `install-fabric.sh` 脚本，镜像应该已经自动下载完毕。

你也可以手动拉取指定版本的镜像：

```bash
# 设置版本变量
export FABRIC_VERSION=2.5.9
export CA_VERSION=1.5.9

# 拉取镜像
docker pull hyperledger/fabric-peer:$FABRIC_VERSION
docker pull hyperledger/fabric-orderer:$FABRIC_VERSION
docker pull hyperledger/fabric-tools:$FABRIC_VERSION
docker pull hyperledger/fabric-ccenv:$FABRIC_VERSION
docker pull hyperledger/fabric-ca:$CA_VERSION
```

## 5. 镜像标签 （Tag）

Fabric 镜像通常提供多种标签：
*   `2.5.9`: 具体版本号（推荐生产环境使用）。
*   `latest`: 指向最新的发布版本（不建议在生产环境使用，因为可能自动升级导致不兼容）。
*   `amd64`, `arm64`: 针对不同 CPU 架构的镜像。

确保你的 `docker-compose.yaml` 文件中引用的镜像版本与你下载的版本一致，以避免启动失败。
