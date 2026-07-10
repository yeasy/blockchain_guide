## 如何升级版本

经典 Fabric 支持在满足版本前置条件时滚动升级组件，但不能据此假定任意历史版本都可直接跨级升级，也不能把 Fabric-X 当作经典 Fabric 的下一主版本。

网络升级主要包括对如下资源进行升级：

* 核心组件：包括 Peer、Orderer、CA 等核心程序；
* 能力配置：更新通道配置中支持的能力集合版本号，以启动新的特性；
* 第三方资源：包括依赖的 CouchDB，以及旧网络中可能仍存在但升级到 Fabric 3.x 前必须迁移掉的 Kafka 排序服务。

### 支持状态与迁移矩阵

起点 | 目标 | 路径性质 | 必做事项 | 结论
--- | --- | --- | --- | ---
经典 Fabric v1.4.x | 经典 Fabric 2.5 LTS | 官方支持滚动升级 | 备份；修正链码 shim/运行时兼容性；对 Peer 执行 `peer node upgrade-dbs`；节点升级完成后再启用 `V2_0`/`V2_5` 能力并迁移到 v2 链码生命周期 | 可升级；不要跳过 v1.4 到 v2.x 的专项步骤
早于 v1.4 的经典 Fabric 1.x | 经典 Fabric 2.5 LTS | 没有当前文档覆盖的直接跨级路径 | 先按对应历史版本文档升级到 v1.4.x，或经业务停机、账本导出与重新部署评估后重建 | 不应直接套用 v1.4→2.5 命令
经典 Fabric 2.x/2.5 | 经典 Fabric 3.x | 官方支持滚动升级组件 | 在升级二进制前移除 Solo/Kafka、迁移全部链码到 v2 生命周期、移除系统通道并启用 Channel Participation API；配置组织级 `OrdererEndpoints` | 前置条件满足后，先滚动升级 Orderer/Peer，再启用 `V3_0` 能力
经典 Fabric 2.5/3.x | Fabric-X 1.x | 新平台迁移，不是经典 Fabric 原地升级 | 建立独立试验网络；重新验证客户端、智能合约、MSP/策略、状态与历史数据导入、性能及灾备；设计双写、停机切换或回滚方案 | **不存在“替换镜像后沿用原账本”的官方直接升级路径**

官方依据如下：

* [Fabric 2.5 文档](https://hyperledger-fabric.readthedocs.io/en/release-2.5/whatsnew.html)将 2.5.x 定位为 LTS，并说明 v1.4.x 可直接滚动升级到 2.5.x；[2.5 升级注意事项](https://hyperledger-fabric.readthedocs.io/en/release-2.5/upgrade_to_newest_version.html)给出了数据库和链码生命周期转换步骤。
* [Fabric 3.x 升级注意事项](https://hyperledger-fabric.readthedocs.io/en/latest/upgrade_to_newest_version.html)明确移除了 Solo、Kafka 和 v1.x 链码生命周期；[能力说明](https://hyperledger-fabric.readthedocs.io/en/latest/capabilities_concept.html)要求仍使用系统通道的网络先迁移到 Channel Participation API。
* LF Decentralized Trust 将[经典 Fabric 和 Fabric-X 描述为两种不同实现](https://www.lfdecentralizedtrust.org/projects/fabric)。[Fabric-X 路线图](https://www.lfdecentralizedtrust.org/blog/the-hyperledger-fabric-x-roadmap)说明其使用新的编程、排序和提交架构，并仍在探索经典 Fabric 链码兼容性。因此，Fabric-X 的 API 或账本格式兼容点不能被解释为经典 Fabric 的升级承诺。

### 能力类型

为了避免网络多个节点运行不同版本组件时出现分叉风险，自 1.1.0 版本起在通道配置中引入了能力（Capabilities），标记节点应当支持和启用的特性。如果某节点程序版本低于能力要求则无法加入或自动退出；同时通道内高版本的节点程序在提交校验时只启用指定的特性集合检查（可参考 core/handlers/validation/builtin）。

目前能力分为三种类型，分别管理不同范围，如下表所示。

类型 | 功能  | 配置路径
--- | --- | ---
通道（Channel）能力| 通道整体相关能力，排序和 Peer 节点都得满足 | /Channel/Capabilities
排序（Orderer）能力| 排序服务能力，只与排序节点有关 | /Channel/Orderer/Capabilities
应用（Application）能力| 应用相关能力，只与 Peer 节点有关 | /Channel/Application/Capabilities

如果要启用相应的能力，需要修改通道配置内对应配置。例如，用户可以指定通道能力为 v1.1.0，排序能力为 v1.1.0 模式下，而应用能力为 v1.3.0。此时，只有不低于 v1.1.0 版本（满足通道和排序能力的较大者）的排序节点，以及不低于 v1.3.0 版本（满足通道和应用能力的较大者）的 Peer 节点可以支持该通道。同时，即使排序节点和 Peer 节点程序版本更新（如 v2.x），仍然只会启用指定的能力集合。

需要注意能力配置只能调整到更新版本而不应回退，例如可以将能力模式 v1.3.0 更新为更高版本的 v1.4.0，反之无意义。这是因为旧版本的节点即便加入到通道内，仍然无法正常处理其中新版本启用阶段的交易。

其中，各能力集合的版本和内容（可参考 common/capabilities）如下表所示，注意并不与程序版本一致。

能力版本 | 起始程序版本 | 类型 | 能力内容
--- | --- | --- | ---
ChannelV1_1 | v1.1.0 | 通道 | 仅供标记，程序版本为 1.1.0+
ChannelV1_3 | v1.3.0 | 通道 | 支持 idemix
OrdererV1_1 | v1.1.0 | 排序 | 重新提交和身份超时检查
OrdererV2_0 | v2.0.0 | 排序 | 排序服务支持从 Kafka 切换到 Raft
ApplicationV1_1 |v1.1.0 | 应用 | 禁止区块内重复交易Id
ApplicationV1_2 |v1.2.0 | 应用 | 正式支持私有数据，支持升级私有数据成员组配置，细粒度的通道资源访问控制（ACL）
ApplicationV1_3 |v1.3.0 | 应用 | 支持基于键值的背书
ApplicationV2_0 |v2.0.0 | 应用 | 新的链码生命周期管理
ApplicationV2_0 |v2.0.0 | 应用 | 支持链码操作

### 推荐升级步骤

#### 升级排序服务

对于不改变排序模式的情况下，升级较为简单。

逐个停止排序节点，并备份本地数据，包括身份文件、账本数据、配置文件等。

升级排序服务程序。重新启动并检查是否工作正常，如获取区块、发送交易等。

*注：自 Fabric v2.3 起，Kafka 排序服务已被正式废弃；升级到 v3.x 之前，应先完成从 Kafka 到 Raft 的迁移。对于新部署网络，当前更稳妥的默认建议仍是优先使用 **Raft**；BFT 更适合在确有拜占庭容错需求时专项评估，而不应简单视为 v3.0+ 的默认首选。*

#### 升级 Peer 节点

逐个停止 Peer 节点，并备份本地数据，包括身份文件、账本数据、链码包、配置文件等。

升级 Peer 程序。重新启动并检查是否工作正常，如查询信息、发送交易提案等。

链码包如果之前有引入旧的第三方库或者 Shim 包，或者需要启用新的 API，则还要执行链码升级操作。

#### 升级 CA 服务

停止 Fabric-CA 服务，备份数据库。

升级 fabric-ca 程序，重新启动并检查是否工作正常，如获取根证书。

```bash
$ fabric-ca-client getcacert -u https://<fabric-ca-server>:7054 --tls.certfiles tls-cert.pem
```

#### 升级通道配置

按照新的格式发送通道更新请求，特别是修改对应能力域值为新的版本。

对于 Fabric 2.x 中仍保留系统通道的旧网络，升级到 Fabric 3.x 前必须先按照官方流程迁移到 Channel Participation API 并移除系统通道；Fabric 3.x 不再支持系统通道，也不支持通过系统通道创建应用通道。

能力配置应按通道逐步更新：先确认排序节点、Peer 节点二进制版本均已满足目标能力，再分别更新各应用通道的 Channel、Orderer 和 Application 能力。更新后测试网络功能，如获取新区块、提交配置更新、调用链码等是否正常。

#### 升级第三方组件

包括 CouchDB 等第三方组件，升级之前最好备份数据文件。

CouchDB 版本自 1.x 版本可以很容易升级到高版本，具体操作可以参考项目文档：https://docs.couchdb.org/en/stable/install/upgrading.html。

如果网络仍停留在 Kafka 排序服务上，更稳妥的做法不是继续升级 Kafka，而是尽快按照官方迁移文档转向 Raft，再继续进行 Fabric 主版本升级。
