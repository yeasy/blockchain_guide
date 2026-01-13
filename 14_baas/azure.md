## 微软 Azure 区块链方案

Azure 是微软推出的云计算平台，向用户提供开放的 IaaS 和 PaaS 服务。

微软曾提供 Azure Blockchain Service，但该服务已于 2021 年 9 月 10 日正式退休。目前，微软推荐用户使用基于 Azure 的 Quorum Blockchain Service (QBS) 或直接在虚拟机/容器服务上自行部署区块链网络。

用户可以在 Azure 应用市场（https://azuremarketplace.microsoft.com/en-us/marketplace/apps）中搜索 “blockchain” 关键字查看最新的区块链解决方案，包括：
*   **Quorum Blockchain Service (QBS)**：由 ConsenSys 管理的完全托管服务；
*   **Hyperledger Fabric on Azure Kubernetes Service (AKS)**：提供在 AKS 上部署 Fabric 网络的模板；
*   **Avalanche 验证节点** 等。

![Azure 上的区块链服务](_images/azure_marketplace.png)

### 替代方案：Hyperledger Fabric on AKS

对于希望在 Azure 上运行 Hyperledger Fabric 的用户，微软提供了基于 Azure Kubernetes Service (AKS) 的部署模板。

该方案具有以下特点：
*   **完全控制**：用户对网络拥有完全的控制权；
*   **云原生**：利用 Kubernetes 的编排能力管理 Fabric 容器；
*   **集成**：可与 Azure Active Directory (AAD)、Azure Monitor 等服务集成。

用户可以通过 Azure CLI 或 ARM 模板快速拉起一套 Fabric 网络，通常包括 Orderer 组织和 Peer 组织。部署完成后，用户可以通过标准的 Fabric CLI 或 SDK 与网络进行交互。

*注：由于云厂商服务策略调整频繁，建议在选择 BaaS 服务前查阅最新的官方文档。*