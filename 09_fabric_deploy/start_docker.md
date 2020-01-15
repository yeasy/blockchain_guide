## 容器方式启动 Fabric 网络

除了上面讲解的手动部署的方式，读者还可以基于容器方式来快速部署 Fabric 网络并验证功能。

首先，按照如下命令下载 Docker-Compose 模板文件，并进入 `hyperledger_fabric` 目录，可以看到有对应多个 Fabric 版本的项目，用户可以根据需求选用特定版本：

```sh
$ git clone https://github.com/yeasy/docker-compose-files
$ cd docker-compose-files/hyperledger_fabric
```

以 Fabric 2.0.0 版本为例，进入到对应目录下，并先下载所需镜像文件：

```bash
$ cd v2.0.0
$ make download
```

查看目录下内容，主要包括若干 Docker-Compose 模板文件，主要包括：

* docker-compose-2orgs-4peer-raft.yaml：包括 4 个 peer 节点（属于两个组织）、3 个 Orderer 节点（Raft 模式）、2 个 CA 节点、1 个客户端节点；
* docker-compose-1orgs-1peers-dev.yaml：包括 1 个 peer 节点、1 个 Orderer 节点、1 个 CA 节点、1 个客户端节点。本地 Fabric 源码被挂载到了客户端节点中，方便进行调试；
* docker-compose-2orgs-4peer-kafka.yaml：包括 4 个 peer 节点（属于两个组织）、3 个 Orderer 节点（Kafka 模式）、2 个 CA 节点、1 个客户端节点；
* docker-compose-2orgs-4peer-couchdb.yaml：包括 4 个 peer 节点（属于两个组织，启用 couchDB 作为状态数据库）、2 个 Orderer 节点、1 个 CA 节点、1 个客户端节点。

使用 Make 命令进行操作。例如使用 HLF_MODE 指定排序服务为 Raft 模式，快速启动网络并执行一系列测试：

```bash
$ HLF_MODE=raft make test
```

`make test` 实际上自动执行了一系列指令：

* make gen_config_crypto：生成网络需要的身份文件；
* make gen_config_channel：生成网络需要的配置文件；
* make start：启动网络；
* make channel_test：执行通道创建和加入通道；
* make update_anchors：更新锚节点信息；
* make cc_test：执行链码相关测试，包括安装、实例化和调用；
* make test_lscc：测试系统链码 LSCC 调用（使用 2.0 中新的链码生命周期则不支持）；
* make test_qscc：测试系统链码 QSCC 调用；
* make test_cscc：测试系统链码 CSCC 调用；
* make test_fetch_blocks：获取通道内的区块；
* make test_config_update：生成新版本配置；
* make test_channel_update：测试更新通道配置；
* make test_configtxlator：测试 configtxlator 转换；
* make test_channel_list：测试列出 Peer 加入的通道；
* make test_channel_getinfo：测试获取通道信息；
* make stop：停止网络。

运行过程中会自动创建网络并逐个完成通道和链码的相关测试，注意查看输出日志中无错误信息。

网络启动后，可以通过 `docker ps` 命令查看本地系统中运行的容器信息：

```bash
$ docker ps
CONTAINER ID        IMAGE                                     COMMAND                  CREATED             STATUS              PORTS                               NAMES
1ee7db027b3f        yeasy/hyperledger-fabric-peer:2.0.0      "peer node start"        27 seconds ago      Up 22 seconds       9443/tcp, 0.0.0.0:8051->7051/tcp    peer1.org1.example.com
8f7bffcd14b3        yeasy/hyperledger-fabric-peer:2.0.0      "peer node start"        27 seconds ago      Up 22 seconds       9443/tcp, 0.0.0.0:10051->7051/tcp   peer1.org2.example.com
8a4e9aaec7ba        yeasy/hyperledger-fabric-peer:2.0.0      "peer node start"        27 seconds ago      Up 22 seconds       9443/tcp, 0.0.0.0:9051->7051/tcp    peer0.org2.example.com
7b9d394f26c0        yeasy/hyperledger-fabric-peer:2.0.0      "peer node start"        27 seconds ago      Up 23 seconds       0.0.0.0:7051->7051/tcp, 9443/tcp    peer0.org1.example.com
ce9ca6c7b672        yeasy/hyperledger-fabric-orderer:2.0.0   "orderer start"          30 seconds ago      Up 27 seconds       8443/tcp, 0.0.0.0:8050->7050/tcp    orderer1.example.com
2646b7f0e462        yeasy/hyperledger-fabric:2.0.0           "bash -c 'cd /tmp; s…"   30 seconds ago      Up 15 seconds       7050-7054/tcp                       fabric-cli
c35e8694c634        yeasy/hyperledger-fabric-orderer:2.0.0   "orderer start"          30 seconds ago      Up 27 seconds       8443/tcp, 0.0.0.0:9050->7050/tcp    orderer2.example.com
1d6dd5009141        yeasy/hyperledger-fabric-orderer:2.0.0   "orderer start"          30 seconds ago      Up 27 seconds       0.0.0.0:7050->7050/tcp, 8443/tcp    orderer0.example.com
```

用户如果希望在客户端、Peer 或 Orderer 容器内执行命令，可以通过 `make cli|peer|orderer` 命令进入到容器中。

例如，如下命令可以让用户登录到客户端节点，在其中以指定身份发送网络请求：

```bash
$ make cli
```

用户也可以通过如下命令来查看日志输出：

```bash
$ make logs
```

更多操作命令用户可以参考 Makefile 内容，在此不再赘述。

