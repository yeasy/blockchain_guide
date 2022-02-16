## 链码概念与结构

超级账本 Fabric 中的链码（Chaincode）延伸自智能合约的概念，负责对应用程序发送的请求做出响应，执行代码逻辑，实现与账本进行交互。

区块链网络中成员协商好业务逻辑后，可将其编程到链码中，所有业务流程将遵循合约代码自动执行。

链码执行过程中可以创建状态（State）并写入账本中。状态绑定到链码的命名空间，仅限该链码访问。在合适许可下，链码可以调用另一个链码，间接访问其状态。在一些场景下，不仅需要访问状态当前值，还需要查询状态的历史值。

原生链码默认在 Docker 容器中执行，2.0 版本中开始支持外部链码。链码通过 gRPC 协议与 Peer 节点进行交互，包括读写账本、返回响应结果等。

Fabric 支持多种语言实现的链码，包括 Go、JavaScript、Java 等。下面以 Go 语言为例介绍链码接口和相关结构。

### Chaincode 接口

每个链码都需要实现以下 Chaincode 接口，包括 Init 和 Invoke 两个方法。

```go
type Chaincode interface {
	Init(stub ChaincodeStubInterface) pb.Response
	Invoke(stub ChaincodeStubInterface) pb.Response
}
```

* Init：当链码收到初始化指令时，Init 方法会被调用。
* Invoke：当链码收到调用（invoke）或查询（query）调用时，Invoke 方法会被调用。

### 链码结构

一个链码的必要结构如下所示：

```go
package main

// 引入必要的包
import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// 声明一个结构体
type SimpleChaincode struct {}

// 为结构体添加 Init 方法
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// 在该方法中实现链码初始化或升级时的处理逻辑
	// 编写时可灵活使用 stub 中的 API
}

// 为结构体添加 Invoke 方法
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// 在该方法中实现链码运行中被调用或查询时的处理逻辑
	// 编写时可灵活使用 stub 中的 API
}

// 主函数，需要调用 shim.Start() 方法
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
```

#### 依赖包

从 `import` 代码段可以看到，链码需要引入如下的依赖包。

* `"github.com/hyperledger/fabric-chaincode-go/shim"`：shim 包提供了链码与账本交互的中间层。链码通过 `shim.ChaincodeStub` 提供的方法来读取和修改账本状态。
* `pb "github.com/hyperledger/fabric-protos-go/peer"`: Init 和 Invoke 方法需要返回 `pb.Response` 类型。

#### Init 和 Invoke 方法

编写链码，关键是要实现 Init 和 Invoke 这两个方法。

当初始化链码时，Init 方法会被调用。如同名字所描述的，该方法用来完成一些初始化的工作。当调用链码时，Invoke 方法被调用，主要业务逻辑都需要在该方法中实现。

Init 或 Invoke 方法以 `stub shim.ChaincodeStubInterface` 作为传入参数，`pb.Response` 作为返回类型。其中，`stub` 封装了丰富的 API，功能包括对账本进行操作、读取交易参数、调用其它链码等。
