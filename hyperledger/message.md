## 消息协议

节点之间通过消息来进行交互，所有消息都由下面的数据结构来实现。

message Message {
   enum Type {
        UNDEFINED = 0;

        DISC_HELLO = 1;
        DISC_DISCONNECT = 2;
        DISC_GET_PEERS = 3;
        DISC_PEERS = 4;
        DISC_NEWMSG = 5;

        CHAIN_STATUS = 6;
        CHAIN_TRANSACTION = 7;
        CHAIN_GET_TRANSACTIONS = 8;
        CHAIN_QUERY = 9;

        SYNC_GET_BLOCKS = 11;
        SYNC_BLOCKS = 12;
        SYNC_BLOCK_ADDED = 13;

        SYNC_STATE_GET_SNAPSHOT = 14;
        SYNC_STATE_SNAPSHOT = 15;
        SYNC_STATE_GET_DELTAS = 16;
        SYNC_STATE_DELTAS = 17;

        RESPONSE = 20;
        CONSENSUS = 21;
    }
    Type type = 1;
    bytes payload = 2;
    google.protobuf.Timestamp timestamp = 3;
}

消息分为四大类：Discovery（探测）、Transaction（交易）、Synchronization（同步）、Consensus（一致性）。

不同消息类型，对应到 payload 中数据不同，分为为对应的子类消息结构。

### Discovery

包括 DISC_HELLO、DISC_GET_PEERS、DISC_PEERS。

节点新加入网络时，会向 `CORE_PEER_DISCOVERY_ROOTNODE` 发送 `DISC_HELLO` 消息，汇报本节点的信息（id、地址、block 数、类型等），开始探测过程。

探测后发现 block 数落后对方，则会触发同步过程。

之后，定期发送 `DISC_GET_PEERS` 消息，获取新加入的节点信息。收到 `DISC_GET_PEERS` 消息的节点会通过 `DISC_PEERS` 消息返回自己知道的节点列表。

### Transaction

包括 Deploy、Invoke、Query。

### Synchronization
SYNC_GET_BLOCKS 和对应的 SYNC_BLOCKS。

SYNC_STATE_GET_SNAPSHOT 和对应的 SYNC_STATE_SNAPSHOT。

SYNC_STATE_GET_DELTAS 和对应的 SYNC_STATE_DELTAS。

### Consensus

CONSENSUS 消息。