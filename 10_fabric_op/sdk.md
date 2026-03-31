## SDK 支持

除了基于命令行的客户端之外，超级账本 Fabric 历史上还提供过多种语言的 SDK，包括 Node.js、Python、Java、Go 等。但对新项目来说，更重要的是区分**历史 SDK**与**当前推荐路径**。

这些客户端能力允许用户和应用跟 Fabric 网络进行交互，但现代应用开发的推荐主路径已经转向 **Fabric Gateway**。因此，本节后半部分提到的传统 SDK 应更多理解为**历史项目或兼容性参考**，而不是新项目的默认选型。

**特别说明**：自 Fabric v2.4+ 起，Fabric Gateway 成为推荐的应用开发方式，它提供了更简洁的 API，简化了交易构造和背书收集的过程。建议新项目优先使用 Fabric Gateway API 而非传统的 SDK 接口。

### Fabric Gateway（推荐）

自 Fabric v2.4 起，**Fabric Gateway** 成为应用开发的首选方式。Gateway 将交易提案构造、背书收集和提交等复杂逻辑从客户端移至 Peer 节点内部，开发者只需连接到一个可信的 Gateway Peer 即可完成所有操作。

Gateway API 支持 Go、Node.js 和 Java 三种语言，源码仓库地址在 github.com/hyperledger/fabric-gateway。

以下是使用 Go 语言 Gateway API 提交交易的简化示例：

```go
package main

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 1. 加载 TLS 证书，建立 gRPC 连接到 Gateway Peer
	certPEM, _ := os.ReadFile("path/to/tls-ca-cert.pem")
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certPEM)
	tlsCreds := credentials.NewClientTLSFromCert(certPool, "peer0.org1.example.com")

	clientConnection, _ := grpc.NewClient("peer0.org1.example.com:7051",
		grpc.WithTransportCredentials(tlsCreds))
	defer clientConnection.Close()

	// 2. 创建客户端身份（从 MSP 证书和私钥）
	clientCert, _ := os.ReadFile("path/to/client-cert.pem")
	clientID, _ := identity.NewX509Identity("Org1MSP", identity.CertificateToPEM(clientCert))
	// ... 此处省略私钥加载和 Sign 函数创建

	// 3. 创建 Gateway 连接
	gw, _ := client.Connect(clientID, client.WithSign(sign),
		client.WithClientConnection(clientConnection))
	defer gw.Close()

	// 4. 获取网络和合约引用
	network := gw.GetNetwork("mychannel")
	contract := network.GetContract("mychaincode")

	// 5. 提交交易（自动完成背书、排序、提交）
	result, _ := contract.SubmitTransaction("CreateAsset", "asset1", "blue", "10")
	fmt.Printf("Transaction result: %s\n", string(result))

	// 6. 查询账本（仅背书，不提交）
	result, _ = contract.EvaluateTransaction("ReadAsset", "asset1")
	fmt.Printf("Query result: %s\n", string(result))
}
```

*注：以上为简化示例，省略了错误处理和私钥加载逻辑。完整可运行的样例参见 [fabric-samples/asset-transfer-basic](https://github.com/hyperledger/fabric-samples/tree/main/asset-transfer-basic/application-gateway-go)。*

相比传统 SDK，Fabric Gateway 的优势在于：客户端逻辑大幅简化（无需手动管理背书策略和节点发现）、自动重试机制、更好的错误处理，以及显著减少的网络往返次数。

### 历史 Node.js SDK（兼容性参考）

作为早期创建的 SDK 项目之一，Node.js SDK 曾支持对 Fabric 链码的主要操作以及访问 Fabric CA 服务。但在现代 Fabric 版本中，它更适合作为历史参考；新开发建议优先选择 Fabric Gateway。

源码仓库地址在 github.com/hyperledger/fabric-sdk-node。

源码的 test/integration/e2e 目录下包括了大量应用的示例代码，可供参考。

### 历史 Python SDK（历史项目）

Python 实现的 SDK 是早期项目，当前更多适合作为历史资料参考；如果是新的生产项目，不应把它当作官方推荐主线。

源码仓库地址在 github.com/hyperledger/fabric-sdk-py。

源码的 test/integration 目录下包括了大量应用的示例代码，可供参考。

### 历史 Java SDK（兼容性参考）

Java SDK 曾支持对 Fabric 中链码的主要操作以及访问 Fabric CA 服务，但现代 Java 应用开发同样更推荐 Fabric Gateway。

源码仓库地址在 github.com/hyperledger/fabric-sdk-java。

### 历史 Go SDK（已被 Gateway 路线取代）

Go SDK 也属于较早期的客户端路线。对新项目而言，更稳妥的建议仍然是优先使用 Fabric Gateway，而不是把传统 Go SDK 视为首选方案。

源码仓库地址在 github.com/hyperledger/fabric-sdk-go。
