# 智能合约漏洞

智能合约一旦部署在区块链上，通常是不可变的（除非设计了升级机制）且代码是公开可见的。这种特性要求编写合约时必须极度谨慎。以下是一些常见的智能合约漏洞。

## 1. 重入攻击 (Reentrancy Attack)

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

## 2. 整数溢出与下溢 (Integer Overflow/Underflow)

**原理**：
在 Solidity 0.8.0 之前，如果一个 `uint8` 变量值为 255，加 1 后会变为 0（溢出）；如果值为 0，减 1 后会变为 255（下溢）。攻击者利用这一点可以修改余额或绕过限制。

**防范**：
*   使用 Solidity 0.8.0 及以上版本（内置了溢出检查）。
*   在旧版本中使用 OpenZeppelin 的 `SafeMath` 库。

## 3. 短地址攻击 (Short Address Attack)

**原理**：
早期以太坊节点在处理交易数据时，如果参数长度不足，会自动补零。如果攻击者构造一个恶意的短地址，可能导致合约解析参数时发生位移，从而错误计算转账金额。

**防范**：
*   在合约层面检查输入参数的长度。
*   现代钱包和节点通常已修复此问题。

## 4. 依赖时间戳 (Timestamp Dependence)

**原理**：
`block.timestamp` 可以被矿工小幅度操纵（通常在十几秒内）。如果合约逻辑（如随机数生成、彩票开奖）完全依赖时间戳，矿工可以尝试操控出块时间来获利。

**防范**：
*   避免在关键逻辑中仅依赖 `block.timestamp`。
*   对于需要随机数的场景，使用像 Chainlink VRF 这样的预言机服务。

## 5. 权限控制不当 (Access Control Issues)

**原理**：
关键函数（如修改拥有者、铸造代币、提取资金）没有添加适当的权限修饰符（如 `onlyOwner`），导致任何人都可以调用。

**防范**：
*   严格检查每个 public/external 函数的权限需求。
*   使用 OpenZeppelin 的 `Ownable` 或 `AccessControl` 库。

## 6. 其它常见漏洞

*   **TX.Origin 攻击**：使用 `tx.origin` 进行鉴权而不是 `msg.sender`，可能导致钓鱼攻击。
*   **未处理的返回值**：低级调用（如 `call`）失败时会返回 false 但不会抛出异常，如果未检查返回值，后续逻辑会继续执行。
*   **前端抢跑 (Front-Running)**：攻击者通过观察内存池中的未确认交易，提高 Gas 费抢先执行自己的交易（如抢购 NFT、夹子攻击）。
