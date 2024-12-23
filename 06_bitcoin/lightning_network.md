## 闪电网络

比特币的交易网络最为人诟病的一点便是交易性能：全网每秒 7 笔左右的交易速度，远低于传统的金融交易系统；同时，等待 6 个块的可信确认将导致约 1 个小时的最终确认时间。

为了提升性能，社区提出了闪电网络等创新的设计。

闪电网络的主要思路十分简单——将大量交易放到比特币区块链之外进行，只把关键环节放到链上进行确认。该设计最早于 2015 年 2 月在论文《The Bitcoin Lightning Network: Scalable Off-Chain Instant Payments》中提出。

比特币的区块链机制自身已经提供了很好的可信保障，但是相对较慢；另一方面考虑，对于大量的小额交易来说，是否真需要这么高的可信性？

闪电网络主要通过引入智能合约的思想来完善链下的交易渠道。核心的概念主要有两个：RSMC（Recoverable Sequence Maturity Contract）和 HTLC（Hashed Timelock Contract）。前者解决了链下交易的确认问题，后者解决了支付通道的问题。

### RSMC

Recoverable Sequence Maturity Contract，即“可撤销的顺序成熟度合同”。这个词很绕，其实主要原理很简单，类似资金池机制。

首先假定交易双方之间存在一个“微支付通道”（资金池）。交易双方先预存一部分资金到“微支付通道”里，初始情况下双方的分配方案等于预存的金额。每次发生交易，需要对交易后产生资金分配结果共同进行确认，同时签字把旧版本的分配方案作废掉。任何一方需要提现时，可以将他手里双方签署过的交易结果写到区块链网络中，从而被确认。从这个过程中可以可以看到，只有在提现时候才需要通过区块链。

任何一个版本的方案都需要经过双方的签名认证才合法。任何一方在任何时候都可以提出提现，提现时需要提供一个双方都签名过的资金分配方案（意味着肯定是某次交易后的结果，被双方确认过，但未必是最新的结果）。在一定时间内，如果另外一方拿出证明表明这个方案其实之前被作废了（非最新的交易结果），则资金罚没给质疑方；否则按照提出方的结果进行分配。罚没机制可以确保了没人会故意拿一个旧的交易结果来提现。

另外，即使双方都确认了某次提现，首先提出提现一方的资金到账时间要晚于对方，这就鼓励大家尽量都在链外完成交易。通过 RSMC，可以实现大量中间交易发生在链外。

### HTLC

微支付通道是通过 Hashed Timelock Contract 来实现的，中文意思是“哈希的带时钟的合约”。这个其实就是限时转账。理解起来也很简单，通过智能合约，双方约定转账方先冻结一笔钱，并提供一个哈希值，如果在一定时间内有人能提出一个字符串，使得它哈希后的值跟已知值匹配（实际上意味着转账方授权了接收方来提现），则这笔钱转给接收方。

不太恰当的例子，约定一定时间内，有人知道了某个暗语（可以生成匹配的哈希值），就可以拿到这个指定的资金。

假设 Alice 和 Bob 想进行跨链或链上资产交换：

* Alice 是交易的发起者，生成 secret。Bob 是交易的参与者，需要使用 Alice 的 secret 解锁资产。
* Alice 生成一个随机密钥（secret），并计算其哈希值 H。
* Alice 部署 HTLC 合约，将资金锁定到合约中，解锁条件是：任何人提供满足 Hash(secret) = H 的 secret，即可解锁资产。同时设置时间锁，超时后资金退还给 Alice。
* Bob 观察到 Alice 提供的哈希值 H，并在链上（或跨链）部署一个自己的 HTLC 合约，将他的资金锁定，并使用同一个哈希值 H 作为解锁条件。
* Alice 提供 secret 原像（例如 "abc123"）到 Bob 的 HTLC 合约上，满足了哈希值要求，获得资产。
* Bob 通过链上记录看到 Alice 提供的 secret，在 Alice 部署的 HTLC 合约上，使用 secret 解锁资金，获得资产。

如果任意一人不提供正确的信息，则智能合约在超时后，会将资产退回。从上面过程可以看出，Bob 的合约超时时间，要短于 Alice 的合约超时时间。

HTLC 机制还可以扩展到多个人的场景。例如三个人的场景，甲想转账给丙，丙先发给甲一个哈希值。甲可以先跟乙签订一个合同，如果你在一定时间内能告诉我一个暗语，我就给你多少钱。乙于是跑去跟丙签订一个合同，如果你告诉我那个暗语，我就给你多少钱。丙于是告诉乙暗语，拿到乙的钱，乙又从甲拿到钱。最终达到结果是甲转账给丙。这样甲和丙之间似乎构成了一条完整的虚拟的“支付通道”。

### 闪电网络

RSMC 保障了两个人之间的直接交易可以在链下完成，HTLC 保障了任意两个人之间的转账都可以通过一条“支付”通道来完成。闪电网络整合这两种机制，就可以实现任意两个人之间的交易都在链下完成了。

在整个交易中，智能合约起到了中介的重要角色，而区块链网络则确保最终的交易结果被确认。

