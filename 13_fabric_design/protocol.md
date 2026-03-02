## 网络协议与交互机制

Fabric 网络的各个节点以及客户端为了能够正确处理交易、同步数据和管理集群，依赖于一套高效且安全的通信机制。这套机制主要基于 **gRPC** 远程调用框架以及用于实现去中心化数据发现与同步的 **Gossip 数据分发协议** 构建。

不再使用早期版本中那种全网广播、高度耦合的通信模式，现代 Fabric 将节点间的通信职责进行了极大的细化和专业化。主要包含了三大核心通信接口与协议层：**Peer 节点间数据交互（Gossip）**、**背书交互（gRPC Endorsement）**以及**排序交付交互（gRPC Deliver/Broadcast）**。

### 1. 基于 Gossip 的数据流协议

Fabric 采用了极其高效且具备很好容错率的 **Gossip 协议**（流言协议）来实现 Peer 节点之间的状态同步、通道数据分发以及服务成员发现。Gossip 极大地减轻了排序服务（Orderer）的广播压力，因为它不需要中心节点将庞大的区块数据一对一地发送给通道里的成千上万个 Peer，而是利用节点间的“口口相传”自动完成全网覆盖。

对于同一个通道内的各个组织（Organizations）而言，Gossip 重点承担了以下职责：

#### a. 成员发现服务（Service Discovery）
节点刚启动或网络拓扑发生变动时，如何知道哪些友方节点在线？
- Peer 节点会定期广播包含其自身标识与证书的在线探测消息。
- 接收到消息的活动节点可以利用这些信息构建出所在通道的局部拓扑网络。为了防止跨组织的信息泄露，Fabric 严格限定了 Gossip 消息的接收者必须拥有相应的通道只读权限并校验身份，不同组织或不同通道之间的 Gossip 是隔离的。

#### b. 交易和区块数据分发（Data Dissemination）
当 Orderer 将一批交易打包成区块后，它并不是将庞大的区块直接发给每个 Peer 节点。
- Orderer 仅仅向该组织预定义的少数几个被称为 **领导节点（Leader Peer）** 的代表发送这个新区块。
- Leader 节点接收到后，再通过 Gossip 以点对点（P2P）的方式迅速在组织内部的其他所有 **提交节点（Committer）** 间泛洪（Fan-out）广播。
- 如果某个节点因为网络抖动落后了，它也可以直接通过与相邻节点对比“区块高度”，自动向对方请求缺失的连续区块，这就是状态同步的过程（State Synchronization）。

#### c. 私有数据（Private Data）安全分发
在涉及企业机密的场景中，交易的具体内容不能被广播并被所有节点看到，只有被授权的组织才能解密。Gossip 被用于在授权节点的内存中私密地分发这些敏感数据，甚至这部分数据都不会被包含在公共的区块里提交给 Orderer 排序。

### 2. 客户端与背书节点的交互：背书服务接口

客户端向 Peer 节点请求执行智能合约（阶段一），这依赖于 Peer 对外暴露的基于 gRPC 的 `Endorser` 服务接口：

```protobuf
service Endorser {
    // 客户端向被选中的某几个背书节点发送交易提案（Proposal）
    rpc ProcessProposal(SignedProposal) returns (ProposalResponse) {}
}
```

其通信流程非常简洁：
1. **客户端发起调用**：发送封装了链码名称、函数、调用参数以及客户端签名信息的 `SignedProposal` 结构数据。
2. **节点同步响应**：Peer 在沙盒中执行完毕后，同步返回包含了模拟执行的读写集以及最重要的——**Peer自身数字签名的结果（ProposalResponse）**。这构成了整个“执行-排序-验证”模型的基础。

客户端可以并行地向多个不同的节点并行发送请求，当客户端集齐了满足背书策略所必须的所有签名后，才会进入下一个通信阶段。

### 3. 客户端（以及Peer）与排序节点的交互：AtomicBroadcast 接口

无论是客户端向 Orderer 提交已经背书完的合法交易（阶段二），还是 Peer 节点乃至客户端向 Orderer 订阅接收新打包区块的数据流（阶段三），统一由排序服务提供的 `AtomicBroadcast` 接口来完成：

```protobuf
service AtomicBroadcast {
    // 1. 客户端使用此接口向排序节点发送交易进行排序
    rpc Broadcast(stream common.Envelope) returns (stream BroadcastResponse) {}

    // 2. Peer节点/客户端使用此接口订阅或请求新生成的区块
    rpc Deliver(stream common.Envelope) returns (stream DeliverResponse) {}
}
```

* **Broadcast 操作**：是一个基于 gRPC Stream（双向数据流）的单向接收端接口。排序节点接收到源源不断的 `Envelope`（封装了带有全套签名的交易），只负责校验发起者是否有向该通道写入的权限（不做深层次交易内容检查），校验通过后将包暂存进内存池等待切分打包并排序。
* **Deliver 操作**：该接口允许具有合法读取通道权限的消费者（大部分是 Leader Peer，或者是某些为了更新 UI 需要监听交易状态的客户端）拉取指定范围的区块数据流。当配置为持续监听状态时，一旦新的区块被 Orderer 创建，Orderer 便会通过此数据流通告这些监听者。

这两个高度精简的接口保证了排序层可以无缝地替换背后的共识引擎（从过去的 Kafka 到现代标准的 Raft），而不必去修改外围组件的任何通信协议层。
