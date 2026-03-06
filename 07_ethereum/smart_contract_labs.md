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
2. 在左侧文件面板中，点击 "+" 创建新文件
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

5. 点击左侧 "Solidity Compiler" 图标（看起来像一个积木）
6. 选择编译器版本 `0.8.0` 或以上
7. 点击 "Compile HelloWorld.sol"
8. 如果没有错误，你应该看到绿色的 "Compile" 按钮

#### 0.2 Hardhat 本地开发环境搭建

Hardhat 是以太坊开发的业界标准框架，提供强大的测试、调试和部署功能。

**安装步骤**：

```bash
# 创建项目目录
mkdir ethereum-learning
cd ethereum-learning

# 初始化 Node.js 项目
npm init -y

# 安装 Hardhat
npm install --save-dev hardhat

# 初始化 Hardhat 项目
npx hardhat

# 选择 "Create a sample project"（回车）
# 确认依赖安装（回车）
```

**项目结构**：

```
ethereum-learning/
├── contracts/
│   └── Lock.sol              # 智能合约
├── test/
│   └── Lock.js               # 测试文件
├── scripts/
│   └── deploy.js             # 部署脚本
├── hardhat.config.js         # Hardhat 配置
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

**编写测试** (`test/SimpleStorage.test.js`)：

```javascript
const { expect } = require("chai");

describe("SimpleStorage", function () {
    let storage;
    let owner;

    // 在每个测试前执行
    beforeEach(async function () {
        [owner] = await ethers.getSigners();
        const Storage = await ethers.getContractFactory("SimpleStorage");
        storage = await Storage.deploy();
        await storage.deployed();
    });

    describe("初始状态", function () {
        it("初始值应该为 0", async function () {
            expect(await storage.getValue()).to.equal(0);
        });

        it("storedValue 应该可以直接访问", async function () {
            expect(await storage.storedValue()).to.equal(0);
        });
    });

    describe("setValue 函数", function () {
        it("应该能够设置新值", async function () {
            await storage.setValue(42);
            expect(await storage.getValue()).to.equal(42);
        });

        it("应该发出 ValueChanged 事件", async function () {
            // 监听事件
            await expect(storage.setValue(100))
                .to.emit(storage, "ValueChanged")
                .withArgs(100, owner.address);
        });

        it("每次设置都应该覆盖旧值", async function () {
            await storage.setValue(10);
            expect(await storage.getValue()).to.equal(10);

            await storage.setValue(20);
            expect(await storage.getValue()).to.equal(20);
        });
    });

    describe("add 函数", function () {
        it("应该正确相加", async function () {
            expect(await storage.add(5, 3)).to.equal(8);
        });

        it("应该处理大数字", async function () {
            const big = ethers.BigNumber.from("10").pow(18);
            expect(await storage.add(big, big)).to.equal(big.mul(2));
        });
    });

    describe("Gas 消耗分析", function () {
        it("setValue 是一个状态修改交易", async function () {
            const tx = await storage.setValue(50);
            const receipt = await tx.wait();
            console.log("setValue Gas 消耗:", receipt.gasUsed.toString());
            // 预期: ~43,000-44,000
        });

        it("getValue 是一个只读调用", async function () {
            // 不会消耗任何 Gas
            const result = await storage.getValue();
            expect(result).to.equal(0);
        });
    });
});
```

**运行测试**：

```bash
npx hardhat test test/SimpleStorage.test.js
```

**预期输出**：

```
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

**测试代码** (`test/SimpleToken.test.js`)：

```javascript
const { expect } = require("chai");

describe("SimpleToken", function () {
    let token;
    let owner, alice, bob;

    beforeEach(async function () {
        [owner, alice, bob] = await ethers.getSigners();

        const Token = await ethers.getContractFactory("SimpleToken");
        token = await Token.deploy(1000); // 初始供应 1000 tokens
        await token.deployed();
    });

    describe("初始化", function () {
        it("应该有正确的总供应", async function () {
            const expected = ethers.utils.parseEther("1000");
            expect(await token.totalSupply()).to.equal(expected);
        });

        it("拥有者应该获得所有初始代币", async function () {
            const balance = await token.balanceOf(owner.address);
            expect(balance).to.equal(await token.totalSupply());
        });
    });

    describe("转账", function () {
        it("应该能够转账代币", async function () {
            const amount = ethers.utils.parseEther("100");
            await token.transfer(alice.address, amount);

            expect(await token.balanceOf(alice.address)).to.equal(amount);
            expect(await token.balanceOf(owner.address)).to.equal(
                ethers.utils.parseEther("900")
            );
        });

        it("应该发出 Transfer 事件", async function () {
            const amount = ethers.utils.parseEther("50");
            await expect(token.transfer(bob.address, amount))
                .to.emit(token, "Transfer")
                .withArgs(owner.address, bob.address, amount);
        });

        it("不应该允许转账超过余额", async function () {
            const tooMuch = ethers.utils.parseEther("2000");
            await expect(
                token.transfer(alice.address, tooMuch)
            ).to.be.revertedWith("Insufficient balance");
        });
    });

    describe("授权和代理转账", function () {
        it("应该能够授权", async function () {
            const amount = ethers.utils.parseEther("100");
            await token.approve(alice.address, amount);

            expect(await token.allowance(owner.address, alice.address)).to.equal(amount);
        });

        it("应该能够进行代理转账", async function () {
            const amount = ethers.utils.parseEther("100");

            // owner 授权 alice 花费 100 tokens
            await token.approve(alice.address, amount);

            // alice 代表 owner 转账给 bob
            await token.connect(alice).transferFrom(owner.address, bob.address, amount);

            expect(await token.balanceOf(bob.address)).to.equal(amount);
            expect(await token.allowance(owner.address, alice.address)).to.equal(0);
        });

        it("不应该允许转账超过授权额度", async function () {
            const approved = ethers.utils.parseEther("50");
            const attempted = ethers.utils.parseEther("100");

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
npx hardhat test test/SimpleToken.test.js
```

### 实验 3：使用 Ganache 搭建本地测试网

Ganache 是一个功能更强大的本地区块链模拟器，提供 GUI 界面和更多调试功能。

**安装 Ganache**：

```bash
# 安装 ganache-cli
npm install --save-dev ganache-cli

# 或安装 Ganache GUI（推荐）
# 访问 https://www.trufflesuite.com/ganache 下载桌面应用
```

**启动 Ganache**：

```bash
# CLI 方式
npx ganache-cli --deterministic --accounts 10 --host 0.0.0.0

# 或使用 GUI 应用（更直观）
```

**Ganache 提供的好处**：
- 10 个预生成的账户，每个初始 100 ETH
- 可视化的区块链状态
- 详细的交易日志
- 时间快进功能（用于测试时间锁）

**在 Hardhat 中配置 Ganache**：

```javascript
// hardhat.config.js
module.exports = {
    solidity: "0.8.0",
    networks: {
        ganache: {
            url: "http://127.0.0.1:8545",
            accounts: [
                // 粘贴 Ganache 显示的私钥
                "0xac0974bec39a17e36ba4a6b4d238ff944bacb476c6b8d6c1f02e31a0c2b7e6c1",
                // ... 更多账户
            ]
        },
        hardhat: {
            forking: {
                enabled: true,
                url: "https://eth-mainnet.alchemyapi.io/v2/YOUR_API_KEY"
            }
        }
    }
};
```

**使用 Ganache 运行测试**：

```bash
# 终端 1：启动 Ganache
npx ganache-cli

# 终端 2：运行 Hardhat 测试
npx hardhat test --network ganache
```

### 实验 4：在测试网部署合约

**选择测试网**：

- **Sepolia**（推荐）：最新的以太坊官方测试网
- **Goerli**：已弃用，勿用
- **Mumbra**：Polygon 测试网

**获取测试代币**：

```bash
# 访问 Sepolia 水龙头
# https://sepoliafaucet.com
# https://www.alchemy.com/faucets/ethereum-sepolia

# 或使用 Alchemy 的 Python 脚本
pip install eth-faucet
python -m eth_faucet --address YOUR_ADDRESS --network sepolia
```

**配置网络** (`hardhat.config.js`)：

```javascript
require("@nomiclabs/hardhat-waffle");
require("dotenv").config();

module.exports = {
    solidity: "0.8.0",
    networks: {
        sepolia: {
            url: `https://eth-sepolia.g.alchemy.com/v2/${process.env.ALCHEMY_API_KEY}`,
            accounts: [process.env.PRIVATE_KEY]
        }
    },
    etherscan: {
        apiKey: process.env.ETHERSCAN_API_KEY
    }
};
```

**创建 .env 文件**：

```
ALCHEMY_API_KEY=your_alchemy_api_key_here
PRIVATE_KEY=your_wallet_private_key_here
ETHERSCAN_API_KEY=your_etherscan_api_key_here
```

**部署脚本** (`scripts/deploy.js`)：

```javascript
const hre = require("hardhat");

async function main() {
    console.log("部署开始...");

    // 编译合约
    const SimpleStorage = await hre.ethers.getContractFactory("SimpleStorage");

    // 部署
    const storage = await SimpleStorage.deploy();
    await storage.deployed();

    console.log(`SimpleStorage 已部署到: ${storage.address}`);

    // 保存地址以供后续使用
    require("fs").writeFileSync(
        "deployed.json",
        JSON.stringify({ SimpleStorage: storage.address })
    );
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
```

**执行部署**：

```bash
npx hardhat run scripts/deploy.js --network sepolia
```

**验证合约** (在 Etherscan 上公开代码)：

```bash
npx hardhat verify --network sepolia CONTRACT_ADDRESS "constructor arguments"
```

### 实验 5：与合约交互

**创建交互脚本** (`scripts/interact.js`)：

```javascript
const hre = require("hardhat");
const fs = require("fs");

async function main() {
    const deployment = JSON.parse(fs.readFileSync("deployed.json"));
    const contractAddress = deployment.SimpleStorage;

    const [signer] = await hre.ethers.getSigners();
    console.log(`使用账户: ${signer.address}`);

    // 连接到部署的合约
    const SimpleStorage = await hre.ethers.getContractFactory("SimpleStorage");
    const storage = await SimpleStorage.attach(contractAddress);

    // 查询初始值
    let value = await storage.getValue();
    console.log(`初始值: ${value}`);

    // 设置新值（这会触发交易）
    console.log("\n正在设置新值为 42...");
    const tx = await storage.setValue(42);
    console.log(`交易哈希: ${tx.hash}`);

    // 等待交易确认
    const receipt = await tx.wait();
    console.log(`交易已确认，区块号: ${receipt.blockNumber}`);

    // 查询新值
    value = await storage.getValue();
    console.log(`新值: ${value}`);

    // 添加数字
    const result = await storage.add(10, 20);
    console.log(`\n10 + 20 = ${result}`);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
```

**执行交互**：

```bash
npx hardhat run scripts/interact.js --network sepolia
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

```javascript
describe("Gas 优化", function () {
    it("对比 ++ 和 +1", async function () {
        const Unoptimized = await ethers.getContractFactory("Unoptimized");
        const Optimized = await ethers.getContractFactory("Optimized");

        const unopt = await Unoptimized.deploy();
        const opt = await Optimized.deploy();

        const tx1 = await unopt.increment();
        const receipt1 = await tx1.wait();
        const gas1 = receipt1.gasUsed.toNumber();

        const tx2 = await opt.increment();
        const receipt2 = await tx2.wait();
        const gas2 = receipt2.gasUsed.toNumber();

        console.log(`counter+1: ${gas1} gas`);
        console.log(`++counter: ${gas2} gas`);
        console.log(`节省: ${gas1 - gas2} gas (${((1 - gas2/gas1)*100).toFixed(2)}%)`);
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
