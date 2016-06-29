##编写智能合约代码

### 什么是chaincode
chaincode是部署在hyperledger fabric网络节点上并与全网络的共享账本进行交互的一段代码，也即智能合约。

### 实现chaincode的接口
首先你需要做的事情是使用golang来实现chaincode的shim接口。在这之中三个核心的函数是**Init**, **Invoke**， and **Query**。三个函数都以函数名和字符串结构作为输入，主要的区别在于三个函数被调用的时机。

#### 依赖
- `fmt` - 包含了 `Println` 函数来为了 debugging/logging.
- `errors` - 标准go errors 包
- `github.com/hyperledger/fabric/core/chaincode/shim` - 与chaincode节点交互的接口代码。shim提供了·stub.PutState·与`stub.GetState`来写入和查询链上键值对的状态。

#### Init()
当你第一次部署你的chaincode代码时，init函数被调用。如同名字所描述的，该函数用来做一些初始化的工作。

#### Invoke()
当你想要调用chaincode代码来做一些实际性的工作时，你可以使用invoke函数。发起的交易将会被链上的区块获取。
它以一个函数名作为参数并基于该参数去调用chaincode中的go函数。

#### Query()
就如同名字所描述的，当你想查询chaincode的状态时，你可以调用`Query`函数

#### Main() 
最后，你需要创建一个`main`函数，当每个节点部署chaincode的实例时，该函数会被调用。
它仅仅在chaincode在某节点上注册时会被调用。


###与你的chaincode代码交互
与chaincode交互的主要方法有cli命令行与rest api，关于rest api的使用请查看该目录下的例子。
