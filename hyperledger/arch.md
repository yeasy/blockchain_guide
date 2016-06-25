## 架构设计

整个架构如下图所示。

![](_images/refarch.png)

包括三大组件：成员权限管理（Membership）、区块链服务（Blockchain）、链上代码服务（Chaincode）。

### 基本术语
* 交易处理（Transaction）：执行账本上的某个函数调用。函数在 chaincode 中实现；
* 交易员（Transactor）：作为客户端发起交易调用；
* 账本（Ledger）：即区块链，带有所有的交易信息和当前的世界观（world state）；
* 世界观（world state）：当前账本的一个（稳定）状态，包括各种 chaincode 的状态，执行交易的结果等；
* 链上代码（Chaincode）：区块链上的应用代码，延伸自“智能合约”，支持 golang、nodejs 等；
* 验证节点（Validating Peer）：维护账本的核心节点，参与一致性维护、对交易的验证和执行；
* 非验证节点（Non-validating Peer）：不参与账本维护，仅作为交易代理响应客户端的 REST 请求，并对交易进行检验，之后转发给验证节点；
* 带许可的账本（Permissioned Ledger）：网络中所有节点必须是经过许可的，非许可过的节点则无法加入网络；
* 隐私保护（Privacy）：交易员可以隐藏交易的身份，其它成员在无特殊权限的情况下，只能对交易进行验证，而无法获知身份信息；
* 秘密保护（Confidentiality）：只有交易双方可以看到交易内容，其它人未经授权则无法看到；
* 审计性（Auditability）：在一定权限和许可下，可以对链上的交易进行审计和检查。

### 成员权限管理

### 区块链服务

区块内容示例：

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

### 链上代码服务
