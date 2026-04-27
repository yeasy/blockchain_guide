## 链码概念与结构

超级账本 Fabric 中的链码（Chaincode）负责响应应用程序提交的交易请求，执行合约逻辑，并与账本状态交互。

链码执行过程中可以创建状态并写入账本。状态绑定到链码命名空间，默认仅由该链码访问；在合适的权限和业务设计下，链码也可以调用其他链码。Fabric 2.x 支持内置链码运行时，也支持外部链码服务模式。

Fabric 支持多种链码开发语言，包括 Go、JavaScript/TypeScript、Java 等。Go 新项目优先使用 Contract API。

### 推荐：Contract API

Contract API 用普通 Go 方法描述交易函数，框架负责参数映射、返回值序列化、错误处理和交易上下文注入。当前示例使用 v2 模块：

```go
require github.com/hyperledger/fabric-contract-api-go/v2 v2.2.1
```

最小结构如下。

```go
package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return ctx.GetStub().PutState("hello", []byte("world"))
}

func (s *SmartContract) Read(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	value, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", key, err)
	}
	if value == nil {
		return "", fmt.Errorf("state %s does not exist", key)
	}
	return string(value), nil
}

func (s *SmartContract) Write(ctx contractapi.TransactionContextInterface, key string, value string) error {
	return ctx.GetStub().PutState(key, []byte(value))
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err)
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
```

### 可选：底层 shim API

需要直接控制 Fabric 响应对象时，可以使用 shim API。当前 Go v2 包路径如下：

```go
import (
	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	pb "github.com/hyperledger/fabric-protos-go-apiv2/peer"
)
```

底层接口只包含 `Init` 和 `Invoke`：

```go
type Chaincode interface {
	Init(stub shim.ChaincodeStubInterface) pb.Response
	Invoke(stub shim.ChaincodeStubInterface) pb.Response
}
```

Fabric v1 早期示例中的独立 `Query` 回调已经不是当前链码接口的一部分。读取类交易也通过 `Invoke` 路径进入链码，再由函数名路由到只读逻辑；客户端侧则通常用 Gateway 的 evaluate 类调用表达“只读查询”。
