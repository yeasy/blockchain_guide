## 监听网络事件

用户可能注意到，发往网络的请求是异步处理模式，这就意味着客户端无法获知提交的交易是否最终接受。Fabric 在 Peer 节点上提供了事件 gRPC 服务，用户可以通过客户端来监听。

下面通过 eventsclient 工具来监听网络中的事件。

首先通过如下命令安装 eventsclient 工具。

```bash
$ cd $GOPATH/src/hyperledger/fabric/examples/events/eventsclient
$ go install && go clean
```

该工具自动封装对 Peer 事件的 gRPC 请求，支持的选项主要包括如下几个：

* -server "localhost:7053"：监听服务地址，一般指定为 Peer 节点的 7053 端口；
* -channelID string：监听指定通道信息，默认为 testchainid；
* -seek int：指定从哪个区块开始监听。-2代表从初始区块（默认），-1代表从当前最新区块；
* -filtered=true：只获取过滤的区块内容，不显示完整内容。
* -quiet：不打印区块内容，只显示区块号；
* -tls：是否启用 TLS，默认关闭；
* -rootCert string：启用 TLS 时指定信任的根 CA 证书路径；
* -mTls：是否开启双向验证（即服务端也同时验证客户端身份），默认关闭。
* -clientCert string：：启用 TLS 时候客户端证书路径；
* -clientKey string：启用 TLS 时候客户端私钥路径；

典型地，用户可以通过环境变量指定所需的参数值，使用如下命令启动监听。

```bash
$ eventsclient \
   -server=${PEER_URL} \
   -channelID=${APP_CHANNEL} \
   -filtered=true \
   -tls=true \
   -clientKey=${TLS_CLIENT_KEY} \
   -clientCert=${TLS_CLIENT_CERT} \
   -rootCert=${TLS_CA_CERT}
```

启动后，该工具会持续监听来自指定通道的事件，并打印出来。

例如，监听 businesschannel 通道内区块信息，并对结果进行过滤输出，命令和结果如下所示。

```bash
$ CORE_PEER_LOCALMSPID=Org1MSP \
CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp \
eventsclient \
  -server=peer0.org1.example.com:7051 \
  -channelID=businesschannel \
  -filtered=true \
  -tls=true \
  -clientKey=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@Org1.example.com/tls/client.key \
  -clientCert=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@Org1.example.com/tls/client.crt \
  -rootCert=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@Org1.example.com/tls/ca.crt

UTC [eventsclient] readEventsStream -> INFO 001 Received filtered block:
{
  "channel_id": "businesschannel",
  "filtered_transactions": [
    {
      "tx_validation_code": "VALID",
      "txid": "",
      "type": "CONFIG"
    }
  ],
  "number": "0"
}
UTC [eventsclient] readEventsStream -> INFO 002 Received filtered block:
{
  "channel_id": "businesschannel",
  "filtered_transactions": [
    {
      "tx_validation_code": "VALID",
      "txid": "",
      "type": "CONFIG"
    }
  ],
  "number": "1"
}
UTC [eventsclient] readEventsStream -> INFO 003 Received filtered block:
{
  "channel_id": "businesschannel",
  "filtered_transactions": [
    {
      "tx_validation_code": "VALID",
      "txid": "",
      "type": "CONFIG"
    }
  ],
  "number": "2"
}
UTC [eventsclient] readEventsStream -> INFO 004 Received filtered block:
{
  "channel_id": "businesschannel",
  "filtered_transactions": [
    {
      "transaction_actions": {
        "chaincode_actions": []
      },
      "tx_validation_code": "VALID",
      "txid": "2832892094f612237b06950b77a6afc13ca9226176e99c2a8577cf4be2074c0a",
      "type": "ENDORSER_TRANSACTION"
    }
  ],
  "number": "3"
}
UTC [eventsclient] readEventsStream -> INFO 005 Received filtered block:
{
  "channel_id": "businesschannel",
  "filtered_transactions": [
    {
      "transaction_actions": {
        "chaincode_actions": []
      },
      "tx_validation_code": "VALID",
      "txid": "fec547335060bb324e8e4a08067c7fa24092e1295cb62dffb14a93bc77b2fbcf",
      "type": "ENDORSER_TRANSACTION"
    }
  ],
  "number": "4"
}
...
```

