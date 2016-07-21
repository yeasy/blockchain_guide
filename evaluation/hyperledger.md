## Hyperledger fabric 性能评测

### 环境配置
| 类型  |     操作系统     | 内核版本 | CPU(GHz) | 内存(GB) |
| :--: | :-------------: | :-----: | :------: | :-----: |
| 物理机 | Ubuntu 14.04.1 | 3.16.0-71-generic | 4x2.0 | 8 |

每个集群启动后等待 10s 以上，待状态稳定。

仅测试单客户端、单服务端的连接性能情况。

### 评测指标

一般评测系统性能指标包括吞吐量（throughput）和延迟（latency）。对于区块链平台系统来说，实际交易延迟包括客户端到系统延迟（往往经过互联网），再加上系统处理反馈延迟（跟不同 consensus 算法关系很大，跟集群之间互联系统关系也很大）。

本次测试仅给出大家最为关注的交易吞吐量（tps）。

### 结果

#### query 交易

##### noops
| clients | VP Nodes | iteration |   tps  |
| -------- | ------- | --------- | ------ |
|    1     |    1    |    2000   | 195.50 |
|    1     |    4    |    2000   | 187.09 |

##### pbft:classic
| clients | VP Nodes | iteration |   tps  |
| -------- | ------- | --------- | ------ |
|    1     |    4    |    2000   | 193.05 |

##### pbft:batch
| clients | VP Nodes | batch size | iteration |   tps  |
| -------- | ------- | --------  | ---------- | ------ |
|    1     |    4    |    2      |    2000    | 193.99 |
|    1     |    4    |    4      |    2000    | 192.49 |
|    1     |    4    |    8      |    2000    | 192.68 |

##### pbft:sieve
| clients | VP Nodes | iteration |   tps  |
| -------- | ------- | --------- | ------ |
|    1     |    4    |    2000   | 192.86 |

#### invoke 交易

##### noops

| clients | VP Nodes | iteration |   tps  |
| -------- | ------- | --------- | ------ |
|   1      |    1    |    2000   | 298.51 |
|   1      |    4    |    2000   | 205.76 |

##### pbft:classic
| clients | VP Nodes | iteration |  tps   |
| -------- | ------- | --------- | ------ |
|    1     |    4    |    2000   | 141.34 |


##### pbft:batch
| clients | VP Nodes | batch size | iteration |   tps  |
| -------- | ------- | ---------  | --------- | ------ |
|    1     |    4    |     2      |    2000   | 214.36 |
|    1     |    4    |     4      |    2000   | 227.53 |
|    1     |    4    |     8      |    2000   | 237.81 |


##### pbft:sieve
| clients | VP Nodes | iteration |   tps  |
| -------- | ------- | --------- | ------ |
|    1     |    4    |    2000   | 253.49* |

*注：sieve 算法目前在所有交易完成后较长时间内并没有取得最终的结果，出现大量类似“vp0_1  | 07:49:26.388 [consensus/obcpbft] main -> WARN 23348 Sieve replica 0 custody expired, complaining: 3kwyMkdCSL4rbajn65v+iYWyJ5aqagXvRR9QU8qezpAZXY4y6uy2MB31SGaAiaSyPMM77TYADdBmAaZveM38zA==”警告信息。*

### 结论
单客户端连接情况下，tps 基本在 190 ~ 300 范围内。