## 什么是区块链

### 定义
区块链技术自身仍然在飞速发展中，也缺乏统一的规范和标准。

[wikipedia](https://en.wikipedia.org/wiki/Block_chain_(database) 给出的定义为：

> A blockchain[1][2][3]—originally, block chain[4][5][6]—is a distributed database that maintains a continuously-growing list of data records hardened against tampering and revision. It consists of data structure blocks—which hold exclusively data in initial blockchain implementations,[7] and both data and programs in some of the more recent implementations—with each block holding batches of individual transactions and the results of any blockchain executables.[8][better source needed] Each block contains a timestamp and information linking it to a previous block.[9]

最早出现区块链是在比特币项目。作为比特币背后的分布式记账平台，区块链在无集中式监管的情况下，稳定运行了八年，支持了海量的交易记录，并未出现严重的漏洞。

公认的最早关于区块链的描述性文献是中本聪所撰写的 [比特币：一种点对点的电子现金系统](https://bitcoin.org/bitcoin.pdf)，但该文献重点在于讨论比特币系统，实际上并没有明确提出区块链的定义和概念。

客观地看，区块链属于一种分布式的记录技术。

跟传统的数据库技术相比，其特点应该包括：

* 维护一条不断增长的链，只可能添加记录，而发生过的记录都不可篡改；
* 去中心化，或者说多中心化，无集中的控制，实现上尽量分布式；
* 可以通过密码学的机制来尽量保护用户信息和记录信息的隐私性。

### 原理
区块链的基本原理理解起来并不难。

首先假设存在一个 P2P 的数据库（这方面的技术相对成熟），剩下来就是大家如何决策去添加数据上来。只允许添加、不允许删除避免了作伪的可能性。这个数据库的结构是一个链，由一个个块组成，这也是其名字的来源。新的数据要加入，必须作为一个新的块来加入。而这个块能否加入，可以通过一些手段来检验出来。

具体到比特币如何使用了区块链技术。比特币将每十分钟内所有的交易都打包在一起，这些信息组成一个块。然后，网络中所有的成员都可以试图来找到一个合法的块（比如基于当前的块的信息，加上时间、id，加上某些其它有用信息等），然后进行一些 hash 计算，并且找到的结果还得满足一定条件（比如小于某个值）。一旦算出来就可以进行全网广播，大家拿到这个算出来的结果，进行正向验证，发现确实符合条件了，就承认你算出来了。

因为算出来的概率要从数学上进行保证，比如每十分钟内大概就刚好算出来一个。所以保证了区块链每十分钟增加一个块。算出来的这个人将获取得到这个时间内所有交易产生的管理费和协议固定发放的奖励费（目前是 25 比特币）。也即俗称的挖矿。