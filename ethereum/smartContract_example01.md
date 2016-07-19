## 链码示例一：Hello World!
### 简介

[smartContract_example01.sol](smartContract_example01.sol)

合约greeter是一个简单的智能合约，你可以使用这个合约来和其他人交流，它的回复会同你的输入完全一样，当输入为“Hello World!”的时候，合约也会回复“Hello World!”。

###目的:
该合约主要面向第一次接触solidity和ethereum的初学者,旨在让大家能够了解如何编写一个简单的智能合约程序,
掌握基本流程。
###主要实现如下的功能：
* 返回你预先设置的字符串

### 主要函数
* `kill`：selfdestruct 是 ethereum 智能合约自带的自毁程序,kill对此方法进行了封装,只有合约的拥有者才可以调用该方法；
* `greet`：返回合约 greeter 里的 greeting属性的值；

### 代码运行分析

####第一步 生成智能合约代码对象
我们先把合约代码[smartContract_example01.sol](smartContract_example01.sol)
压缩为一行．新建一个ssh session, 切换到geth用户环境`su - geth`, 然后输入：`cat smartContract_example01.sol | tr '\n' ' '`.
切换到以太坊控制台，把合约代码保存为一个变量:
```js
var greeterSource = 'contract mortal { address owner; function mortal() { owner = msg.sender; } function kill() { if (msg.sender == owner) selfdestruct(owner); } } contract greeter is mortal { string greeting; function greeter(string _greeting) public { greeting = _greeting; } function greet() constant returns (string) { return greeting; } }'
```

####第二步 编译合约代码
然后编译合约代码：
```js
var greeterCompiled = web3.eth.compile.solidity(greeterSource)
```
`greeterCompiled.Token.code`可以看到编译好的二进制代码
`greeterCompiled.Token.info.abiDefinition`可以看到合约的ABI

####第三步 设置希望返回的字符串
```js
var _greeting = "Hello World!"
```
####第四步 部署合约
接下来我们要把编译好的合约部署到网络上去．

首先我们用ABI来创建一个javascript环境中的合约对象：
```sol
var greeterContract = web3.eth.contract(greeterCompiled.greeter.info.abiDefinition);
```
我们通过合约对象来部署合约：
```js
var greeter = greeterContract.new(_greeting,{from:web3.eth.accounts[0], data: greeterCompiled.greeter.code, gas: 300000}, function(e, contract){
    if(!e) {
      if(!contract.address) {
        console.log("Contract transaction send: TransactionHash: " + contract.transactionHash + " waiting to be mined...");
      } else {
        console.log("Contract mined! Address: " + contract.address);
        console.log(contract);
      }
    }
})
```

- greeterContract.new方法的第一个参数设置了这个新合约的构造函数初始化的值
- greeterContract.new方法的第二个参数设置了这个新合约的创建者地址from,
这个新合约的代码data, 和用于创建新合约的费用gas．gas是一个估计值，只要比所需要的gas多就可以
，合约创建完成后剩下的gas会退还给合约创建者．
- greeterContract.new方法的第三个参数设置了一个回调函数，可以告诉我们部署是否成功．

contract.new执行时会提示输入钱包密码．执行成功后，我们的合约Token就已经广播到网络上了．
此时只要等待矿工把我们的合约打包保存到以太坊区块链上，部署就完成了．
####第五步 挖矿
在公有链上，矿工打包平均需要15秒，在私有链上，我们需要自己来做这件事情．首先开启挖矿：
```js
miner.start(1)
```
此时需要等待一段时间，以太坊节点会生成挖矿必须的数据，这些数据都会放到内存里面．
在数据生成好之后，挖矿就会开始，稍后就能在控制台输出中看到类似：
```
...
I0714 22:00:19.694219 ethash.go:291] Generating DAG: 97%
I0714 22:00:22.987934 ethash.go:291] Generating DAG: 98%
I0714 22:00:26.543035 ethash.go:291] Generating DAG: 99%
I0714 22:00:29.912655 ethash.go:291] Generating DAG: 100%
I0714 22:00:29.915580 ethash.go:276] Done generating DAG for epoch 2, it took 5m34.983289765s
```
####第六步 停止挖矿(可选)
当生成DAG结束,提示已经挖出至少一个矿以后,我们需要停止挖矿(当然,你也可以不停,就是会一直输出)
```js
miner.stop()
```

####第七步 部署在其他节点上
现在,你已经成功部署了一个智能合约,当运行以下代码时:
```js
//由于该命令未改变blockchain,所以不会有任何花费
greeter.greet();
```
命令行上会出现如下返回结果:
```js
'Hello World!'
```
好了,我们的第一个智能合约程序 "Hello World!" 已经完成了,但是目前它只有一个节点!

####第八步 部署在其他节点上
为了使得其他人可以运行你的智能合约，你需要两个信息：
1. 智能合约地址Address
2. 智能合约ABI（Application Binary Interface），ABI其实就是一个有序的用户手册，描述了所有方法的名字和如何调用它们。
我们可以使用如下代码获得其ABI和智能合约地址:
```js
greeterCompiled.greeter.info.abiDefinition;
greeter.address;
```

然后你可以实例化一个JavaScript对象，该对象可以用来在任意联网机器上调用该合约，此处***ABI***和***Address***是上述代码返回值。
```js
var greeter = eth.contract(ABI).at(Address);
```

####第九步 自毁程序
一个交易需要被发送到网络需要支付费用，自毁程序是对网络的补充，花费的费用远小于一次常用交易。

你可以通过以下代码来检验是否成功，如果自毁程序运行成功以下代码会返回0：
```js
greeter.kill.sendTransaction({from:eth.accounts[0]})
```

###参考文献
[THE GREETER YOUR DIGITAL PAL WHO'S FUN TO BE WITH](
https://www.ethereum.org/greeter#compiling-your-contract)

[以太坊本地私有链开发环境搭建](
http://ethfans.org/posts/ethereum-private-network-bootstrap)
