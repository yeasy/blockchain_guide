## 消息协议

节点之间通过消息来进行交互，所有消息都由下面的数据结构来实现。

```protobuf
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
```

消息分为四大类：Discovery（探测）、Transaction（交易）、Synchronization（同步）、Consensus（一致性）。

不同消息类型，对应到 payload 中数据不同，分为对应的子类消息结构。

### Discovery

包括 DISC_HELLO、DISC_GET_PEERS、DISC_PEERS。

DISC_HELLO 消息结构如下。

```protobuf
message HelloMessage { PeerEndpoint peerEndpoint = 1; uint64 blockNumber = 2;}message PeerEndpoint { PeerID ID = 1; string address = 2; enum Type { UNDEFINED = 0; VALIDATOR = 1; NON_VALIDATOR = 2; } Type type = 3; bytes pkiID = 4;}

message PeerID { string name = 1;}
```

节点新加入网络时，会向 `CORE_PEER_DISCOVERY_ROOTNODE` 发送 `DISC_HELLO` 消息，汇报本节点的信息（id、地址、block 数、类型等），开始探测过程。

探测后发现 block 数落后对方，则会触发同步过程。

之后，定期发送 `DISC_GET_PEERS` 消息，获取新加入的节点信息。收到 `DISC_GET_PEERS` 消息的节点会通过 `DISC_PEERS` 消息返回自己知道的节点列表。

### Transaction

包括 Deploy、Invoke、Query。消息结构如下：

```protobuf
message Transaction { enum Type { UNDEFINED = 0; CHAINCODE_DEPLOY = 1; CHAINCODE_INVOKE = 2; CHAINCODE_QUERY = 3; CHAINCODE_TERMINATE = 4; } Type type = 1; string uuid = 5; bytes chaincodeID = 2; bytes payloadHash = 3;

 ConfidentialityLevel confidentialityLevel = 7; bytes nonce = 8; bytes cert = 9; bytes signature = 10;

 bytes metadata = 4; google.protobuf.Timestamp timestamp = 6;}

message TransactionPayload { bytes payload = 1;}

enum ConfidentialityLevel { PUBLIC = 0; CONFIDENTIAL = 1;}
```

### Synchronization
当节点发现自己 block 落后网络中最新状态，则可以通过发送如下消息（由 consensus 策略决定）来获取对应的返回。

* SYNC_GET_BLOCKS（对应 SYNC_BLOCKS）：获取给定范围内的 block 数据；
* SYNC_STATE_GET_SNAPSHOT（对应 SYNC_STATE_SNAPSHOT）：获取最新的世界观快照；
* SYNC_STATE_GET_DELTAS（对应 SYNC_STATE_DELTAS）：获取某个给定范围内的 block 对应的状态变更。

### Consensus

consensus 组件收到 `CHAIN_TRANSACTION` 类消息后，将其转换为 `CONSENSUS` 消息，然后向所有的 VP 节点广播。

收到 `CONSENSUS` 消息的节点会按照预定的 consensus 算法进行处理。
