## 顶级项目介绍

Hyperledger 所有项目代码托管在 [Github](https://github.com/hyperledger/)上。目前，主要包括如下顶级项目（按时间顺序）。

* [Fabric](https://github.com/hyperledger/fabric)：包括 [Fabric](https://github.com/hyperledger/fabric)、[Fabric CA](https://github.com/hyperledger/fabric-ca)、Fabric SDK（包括 Node.Js、Java、Python 和 Go 语言）等，目标是区块链的基础核心平台，支持 PBFT 等新的共识机制，支持权限管理，最早由 IBM 和 DAH 于 2015 年底发起；
* [Sawtooth](https://github.com/hyperledger/sawtooth-core)：包括 arcade、[core](https://github.com/hyperledger/sawtooth-core)、dev-tools、[validator](https://github.com/hyperledger/sawtooth-validator)、mktplace 等。是 Intel 主要发起和贡献的区块链平台，支持全新的基于硬件芯片的共识机制 Proof of Elapsed Time（PoET）， 2016 年 4 月贡献到社区；
* [Blockchain Explorer](https://github.com/hyperledger/blockchain-explorer)：提供 Web 操作界面，通过界面快速查看查询绑定区块链的状态（区块个数、交易历史）信息等，由 DTCC、IBM、Intel 等开发支持，2016 年 8 月贡献到社区；
* [Iroha](https://github.com/hyperledger/Iroha)：账本平台项目，基于 C++ 实现，带有不少面向 Web 和 Mobile 的特性，主要由 Soramitsu 于 2016 年 10 月发起和贡献；
* [Cello](https://github.com/hyperledger/cello)：提供区块链平台的部署和运行时管理功能。使用 Cello，管理员可以轻松部署和管理多条区块链；应用开发者可以无需关心如何搭建和维护区块链，由 IBM 团队于 2017 年 1 月贡献到社区；
* [Indy](https://github.com/hyperledger/indy)：提供基于分布式账本技术的数字身份管理机制，由 Sovrin 基金会发起，2017 年 3 月底正式贡献到社区；
* [Composer](https://github.com/hyperledger/composer)：[End of Life] 提供面向链码开发的高级语言支持，自动生成链码代码等，由 IBM 团队发起并维护，2017 年 3 月底贡献到社区。目前已经停止维护；
* [Burrow](https://github.com/hyperledger/burrow)：提供以太坊虚拟机的支持，实现支持高效交易的带权限的区块链平台，由 Monax 公司发起支持，2017 年 4 月贡献到社区；
* [Quilt](https://github.com/hyperledger/quilt)：对 W3C 支持的跨账本协议 Interledger 的 Java 实现。2017 年 10 月正式贡献到社区；
* [Caliper](https://github.com/hyperledger/burrow)：提供对区块链平台性能的测试工具，由华为公司发起支持，2018 年 3 月正式贡献到社区。
* [Ursa](https://github.com/hyperledger/ursa)：[End of Life] 提供一套密码学相关组件，初始贡献者包括来自 Fujitsu、Sovrin、Intel、DFINITY、State Street、IBM、Bitwise IO 等企业的开发者，2018 年 11 月正式被接收到社区，2023 年 4 月结束生命周期；
* [Grid](https://github.com/hyperledger/burrow)：提供帮助快速构建供应链应用的框架，由 Cargill、Intel 和 Bitwise IO 公司发起支持，2018 年 12 月正式贡献到社区；
* [Transact](https://github.com/hyperledger/transact)：提供运行交易的引擎和框架，由 Bitwise IO、Cargill、Intel、IBM、HACERA 等公司发起支持，2019 年 5 月正式贡献到社区；
* [Aries](https://github.com/hyperledger/aries)：为客户端提供共享的密码学钱包，由 Sovrin、C3I 和 Evernym 等公司发起支持，2019 年 5 月正式贡献到社区；
* [Besu](https://github.com/hyperledger/besu)：作为企业级的以太坊客户端支持，由 Consensys、Hacera、JPM 和 Redhat 等公司发起支持，2019 年 8 月正式贡献到社区；
* [Avalon](https://github.com/hyperledger/avalon)：提供链下计算支持，增强安全性和可扩展性，由 Intel、IEX、IBM 和 Consensys 等公司发起支持，2019 年 9 月正式贡献到社区。

这些顶级项目分别从平台、工具和类库三个层次相互协作，构成了完善的生态系统，如下图所示。

![Hyperledger 顶级项目](_images/top_projects.png)

所有项目一般都需要经历提案（Proposal）、孵化（Incubation）、活跃（Active）、退出（Deprecated）、终结（End of Life）等 5 个生命周期。

任何希望加入到 Hyperledger 社区中的项目，必须首先由发起人编写提案。描述项目的目的、范围、参与者和开发计划等重要信息，并由全球技术委员会来进行评审投票，评审通过则可以进入到社区内进行孵化。项目成熟后可以申请进入到活跃状态，发布正式的版本。项目不再活跃后可以进入维护阶段，最终结束生命周期。

### Fabric 项目

![Hyperledger Fabric 项目](_images/fabric.png)

作为最早加入到超级账本项目中的顶级项目，Fabric 由 IBM、DAH 等企业于 2015 年底联合贡献到社区。项目在 Github 上地址为 https://github.com/hyperledger/fabric。

该项目的定位是面向企业的分布式账本平台，其创新地引入了权限管理支持，设计上支持可插拔、可扩展，是首个面向联盟链场景的开源项目。

Fabric 项目基于 Go 语言实现，贡献者超过 200 人，总提交次数已经超过 15000 次，核心代码数超过 15 万行。

Fabric 项目目前处于活跃状态，当前长期支持版本（LTS）为 2.5.x 系列，已发布至 v2.5.9（2024年6月）。Fabric 3.0 已发布 Beta 版本，引入了拜占庭容错（BFT）排序服务等新特性。项目同时包括 Fabric CA、Fabric Gateway、多语言 SDK 等子项目。

项目的邮件列表地址为 fabric@lists.hyperledger.org。

### Sawtooth 项目

![Hyperledger Sawtooth 项目](_images/sawtooth.png)

Sawtooth 项目由 Intel 等企业于 2016 年 4 月提交到社区，包括 sawtooth-core、sawtooth-supply-chain、sawtooth-marketplace、sawtooth-seth、sawtooth-next-directory、sawtooth-explorer 等数十个子项目。核心代码在 Github 上地址为 https://github.com/hyperledger/sawtooth-core。

该项目的定位也是分布式账本平台，基于 Python 语言实现。项目目前处于 Active 阶段，核心项目的贡献者超过 70 人，提交次数已经超过 8000 次。

Sawtooth 项目利用 Intel 芯片的专属功能，实现了低功耗的 Proof of Elapsed Time（PoET）共识机制，并支持交易族（Transaction Family），方便用户使用它来快速开发应用。

项目的邮件列表地址为 sawtooth@lists.hyperledger.org。

### Iroha 项目

![Hyperledger Iroha 项目](_images/iroha.png)

Iroha 项目由 Soramitsu 等企业于 2016 年 10 月提交到社区，包括 iroha、iroha-android、iroha-ios、iroha-python、iroha-javascript 等子项目。核心代码在 Github 上地址为 https://github.com/hyperledger/iroha。

该项目的定位是分布式账本平台框架，基于 C++ 语言实现。项目目前处于 Active 阶段，贡献者超过 50 人，提交次数已经超过 7000 次。

Iroha 项目在设计上类似 Fabric，同时提供了基于 C++ 的区块链开发环境，并考虑了移动端和 Web 端的一些需求。

项目的邮件列表地址为 iroha@lists.hyperledger.org。

###  Explorer 项目

![Hyperledger Explorer 项目](_images/explorer.png)

Explorer 项目由 Intel、DTCC、IBM 等企业于 2016 年 8 月提交到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/blockchain-explorer，目前贡献者超过 40 人，提交次数超过 350 次。

该项目的定位是区块链平台的浏览器，基于 Node.js 语言实现，提供 Web 操作界面。用户可以使用它来快速查看底层区块链平台的运行信息，如区块个数、交易情况、网络状况等。

项目的邮件列表地址为 explorer@lists.hyperledger.org。

### Cello 项目

![Hyperledger Cello 项目](_images/cello.png)

Cello 项目由笔者领导的技术团队于 2017 年 1 月贡献到社区。Github 上仓库地址为 https://github.com/hyperledger/cello（核心代码）和 https://github.com/hyperledger/cello-analytics（侧重数据分析）。

该项目的定位为区块链网络的操作系统，实现区块链网络自动化部署，以及对区块链网络的运行时管理。使用 Cello，可以让区块链应用人员专注到应用开发，而无需关心底层平台的管理和维护。已有一些企业基于 Cello 项目代码构建了区块链即服务（Blockchain-as-a-Service）平台。

Cello 的主要开发语言为 Python 和 JavaScript 等，底层支持包括裸机、虚拟机、容器云（包括 Swarm、Kubernetes）等多种基础架构。目前贡献者超过 40 人，提交次数超过 1000 次。

项目的邮件列表地址为 cello@lists.hyperledger.org。

### Indy 项目

![Hyperledger Indy 项目](_images/indy.png)

Indy 项目由 Sovrin 基金会牵头进行开发，致力于打造一个基于区块链和分布式账本技术的数字身份管理平台。该平台支持去中心化，支持跨区块链和跨应用的操作，实现全球化的身份管理。Indy 项目于 2017 年 3 月底正式加入到超级账本项目。目前包括 indy-node、indy-sdk、indy-plenum、indy-hipe、indy-crypto、indy-agent 等项目。

该项目主要由 Python 语言开发，包括服务节点、客户端和通用库等。目前处于 Active 阶段，贡献者超过 60人，已有超过 5000 次提交。

项目的邮件列表地址为 indy@lists.hyperledger.org。

### Composer 项目

![Hyperledger Composer 项目](_images/composer.png)

Composer 项目由 IBM 团队于 2017 年 3 月底贡献到社区，试图提供一个 Hyperledger Fabric 的开发辅助框架。使用 Composer，开发人员可以使用 Javascript 语言定义应用逻辑，再加上资源、参与者、交易等模型和访问规则，生成 Hyperledger Fabric 支持的链码。

该项目主要由 Node.Js 语言开发，贡献者超过 80人，已有超过 5000 次提交。该项目已经停止维护（End of Life）。

项目的邮件列表地址为 composer@lists.hyperledger.org。

### Burrow 项目

![Hyperledger Burrow 项目](_images/burrow.png)

Burrow 项目由 Monax、Intel 等企业于 2017 年 4 月提交到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/burrow。

该项目的前身为 eris-db，基于 Go 语言实现的以太坊虚拟机，目前贡献者超过 20人，提交次数已经超过 2000 次。

Burrow 项目提供了支持以太坊虚拟机的智能合约区块链平台，并支持 Proof-of-Stake 共识机制（Tendermint）和权限管理，可以提供快速的区块链交易。

项目的邮件列表地址为 burrow@lists.hyperledger.org。

### Quilt 项目

![Hyperledger Quilt 项目](_images/quilt.png)

Quilt 项目由 NTT、Ripple 等企业于 2017 年 10 月提交到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/quilt。

Quilt 项目前身为 W3C 支持的 Interledger 协议的 Java 实现，主要试图为转账服务提供跨多个区块链平台的支持。目前贡献者超过 10人，提交次数已经超过 100 次。

项目的邮件列表地址为 quilt@lists.hyperledger.org。

### Caliper 项目

![Hyperledger Caliper 项目](_images/caliper.png)

Caliper 项目由华为于 2018 年 3 月提交到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/caliper。

Caliper 项目希望能为评测区块链的性能（包括吞吐、延迟、资源使用率等）提供统一的工具套装，主要基于 Node.js 语言实现，支持对 Fabric、Sawtooth、Burrow 等项目进行性能测试。目前贡献者超过 20人，提交次数超过 400 次。

项目的邮件列表地址为 caliper@lists.hyperledger.org。

### Ursa 项目

![Hyperledger Ursa 项目](_images/ursa.png)

Ursa 项目前身为加密实现库项目，由 Fujitsu、Sovrin、Intel、DFINITY、State Street、IBM、Bitwise IO 等企业的开发者于 2018 年 11 月正式贡献到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/ursa。

Ursa 项目曾经提供一套方便、安全的密码学软件库（包括加解密、零知识证明等），为实现区块链平台实现提供便利。主要基于 Rust 语言实现，目前包括两个子组件（基础密码实现库 Base Crypto 和零知识证明库 Z-Mix）。参与贡献者超过 10人，提交次数超过 400 次。项目已于 2023 年 4 月结束生命周期。

项目的邮件列表地址为 ursa@lists.hyperledger.org。

### Grid 项目

![Hyperledger Grid 项目](_images/grid.png)

Grid 项目由 Cargill、Intel 和 Bitwise IO 公司于 2018 年 12 月提交到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/grid。

Grid 项目为开发基于区块链的供应链场景应用提供框架支持和参考实现，包括智能合约、数据模型、领域模型、样例应用等。
，主要基于 Python 语言实现，并使用 Sabre（基于 WebAssembly/WASM 的智能合约引擎）来运行智能合约。目前贡献者超过 40人，提交次数超过 5000 次。

项目的邮件列表地址为 grid@lists.hyperledger.org。

### Transact 项目

![Hyperledger Transact 项目](_images/transact.png)

Transact 项目由 Bitwise IO、Cargill、Intel、IBM、HACERA 等公司于 2019 年 5 月提交到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/transact。

Transact 项目为区块链提供交易执行的平台和代码库，其他的框架性项目可以基于 Transact 来管理交易的执行过程和状态。Transact 项目试图打造一个通用的智能合约引擎来支持包括 EVM、WebAssembly 等合约的运行。目前包括 transact、transact-rfcs、transact-contrib 等子项目。

项目的邮件列表地址为 transact@lists.hyperledger.org。

### Aries 项目

![Hyperledger Aries 项目](_images/aries.png)

Aries 项目由 Sovrin、C3I、和 Evernym 等公司于 2019 年 5 月提交到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/aries。

Aries 项目希望能为客户端提供共享的密码学钱包和相关的代码库（包括零知识证明），以及对于链下交互的消息支持，以简化区块链客户端的开发。

项目的邮件列表地址为 aries@lists.hyperledger.org。

### Besu 项目

![Hyperledger Besu 项目](_images/besu.png)

* Besu 项目由 Consensys、Hacera、JPM 和 Redhat 等公司于 2019 年 8 月正式贡献到社区。核心代码在 Github 上地址为 https://github.com/hyperledger/besu。

Besu 提供对以太坊协议的支持，由 Java 实现。

项目的邮件列表地址为 besu@lists.hyperledger.org。

### Avalon 项目

![Hyperledger Avalon 项目](_images/avalon.png)

* Avalon 项目由 Intel、IEX、IBM 和 Consensys 等公司于 2019 年 9 月正式贡献到社区。主要提供链下的安全计算支持，重点考虑了安全性和可扩展性。项目核心代码在 https://github.com/hyperledger/avalon。

项目的邮件列表地址为 avalon@lists.hyperledger.org。
