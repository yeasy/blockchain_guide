##membersrcv应用案例

###下载相关镜像
```sh
$ docker pull yeasy/hyperledger:latest
$ docker tag yeasy/hyperledger:latest hyperledger/fabric-baseimage:latest
$ docker pull yeasy/hyperledger-peer:latest
$ docker pull yeasy/hyperledger-peer:noops
$ docker pull yeasy/hyperledger-peer:pbft
$ docker pull yeasy/hyperledger-membersrvc:latest
```

进入hyperledger项目，启动带成员管理的pbft集群
```sh
$ git clone https://github.com/yeasy/docker-compose-files
$ cd docker-compose-files/hyperledger
$ docker-compose up
```

###用户注册



###chaincode deploy


###chaincode invoke

###chaincode query
