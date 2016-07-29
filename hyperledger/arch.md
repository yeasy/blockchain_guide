## 架构设计

整个架构如下图所示。

![](_images/refarch.png)

包括三大组件：区块链服务（Blockchain）、链码服务（Chaincode）、成员权限管理（Membership）。

### 基本术语

* 交易处理（Transaction）：执行账本上的某个函数调用。函数在 chaincode 中实现；
* 交易员（Transactor）：作为客户端发起交易调用；
* 账本（Ledger）：即区块链，带有所有的交易信息和当前的世界状态（world state）；
* 世界状态（World State）：当前账本的一个（稳定）状态，包括所有 chaincode 中所有键值对的集合。是一个键值集合，一般用 `{chaincodeID, ckey}` 代表键；
* 链码（Chaincode）：区块链上的应用代码，延伸自“智能合约”，支持 golang、nodejs 等；
* 验证节点（Validating Peer）：维护账本的核心节点，参与一致性维护、对交易的验证和执行；
* 非验证节点（Non-validating Peer）：不参与账本维护，仅作为交易代理响应客户端的 REST 请求，并对交易进行一些基本的有效性检查，之后转发给验证节点；
* 带许可的账本（Permissioned Ledger）：网络中所有节点必须是经过许可的，非许可过的节点则无法加入网络；
* 隐私保护（Privacy）：交易员可以隐藏交易的身份，其它成员在无特殊权限的情况下，只能对交易进行验证，而无法获知身份信息；
* 秘密保护（Confidentiality）：只有交易双方可以看到交易内容，其它人未经授权则无法看到；
* 审计性（Auditability）：在一定权限和许可下，可以对链上的交易进行审计和检查。

### 区块链服务

区块链服务提供一个分布式账本平台。一般地，多个交易被打包进区块中，多个区块构成一条区块链。

#### 交易

交易意味着围绕着某个链码进行操作。

交易可以改变世界状态。

交易中包括的内容主要有：

* 交易类型：目前包括 Deploy、Invoke、Query、Terminate 四种；
* uuid：代表交易的唯一编号；
* 链码编号 chaincodeID：交易针对的链码；
* 负载内容的 hash 值：Deploy 或 Invoke 时候可以指定负载内容；
* 交易的保密等级 ConfidentialityLevel；
* 交易相关的 metadata 信息；
* 临时生成值 nonce：跟安全机制相关；
* 交易者的证书信息 cert；
* 签名信息 signature；
* metadata 信息；
* 时间戳 timestamp。

交易的数据结构（Protobuf 格式）定义为

```protobuf
message Transaction {
    enum Type {
        UNDEFINED = 0;
        // deploy a chaincode to the network and call `Init` function
        CHAINCODE_DEPLOY = 1;
        // call a chaincode `Invoke` function as a transaction
        CHAINCODE_INVOKE = 2;
        // call a chaincode `query` function
        CHAINCODE_QUERY = 3;
        // terminate a chaincode; not implemented yet
        CHAINCODE_TERMINATE = 4;
    }
    Type type = 1;
    //store ChaincodeID as bytes so its encrypted value can be stored
    bytes chaincodeID = 2;
    bytes payload = 3;
    bytes metadata = 4;
    string uuid = 5;
    google.protobuf.Timestamp timestamp = 6;

    ConfidentialityLevel confidentialityLevel = 7;
    string confidentialityProtocolVersion = 8;
    bytes nonce = 9;

    bytes toValidators = 10;
    bytes cert = 11;
    bytes signature = 12;
}
```

#### 区块

区块打包交易，确认交易后的世界状态。

一个区块中包括的内容主要有：

* 版本号 version：协议的版本信息；
* 时间戳 timestamp：由区块提议者设定；
* 交易信息的默克尔树的根 hash 值：由区块所包括的交易构成；
* 世界观的默克尔树的根 hash 值：由当前整个世界的状态值构成；
* 前一个区块的 hash 值：构成链所必须；
* 共识相关的元数据：可选值；
* 非 hash 数据：不参与 hash 过程，各个 peer 上的值可能不同，例如本地提交时间、交易处理的返回值等；

_注意具体的交易信息并不存放在区块中。_

交易的数据结构（Protobuf 格式）定义为

```protobuf
message Block {
    uint32 version = 1;
    google.protobuf.Timestamp timestamp = 2;
    repeated Transaction transactions = 3;
    bytes stateHash = 4;
    bytes previousBlockHash = 5;
    bytes consensusMetadata = 6;
    NonHashData nonHashData = 7;
}
```

一个真实的区块内容示例：

```json
{
    "nonHashData": {
        "localLedgerCommitTimestamp": {
            "nanos": 975295157,
                "seconds": 1466057539
        },
            "transactionResults": [
            {
                "uuid": "7be1529ee16969baf9f3156247a0ee8e7eee99a6a0a816776acff65e6e1def71249f4cb1cad5e0f0b60b25dd2a6975efb282741c0e1ecc53fa8c10a9aaa31137"
            }
            ]
    },
        "previousBlockHash": "RrndKwuojRMjOz/rdD7rJD/NUupiuBuCtQwnZG7Vdi/XXcTd2MDyAMsFAZ1ntZL2/IIcSUeatIZAKS6ss7fEvg==",
        "stateHash": "TiIwROg48Z4xXFFIPEunNpavMxnvmZKg+yFxKK3VBY0zqiK3L0QQ5ILIV85iy7U+EiVhwEbkBb1Kb7w1ddqU5g==",
        "transactions": [
        {
            "chaincodeID": "CkdnaXRodWIuY29tL2h5cGVybGVkZ2VyL2ZhYnJpYy9leGFtcGxlcy9jaGFpbmNvZGUvZ28vY2hhaW5jb2RlX2V4YW1wbGUwMhKAATdiZTE1MjllZTE2OTY5YmFmOWYzMTU2MjQ3YTBlZThlN2VlZTk5YTZhMGE4MTY3NzZhY2ZmNjVlNmUxZGVmNzEyNDlmNGNiMWNhZDVlMGYwYjYwYjI1ZGQyYTY5NzVlZmIyODI3NDFjMGUxZWNjNTNmYThjMTBhOWFhYTMxMTM3",
            "payload": "Cu0BCAESzAEKR2dpdGh1Yi5jb20vaHlwZXJsZWRnZXIvZmFicmljL2V4YW1wbGVzL2NoYWluY29kZS9nby9jaGFpbmNvZGVfZXhhbXBsZTAyEoABN2JlMTUyOWVlMTY5NjliYWY5ZjMxNTYyNDdhMGVlOGU3ZWVlOTlhNmEwYTgxNjc3NmFjZmY2NWU2ZTFkZWY3MTI0OWY0Y2IxY2FkNWUwZjBiNjBiMjVkZDJhNjk3NWVmYjI4Mjc0MWMwZTFlY2M1M2ZhOGMxMGE5YWFhMzExMzcaGgoEaW5pdBIBYRIFMTAwMDASAWISBTIwMDAw",
            "timestamp": {
                "nanos": 298275779,
                "seconds": 1466057529
            },
            "type": 1,
            "uuid": "7be1529ee16969baf9f3156247a0ee8e7eee99a6a0a816776acff65e6e1def71249f4cb1cad5e0f0b60b25dd2a6975efb282741c0e1ecc53fa8c10a9aaa31137"
        }
    ]
}
```

### 链码服务

链码包含所有的处理逻辑，并对外提供接口，外部通过调用链码接口来改变世界观。

链码目前支持的交易类型包括：部署（Deploy）、调用（Invoke）和查询（Query）。

* 部署：VP 节点利用链码创建沙盒，沙盒启动后，处理 protobuf 协议的 shim 层一次性发送包含 ChaincodeID 信息的 REGISTER 消息给 VP 节点，进行注册，注册完成后，VP 节点通过 gRPC 传递参数并调用链码 Invoke 函数完成初始化；
* 调用：VP 节点发送 TRANSACTION 消息给链码沙盒的 shim 层，shim 层用传过来的参数调用链码的 Invoke 函数完成调用；
* 查询：VP 节点发送 QUERY 消息给链码沙盒的 shim 层，shim 层用传过来的参数调用链码的 Query 函数完成查询。

不同链码之间可能互相调用和查询。

### 成员权限管理

通过基于 PKI 的成员权限管理，平台可以对接入的节点和客户端的能力进行限制。

证书有三种，Enrollment，Transaction，以及确保安全通信的 TLS 证书。

* 注册证书 ECert：颁发给提供了注册凭证的用户或节点，一般长期有效；
* 交易证书 TCert：颁发给用户，控制每个交易的权限，一般针对某个交易，短期有效。
* 通信证书 TLSCert：控制对网络的访问，并且防止窃听。

![](_images/memserv-components.png)

### 消息类型

节点之间通过消息来进行交互，所有消息都由下面的数据结构来实现。

```protobuf
message Message {
   enum Type {
        UNDEFINED = 0;

        DISC_HELLO = 1;
        DISC_DISCONNECT = 2;
        DISC_GET_PEERS = 3;
        DISC_PEERS = 4;
        DISC_NEWMSG = 5;

        CHAIN_STATUS = 6;
        CHAIN_TRANSACTION = 7;
        CHAIN_GET_TRANSACTIONS = 8;
        CHAIN_QUERY = 9;

        SYNC_GET_BLOCKS = 11;
        SYNC_BLOCKS = 12;
        SYNC_BLOCK_ADDED = 13;

        SYNC_STATE_GET_SNAPSHOT = 14;
        SYNC_STATE_SNAPSHOT = 15;
        SYNC_STATE_GET_DELTAS = 16;
        SYNC_STATE_DELTAS = 17;

        RESPONSE = 20;
        CONSENSUS = 21;
    }
    Type type = 1;
    bytes payload = 2;
    google.protobuf.Timestamp timestamp = 3;
}
```

消息分为四大类：Discovery（探测）、Transaction（交易）、Synchronization（同步）、Consensus（一致性）。

不同消息类型，payload 中数据不同。

#### Discovery
包括 DISC_HELLO、DISC_GET_PEERS、DISC_PEERS。

#### Transaction
包括 Deploy、Invoke、Query。

#### Synchronization
SYNC_GET_BLOCKS 和对应的 SYNC_BLOCKS。

SYNC_STATE_GET_SNAPSHOT 和对应的 SYNC_STATE_SNAPSHOT。

SYNC_STATE_GET_DELTAS 和对应的 SYNC_STATE_DELTAS。

#### Consensus
CONSENSUS 消息。