## 术语

### 通用术语

* Blockchain：区块链，基于密码学的可实现信任化的信息存储和处理技术。
* CA：Certificate Authority，负责证书的创建、颁发，在 PKI 体系中最为核心的角色。
* Chaincode：链上代码，运行在区块链上提前约定的智能合约，支持多种语言实现。
* Distributed Ledger：分布式记账本，大家都认可的去中心化的账本记录平台。
* DLT：Distributed Ledger Technology。
* DTCC：Depository Trust and Clearing Corporation，存托和结算公司，全球最大的金融交易后台服务机构。
* Fintech：Financial Technology，跟金融相关的（信息）技术。
* Hash：哈希算法，任意长度的二进制值映射为较短的固定长度的二进制值的算法。
* Lightning Network：闪电网络，通过链外的微支付通道来增大交易吞吐量的技术。
* Nonce：密码学术语，表示一个临时的值，多为随机字符串。
* P2P：点到点的通信网络，网络中所有节点地位均等，不存在中心化的控制机制。
* PKI：Public key infrastructure，基于公钥体系的安全基础架构。
* Smart Contract：智能合约，运行在区块链上提前约定的合同；
* Sybil Attack（女巫攻击）：少数节点通过伪造或盗用身份伪装成大量节点，进而对分布式系统系统进行破坏。
* SWIFT：Society for Worldwide Interbank Financial Telecommunication，环球银行金融电信协会，运营世界金融电文网络，服务银行和金融机构。
* 挖矿：通过暴力尝试来找到一个字符串，使得它加上一组交易信息后的 hash 值符合特定规则（例如前缀包括若干个 0），找到的人可以宣称新区块被发现，并获得系统奖励的比特币。
* 矿工：参与挖矿的人或组织。
* 矿机：专门为比特币挖矿而设计的设备，包括 GPU、专用芯片等。
* 矿池：采用团队协作方式来集中算力进行挖矿，对产出的比特币进行分配。
* 市场深度：未成交的交易，衡量市场承受大额交易后汇率的稳定能力。
* 图灵完备：指一个机器或装置能用来模拟图灵机（现代通用计算机的雏形）的功能，图灵完备的机器在可计算性上等价。

### 比特币、以太坊相关术语
* Bitcoin：比特币，中本聪发起的数字货币技术。
* DAO：Decentralized Autonomous Organization，分布式自治组织，基于区块链的按照智能合约联系起来的松散众筹群体。
* PoW：Proof of Work，工作量证明，在一定难题前提下求解一个 SHA256 的 hash 问题。

### Hyperledger 相关术语

* Auditability（审计性）：在一定权限和许可下，可以对链上的交易进行审计和检查。
* Block（区块）：代表一批得到确认的交易信息的整体，准备被共识加入到区块链中。
* Blockchain（区块链）：由多个区块链接而成的链表结构，除了首个区块，每个区块都包括前继区块内容的 hash 值。
* Chaincode（链码）：区块链上的应用代码，扩展自“智能合约”概念，支持 golang、nodejs 等。
* Committer（提交节点）：1.0 架构中一种 peer 节点角色，负责对 orderer 排序后的交易进行检查，选择合法的交易执行并写入存储。
* Confidentiality（保密）：只有交易相关方可以看到交易内容，其它人未经授权则无法看到。
* Endorser（推荐节点）：1.0 架构中一种 peer 节点角色，负责检验某个交易是否合法，是否愿意为之背书、签名。
* Ledger（账本）：包括区块链结构（带有所有的交易信息）和当前的世界观（world state）。
* MSP（Member Service Provider，成员服务提供者）：成员服务的抽象访问接口，实现对不同成员服务的可拔插支持。
* Non-validating Peer（非验证节点）：不参与账本维护，仅作为交易代理响应客户端的 REST 请求，并对交易进行一些基本的有效性检查，之后转发给验证节点。
* Orderer（排序节点）：1.0 架构中的共识服务角色，负责排序看到的交易，提供全局确认的顺序。
* Permissioned Ledger（带权限的账本）：网络中所有节点必须是经过许可的，非许可过的节点则无法加入网络。
* Privacy（隐私保护）：交易员可以隐藏交易的身份，其它成员在无特殊权限的情况下，只能对交易进行验证，而无法获知身份信息。
* Transaction（交易）：执行账本上的某个函数调用。具体函数在 chaincode 中实现。
* Transactor（交易者）：发起交易调用的客户端。
* Validating Peer（验证节点）：维护账本的核心节点，参与一致性维护、对交易的验证和执行。
* World State（世界观）：是一个键值数据库，chaincode 用它来存储交易相关的状态。