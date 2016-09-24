## 应用案例

### 双方交易案例
两方（如 a 和 b）之间进行价值的转移。

集群启动后，进入一个 VP 节点。以 pbft 模式为例，节点名称为 `pbft_vp0_1`。

```sh
$ docker exec -it pbft_vp0_1 bash
```

部署 chaincode example02。

```sh
$ peer chaincode deploy -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02 -c '{"Function":"init", "Args": ["a","100", "b", "200"]}'
03:08:44.740 [chaincodeCmd] chaincodeDeploy -> INFO 001 Deploy result: type:GOLANG chaincodeID:<path:"github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02" name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"init" args:"a" args:"100" args:"b" args:"200" >
Deploy chaincode: ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539
03:08:44.740 [main] main -> INFO 002 Exiting.....
```

返回 chaincode id 为 `ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539`，后面将用这个 id 来标识这次交易。为了方便，把它记录到环境变量 CC_ID 中。

```sh
$ CC_ID="ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539"
```

部署成功后，系统中会自动生成几个 chaincode 容器，例如

```sh
CONTAINER ID        IMAGE                                                                                                                                      COMMAND                  CREATED             STATUS              PORTS                                   NAMES
e86c26bad76f        dev-vp1-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539   "/opt/gopath/bin/ee5b"   2 minutes ago       Up 2 minutes                                                dev-vp1-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539
597ebaf929a0        dev-vp2-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539   "/opt/gopath/bin/ee5b"   2 minutes ago       Up 2 minutes                                                dev-vp2-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539
8748a3b47312        dev-vp3-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539   "/opt/gopath/bin/ee5b"   2 minutes ago       Up 2 minutes                                                dev-vp3-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539
cf6e762f6a2e        dev-vp0-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539   "/opt/gopath/bin/ee5b"   2 minutes ago       Up 2 minutes                                                dev-vp0-ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539
```

查询 a 手头的价值，为初始值 100。

```sh
$ peer chaincode query -n ${CC_ID} -c '{"Function": "query", "Args": ["a"]}'
03:22:31.420 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001 Successfully queried transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"query" args:"a" > >
Query Result: 100
03:22:31.420 [main] main -> INFO 002 Exiting.....
```

a 向 b 转账 10 元。

```sh
$ peer chaincode invoke -n ${CC_ID} -c '{"Function": "invoke", "Args": ["a", "b", "10"]}'
03:22:57.345 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001 Successfully invoked transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"invoke" args:"a" args:"b" args:"10" > > (fc298ffb-c763-4ed0-9da2-072de2ab20b1)
03:22:57.345 [main] main -> INFO 002 Exiting.....
```

查询 a 手头的价值，为新的值 90。

```sh
$ peer chaincode query -n ${CC_ID} -c '{"Function": "query", "Args": ["a"]}'
03:23:33.045 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001 Successfully queried transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"query" args:"a" > >
Query Result: 90
03:23:33.045 [main] main -> INFO 002 Exiting.....
···
