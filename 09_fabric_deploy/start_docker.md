## 容器方式启动 Fabric 网络

除了上面讲解的手动部署方式，读者还可以使用官方 `fabric-samples` 中的 test network 快速验证 Fabric 网络。该网络基于 Docker Compose，适合学习、开发和测试，不是生产部署模板。

完成 `install-fabric.sh` 安装后，进入 test network 目录：

```bash
$ cd $HOME/hyperledger/fabric-samples/test-network
$ ./network.sh down
```

启动网络并创建通道：

```bash
$ ./network.sh up createChannel -ca -c businesschannel
```

其中：

* `-ca` 表示使用 Fabric CA 生成测试网络的身份材料；
* `-c businesschannel` 指定要创建并加入的通道名称；
* 如需使用 CouchDB 状态数据库，可增加 `-s couchdb`；
* 如需在 Fabric 3.x 中体验 BFT 排序服务，可使用 `-bft`；该选项不适用于 Fabric 2.x。

test network 默认配置较简化：两个 Peer 组织、一个 Orderer 组织，并以本机 Docker Compose 网络隔离运行。生产网络需要独立设计 CA、MSP、TLS、排序服务法定人数、持久化存储、operations endpoint、监控和证书轮换等配置。

网络启动后，可以通过如下命令查看容器：

```bash
$ docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
```

如需单独创建通道，可先启动节点，再执行：

```bash
$ ./network.sh up -ca
$ ./network.sh createChannel -c businesschannel
```

如需部署官方示例链码，可执行：

```bash
$ ./network.sh deployCC \
    -ccn basic \
    -ccp ../asset-transfer-basic/chaincode-go \
    -ccl go \
    -c businesschannel
```

完成测试后清理网络和生成物：

```bash
$ ./network.sh down
```
