## Paxos 与 Raft

Paxos 最初设计为解决存在故障，但不存在恶意节点（无伪造消息，但可能丢失或重复）场景下的一致性问题。

Raft 是对 Paxos 的重新设计和实现。

### Paxos

1990 年由 Leslie Lamport 提出的 [Paxos](http://research.microsoft.com/users/lamport/pubs/lamport-paxos.pdf) 一致性算法，在工程角度实现了一种最大化保障一致性（存在极小的概率无法实现一致性）的机制。

故事背景是古希腊 Paxon 岛上的多个法官在一个大厅内对一个议案进行表决，如何达成统一的结果。他们之间通过服务人员来传递纸条，但法官可能离开或进入大厅，服务人员可能偷懒去睡觉。

Paxos 是第一个被证明的一致性算法，其原理基于两阶段提交并进行扩展。

作为现在一致性算法设计的鼻祖，算法以复杂难懂出名。算法中将节点分为三种类型：

* proposer：提出一个提案，等待大家批准为结案；
* acceptor：负责对提案进行投票；
* learner：被告知结案结果，并与之统一。

基本过程包括 proposer 提出提案，先争取多个 acceptor 的支持，超过一半支持时，则发送结案结果给所有人。一个潜在的问题是 proposer 在此过程中出现故障，可以通过超时机制来解决。极为凑巧的情况下，每次新的一轮提案的 proposer 都恰好故障，系统则永远无法达成一致（概率很小）。

Paxos 能保证在超过一半的正常节点存在时，系统能达成一致。

Paxos 被应用在 Chubby、ZooKeeper 这样的系统中。

### Raft

[Raft](https://ramcloud.atlassian.net/wiki/download/attachments/6586375/raft.pdf) 算法是Paxos 算法的一种简化实现。

包括三种角色：leader、candiate 和 follower，其基本过程为：

* Leader 选举：每个 candidate 随机经过一定时间都会提出选举方案，最近阶段中得票最多者被选为 leader；
* 同步 log：leader 会找到系统中 log 最新的记录，并强制所有的 follower 来刷新到这个记录；

*注：此处 log 并非是指日志消息，而是各种事件的发生记录。*