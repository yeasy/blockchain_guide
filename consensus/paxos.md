## Paxos 与 Raft

### Paxos

1990 年由 Leslie Lamport 提出的 [Paxos](http://research.microsoft.com/users/lamport/pubs/lamport-paxos.pdf) 一致性算法，在工程角度实现了一种最大化保障一致性（极小的概率无法实现一致性）的机制。

Paxos 是第一个被证明的一致性算法，其原理是现在一致性算法设计的鼻祖，然而以复杂难懂出名。

Paxos 被应用在 Chubby、ZooKeeper 这样的系统中。

### Raft

[Raft](https://ramcloud.atlassian.net/wiki/download/attachments/6586375/raft.pdf) 算法是Paxos 算法的一种简化实现。

包括三种角色：leader、candiate 和 follower，其基本过程为：

* Leader 选举：每个 candidate 随机经过一定时间都会提出选举方案，最近阶段中得票最多者被选为 leader；
* 同步 log：leader 会找到系统中 log 最新的记录，并强制所有的 follower 来刷新到这个记录；

*注：此处 log 并非是指日志消息，而是各种事件的发生记录。*