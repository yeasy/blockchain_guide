# Web3 概念与架构

互联网的发展经历的三次重要迭代。Web3 是对当前互联网模式（Web2）的一次深刻变革，其核心目标是将数字资产的所有权和数据的控制权从平台巨头手中回归给用户。

## 从 Web1 到 Web3

### Web1 (1990-2005): Read (只读)
Web1 是静态信息的展示平台。如早期的 Yahoo、新浪门户。
*   **特征**：去中心化协议（HTTP, SMTP），但主要由公司制作内容，用户被动消费。
*   **模式**：平台 -> 用户。

### Web2 (2005-2020): Read-Write (读写)
Web2 是社交网络和平台经济的时代。如 Facebook, Twitter, 微信。
*   **特征**：用户创造内容（UGC），交互性强。但数据和价值被中心化平台垄断和变现。
*   **模式**：用户 <-> 平台 <-> 用户。
*   **痛点**：隐私泄露、数据垄断、算法杀熟、平台审查。

### Web3 (2020-至今): Read-Write-Own (读写拥有)
Web3 是基于区块链的价值互联网。
*   **特征**：基于密码学和分布式账本，用户拥有自己的身份（私钥）、数据和资产（Token/NFT）。
*   **核心理念**：**无需许可 (Permissionless)**、**抗审查 (Censorship Resistant)**、**用户主权 (User Sovereignty)**。

## Web3 技术架构

Web3 的技术栈通常被划分为四层架构：

### 1. 基础设施层 (Layer 0 & Layer 1)
这是 Web3 的物理根基和信任之源。
*   **Layer 0**：负责异构链之间的通信和安全性共享，如 **Polkadot** (Relay Chain) 和 **Cosmos** (IBC)。
*   **Layer 1**：基础公链，负责共识机制、账本维护和智能合约执行。如 **Ethereum**、**Solana**、**Bitcoin**。

### 2. 扩展层 (Layer 2)
为了解决 Layer 1 的性能瓶颈而生（详见以太坊章节）。
*   主要技术：**Rollups** (Optimism, ZK-Rollups)。
*   作用：提供高吞吐和低成本的执行环境。

### 3. 中间件与协议层
连接底层区块链与上层应用的关键组件。
*   **去中心化存储**：IPFS, Arweave, Filecoin。解决“区块链存不下图片和视频”的问题。
*   **预言机 (Oracle)**：Chainlink。将现实世界的数据（如币价、天气）可信地喂给智能合约。
*   **数字身份 (DID)**：ENS (Ethereum Name Service)。将复杂的哈希地址解析为人类可读的名称（如 alice.eth）。
*   **索引查询**：The Graph。方便 DApp 快速查询链上复杂数据。

### 4. 应用层 (DApp)
用户直接交互的界面。
*   **DeFi**：Uniswap, Aave。
*   **NFT Market**：OpenSea, Blur。
*   **GameFi**：Axie Infinity。
*   **SocialFi**：Lens Protocol。

## Web3 账户体系：钱包

在 Web3 中，**钱包 (Wallet)** 不仅仅是管理资金的工具，更是用户的**通用数字身份**。
*   **非托管 (Non-custodial)**：私钥完全由用户自己掌管，“Not your keys, not your coins”。
*   **单点登录 (Sign-in with Ethereum)**：使用通过一个钱包地址，即可登录所有 Web3 应用，无需重复注册账号，数据也可以跨应用携带。

这一架构彻底打破了 Web2 时代不同应用之间的数据孤岛（Walled Gardens），为创新提供了无限可能。
