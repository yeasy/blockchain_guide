# 管理节点

Fabric 网络由多种不同类型的节点组成，每种节点负责不同的功能。了解如何管理这些节点对于维护一个健康的区块链网络至关重要。

## 节点类型

### Peer 节点

Peer 节点是 Fabric 网络的核心组件，主要负责：
*   **维护账本**：存储区块链数据和世界状态（State DB）。
*   **执行链码**：运行智能合约并生成交易提案的背书。
*   **参与共识**：验证交易并提交区块。

根据角色不同，Peer 可以分为：
*   **背书节点 (Endorser)**：执行链码并对结果签名背书。
*   **提交节点 (Committer)**：验证交易并将区块写入账本（所有 Peer 都是 Committer）。
*   **锚节点 (Anchor Peer)**：用于跨组织的 Gossip 通信发现。

### Orderer 节点

Orderer 节点负责交易排序和区块生成：
*   接收来自客户端的交易。
*   对交易进行全局排序。
*   将交易打包成区块并分发给 Peer 节点。
*   Fabric 2.x 默认使用 **Raft** 共识协议。

### CA 节点

证书授权（Certificate Authority）节点负责身份管理：
*   为网络参与者签发数字证书。
*   管理证书的注册、吊销等生命周期。

## 常用管理操作

### 查看节点状态

```bash
# 查看 Peer 节点加入的通道
peer channel list

# 查看 Peer 节点安装的链码
peer lifecycle chaincode queryinstalled
```

### 节点日志管理

```bash
# 查看 Docker 容器日志
docker logs peer0.org1.example.com

# 调整日志级别
peer node logging setlevel gossip warning
```

### Orderer 集群管理

Raft 模式下的 Orderer 集群支持动态添加和移除节点。通过更新系统通道配置可以实现节点的增减，无需停机。

## 最佳实践

1.  **资源隔离**：生产环境中，Peer、Orderer 和 CA 应部署在不同的物理机或虚拟机上。
2.  **备份策略**：定期备份 Peer 节点的账本数据和配置文件。
3.  **监控告警**：使用 Prometheus + Grafana 监控节点健康状态。
4.  **证书管理**：设置证书到期提醒，避免因证书过期导致网络故障。
