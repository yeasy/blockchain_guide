
# 区块链技术指南
0.7.9

区块链技术是金融科技（Fintech）领域的一项重要技术创新。

作为去中心化记账（Decentrialized Ledger Technology，DLT）平台的核心技术，区块链被认为在金融、征信、物联网、经济贸易结算、资产管理等众多领域都拥有广泛的应用前景。

区块链技术自身尚处于快速发展的初级阶段，现有区块链系统在设计和实现中利用了分布式系统、密码学、博弈论、网络协议等诸多学科的知识，为学习原理和实践应用都带来了不小的挑战。

目前该领域尚缺乏一本较为系统的技术资料。本书希望可以探索区块链概念的来龙去脉，剥茧抽丝，剖析关键技术原理，同时讲解实践应用。

在参与相关开源项目，以及编写区块链云服务平台的过程中，笔者积累了一些实践经验，也通过本书一并分享出来，希望能推动区块链技术的早日成熟和更多应用场景的出现。

本书适用于对区块链技术感兴趣，且具备一定信息和金融基础知识的读者；无技术背景的读者也可以从中了解到区块链的应用现状。

在线阅读：[GitBook](https://www.gitbook.com/book/yeasy/blockchain_guide) 或 [GitHub](https://github.com/yeasy/blockchain_guide/blob/master/SUMMARY.md)。

* pdf 版本 [下载](https://www.gitbook.com/download/pdf/book/yeasy/blockchain_guide)
* epub 版本 [下载](https://www.gitbook.com/download/epub/book/yeasy/blockchain_guide)

欢迎大家加入区块链技术讨论群：

* QQ 群   I：335626996（已满）
* QQ 群  II：523889325（已满）
* QQ 群 III：414919574（已满）
* QQ 群  IV：364824846（可加）

## 版本历史

* 0.8.0: 2016-XX-YY
  * 完善应用场景等；
  * 完善分布式系统技术；
  * 完善密码学技术；
  * 根据最新代码更新 Hyperledger 使用。
* 0.7.0: 2016-09-10
  * 完善一致性技术等；
  * 修正文字。
* 0.6.0: 2016-08-05
  * 修改文字；
  * 增加更多智能合约；
  * 增加更多业务场景。
* 0.5.0: 2016-07-10
  * 增加 Hyperledger 项目的内容；
  * 增加以太坊项目内容；
  * 增加闪电网络介绍、关键技术剖析；
  * 补充区块链即服务；
  * 增加比特币项目。
* 0.4.0: 2016-06-02
    * 添加应用场景分析。
* 0.3.0: 2016-05-12
    * 添加数字货币问题分析。
* 0.2.0: 2016-04-07
    * 添加 Hyperledger 项目简介。
* 0.1.0: 2016-01-17
    * 添加区块链简介。

## 参与贡献
贡献者 [名单](https://github.com/yeasy/blockchain_guide/graphs/contributors)。

区块链技术自身仍在快速发展中，生态环境也在蓬勃成长。

本书源码开源托管在 Github 上，欢迎参与维护：[github.com/yeasy/blockchain_guide](https://github.com/yeasy/blockchain_guide)。

首先，在 GitHub 上 `fork` 到自己的仓库，如 `docker_user/blockchain_guide`，然后 `clone` 到本地，并设置用户信息。

```sh
$ git clone git@github.com:docker_user/blockchain_guide.git
$ cd blockchain_guide
$ git config user.name "yourname"
$ git config user.email "your email"
```

更新内容后提交，并推送到自己的仓库。

```sh
$ #do some change on the content
$ git commit -am "Fix issue #1: change helo to hello"
$ git push
```

最后，在 GitHub 网站上提交 pull request 即可。

另外，建议定期使用项目仓库内容更新自己仓库内容。
```sh
$ git remote add upstream https://github.com/yeasy/blockchain_guide
$ git fetch upstream
$ git checkout master
$ git rebase upstream/master
$ git push -f origin master
```
