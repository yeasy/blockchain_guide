## 智能合约开发实验指南

本章提供了一系列循序渐进的实验教程，帮助开发者从零开始掌握以太坊智能合约的开发、部署和交互。所有实验都可以在本地环境中完成，无需花费真实的以太币。

### 实验 0：搭建开发环境

#### 0.1 使用 Remix IDE（无需本地环装）

Remix 是官方提供的浏览器 IDE，无需任何安装，立即可用。

**访问地址**：https://remix.ethereum.org

**优点**：
- 无需安装，打开浏览器即可
- 内置编译器和调试器
- 支持直接连接本地区块链和测试网

**创建第一个合约**：

1. 打开 Remix IDE
2. 在左侧文件面板中，点击 “+” 创建新文件
3. 命名为 `HelloWorld.sol`
4. 输入以下代码：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract HelloWorld {
    string public message = "Hello, Blockchain!";

    function updateMessage(string memory newMessage) public {
        message = newMessage;
    }

    function getMessage() public view returns (string memory) {
        return message;
    }
}
```

5. 点击左侧 “Solidity Compiler” 图标（看起来像一个积木）
6. 选择编译器版本 `0.8.0` 或以上
7. 点击 “Compile HelloWorld.sol”
8. 如果没有错误，你应该看到绿色的 “Compile” 按钮

#### 0.2 Hardhat 本地开发环境搭建

Hardhat 是以太坊开发的主流框架，提供测试、调试、部署和本地链模拟能力。本实验使用 Hardhat 3 的 TypeScript + Mocha + ethers.js 工具链。

**安装步骤**：

```bash
# 创建项目目录
mkdir ethereum-learning
cd ethereum-learning

# 初始化 Node.js 项目
npm init -y

# 安装 Hardhat（Hardhat 3 要求 Node.js 22 或以上）
npm install --save-dev hardhat

# 初始化 Hardhat 项目
npx hardhat --init

# 选择 "A TypeScript Hardhat project using Mocha and Ethers.js"
# 按提示安装 @nomicfoundation/hardhat-toolbox-mocha-ethers 等依赖
```

**项目结构**：

```text
ethereum-learning/
├── contracts/
│   ├── Counter.sol           # 示例合约
│   └── Counter.t.sol         # Solidity 测试（可选）
├── ignition/
│   └── modules/
│       └── Counter.ts        # Hardhat Ignition 部署模块
├── test/
│   └── Counter.ts            # TypeScript 测试文件
├── scripts/
│   └── deploy.ts             # 脚本文件
├── hardhat.config.ts         # Hardhat 配置
└── package.json
```

#### 0.3 验证环境

```bash

# 编译合约
npx hardhat compile

# 运行测试
npx hardhat test

# 启动本地区块链
npx hardhat node
```

如果以上命令都成功执行，恭喜！开发环境已搭建完成。

### 实验 1：部署和调用简单合约

**目标**：理解合约的基本结构、部署过程和状态查询。

**合约代码** (`contracts/SimpleStorage.sol`)：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleStorage {
    // 状态变量：永久存储在区块链上
    uint256 public storedValue;

    // 事件：用于链外应用监听
    event ValueChanged(uint256 indexed newValue, address indexed changer);

    // 函数：修改状态（花费 Gas）
    function setValue(uint256 newValue) public {
        storedValue = newValue;
        emit ValueChanged(newValue, msg.sender);
    }

    // 函数：查询状态（不花费 Gas）
    function getValue() public view returns (uint256) {
        return storedValue;
    }

    // 函数：计算但不修改状态
    function add(uint256 a, uint256 b) public pure returns (uint256) {
        return a + b;
    }
}
```

**编写测试** (`test/SimpleStorage.test.ts`)：

```typescript
import { expect } from "chai";
import { network } from "hardhat";

describe("SimpleStorage", function () {
    let ethers: any;
    let storage: any;
    let owner: any;

    // 在每个测试前执行
    beforeEach(async function () {
        ({ ethers } = await network.create());
        [owner] = await ethers.getSigners();
        storage = await ethers.deployContract("SimpleStorage");
        await storage.waitForDeployment();
    });

    describe("初始状态", function () {
        it("初始值应该为 0", async function () {
            expect(await storage.getValue()).to.equal(0n);
        });

        it("storedValue 应该可以直接访问", async function () {
            expect(await storage.storedValue()).to.equal(0n);
        });
    });

    describe("setValue 函数", function () {
        it("应该能够设置新值", async function () {
            await storage.setValue(42);
            expect(await storage.getValue()).to.equal(42n);
        });

        it("应该发出 ValueChanged 事件", async function () {
            // 监听事件
            await expect(storage.setValue(100))
                .to.emit(storage, "ValueChanged")
                .withArgs(100n, owner.address);
        });

        it("每次设置都应该覆盖旧值", async function () {
            await storage.setValue(10);
            expect(await storage.getValue()).to.equal(10n);

            await storage.setValue(20);
            expect(await storage.getValue()).to.equal(20n);
        });
    });

    describe("add 函数", function () {
        it("应该正确相加", async function () {
            expect(await storage.add(5, 3)).to.equal(8n);
        });

        it("应该处理大数字", async function () {
            const big = 10n ** 18n;
            expect(await storage.add(big, big)).to.equal(big * 2n);
        });
    });

    describe("Gas 消耗分析", function () {
        it("setValue 是一个状态修改交易", async function () {
            const tx = await storage.setValue(50);
            const receipt = await tx.wait();
            if (receipt === null) throw new Error("交易未确认");
            console.log("setValue Gas 消耗:", receipt.gasUsed.toString());
            // 预期: ~43,000-44,000
        });

        it("getValue 是一个只读调用", async function () {
            // 不会消耗任何 Gas
            const result = await storage.getValue();
            expect(result).to.equal(0n);
        });
    });
});
```

**运行测试**：

```bash
npx hardhat test test/SimpleStorage.test.ts
```

**预期输出**：

```text
SimpleStorage
  初始状态
    ✓ 初始值应该为 0
    ✓ storedValue 应该可以直接访问
  setValue 函数
    ✓ 应该能够设置新值
    ✓ 应该发出 ValueChanged 事件
    ✓ 每次设置都应该覆盖旧值
  add 函数
    ✓ 应该正确相加
    ✓ 应该处理大数字
  Gas 消耗分析
    ✓ setValue 是一个状态修改交易
    setValue Gas 消耗: 43210
    ✓ getValue 是一个只读调用

9 passing
```

### 实验 2：代币合约（ERC-20）

**目标**：理解标准合约接口、代币的基本操作（转账、授权）。

**合约代码** (`contracts/SimpleToken.sol`)：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleToken {
    string public name = "Simple Token";
    string public symbol = "STK";
    uint8 public decimals = 18;
    uint256 public totalSupply;

    // 账户余额
    mapping(address => uint256) public balanceOf;

    // 授权额度: owner => spender => amount
    mapping(address => mapping(address => uint256)) public allowance;

    // 事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    // 构造函数：初始化代币供应
    constructor(uint256 initialSupply) {
        totalSupply = initialSupply * 10 ** uint256(decimals);
        balanceOf[msg.sender] = totalSupply;
    }

    // 转账函数
    function transfer(address to, uint256 value) public returns (bool) {
        require(to != address(0), "Invalid address");
        require(balanceOf[msg.sender] >= value, "Insufficient balance");

        balanceOf[msg.sender] -= value;
        balanceOf[to] += value;

        emit Transfer(msg.sender, to, value);
        return true;
    }

    // 授权函数：允许 spender 花费你的代币
    function approve(address spender, uint256 value) public returns (bool) {
        allowance[msg.sender][spender] = value;
        emit Approval(msg.sender, spender, value);
        return true;
    }

    // 代理转账：spender 转账 owner 的代币
    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        require(from != address(0), "Invalid address");
        require(to != address(0), "Invalid address");
        require(balanceOf[from] >= value, "Insufficient balance");
        require(allowance[from][msg.sender] >= value, "Insufficient allowance");

        balanceOf[from] -= value;
        balanceOf[to] += value;
        allowance[from][msg.sender] -= value;

        emit Transfer(from, to, value);
        return true;
    }
}
```

**测试代码** (`test/SimpleToken.test.ts`)：

```typescript
import { expect } from "chai";
import { network } from "hardhat";

describe("SimpleToken", function () {
    let ethers: any;
    let token: any;
    let owner: any, alice: any, bob: any;

    beforeEach(async function () {
        ({ ethers } = await network.create());
        [owner, alice, bob] = await ethers.getSigners();

        token = await ethers.deployContract("SimpleToken", [1000]); // 初始供应 1000 tokens
        await token.waitForDeployment();
    });

    describe("初始化", function () {
        it("应该有正确的总供应", async function () {
            const expected = ethers.parseEther("1000");
            expect(await token.totalSupply()).to.equal(expected);
        });

        it("拥有者应该获得所有初始代币", async function () {
            const balance = await token.balanceOf(owner.address);
            expect(balance).to.equal(await token.totalSupply());
        });
    });

    describe("转账", function () {
        it("应该能够转账代币", async function () {
            const amount = ethers.parseEther("100");
            await token.transfer(alice.address, amount);

            expect(await token.balanceOf(alice.address)).to.equal(amount);
            expect(await token.balanceOf(owner.address)).to.equal(
                ethers.parseEther("900")
            );
        });

        it("应该发出 Transfer 事件", async function () {
            const amount = ethers.parseEther("50");
            await expect(token.transfer(bob.address, amount))
                .to.emit(token, "Transfer")
                .withArgs(owner.address, bob.address, amount);
        });

        it("不应该允许转账超过余额", async function () {
            const tooMuch = ethers.parseEther("2000");
            await expect(
                token.transfer(alice.address, tooMuch)
            ).to.be.revertedWith("Insufficient balance");
        });
    });

    describe("授权和代理转账", function () {
        it("应该能够授权", async function () {
            const amount = ethers.parseEther("100");
            await token.approve(alice.address, amount);

            expect(await token.allowance(owner.address, alice.address)).to.equal(amount);
        });

        it("应该能够进行代理转账", async function () {
            const amount = ethers.parseEther("100");

            // owner 授权 alice 花费 100 tokens
            await token.approve(alice.address, amount);

            // alice 代表 owner 转账给 bob
            await token.connect(alice).transferFrom(owner.address, bob.address, amount);

            expect(await token.balanceOf(bob.address)).to.equal(amount);
            expect(await token.allowance(owner.address, alice.address)).to.equal(0n);
        });

        it("不应该允许转账超过授权额度", async function () {
            const approved = ethers.parseEther("50");
            const attempted = ethers.parseEther("100");

            await token.approve(alice.address, approved);

            await expect(
                token.connect(alice).transferFrom(owner.address, bob.address, attempted)
            ).to.be.revertedWith("Insufficient allowance");
        });
    });
});
```

**运行测试**：

```bash
npx hardhat test test/SimpleToken.test.ts
```

### 实验 3：使用 Hardhat Network 或 Anvil 搭建本地测试网

Ganache 及其旧 CLI 已经 sunset，不再作为新项目的默认本地链。新项目优先使用 Hardhat Network；如果需要独立的高性能本地节点，可使用 Foundry Anvil。

**Hardhat Network（推荐默认选择）**：

```bash
# 在当前 Hardhat 项目内启动本地 JSON-RPC 节点
npx hardhat node

# 在另一个终端连接这个节点运行脚本或测试
npx hardhat run scripts/interact.ts --network localhost
```

Hardhat Network 的特点：
- 默认提供 20 个测试账户，每个账户有充足测试 ETH
- 和 Hardhat 测试、调试、stack trace、console.log 深度集成
- 支持主网 fork、快照、时间调整等测试辅助能力

**Foundry Anvil（独立本地节点）**：

```bash
# 安装 Foundry 后，启动 Anvil
anvil

# 使用远程 RPC fork 主网或测试网
anvil --fork-url https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY
```

Anvil 的特点：
- 启动快，适合和 Foundry/Forge/Cast 一起使用
- 默认提供 10 个测试账户，每个账户有 10000 ETH
- 默认 RPC 地址通常为 `http://127.0.0.1:8545`
- 支持 fork、固定区块、自动出块和间隔出块

**在 Hardhat 中连接 Anvil**：

```typescript
// hardhat.config.ts
import { configVariable, defineConfig } from "hardhat/config";
import hardhatToolboxMochaEthers from "@nomicfoundation/hardhat-toolbox-mocha-ethers";

export default defineConfig({
    plugins: [hardhatToolboxMochaEthers],
    solidity: "0.8.28",
    networks: {
        anvil: {
            type: "http",
            url: "http://127.0.0.1:8545",
            chainId: 31337
        },
        hardhatMainnet: {
            type: "edr-simulated",
            chainType: "l1",
            forking: {
                url: configVariable("MAINNET_RPC_URL")
            }
        }
    }
});
```

**使用 Anvil 运行测试**：

```bash
# 终端 1：启动 Anvil
anvil

# 终端 2：运行 Hardhat 测试
npx hardhat test --network anvil
```

### 实验 4：在测试网部署合约

**选择测试网**：

- **Sepolia**（推荐）：以太坊应用和智能合约开发的默认测试网
- **Hoodi**：验证者、质押和协议升级测试网，不是普通 dApp 的默认选择
- **Polygon Amoy**：Polygon PoS 当前测试网，替代旧 Polygon PoS 测试网

**获取测试代币**：

```bash
# 访问 Sepolia 水龙头
# https://sepoliafaucet.com
# https://www.alchemy.com/faucets/ethereum-sepolia

# Polygon Amoy 测试 POL 可使用 Polygon 文档列出的第三方水龙头
# https://docs.polygon.technology/tools/gas/matic-faucet/
```

**配置网络** (`hardhat.config.ts`)：

```typescript
import { configVariable, defineConfig } from "hardhat/config";
import hardhatToolboxMochaEthers from "@nomicfoundation/hardhat-toolbox-mocha-ethers";

export default defineConfig({
    plugins: [hardhatToolboxMochaEthers],
    solidity: "0.8.28",
    networks: {
        sepolia: {
            type: "http",
            chainType: "l1",
            url: configVariable("SEPOLIA_RPC_URL"),
            accounts: [configVariable("SEPOLIA_PRIVATE_KEY")]
        },
        polygonAmoy: {
            type: "http",
            chainId: 80002,
            url: configVariable("POLYGON_AMOY_RPC_URL"),
            accounts: [configVariable("POLYGON_AMOY_PRIVATE_KEY")]
        }
    },
    verify: {
        etherscan: {
            apiKey: configVariable("ETHERSCAN_API_KEY")
        }
    }
});
```

Hardhat 3 的 `configVariable` 默认读取环境变量；敏感值也可以用 `hardhat-keystore` 加密保存。

```bash
# macOS/Linux 示例
export SEPOLIA_RPC_URL="https://eth-sepolia.g.alchemy.com/v2/YOUR_API_KEY"
export SEPOLIA_PRIVATE_KEY="your_private_key_here"
export MAINNET_RPC_URL="https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY"
export POLYGON_AMOY_RPC_URL="https://polygon-amoy.g.alchemy.com/v2/YOUR_API_KEY"
export POLYGON_AMOY_PRIVATE_KEY="your_private_key_here"
export ETHERSCAN_API_KEY="your_etherscan_api_key_here"

# 或使用 Hardhat keystore
npx hardhat keystore set SEPOLIA_RPC_URL
npx hardhat keystore set SEPOLIA_PRIVATE_KEY
npx hardhat keystore set MAINNET_RPC_URL
npx hardhat keystore set ETHERSCAN_API_KEY
```

**部署脚本** (`scripts/deploy.ts`)：

```typescript
import fs from "node:fs";
import { network } from "hardhat";

console.log("部署开始...");

const { ethers } = await network.create();

// 部署
const storage = await ethers.deployContract("SimpleStorage");
await storage.waitForDeployment();

const address = await storage.getAddress();
console.log(`SimpleStorage 已部署到: ${address}`);

// 保存地址以供后续使用
fs.writeFileSync(
    "deployed.json",
    JSON.stringify({ SimpleStorage: address }, null, 2)
);
```

**执行部署**：

```bash
npx hardhat run scripts/deploy.ts --network sepolia
```

**验证合约** (在 Etherscan 上公开代码)：

```bash
npx hardhat verify --network sepolia CONTRACT_ADDRESS "constructor arguments"
```

### 实验 5：与合约交互

**创建交互脚本** (`scripts/interact.ts`)：

```typescript
import fs from "node:fs";
import { network } from "hardhat";

const deployment = JSON.parse(fs.readFileSync("deployed.json", "utf8"));
const contractAddress = deployment.SimpleStorage;

const { ethers } = await network.create();
const [signer] = await ethers.getSigners();
console.log(`使用账户: ${signer.address}`);

// 连接到部署的合约
const storage = await ethers.getContractAt("SimpleStorage", contractAddress, signer);

// 查询初始值
let value = await storage.getValue();
console.log(`初始值: ${value}`);

// 设置新值（这会触发交易）
console.log("\n正在设置新值为 42...");
const tx = await storage.setValue(42);
console.log(`交易哈希: ${tx.hash}`);

// 等待交易确认
const receipt = await tx.wait();
if (receipt === null) throw new Error("交易未确认");
console.log(`交易已确认，区块号: ${receipt.blockNumber}`);

// 查询新值
value = await storage.getValue();
console.log(`新值: ${value}`);

// 添加数字
const result = await storage.add(10, 20);
console.log(`\n10 + 20 = ${result}`);
```

**执行交互**：

```bash
npx hardhat run scripts/interact.ts --network sepolia
```

### 实验 6：Gas 优化分析

**合约优化前后对比**：

```solidity
// 不优化的版本
contract Unoptimized {
    uint256 public counter;

    function increment() public {
        counter = counter + 1;  // SSTORE + SLOAD
    }
}

// 优化的版本
contract Optimized {
    uint256 public counter;

    function increment() public {
        unchecked {  // 0.8.0 起默认检查，如果确定不溢出可禁用
            ++counter;  // ++counter 比 counter+1 少一次 SLOAD
        }
    }
}
```

**测试脚本分析 Gas**：

```typescript
import { network } from "hardhat";

describe("Gas 优化", function () {
    it("对比 ++ 和 +1", async function () {
        const { ethers } = await network.create();

        const unopt = await ethers.deployContract("Unoptimized");
        const opt = await ethers.deployContract("Optimized");
        await unopt.waitForDeployment();
        await opt.waitForDeployment();

        const tx1 = await unopt.increment();
        const receipt1 = await tx1.wait();
        if (receipt1 === null) throw new Error("交易未确认");
        const gas1 = receipt1.gasUsed;

        const tx2 = await opt.increment();
        const receipt2 = await tx2.wait();
        if (receipt2 === null) throw new Error("交易未确认");
        const gas2 = receipt2.gasUsed;

        const saved = gas1 - gas2;
        const percent = Number(saved * 10000n / gas1) / 100;

        console.log(`counter+1: ${gas1} gas`);
        console.log(`++counter: ${gas2} gas`);
        console.log(`节省: ${saved} gas (${percent.toFixed(2)}%)`);
    });
});
```

### 总结

通过以上实验，你应该掌握了：

1. ✓ Remix 和 Hardhat 开发环境配置
2. ✓ 基本合约编写和编译
3. ✓ 单元测试的编写和运行
4. ✓ 合约部署到本地和测试网
5. ✓ 与合约的交互和事件监听
6. ✓ Gas 优化和性能分析

**下一步建议**：
- 学习高级 Solidity 特性（继承、接口、库）
- 阅读和分析现有项目的合约代码
- 参与 Hackathon，实现真实项目
- 进行代码审计，学习安全最佳实践

所有代码都已经过测试，可直接使用。祝你开发愉快！
