## 网络协议与交互机制

Fabric 网络的各个节点以及客户端为了能够正确处理交易、同步数据和管理集群，依赖于一套高效且安全的通信机制。这套机制主要基于 **gRPC** 远程调用框架以及用于实现去中心化数据发现与同步的 **Gossip 数据分发协议** 构建。

不再使用早期版本中那种全网广播、高度耦合的通信模式，现代 Fabric 将节点间的通信职责进行了极大的细化和专业化。主要包含了三大核心通信接口与协议层：**Peer 节点间数据交互（Gossip）**、**客户端到 Peer 的 Gateway/Endorsement 交互（gRPC）**以及**排序交付交互（gRPC Deliver/Broadcast）**。

### 1. 基于 Gossip 的数据流协议

Fabric 采用了极其高效且具备很好容错率的 **Gossip 协议**（流言协议）来实现 Peer 节点之间的成员发现、服务发现元数据传播、私有数据分发，以及可选的区块传播和状态补齐。需要注意的是，Fabric v2.2 之后的默认和推荐生产方向通常是让 Peer 直接从排序服务（Orderer）拉取区块；使用 Gossip 在组织内部级联传播区块属于可配置能力，而不是现代部署的必然路径。

对于同一个通道内的各个组织（Organizations）而言，Gossip 重点承担了以下职责：

#### a. 成员发现服务（Service Discovery）

节点刚启动或网络拓扑发生变动时，如何知道哪些友方节点在线？
- Peer 节点会定期广播包含其自身标识与证书的在线探测消息。
- 接收到消息的活动节点可以利用这些信息构建出所在通道的局部拓扑网络。为了防止跨组织的信息泄露，Fabric 严格限定了 Gossip 消息的接收者必须拥有相应的通道只读权限并校验身份，不同组织或不同通道之间的 Gossip 是隔离的。

#### b. 交易和区块数据分发（Data Dissemination）

当 Orderer 将一批交易打包成区块后，Peer 通过 Deliver 服务拉取新区块。现代配置中，多个 Peer 甚至所有 Peer 都可以直接从排序服务拉取区块，以减少组织内的级联转发依赖。
- 需要节省排序服务连接或兼容旧部署时，可以配置一个或多个 **领导节点（Leader Peer）** 从排序服务拉取区块。
- 只有在启用对应配置时，Leader 节点才会通过 Gossip 以点对点（P2P）的方式在组织内部继续分发新区块。
- 如果某个节点因为网络抖动落后了，可以在启用状态传输配置时向其他 Peer 拉取缺失区块；也可以直接从排序服务补齐。

#### c. 私有数据（Private Data）安全分发

在涉及企业机密的场景中，私有数据集合（Private Data Collection）不会把明文私有数据写入通道公共区块，也不会提交给 Orderer 排序。实际私有数据通过 Gossip 点对点分发给集合策略授权的组织 Peer，并存入这些 Peer 的私有状态数据库；公共账本中只记录私有数据键和值的哈希，用于全通道一致验证、审计和后续证明。

### 2. 客户端与 Peer 的交互：Gateway 与背书服务接口

Fabric v2.4 之后，客户端推荐通过 Gateway API 连接一个受信任的 Peer。Gateway 会调用发现服务确定满足背书策略所需的 Peer，并代表客户端发起背书请求、收集足够背书、提交交易并等待提交状态。底层背书执行仍由 Peer 的 `Endorser` 服务完成：

```protobuf
service Endorser {
    // Gateway 或兼容旧客户端向被选中的背书节点发送交易提案（Proposal）
    rpc ProcessProposal(SignedProposal) returns (ProposalResponse) {}
}
```

其通信流程非常简洁：
1. **Gateway 发起调用**：发送封装了链码名称、函数、调用参数以及客户端签名信息的 `SignedProposal` 结构数据。兼容旧 SDK 的客户端也可以直接调用该接口。
2. **节点同步响应**：Peer 在沙盒中执行完毕后，同步返回包含了模拟执行的读写集以及最重要的——**Peer自身数字签名的结果（ProposalResponse）**。这构成了整个“执行-排序-验证”模型的基础。

Gateway 会根据发现服务和背书策略向足够的 Peer 请求背书，并在收集到满足策略的签名后组装交易信封返回给客户端签名。客户端直接并行请求多个背书 Peer 是旧 SDK 兼容路径，不是现代 Gateway API 的默认职责划分。

### 3. 客户端（以及Peer）与排序节点的交互：AtomicBroadcast 接口

无论是 Gateway/兼容客户端向 Orderer 提交已经背书完的合法交易（阶段二），还是 Peer 节点向 Orderer 订阅接收新打包区块的数据流（阶段三），统一由排序服务提供的 `AtomicBroadcast` 接口来完成：

```protobuf
service AtomicBroadcast {
    // 1. Gateway 或兼容客户端使用此接口向排序节点发送交易进行排序
    rpc Broadcast(stream common.Envelope) returns (stream BroadcastResponse) {}

    // 2. Peer 节点使用此接口订阅或请求新生成的区块
    rpc Deliver(stream common.Envelope) returns (stream DeliverResponse) {}
}
```

* **Broadcast 操作**：是一个基于 gRPC Stream（双向数据流）的单向接收端接口。排序节点接收到 `Envelope`（封装了带有全套签名的交易）后，校验发起者是否有向该通道写入的权限，校验通过后等待切分打包并排序。
* **Deliver 操作**：该接口允许具有合法读取通道权限的 Peer 拉取指定范围的区块数据流。Peer 可以作为组织 Leader，也可以作为普通 Peer 直接从排序服务拉取区块；是否再通过 Gossip 传播新区块取决于配置。

这两个高度精简的接口保证了排序层可以无缝地替换背后的共识引擎（从过去的 Kafka 到现代标准的 Raft），而不必去修改外围组件的任何通信协议层。
