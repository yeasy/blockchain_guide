## 应用案例

### 两人交易案例
两个人 a 和 b 之间进行价值的转移。

集群启动后，进入一个 VP 节点。

```sh
$ docker exec -it vp0 bash
```

部署 chaincode example02。

```sh
$ peer chaincode deploy -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02 -c '{"Function":"init", "Args": ["a","100", "b", "200"]}'
13:16:35.643 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
5844bc142dcc9e788785e026e22c855957b2c754c912702c58d997dedbc9a042f05d152f6db0fbd7810d95c1b880c210566c9de3093aae0ab76ad2d90e9cfaa5
```

查询 a 手头的价值，为初始值 100。

```sh
$ peer chaincode query -n 5844bc142dcc9e788785e026e22c855957b2c754c912702c58d997dedbc9a042f05d152f6db0fbd7810d95c1b880c210566c9de3093aae0ab76ad2d90e9cfaa5 -c '{"Function": "query", "Args": ["a"]}'
13:20:07.952 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
100
```

a 向 b 转账 10 元。

```sh
$ peer chaincode invoke -n 5844bc142dcc9e788785e026e22c855957b2c754c912702c58d997dedbc9a042f05d152f6db0fbd7810d95c1b880c210566c9de3093aae0ab76ad2d90e9cfaa5 -c '{"Function": "invoke", "Args": ["a", "b", "10"]}'
13:20:31.028 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
ec3c675b-a2fe-4429-ab44-7f389e454657
```

查询 a 手头的价值，为新的值 90。

```sh
$ peer chaincode query -n 5844bc142dcc9e788785e026e22c855957b2c754c912702c58d997dedbc9a042f05d152f6db0fbd7810d95c1b880c210566c9de3093aae0ab76ad2d90e9cfaa5 -c '{"Function": "query", "Args": ["a"]}'
13:20:35.725 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
90
···