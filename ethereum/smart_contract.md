## 使用智能合约

以太坊社区有不少提供智能合约编写、编译、发布、调用等功能的工具，用户和开发者可以根据需求或开发环境自行选择。

本节将向开发者介绍使用 Geth 客户端搭建测试用的本地区块链，以及如何在链上部署和调用智能合约。

### 搭建测试用区块链

由于在以太坊公链上测试智能合约需要消耗以太币，所以对于开发者开发测试场景，可以选择本地自行搭建一条测试链。开发好的智能合约可以很容易的切换接口部署到公有链上。注意测试链不同于以太坊公链，需要给出一些非默认的手动配置。

#### 配置初始状态

首先配置私有区块链网络的初始状态。新建文件 `genesis.json`，内容如下。

```json
{
  "config": {
        "chainId": 22,
        "homesteadBlock": 0,
        "eip155Block": 0,
        "eip158Block": 0
  },
  "alloc"      : {},
  "coinbase"   : "0x0000000000000000000000000000000000000000",
  "difficulty" : "0x400",
  "extraData"  : "",
  "gasLimit"   : "0x2fefd8",
  "nonce"      : "0x0000000000000038",
  "mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp"  : "0x00"
}
```

其中，`chainId` 指定了独立的区块链网络 ID，不同 ID 网络的节点无法互相连接。配置文件还对当前挖矿难度 `difficulty`、区块 Gas 消耗限制 `gasLimit` 等参数进行了设置。

#### 启动区块链

用以下命令初始化区块链，生成创世区块和初始状态。

```bash
$ geth --datadir /path/to/datadir init /path/to/genesis.json
```

其中，`--datadir` 指定区块链数据的存储位置，可自行选择一个目录地址。

接下来用以下命令启动节点，并进入 Geth 命令行界面。

```bash
$ geth --identity "TestNode" --rpc --rpcport "8545" --datadir /path/to/datadir --port "30303" --nodiscover console
```

各选项的含义如下。

* `--identity`：指定节点 ID；
* `--rpc`： 表示开启 HTTP-RPC 服务；
* `--rpcport`： 指定 HTTP-RPC 服务监听端口号（默认为 8545）；
* `--datadir`： 指定区块链数据的存储位置；
* `--port`： 指定和其他节点连接所用的端口号（默认为 30303）；
* `--nodiscover`： 关闭节点发现机制，防止加入有同样初始配置的陌生节点；

#### 创建账号

用上述 `geth console` 命令进入的命令行界面采用 JavaScript 语法。可以用以下命令新建一个账号。

```
> personal.newAccount()

Passphrase:
Repeat passphrase:
"0x1b6eaa5c016af9a3d7549c8679966311183f129e"
```

输入两遍密码后，会显示生成的账号，如`"0x1b6eaa5c016af9a3d7549c8679966311183f129e"`。可以用以下命令查看该账号余额。

```
> myAddress = "0x1b6eaa5c016af9a3d7549c8679966311183f129e"
> eth.getBalance(myAddress)
0
```

看到该账号当前余额为 0。可用 `miner.start()` 命令进行挖矿，由于初始难度设置的较小，所以很容易就可挖出一些余额。`miner.stop()` 命令可以停止挖矿。

### 创建和编译智能合约

以 Solidity 编写的智能合约为例。为了将合约代码编译为 EVM 二进制，需要安装 Solidity 编译器 solc。

```bash
$ apt-get install solc
```

新建一个 Solidity 智能合约文件，命名为 `testContract.sol`，内容如下。该合约包含一个方法 multiply，作用是将输入的整数乘以 7 后输出。

```
pragma solidity ^0.4.0;
contract testContract {
  function multiply(uint a) returns(uint d) {
    d = a * 7;
  }
}
```

用 solc 获得合约编译后的 EVM 二进制码。

```bash
$ solc --bin testContract.sol

======= testContract.sol:testContract =======
Binary:
6060604052341561000c57fe5b5b60a58061001b6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063c6888fa114603a575bfe5b3415604157fe5b60556004808035906020019091905050606b565b6040518082815260200191505060405180910390f35b60006007820290505b9190505600a165627a7a72305820748467daab52f2f1a63180df2c4926f3431a2aa82dcdfbcbde5e7d036742a94b0029
```

再用 solc 获得合约的 JSON ABI（Application Binary Interface），其中指定了合约接口，包括可调用的合约方法、变量、事件等。

```bash
$ solc --abi testContract.sol

======= testContract.sol:testContract =======
Contract JSON ABI
[{"constant":false,"inputs":[{"name":"a","type":"uint256"}],"name":"multiply","outputs":[{"name":"d","type":"uint256"}],"payable":false,"type":"function"}]
```

下面回到 Geth 的 JavaScript 环境命令行界面，用变量记录上述两个值。注意在 code 前加上 `0x` 前缀。

```
> code = "0x6060604052341561000c57fe5b5b60a58061001b6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063c6888fa114603a575bfe5b3415604157fe5b60556004808035906020019091905050606b565b6040518082815260200191505060405180910390f35b60006007820290505b9190505600a165627a7a72305820748467daab52f2f1a63180df2c4926f3431a2aa82dcdfbcbde5e7d036742a94b0029"
> abi = [{"constant":false,"inputs":[{"name":"a","type":"uint256"}],"name":"multiply","outputs":[{"name":"d","type":"uint256"}],"payable":false,"type":"function"}]
```

### 部署智能合约

在 Geth 的 JavaScript 环境命令行界面，首先用以下命令解锁自己的账户，否则无法发送交易。

```
> personal.unlockAccount(myAddress)

Unlock account 0x1b6eaa5c016af9a3d7549c8679966311183f129e
Passphrase:
true
```

接下来发送部署合约的交易。

```
> myContract = eth.contract(abi)
> contract = myContract.new({from:myAddress,data:code,gas:1000000})
```

如果此时没有在挖矿，用 `txpool.status` 命令可看到本地交易池中有一个待确认的交易。可用以下命令查看当前待确认的交易。

```
> eth.getBlock("pending",true).transactions

[{
    blockHash: "0xbf0619ca48d9e3cc27cd0ab0b433a49a2b1bed90ab57c0357071b033aca1f2cf",
    blockNumber: 17,
    from: "0x1b6eaa5c016af9a3d7549c8679966311183f129e",
    gas: 90000,
    gasPrice: 20000000000,
    hash: "0xa019c2e5367b3ad2bbfa427b046ab65c81ce2590672a512cc973b84610eee53e",
    input: "0x6060604052341561000c57fe5b5b60a58061001b6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063c6888fa114603a575bfe5b3415604157fe5b60556004808035906020019091905050606b565b6040518082815260200191505060405180910390f35b60006007820290505b9190505600a165627a7a72305820748467daab52f2f1a63180df2c4926f3431a2aa82dcdfbcbde5e7d036742a94b0029",
    nonce: 1,
    r: "0xbcb2ba94f45dfb900a0533be3c2c603c2b358774e5fe89f3344031b202995a41",
    s: "0x5f55fb1f76aa11953e12746bc2d19fbea6aeb1b9f9f1c53a2eefab7058515d99",
    to: null,
    transactionIndex: 0,
    v: "0x4f",
    value: 0
}]
```

可以用 `miner.start()` 命令挖矿，一段时间后，交易会被确认，即随新区块进入区块链。

### 调用智能合约

用以下命令可以发送交易，其中 sendTransaction 方法的前几个参数与合约中 multiply 方法的输入参数对应。这种方式，交易会通过挖矿记录到区块链中，如果涉及状态改变也会获得全网共识。

```
> contract.multiply.sendTransaction(10, {from:myAddress})
```

如果只是想本地运行该方法查看返回结果，可采用如下方式获取结果。

```
> contract.multiply.call(10)
70
```
