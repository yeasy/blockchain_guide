## 智能合约漏洞

智能合约一旦部署在区块链上，通常是不可变的（除非设计了升级机制）且代码是公开可见的。这种特性要求编写合约时必须极度谨慎。以下是一些常见的智能合约漏洞。

### 1. 重入攻击 （Reentrancy Attack）

这是最著名也是破坏性最大的漏洞之一（如 The DAO 事件）。

**原理**：
攻击者创建一个恶意合约，该合约调用受害合约的提款函数。如果受害合约在更新用户余额之前就发送资金，恶意合约可以在其 `fallback` 函数中再次调用受害合约的提款函数，形成递归调用，直到受害合约余额被耗尽。

**防范**：
*   **Checks-Effects-Interactions 模式**：先检查条件，再更新状态（扣除余额），最后进行交互（发送以太币）。
*   使用重入锁（Reentrancy Guard）。

```solidity
// 易受攻击的代码
function withdraw() public {
    uint balance = userBalances[msg.sender];
    require(balance > 0);
    (bool success, ) = msg.sender.call{value: balance}(""); // Interaction
    require(success);
    userBalances[msg.sender] = 0; // Effect (Too late!)
}

// 安全代码
function withdraw() public {
    uint balance = userBalances[msg.sender];
    require(balance > 0);
    userBalances[msg.sender] = 0; // Effect
    (bool success, ) = msg.sender.call{value: balance}(""); // Interaction
    require(success);
}
```

### 2. 整数溢出与下溢 (Integer Overflow/Underflow)

**原理**：
在 Solidity 0.8.0 之前，如果一个 `uint8` 变量值为 255，加 1 后会变为 0（溢出）；如果值为 0，减 1 后会变为 255（下溢）。攻击者利用这一点可以修改余额或绕过限制。

**防范**：
*   使用 Solidity 0.8.0 及以上版本（内置了溢出检查）。
*   在旧版本中使用 OpenZeppelin 的 `SafeMath` 库。

### 3. 短地址攻击 （Short Address Attack）

**原理**：
早期以太坊节点在处理交易数据时，如果参数长度不足，会自动补零。如果攻击者构造一个恶意的短地址，可能导致合约解析参数时发生位移，从而错误计算转账金额。

**防范**：
*   在合约层面检查输入参数的长度。
*   现代钱包和节点通常已修复此问题。

### 4. 依赖时间戳 （Timestamp Dependence）

**原理**：
`block.timestamp` 可以被矿工小幅度操纵（通常在十几秒内）。如果合约逻辑（如随机数生成、彩票开奖）完全依赖时间戳，矿工可以尝试操控出块时间来获利。

**防范**：
*   避免在关键逻辑中仅依赖 `block.timestamp`。
*   对于需要随机数的场景，使用像 Chainlink VRF 这样的预言机服务。

### 5. 权限控制不当 （Access Control Issues）

**原理**：
关键函数（如修改拥有者、铸造代币、提取资金）没有添加适当的权限修饰符（如 `onlyOwner`），导致任何人都可以调用。

**防范**：
*   严格检查每个 public/external 函数的权限需求。
*   使用 OpenZeppelin 的 `Ownable` 或 `AccessControl` 库。

### 6. 其它常见漏洞

*   **TX.Origin 攻击**：使用 `tx.origin` 进行鉴权而不是 `msg.sender`，可能导致钓鱼攻击。
*   **未处理的返回值**：低级调用（如 `call`）失败时会返回 false 但不会抛出异常，如果未检查返回值，后续逻辑会继续执行。
*   **前端抢跑 (Front-Running)**：攻击者通过观察内存池中的未确认交易，提高 Gas 费抢先执行自己的交易（如抢购 NFT、夹子攻击）。

### 7. DeFi 时代的新型攻击

随着 DeFi 生态的成熟，出现了若干新型攻击向量：

*   **闪电贷攻击 (Flash Loan Attack)**：攻击者利用闪电贷在单笔交易中借入大量资金，通过操纵预言机价格或流动性池来套利，交易结束前归还贷款。由于无需抵押，攻击成本极低。典型案例包括 bZx（2020）、Cream Finance（2021）、Euler Finance（2023，损失约 1.97 亿美元）。防范措施包括使用时间加权平均价格（TWAP）预言机、限制单笔交易的价格影响、添加闪电贷检测机制。

*   **MEV（最大可提取价值）攻击**：矿工/验证者通过重排、插入或审查交易来提取额外价值。常见形式包括三明治攻击（在用户交易前后插入交易来套利）、即时清算（抢先执行清算交易获取奖励）等。以太坊社区通过 PBS（提议者-构建者分离）和 Flashbots 等方案来缓解 MEV 问题。

*   **跨链桥攻击**：跨链桥由于管理大量锁定资产且涉及多链交互，成为高价值攻击目标。典型案例包括 Ronin Bridge（2022，6.2 亿美元）、Wormhole（2022，3.2 亿美元）。防范需要多重签名、零知识证明验证、延迟提款等机制。

*   **治理攻击**：攻击者通过闪电贷临时获取大量治理代币，在 DAO 投票中通过恶意提案。防范措施包括投票锁定期、时间锁（Timelock）和快照机制。
