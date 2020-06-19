## SDK 支持

除了基于命令行的客户端之外，超级账本 Fabric 提供了多种语言的 SDK，包括 Node.Js、Python、Java、Go 等。它们封装了 Fabric 网络中节点提供的 gRPC 服务接口，可以实现更方便的调用。

这些客户端 SDK 允许用户和应用跟 Fabric 网络进行交互，还可以实现更为复杂的操作，实现包括节点的启停、通道的创建和加入、链码的生命周期管理等操作。SDK 项目目前已经初步成熟，更多特性仍在开发中，感兴趣的读者可以通过如下途径获取到 SDK 的源码并进行尝试。

### 基于 Node.Js 实现的 SDK

作为早期创建的 SDK 项目之一，Node.Js 实现的 SDK 目前已经支持了对 Fabric 链码的主要操作，包括安装链码、实例化并进行调用等，以及访问 Fabric CA 服务。内带了不少操作的例子可供参考。

源码仓库地址在 github.com/hyperledger/fabric-sdk-node。

源码的 test/integration/e2e 目录下包括了大量应用的示例代码，可供参考。

### 基于 Python 实现的 SDK

早期创建的 SDK 项目之一。Python 实现的 SDK 目前已经完成了对 Fabric 链码的主要操作，包括安装链码、实例化并进行调用等，以及使用 Fabric CA 的基础功能。

源码仓库地址在 github.com/hyperledger/fabric-sdk-py。

源码的 test/integration 目录下包括了大量应用的示例代码，可供参考。

### 基于 Java 实现的 SDK

属于较新的 SDK 项目。Java SDK 目前支持对 Fabric 中链码的主要操作，以及访问 Fabric CA 服务。

源码仓库地址在 github.com/hyperledger/fabric-sdk-java。

### 基于 Go 实现的 SDK

属于较新的 SDK 项目。Go SDK 提取了原先 Fabric 中的相关代码，目前支持对 Fabric 中链码的主要操作。将来，Fabric 中的命令行客户端将可能基于该 SDK 重新实现。

源码仓库地址在 github.com/hyperledger/fabric-sdk-go。
