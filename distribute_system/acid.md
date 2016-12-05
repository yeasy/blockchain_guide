## ACID 原则
即 Atomicity（原子性）、Consistency（一致性）、Isolation（隔离性）、Durability（持久性）。

ACID 原则描述了对分布式数据库的一致性需求，同时付出了可用性的代价。

* Atomicity：每次操作是原子的，要么成功，要么不执行；
* Consistency：数据库的状态是一致的，无中间状态；
* Isolation：各种操作彼此互相不影响；
* Durability：状态的改变是持久的，不会失效。

一个与之相对的原则是 BASE（Basic Availiability，Soft state，Eventually Consistency），牺牲掉对一致性的约束（最终一致性），来换取一定的可用性。

