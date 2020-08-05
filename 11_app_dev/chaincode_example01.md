## 链码示例一：信息公证
### 简介

[chaincode_example01.go](chaincode_example01.go) 主要实现如下的功能：

* 初始化，以键值形式存放信息；
* 允许读取和修改键值。

代码中，首先初始化了 `hello_world` 的值，并根据请求中的参数创建修改查询链上 `key` 中的值，本质上实现了一个简单的可修改的键值数据库。

### 主要函数

* `read`：读取key `args[0]` 的 value；
* `write`：创建或修改 key `args[0]` 的 value；
* `init`：初始化 key `hello_world` 的 value；
* `invoke`：根据传递参数类型调用执行相应的 `init` 和 `write` 函数；
* `query`：调用 `read` 函数查询 `args[0]` 的 value。

### 代码运行分析

`main` 函数作为程序的入口，调用 shim 包的 start 函数，启动 chaincode 引导程序的入口节点。如果报错，则返回。

```go
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
```

当智能合约部署在区块链上，可以通过 rest api 进行交互。

三个主要的函数是 `init`，`invoke`，`query`。在三个函数中，通过 `stub.PutState`与 `stub.GetState` 存储访问 ledger 上的键值对。

### 通过 REST API 操作智能合约

假设以 jim 身份登录 pbft 集群，请求部署该 chaincode 的 json 请求格式为：
```json
{
    "jsonrpc": "2.0",
    "method": "deploy",
    "params": {
        "type": 1,
        "chaincodeID": {
            "path": "https://github.com/ibm-blockchain/learn-chaincode/finished"
        },
        "ctorMsg": {
            "function": "init",
            "args": [
                "hi there"
            ]
        },
        "secureContext": "jim"
    },
    "id": 1
}
```

目前 path 仅支持 github 上的目录，ctorMsg 中为函数 `init` 的传参。

调用 invoke 函数的 json 格式为：

```json
{
    "jsonrpc": "2.0",
    "method": "invoke",
    "params": {
        "type": 1,
        "chaincodeID": {
            "name": "4251b5512bad70bcd0947809b163bbc8398924b29d4a37554f2dc2b033617c19cc0611365eb4322cf309b9a5a78a5dba8a5a09baa110ed2d8aeee186c6e94431"
        },
        "ctorMsg": {
            "function": "init",
            "args": [
                "swb"
            ]
        },
        "secureContext": "jim"
    },
    "id": 2
}
```

其中 name 字段为 `deploy` 后返回的 message 字段中的字符串。

`query` 的接口也是类似的。
