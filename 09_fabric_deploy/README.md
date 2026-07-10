# Fabric 安装与部署

**纸上得来终觉浅，绝知此事要躬行。**

作为被广泛应用的联盟链项目，Fabric 吸取了来自科技界和金融界的最新成果，提供面向企业场景的开放分布式账本平台实现。

本章将带领读者动手实践，学习如何从源码进行编译、安装 Fabric ，如何使用官方容器镜像，以及在多服务器环境下部署一个典型的 Fabric 网络，同时，还将介绍通道的相关实践操作。

## 版本选择与迁移边界

Hyperledger 项目当前同时提供经典 Fabric 与 Fabric-X 两种不同实现。经典 Fabric 的 2.5.x 是长期支持（LTS）分支，3.x 是引入 `V3_0` 通道能力和 SmartBFT 等特性的后续主版本；Fabric-X 则面向高吞吐数字资产场景，采用 ARMA 排序和解耦的提交架构。它们不是同一条版本号序列。

现状或目标 | 官方定位 | 推荐路径 | 进入下一阶段前的关键检查
--- | --- | --- | ---
经典 Fabric 1.x | 已停止维护的历史分支；仅 v1.4.x 有官方直达 2.5.x 的滚动升级说明 | 先备份；v1.4.x 可按 2.5 文档升级，早于 v1.4 的网络应先按对应历史文档到 v1.4.x，或评估重建网络 | 更新旧链码依赖，执行 Peer 数据库升级，准备 `V2_0` 链码生命周期
经典 Fabric 2.5 LTS | 经典 Fabric 的生产 LTS 分支 | 新部署优先选择最新 2.5.x 补丁；现有 2.x 可滚动升级到 2.5.x | 先升级二进制，再按需启用 `V2_5` 应用能力
经典 Fabric 3.x | 经典 Fabric 的当前功能分支 | 从 2.x 滚动升级；先在 2.x 阶段清除已移除功能 | Kafka/Solo 改为 Raft、旧链码生命周期改为 v2 生命周期、移除系统通道并使用 Channel Participation API、配置组织级 `OrdererEndpoints`
Fabric-X 1.x | 与经典 Fabric 并列的、面向高吞吐数字资产的独立实现 | 作为新平台单独部署和验收，规划应用、身份、治理与数据迁移 | **不是**从经典 Fabric 2.5/3.x 直接替换二进制即可完成的升级；API 或区块格式兼容性不等于受支持的原地升级路径

判断依据以官方资料为准：[Fabric 2.5 LTS 说明](https://hyperledger-fabric.readthedocs.io/en/release-2.5/whatsnew.html)、[从 1.4/2.x 升级到 2.5 的注意事项](https://hyperledger-fabric.readthedocs.io/en/release-2.5/upgrade_to_newest_version.html)、[进入 Fabric 3.x 的前置条件](https://hyperledger-fabric.readthedocs.io/en/latest/upgrade_to_newest_version.html)、[Fabric 与 Fabric-X 的双实现定位](https://www.lfdecentralizedtrust.org/projects/fabric)及 [Fabric-X 路线图](https://www.lfdecentralizedtrust.org/blog/the-hyperledger-fabric-x-roadmap)。具体升级步骤见[运维章节的升级矩阵](../10_fabric_op/upgrade.md)。
