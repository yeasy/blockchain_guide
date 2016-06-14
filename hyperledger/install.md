## 安装部署

社区在很长一段时间内并没有推出比较容易上手的安装部署方案，于是笔者设计了基于 Docker 容器的一键式部署方案，该方案推出后在社区受到了不少人的关注和应用。官方在安装部署方面已有了一些改善，但仍然存在一些问题。

如果你是初次接触 hyperledger fabric 项目，推荐采用如下的步骤。

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
$ docker pull yeasy/hyperledger:latest
$ docker tag yeasy/hyperledger:latest hyperledger/fabric-baseimage:latest
$ docker pull yeasy/hyperledger-peer:latest
$ docker pull yeasy/hyperledger-peer:noops
$ docker pull yeasy/hyperledger-peer:pbft
```

### 采用 noops 测试
noops 默认没有采用 consensus 机制，1 个节点即可，可以用来进行快速测试。

```sh
$ docker run --name=vp0 \
                    --restart=unless-stopped \
                    -it \
                    -p 5000:5000 \
                    -p 30303:30303 \
                    -v /var/run/docker.sock:/var/run/docker.sock \
                    -e CORE_PEER_ID=vp0 \
                    -e CORE_PEER_ADDRESSAUTODETECT=true \
                    -e CORE_NOOPS_BLOCK_TIMEOUT=10 \
                    yeasy/hyperledger-peer:noops peer node start
```

### 启动 PBFT 集群

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

### 服务端口
Hyperledger 默认监听的服务端口包括：

* 5000: REST 服务端口，推荐 NVP 节点开放；
* 30303: Peer 服务监听端口；
* 30304: CLI 从 chaincode 的回调端口；
* 31315: VP 上的事件服务端口。
