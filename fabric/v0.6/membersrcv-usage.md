## 权限管理

权限管理机制是 hyperledger fabric 项目的一大特色。下面给出使用权限管理的一个应用案例。

### 启动集群

首先下载相关镜像。

```sh
$ docker pull yeasy/hyperledger:latest
$ docker tag yeasy/hyperledger:latest hyperledger/fabric-baseimage:latest
$ docker pull yeasy/hyperledger-peer:latest
$ docker pull yeasy/hyperledger-membersrvc:latest
```

进入 hyperledger 项目，启动带成员管理的 PBFT 集群。

```sh
$ git clone https://github.com/yeasy/docker-compose-files
$ cd docker-compose-files/hyperledger/0.6/pbft
$ docker-compose -f 4-peers-with-membersrvc.yml up
```

### 用户登陆

当启用了权限管理后，首先需要登录，例如以内置账户 jim 账户登录。

登录 vp0，并执行登录命令。

```sh
$ docker exec -it pbft_vp0_1 bash
# peer network login jim
06:57:13.603 [networkCmd] networkLogin -> INFO 001 CLI client login...
06:57:13.603 [networkCmd] networkLogin -> INFO 002 Local data store for client loginToken: /var/hyperledger/production/client/
Enter password for user 'jim': 6avZQLwcUe9b
06:57:25.022 [networkCmd] networkLogin -> INFO 003 Logging in user 'jim' on CLI interface...
06:57:25.576 [networkCmd] networkLogin -> INFO 004 Storing login token for user 'jim'.
06:57:25.576 [networkCmd] networkLogin -> INFO 005 Login successful for user 'jim'.
06:57:25.576 [main] main -> INFO 006 Exiting.....
```

也可以用 REST 方式：

```sh
POST  HOST:7050/registrar
```

Request：

```json
{
  "enrollId": "jim",
  "enrollSecret": "6avZQLwcUe9b"
}
```

Response：

```json
{
    "OK": "User jim is already logged in."
}
```

### chaincode 部署

登录之后，chaincode 的部署、调用等操作与之前类似，只是需要通过 -u 选项来指定用户名。

在 vp0 上执行命令：

```sh
# peer chaincode deploy -u jim -p github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02 -c '{"Function":"init", "Args": ["a","100", "b", "200"]}'
06:58:20.099 [chaincodeCmd] getChaincodeSpecification -> INFO 001 Local user 'jim' is already logged in. Retrieving login token.
06:58:22.178 [chaincodeCmd] chaincodeDeploy -> INFO 002 Deploy result: type:GOLANG chaincodeID:<path:"github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02" name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"init" args:"a" args:"100" args:"b" args:"200" >
Deploy chaincode: ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539
06:58:22.178 [main] main -> INFO 003 Exiting.....
```

记录下返回的 chaincode ID。

```sh
# CC_ID=ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539
```

此时，查询账户值应当为初始值。

```sh
# peer chaincode query -u jim -n ${CC_ID} -c '{"Function": "query", "Args": ["a"]}'
07:28:39.925 [chaincodeCmd] getChaincodeSpecification -> INFO 001 Local user 'jim' is already logged in. Retrieving login token.
07:28:40.281 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 002 Successfully queried transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"query" args:"a" > secureContext:"jim" >
Query Result: 100
07:28:40.281 [main] main -> INFO 003 Exiting.....
```

也可以通过 REST 方式进行：

```sh
POST  HOST:7050/chaincode
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
        "message": "ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539"
    },
    "id": 1
}
```

### chaincode 调用

在账户 a、b 之间进行转账 10 元的操作。

```sh
# peer chaincode invoke -u jim -n ${CC_ID} -c '{"Function": "invoke", "Args": ["a", "b", "10"]}'
07:29:25.245 [chaincodeCmd] getChaincodeSpecification -> INFO 001 Local user 'jim' is already logged in. Retrieving login token.
07:29:25.585 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 002 Successfully invoked transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"invoke" args:"a" args:"b" args:"10" > secureContext:"jim" > (f8347e3b-7230-4561-9017-3946756a0bf4)
07:29:25.585 [main] main -> INFO 003 Exiting.....
```

也可以通过 REST 方式进行：

```sh
POST  HOST:7050/chaincode
```

Request：

```json
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d"
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
        "message": "66308740-a2c5-4a60-81f1-778dbed49cc3"
    },
    "id": 3
}
```

### chaincode 查询

查询 a 账户的余额。

```sh
# peer chaincode query -u jim -n ${CC_ID} -c '{"Function": "query", "Args": ["a"]}'
07:29:55.844 [chaincodeCmd] getChaincodeSpecification -> INFO 001 Local user 'jim' is already logged in. Retrieving login token.
07:29:56.198 [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 002 Successfully queried transaction: chaincodeSpec:<type:GOLANG chaincodeID:<name:"ee5b24a1f17c356dd5f6e37307922e39ddba12e5d2e203ed93401d7d05eb0dd194fb9070549c5dc31eb63f4e654dbd5a1d86cbb30c48e3ab1812590cd0f78539" > ctorMsg:<args:"query" args:"a" > secureContext:"jim" >
Query Result: 90
07:29:56.198 [main] main -> INFO 003 Exiting.....
```

也可以通过 REST 方式进行：

```sh
POST  HOST:7050/chaincode
```

Request：

```json
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"980d4bb7f69578592e5775a6da86d81a221887817d7164d3e9d4d4df1c981440abf9a61417eaf8ad6f7fc79893da36de2cf4709131e9af39bca6ebc2e5a1cd9d"
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
GET  HOST:7050/chain/blocks/2
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



