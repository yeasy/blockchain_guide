## 简介

### 历史
区块链已经成为当下最受人关注的开源技术，有人说它将颠覆金融行业的未来。然而对很多人来说，区块链技术难以理解和实现，而且缺乏统一的规范。

2015 年 12 月，[Linux 基金](http://www.linuxfoundation.org/) 会牵头，联合 30 家初始成员（包括各大金融、科技公司和相关开源组织），共同 [宣告](https://www.hyperledger.org/news/announcement/2016/02/hyperledger-project-announces-30-founding-members) 了 [Hyperledger](https://www.hyperledger.org) 项目的成立。该项目试图打造一个透明、公开、去中心化的超级账本项目，作为区块链技术的开源规范和标准，让更多的应用能更容易的建立在区块链技术之上。目前已经有超过 80 家企业和机构（大部分均为各自行业的领导者）宣布加入 Hyperledger 项目。

如果说以比特币为代表的货币区块链技术为 1.0，以以太坊为代表的合同区块链技术为 2.0，那么实现了完备的权限控制和安全保障的 Hyperledger 项目毫无疑问代表着 3.0 时代的到来。

IBM 贡献了数万行已有的 [Open Block Chain](https://github.com/openblockchain) 代码，Digital Asset 则贡献了企业和开发者相关资源，R3 贡献了新的金融交易架构，Intel 也刚贡献了跟分布式账本相关的代码。

首届技术委员会主席由来自 IBM 开源技术部 CTO 的 [Chris Ferris](https://www.linkedin.com/in/chrisfer) 担任，委员会主席则由来自 Digital Asset Holdings 的 CEO Blythe Masters 担任。另外，自 2016 年 5 月起，Apache 基金会创始人 Brian Behlendorf 担任超级账本项目的首位执行董事。

该项目的出现，实际上宣布区块链技术已经不单纯是一个开源技术了，已经正式被主流机构和市场认可；同时，Hyperledger 首次提出和实现的完备权限管理、创新的一致性算法和可拔插的框架，对于区块链相关技术和产业的发展都将产生深远的影响。

项目官方地址托管在 [Linux 基金会网站](https://blockchain.linuxfoundation.org/)，代码托管在 [Github](https://github.com/hyperledger/hyperledger) 上，目前已经获得了不少关注。

![](_images/hyperledger.png)

目前主要包括两大子项目：

* [fabric](https://github.com/hyperledger/fabric)：包括 [fabric](https://github.com/hyperledger/fabric) 和 [fabric-api](https://github.com/hyperledger/fabric-api)，目标是区块链的基础核心平台，支持 pbft 等新的 consensus 机制，支持权限管理，最早由 IBM 和 DAH 发起；
* [sawtooth Lake](https://github.com/hyperledger/sawtooth-core)：包括 arcade、[core]((https://github.com/hyperledger/sawtooth-core)、dev-tools、[validator](https://github.com/hyperledger/sawtooth-validator)、mktplace 等。是 Intel 主要贡献和主导的区块链平台，支持全新的共识机制 Proof of Elapsed Time（PoET）。

目前，所有项目均处于孵化（Incubation）状态。

项目约定共同遵守的 [基本原则](https://github.com/hyperledger/hyperledger) 为：

* 重视模块化设计，包括交易、合同、一致性、身份、存储等技术场景；
* 代码可读性，保障新功能和模块都可以很容易添加和扩展；
* 演化路线，随着需求的深入和更多的应用场景，不断增加和演化新的项目。
