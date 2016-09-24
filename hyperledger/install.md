## 安装部署

社区在很长一段时间内并没有推出比较容易上手的安装部署方案，于是笔者设计了基于 Docker 容器的一键式部署方案，该方案推出后在社区受到了不少人的关注和应用。官方在安装部署方面已有了一些改善，具体可以参考代码 doc 目录下内容，但仍然存在一些问题。

如果你是初次接触 hyperledger fabric 项目，推荐采用如下的步骤，基于 Docker-compose 的一键部署。

*动手前，建议适当了解一些 [Docker 相关知识](https://github.com/yeasy/docker_practice)。*

### 安装 Docker

Docker 支持 Linux 常见的发行版，如 Redhat/Centos/Ubuntu 等。

```sh
$ curl -fsSL https://get.docker.com/ | sh
```

安装成功后，停止默认启动的 Docker 服务。

```sh
$ sudo service docker stop
```

用如下命令手动启动 Docker 服务。

```sh
$ sudo docker daemon --api-cors-header="*" -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock
```

### 安装 docker-compose

首先，安装 python-pip 软件包。

```sh
$ sudo aptitude install python-pip
```

安装 docker-compose。

```sh
$ sudo pip install docker-compose
```

### 下载镜像

下载相关镜像，并进行配置。

```sh
$ docker pull yeasy/hyperledger-fabric:latest
$ docker tag yeasy/hyperledger-fabric:latest hyperledger/fabric-peer:latest
$ docker tag yeasy/hyperledger-fabric:latest hyperledger/fabric-baseimage:latest
$ docker tag yeasy/hyperledger-fabric:latest hyperledger/fabric-membersrvc:latest
```

*注：如果采用某个稳定分支的代码，需要下载对应镜像。例如 0.6-dp 分支，则需要修改下载镜像的 latest 标签为 0.6-dp。*

之后，用户可以选择采用不同的一致性机制，包括 noops、pbft 两类。

### 使用 noops 模式
noops 默认没有采用 consensus 机制，1 个节点即可，可以用来进行快速测试。

```sh
$ docker run --name=vp0 \
                    --restart=unless-stopped \
                    -it \
                    -p 7050:7050 \
                    -p 7051:7051 \
                    -v /var/run/docker.sock:/var/run/docker.sock \
                    -e CORE_PEER_ID=vp0 \
                    -e CORE_PEER_ADDRESSAUTODETECT=true \
                    -e CORE_NOOPS_BLOCK_WAIT=10 \
                    yeasy/hyperledger-peer:latest peer node start
```

### 使用 PBFT 模式

PBFT 是经典的分布式一致性算法，也是 hyperledger 目前最推荐的算法，该算法至少需要 4 个节点。

首先，下载 compose 文件。

```sh
$ git clone https://github.com/yeasy/docker-compose-files
```

进入 hyperledger 项目，并启动集群。

```sh
$ cd docker-compose-files/hyperledger
$ docker-compose up
```

### 多物理节点部署

上述方案的典型场景是单物理节点上部署多个 Peer 节点。如果要扩展到多物理节点，需要容器云平台的支持，如 Swarm 等。

当然，用户也可以分别在各个物理节点上通过手动启动容器的方案来实现跨主机组网，每个物理节点作为一个 peer 节点。

首先，以 4 节点下的 PBFT 模式为例，配置 4 台互相连通的物理机，分别按照上述步骤配置 Docker，下载镜像。

4 台物理机分别命名为 vp0 ~ vp3。

#### vp0

vp0 作为初始的探测节点。

```sh
docker run --name=vp0 \
    --net="host" \
    --restart=unless-stopped \
    -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e CORE_PEER_ID=vp0 \
    -e CORE_PBFT_GENERAL_N=4 \
    -e CORE_LOGGING_LEVEL=debug \
    -e CORE_PEER_ADDRESSAUTODETECT=true \
    -e CORE_PEER_NETWORKID=dev \
    -e CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft \
    -e CORE_PBFT_GENERAL_MODE=batch \
    -e CORE_PBFT_GENERAL_TIMEOUT_REQUEST=10s \
    yeasy/hyperledger-peer:latest peer node start
```

#### vp1 ~ vp3

以 vp1 为例，假如 vp0 的地址为 10.0.0.1。

```sh
NAME=vp1 \
ROOT_NODE=10.0.0.1 \
docker run --name=${NAME} \
    --net="host" \
    --restart=unless-stopped \
    -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -e CORE_PEER_ID=${NAME} \
    -e CORE_PBFT_GENERAL_N=4 \
    -e CORE_LOGGING_LEVEL=debug \
    -e CORE_PEER_ADDRESSAUTODETECT=true \
    -e CORE_PEER_NETWORKID=dev \
    -e CORE_PEER_VALIDATOR_CONSENSUS_PLUGIN=pbft \
    -e CORE_PBFT_GENERAL_MODE=batch \
    -e CORE_PBFT_GENERAL_TIMEOUT_REQUEST=10s \
    -e CORE_PEER_DISCOVERY_ROOTNODE=${ROOT_NODE}:7051 \
    yeasy/hyperledger-peer:latest peer node start
```

### 服务端口
Hyperledger 默认监听的服务端口包括：

* 7050: REST 服务端口，推荐 NVP 节点开放，旧版本中为 5000；
* 7051：peer gRPC 服务监听端口，旧版本中为 30303；
* 7052：peer CLI 端口，旧版本中为 30304；
* 7053：peer 事件服务端口，旧版本中为 31315；
* 7054：eCAP
* 7055：eCAA
* 7056：tCAP
* 7057：tCAA
* 7058：tlsCAP
* 7059：tlsCAA


