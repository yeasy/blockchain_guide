## 链上代码

### 什么是 chaincode
chaincode（链码）是部署在 Hyperledger fabric 网络节点上，可被调用与分布式账本进行交互的一段程序代码，也即狭义范畴上的“智能合约”。链码在 VP 节点上的隔离沙盒（目前为 Docker 容器）中执行，并通过 gRPC 协议来被相应的 VP 节点调用和查询。

Hyperledger 支持多种计算机语言实现的 chaincode，包括 Golang、JavaScript、Java 等。

### 实现 chaincode 接口
下面以 golang 为例来实现 chaincode 的 shim 接口。在这之中三个核心的函数是 **Init**, **Invoke**， 和 **Query**。三个函数都以函数名和字符串结构作为输入，主要的区别在于三个函数被调用的时机。

#### 依赖包

chaincode 需要引入如下的软件包。

* `fmt`：包含了 `Println` 等标准函数.
* `errors`：标准 errors 类型包；
* `github.com/hyperledger/fabric/core/chaincode/shim`：与 chaincode 节点交互的接口代码。shim 包 提供了 `stub.PutState` 与 `stub.GetState` 等方法来写入和查询链上键值对的状态。

比较重要的 shim 包，通过封装 gRPC 消息到 VP 节点来完成操作，如：

* PUT_STATE：修改某个状态（键值）的值；
* GET_STATE：获取某个状态的值；
* DEL_STATE：删除某个键值；
* RANGE_QUERY_STATE：获取某个范围内的键值，需要键的命名可构成规则的范围；
* INVOKE_CHAINCODE：调用其它链码方法；
* QUERY_CHAINCODE：查询同一上下文下的其它链码。

#### Init()函数
当首次部署 chaincode 代码时，init 函数被调用。如同名字所描述的，该函数用来做一些初始化的工作。

#### Invoke()函数
当通过调用 chaincode 代码来做一些实际性的工作时，可以使用 invoke 函数。发起的交易将会被链上的区块获取并记录。

它以被调用的函数名作为参数，并基于该参数去调用 chaincode 中匹配的的 go 函数。

#### Query()函数
顾名思义，当需要查询 chaincode 的状态时，可以调用 `Quer()` 函数。

#### Main() 函数
最后，需要创建一个 `main` 函数，当每个节点部署 chaincode 的实例时，该函数会被调用。

它仅仅在 chaincode 在某节点上注册时会被调用。


### 与 chaincode 代码进行交互
与 chaincode 交互的主要方法有 cli 命令行与 rest api，关于 rest api 的使用请查看该目录下的例子。
