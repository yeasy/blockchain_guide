##智能合约案例一代码分析
###`chaincode_finished.go`主要实现如下的功能：
- 初始化key `hello_world` 的value
- 读取和修改key `args[0]`的value

该部署在fabric网络上的初始化`hello_world`的值，并根据请求中的参数创建修改查询链上`key`中的值，本质上实现了一个简单的可修改的字典。

###function及实现的功能
- `read`  读取key `args[0]` 的value
- `write`  创建或修改key `args[0]` 的value
- `init`  初始化key `hello_world`的value
- `invoke`  根据传递参数类型调用执行相应的 `init`和`write`函数
- `query`  调用`read`函数查询 `args[0]`的value

###智能合约代码运行分析

`main`函数作为程序的入口，调用shim的start函数，启动chaincode引导程序的入口节点。
如果报错，返回。
```
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
```

当智能合约部署在区块链上，可以通过rest api进行交互。
三个主要的函数是`init`,`invoke`,`query`。在三个函数中，通过`stub.PutState`与`stub.GetState`存储访问ledger上的键值对。
###智能合约request json

假设以jim身份登录pbft集群，请求部署该chaincode json格式为：
```
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
目前path仅支持github上的目录，ctorMsg中为函数`init`的传参。

调用invoke函数的json格式为：
```
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
其中name字段为`deploy`后返回的message字段中的字符串。`query`与上诉两个类似。
