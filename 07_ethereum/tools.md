## 相关工具

### 客户端和开发库

以太坊社区提供了多种语言实现的客户端，支持标准的 JSON-RPC 协议。合并后，主网节点通常由执行层客户端（Execution Client）和共识层客户端（Consensus Client）共同组成。

**执行层客户端（Execution Client）：**
*   [Geth (go-ethereum)](https://github.com/ethereum/go-ethereum)：Go 语言实现，最主流的客户端；
*   [Nethermind](https://github.com/NethermindEth/nethermind)：C# 语言实现，性能优异；
*   [Besu](https://github.com/hyperledger/besu)：Java 语言实现，Hyperledger 项目之一，适合企业级应用；
*   [Erigon](https://github.com/ledgerwatch/erigon)：Go 语言实现，专注于存储效率和同步速度。

**共识层客户端（Consensus Client）：**
*   [Prysm](https://github.com/prysmaticlabs/prysm)：Go 语言实现；
*   [Lighthouse](https://github.com/sigp/lighthouse)：Rust 语言实现；
*   [Teku](https://github.com/ConsenSys/teku)：Java 语言实现。

**开发库：**
*   [ethers.js](https://docs.ethers.org/v6/)：轻量级且功能强大的 JavaScript/TypeScript 库，当前 v6 API 使用 `bigint`、`ethers.parseEther()` 等现代接口；
*   [web3.js](https://github.com/web3/web3.js)：历史悠久的 JavaScript 库，常见于 legacy 项目；
*   [viem](https://github.com/wagmi-dev/viem)：基于 TypeScript 构建的新一代库，专注于性能。

### 以太坊钱包

钱包是用户进入 Web3 世界的入口。

*   [MetaMask](https://metamask.io/)：浏览器插件钱包的行业标准，支持几乎所有 EVM 兼容链。
*   [Rabby Wallet](https://rabby.io/)：专为 DeFi 用户设计，提供交易模拟和风险扫描功能，体验更佳。
*   [Frame](https://frame.sh/)：专注于隐私和 macOS 原生体验的桌面钱包。

### 开发框架与 IDE

现代以太坊开发工具栈已经发生了巨大变化。

*   [Hardhat](https://hardhat.org/)：主流以太坊开发环境，支持 TypeScript、Solidity 测试、Hardhat Network、本地调试和插件生态。
*   [Foundry](https://getfoundry.sh/)：基于 Rust 编写的高性能开发框架，包含 Forge、Cast、Anvil 和 Chisel，适合 Solidity 原生测试和脚本化工作流。
*   [Remix](https://remix.ethereum.org)：基于浏览器的 IDE，无需安装，适合快速原型开发和教学。
*   [Ganache](https://archive.trufflesuite.com/ganache/)：Truffle/Ganache 已 sunset，代码和文档作为归档保留；新项目优先使用 Hardhat Network 或 Foundry Anvil。

### 本地测试链

*   [Hardhat Network](https://hardhat.org/docs/guides/hardhat-node)：Hardhat 内置的本地开发链，可通过 `npx hardhat node` 暴露 JSON-RPC，适合与 Hardhat 测试、调试、fork 和脚本集成。
*   [Anvil](https://getfoundry.sh/anvil/overview)：Foundry 提供的快速本地 Ethereum 节点，适合独立本地链、fork 测试和与 Forge/Cast 配合使用。

### 网站资源

*   [Etherscan](https://etherscan.io/)：最权威的区块浏览器，查看所有链上交易、合约代码和账户状态。
*   [DefiLlama](https://defillama.com/)：最全面的 DeFi TVL 和数据分析平台。
*   [Dune Analytics](https://dune.com/)：强大的链上数据可视化分析平台，社区贡献了大量 Dashboard。
