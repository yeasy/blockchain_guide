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
08:23:25.028 [chaincodeCmd] getChaincodeSpecification -> INFO 001 Local user 'jim' is already logged in. Retrieving login token.
08:23:26.162 [chaincodeCmd] chaincodeDeploy -> INFO 002 Deploy result: type:GOLANG chaincodeID:<path:"github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02" name:"980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d" > ctorMsg:<args:"init" args:"a" args:"100" args:"b" args:"200" >
08:23:26.163 [main] main -> INFO 003 Exiting.....
```

部署成功后，系统中会自动生成几个 chaincode 容器，例如

```sh
CONTAINER ID        IMAGE                                                                                                                                      COMMAND                  CREATED             STATUS              PORTS                                   NAMES
07d5c07bbab3        dev-vp1-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d   "/opt/gopath/bin/980d"   15 minutes ago      Up 15 minutes                                               dev-vp1-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d
52b1be5a7bda        dev-vp0-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d   "/opt/gopath/bin/980d"   15 minutes ago      Up 15 minutes                                               dev-vp0-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d
480df639c212        dev-vp2-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d   "/opt/gopath/bin/980d"   15 minutes ago      Up 15 minutes                                               dev-vp2-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d
14ecdae1adbf        dev-vp3-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d   "/opt/gopath/bin/980d"   15 minutes ago      Up 15 minutes                                               dev-vp3-980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d
```

查询 a 手头的价值，为初始值 100。

```sh
$ peer chaincode query -n 980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d -c '{"Function": "query", "Args": ["a"]}'
08:37:12.415 [chaincodeCmd] getChaincodeSpecification -> INFO 001 Local user 'jim' is already logged in. Retrieving login token.
08:37:12.516 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 002 Successfully queried transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d" > ctorMsg:<args:"query" args:"a" > secureContext:"jim" >
Query Result: 100
08:37:12.516 [main] main -> INFO 003 Exiting.....
```

a 向 b 转账 10 元。

```sh
$ peer chaincode invoke -n 980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d -c '{"Function": "invoke", "Args": ["a", "b", "10"]}'
08:37:51.211 [chaincodeCmd] getChaincodeSpecification -> INFO 001 Local user 'jim' is already logged in. Retrieving login token.
08:37:51.309 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 002 Successfully invoked transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d" > ctorMsg:<args:"invoke" args:"a" args:"b" args:"10" > secureContext:"jim" > (66308740-a2c5-4a60-81f1-778dbed49cc3)
08:37:51.309 [main] main -> INFO 003 Exiting.....
```

查询 a 手头的价值，为新的值 90。

```sh
$ peer chaincode query -n 980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d -c '{"Function": "query", "Args": ["a"]}'
08:55:12.961 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001 Successfully queried transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d" > ctorMsg:<args:"query" args:"a" > >
Query Result: 90
08:55:12.962 [main] main -> INFO 002 Exiting.....
···