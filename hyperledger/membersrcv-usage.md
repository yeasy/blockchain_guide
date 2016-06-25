## 权限管理

权限管理机制是 hyperledger fabric 项目的一大特色。下面给出使用权限管理的一个应用案例。

### 下载相关镜像
首先启动相关的环境。

```sh
$ docker pull yeasy/hyperledger:latest
$ docker tag yeasy/hyperledger:latest hyperledger/fabric-baseimage:latest
$ docker pull yeasy/hyperledger-peer:latest
$ docker pull yeasy/hyperledger-membersrvc:latest
```

进入 hyperledger 项目，启动带成员管理的 PBFT 集群。

```sh
$ git clone https://github.com/yeasy/docker-compose-files
$ cd docker-compose-files/hyperledger
$ docker-compose -f docker-compose-with-membersrvc.yml up
```

### 用户登陆
以 jim 账户登录，URL：
```sh
POST  HOST:5000/registrar
```
Request：
```
{
  "enrollId": "jim",
  "enrollSecret": "6avZQLwcUe9b"
}
```

Response：
```
{
    "OK": "User jim is already logged in."
}
```
### chaincode 部署
将 https://github.com//hyperledger/fabric/examples/chaincode/go/chaincode_example02 的chaincode部署到 PBFT 集群上，并初始化 a、b 两个账户。

URL：

```sh
POST  HOST:5000/chaincode
```

Request：

```json
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID":{
        "path":"github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"
    },
    "ctorMsg": {
        "function":"init",
        "args":["a", "1000", "b", "2000"]
    },
    "secureContext": "jim"
  },
  "id": 1
}
```

Response：

```json
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message": "28bb2b2316171a706bb2810ec35d095f430877bf443f1061ef0f60bbe753ed440700a5312c16390d3b30199fe9465c3b75d5944358caae01ca81ef28128a1bfb"
    },
    "id": 1
}
```

### chaincode 调用
在账户 a，b 间进行转账，URL：

```sh
POST  HOST:5000/chaincode
```

Request：

```json
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"28bb2b2316171a706bb2810ec35d095f430877bf443f1061ef0f60bbe753ed440700a5312c16390d3b30199fe9465c3b75d5944358caae01ca81ef28128a1bfb"
      },
      "ctorMsg": {
         "function":"invoke",
         "args":["a", "b", "100"]
      },
    "secureContext": "jim"
  },
  "id": 3
}
```

Response：

```json
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message": "2b3b6cf3-9887-4dd5-8f2e-3634ec9c719a"
    },
    "id": 3
}
```
### chaincode 查询

查询 a 账户的余额 URL：

```sh
POST  HOST:5000/chaincode
```

Request：

```json
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"28bb2b2316171a706bb2810ec35d095f430877bf443f1061ef0f60bbe753ed440700a5312c16390d3b30199fe9465c3b75d5944358caae01ca81ef28128a1bfb"
      },
      "ctorMsg": {
         "function":"query",
         "args":["a"]
      },
      "secureContext": "jim"
  },
  "id": 5
}
```

Response：

```json
{
    "jsonrpc": "2.0",
    "result": {
        "status": "OK",
        "message": "900"
    },
    "id": 5
}
```


### 区块信息查询
URL：

```sh
GET  HOST:5000/chain/blocks/2
```

Response：

```json
{
    "transactions": [
        {
            "type": 2,
            "chaincodeID": "EoABMjhiYjJiMjMxNjE3MWE3MDZiYjI4MTBlYzM1ZDA5NWY0MzA4NzdiZjQ0M2YxMDYxZWYwZjYwYmJlNzUzZWQ0NDA3MDBhNTMxMmMxNjM5MGQzYjMwMTk5ZmU5NDY1YzNiNzVkNTk0NDM1OGNhYWUwMWNhODFlZjI4MTI4YTFiZmI=",
            "payload": "Cp0BCAESgwESgAEyOGJiMmIyMzE2MTcxYTcwNmJiMjgxMGVjMzVkMDk1ZjQzMDg3N2JmNDQzZjEwNjFlZjBmNjBiYmU3NTNlZDQ0MDcwMGE1MzEyYzE2MzkwZDNiMzAxOTlmZTk0NjVjM2I3NWQ1OTQ0MzU4Y2FhZTAxY2E4MWVmMjgxMjhhMWJmYhoTCgZpbnZva2USAWESAWISAzEwMA==",
            "uuid": "2b3b6cf3-9887-4dd5-8f2e-3634ec9c719a",
            "timestamp": {
                "seconds": 1466577447,
                "nanos": 399637431
            },
            "nonce": "5AeA6S1odhPIDiGjFTFG8ttcihOoNNsh",
            "cert": "MIICPzCCAeSgAwIBAgIRAMndnS+Me0G6gs4J9/fb8HcwCgYIKoZIzj0EAwMwMTELMAkGA1UEBhMCVVMxFDASBgNVBAoTC0h5cGVybGVkZ2VyMQwwCgYDVQQDEwN0Y2EwHhcNMTYwNjIyMDYzMzE4WhcNMTYwOTIwMDYzMzE4WjAxMQswCQYDVQQGEwJVUzEUMBIGA1UEChMLSHlwZXJsZWRnZXIxDDAKBgNVBAMTA2ppbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABDLd2W8PxzgB4A85Re2x44BApbOGqP05tnkygbXSctLiqi5HVfwRAACS6znVA9+toni59Yy+XAH3w2offdjFW3mjgdwwgdkwDgYDVR0PAQH/BAQDAgeAMAwGA1UdEwEB/wQCMAAwDQYDVR0OBAYEBAECAwQwDwYDVR0jBAgwBoAEAQIDBDBNBgYqAwQFBgcBAf8EQAfASTE6bZ0P5mrEzTa5r1UyKFv+dKezBiGU0V3l2iWzk9evlGMvaC2pwhEKfKDdKxs7YSMYe/7cLq/oF++GBVowSgYGKgMEBQYIBEBEO3TKXuORl5Geuco8Gnn5TkoIl4+b96aPGDGvKbmDjMXR9vEBuUXTnsbDL53j7kC8/XQs1kZboC1ojLeUSN03MAoGCCqGSM49BAMDA0kAMEYCIQCZqyANMFcu1WiMe2So0pC7eRU95F0+qUXLAKZsPWv/YQIhALmNaglP7CoMOe2qxehucmffDlu0BRLSYDHyV9xcxmkH",
            "signature": "MEYCIQDob3NqdrfwlSGhi+zz+Ypl7S9QQ07RIFr8nV92e8KDNgIhANIljz4tRS8vwQk01hTemNQFJX2zMI6DhSUFZivbbtoR"
        }
    ],
    "stateHash": "7YUoVvYnMLHbLf47uTixLtkjF6xM9DuvgSWC92MbOUzk09xhcRBBLZqe5FvJElgZemELBOcuIFnubL0LiGH0yw==",
    "previousBlockHash": "On4BlpqCYNpugUKluqvOcbvkr3TAQxmlISLdd6qrONtIgmQ4iUDeWxAA9lUCceZfF8tke8A0Wy7m9tksNpKodw==",
    "consensusMetadata": "CAI=",
    "nonHashData": {
        "localLedgerCommitTimestamp": {
            "seconds": 1466577447,
            "nanos": 653618964
        },
        "transactionResults": [
            {
                "uuid": "2b3b6cf3-9887-4dd5-8f2e-3634ec9c719a"
            }
        ]
    }
}
```
