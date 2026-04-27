## 使用智能合约

以太坊智能合约开发通常不再直接使用 Geth 控制台、旧式内置 JS 对象或 PoW 私链挖矿流程。自 The Merge 之后，Geth 是执行层客户端，完整节点需要配合共识层客户端；旧教程中的本地挖矿和 JS 控制台部署流程不适合当前开发。

本节采用 Hardhat 本地开发网络演示合约编写、编译、部署和调用。它更接近当前 DApp 开发流程，也便于切换到 Sepolia 等测试网。

### 初始化项目

先安装 Node.js 22 或更高版本，然后创建 Hardhat 项目。

```bash
mkdir hardhat-example
cd hardhat-example
npx hardhat --init
npx hardhat test
```

按默认模板创建项目即可。默认项目会包含 `contracts/`、`test/`、`ignition/modules/` 和 `hardhat.config.ts` 等文件。

### 创建合约

新建 `contracts/TestContract.sol`。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

contract TestContract {
    function multiply(uint256 a) public pure returns (uint256) {
        return a * 7;
    }
}
```

编译合约。

```bash
npx hardhat compile
```

Hardhat 会根据项目配置管理 Solidity 编译器，不需要手工把 `solc --bin` 和 `solc --abi` 输出复制到控制台。

### 本地部署

使用 Hardhat Ignition 描述部署过程。新建 `ignition/modules/TestContract.ts`。

```typescript
import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

export default buildModule("TestContractModule", (m) => {
  const testContract = m.contract("TestContract");

  return { testContract };
});
```

如果只想在一次性模拟网络中验证部署，直接运行：

```bash
npx hardhat ignition deploy ignition/modules/TestContract.ts
```

如果需要一个持续运行的本地 JSON-RPC 节点，在第一个终端运行：

```bash
npx hardhat node
```

在第二个终端部署到该本地节点：

```bash
npx hardhat ignition deploy ignition/modules/TestContract.ts --network localhost
```

输出中会给出合约地址。切换到 Sepolia 等测试网时，应在 `hardhat.config.ts` 中配置 RPC URL 和部署账户，并把私钥放在环境变量或 Hardhat 配置变量中，不要写入代码仓库。

### 调用合约

调试阶段可以用 Hardhat 控制台在本地模拟网络中部署并调用。

```bash
npx hardhat console
```

```javascript
> const { viem } = await network.create()
> const c = await viem.deployContract("TestContract")
> await c.read.multiply([10n])
70n
```

如果已经部署到本地节点或测试网，可使用前端、脚本或命令行工具通过合约地址和 ABI 调用。Foundry 的 `cast` 也常用于快速读写合约。

```bash
cast call <CONTRACT_ADDRESS> "multiply(uint256)(uint256)" 10 --rpc-url http://127.0.0.1:8545
```

### 使用 Geth 私有网络的场景

如果目标是测试多节点执行层/共识层行为，而不是普通合约开发，应参考 Geth 官方的 Kurtosis 私有网络文档。当前 Geth 私网需要同时启动执行客户端和共识客户端，不应再按旧版挖矿或 PoA 单客户端教程搭建。
