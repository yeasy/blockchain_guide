## 自动探测网络信息

客户端要往 Fabric 网络中发送请求，首先需要知道网络的相关信息，如网络中成员、背书节点的地址、链码安装信息等。

在 Fabric v1.2.0 版本之前这些信息需要使用者来手动采集提供。这种方式下需要提前指定，容易出错，另外是当网络中信息变更后（如节点上下线）还需要再次更新。

为了解决这些问题，社区自 v1.2.0 版本开始在 Peer 节点上提供了 Discovery 服务，并编写了 discover 客户端工具（位于 discovery/cmd）。该工具可以访问 Peer 节点提供的 Discovery 服务，查询获取指定信息。

### 主要功能

discover 工具目前提供如下的查询功能：

* 节点信息查询：使用 `peers` 子命令查询节点的身份、服务等信息；
* 通道配置：使用 `config` 子命令查询通道的配置信息，包括成员组织、排序服务信息等；
* 链码背书信息：使用 `endorsers` 子命令查询对某个链码可以进行背书的节点信息。

命令使用格式为 `discover [全局参数] <子命令> [子命令参数列表]`。

### 全局参数

discover 支持的全局参数和相关说明如下：

```bash
* --help：输出帮助信息；
* --configFile=CONFIGFILE：指定从配置文件中载入参数配置，则无需从命令行指定参数；
* --peerTLSCA=PEERTLSCA：指定校验 peer 端 TLS 的 CA 证书；
* --tlsCert=TLSCERT：指定客户端使用的 TLS 证书（可选，当 Peer 校验客户端 TLS 时）；
* --tlsKey=TLSKEY：指定客户端使用的 TLS 私钥（可选，当 Peer 校验客户端 TLS 时）；
* --userKey=USERKEY：客户端签名私钥；
* --userCert=USERCERT：客户端签名证书；
* --MSP=MSP：指定客户端的 MSP ID。
```

### 子命令

discover 目前支持四个子命令：`peers`、`config`、`endorsers`、`saveConfig`，可以通过 `help <子命令>` 来查看各子命令的功能和使用方法。

#### peers 子命令

显示网络中的 Peer 节点信息，包括它们的 MSP Id、gRPC 服务监听地址和身份证书。

命令格式为 `peers [参数列表]`，支持参数如下：

```bash
* --server=SERVER：指定命令连接的 Peer 节点地址；
* --channel=CHANNEL：指定查询某个特定通道内的节点信息。
```

例如，通过 peer0.org1.example.com 节点查询 businesschannel 通道内的 Peer 节点信息，可以执行如下命令：

```bash
$ discover \
    --peerTLSCA tls/ca.crt \
    --userKey msp/keystore/f76cf3c92dac81103c82d5490c417ac0123c279f93213f65947d8cc69e11fbc5_sk \
    --userCert msp/signcerts/Admin\@org1.example.com-cert.pem \
    --MSP Org1MSP \
    --tlsCert tls/client.crt \
    --tlsKey tls/client.key \
    peers \
    --server peer0.org1.example.com:7051 \
    --channel businesschannel

[
	{
		"MSPID": "Org2MSP",
		"Endpoint": "peer1.org2.example.com:7051",
		"Identity": "-----BEGIN CERTIFICATE-----\nMIICKD...pVTw==\n-----END CERTIFICATE-----\n"
	},
	{
		"MSPID": "Org2MSP",
		"Endpoint": "peer0.org2.example.com:7051",
		"Identity": "-----BEGIN CERTIFICATE-----\nMIICKT...cGaA=\n-----END CERTIFICATE-----\n"
	},
	{
		"MSPID": "Org1MSP",
		"Endpoint": "peer0.org1.example.com:7051",
		"Identity": "-----BEGIN CERTIFICATE-----\nMIICKD...mgaA==\n-----END CERTIFICATE-----\n"
	},
	{
		"MSPID": "Org1MSP",
		"Endpoint": "peer1.org1.example.com:7051",
		"Identity": "-----BEGIN CERTIFICATE-----\nMIICKD...UO+g==\n-----END CERTIFICATE-----\n"
	}
]
```

结果显示 businesschannel 通道内目前包括属于 2 个组织的 4 个 Peer 节点成员：

* Org1MSP
  * peer0.org1.example.com
  * peer1.org1.example.com
* Org2MSP
  * peer0.org2.example.com
  * peer1.org2.example.com

#### config 子命令

显示网络中的通道配置信息，包括各个组织的 MSP 信息和排序节点信息。

命令格式为 `config [<参数>]`，支持参数如下：

```bash
* --server=SERVER：指定命令连接的 Peer 节点地址；
* --channel=CHANNEL：指定查询某个特定通道内的配置信息。
```

例如，通过 peer0.org1.example.com 节点查询 businesschannel 通道内的配置信息，可以执行如下命令：

```bash
$ discover \
    --peerTLSCA tls/ca.crt \
    --userKey msp/keystore/f76cf3c92dac81103c82d5490c417ac0123c279f93213f65947d8cc69e11fbc5_sk \
    --userCert msp/signcerts/Admin\@org1.example.com-cert.pem \
    --MSP Org1MSP \
    --tlsCert tls/client.crt \
    --tlsKey tls/client.key \
    config \
    --server peer0.org1.example.com:7051 \
    --channel businesschannel

{
	"msps": {
		"OrdererMSP": {
			"name": "OrdererMSP",
			"root_certs": [
				"LS0tLS...tLQo="
			],
			"admins": [
				"LS0tLS...LS0K"
			],
			"crypto_config": {
				"signature_hash_family": "SHA2",
				"identity_identifier_hash_function": "SHA256"
			},
			"tls_root_certs": [
				"LS0tLS...0tCg=="
			]
		},
		"Org1MSP": {
			"name": "Org1MSP",
			"root_certs": [
				"LS0tLS...0tCg=="
			],
			"admins": [
				"LS0tLS...LS0K"
			],
			"crypto_config": {
				"signature_hash_family": "SHA2",
				"identity_identifier_hash_function": "SHA256"
			},
			"tls_root_certs": [
				"LS0tLS...LS0K"
			],
			"fabric_node_ous": {
				"enable": true,
				"client_ou_identifier": {
					"certificate": "LS0tLS...0tCg==",
					"organizational_unit_identifier": "client"
				},
				"peer_ou_identifier": {
					"certificate": "LS0tLS...0tCg==",
					"organizational_unit_identifier": "peer"
				}
			}
		},
		"Org2MSP": {
			"name": "Org2MSP",
			"root_certs": [
				"LS0tLS...LS0K"
			],
			"admins": [
				"LS0tLS...0tCg=="
			],
			"crypto_config": {
				"signature_hash_family": "SHA2",
				"identity_identifier_hash_function": "SHA256"
			},
			"tls_root_certs": [
				"LS0tLS...LS0K"
			],
			"fabric_node_ous": {
				"enable": true,
				"client_ou_identifier": {
					"certificate": "LS0tLS...LS0K",
					"organizational_unit_identifier": "client"
				},
				"peer_ou_identifier": {
					"certificate": "LS0tLS...LS0K",
					"organizational_unit_identifier": "peer"
				}
			}
		}
	},
	"orderers": {
		"OrdererMSP": {
			"endpoint": [
				{
					"host": "orderer.example.com",
					"port": 7050
				}
			]
		}
	}
}
```

结果将显示通道内的各个 MSP 的信息和排序服务信息。

#### endorsers 子命令

显示网络中的背书节点信息，包括它们的 MSP Id、账本高度、服务地址和身份证书等。

命令格式为 `endorsers [参数列表]`，支持参数如下：

```bash
* --server=SERVER：指定命令连接的 Peer 节点地址；
* --channel=CHANNEL：指定查询某个特定通道内的节点信息；
* --chaincode=CHAINCODE：指定链码名称列表；
* --collection=CC:C1,C2...：指定链码中集合信息。
```

例如，查询可以对链码 marblesp 的 collectionMarbles 集合进行背书的节点，可以执行如下命令：

```bash
$ discover \
    --peerTLSCA tls/ca.crt \
    --userKey msp/keystore/f76cf3c92dac81103c82d5490c417ac0123c279f93213f65947d8cc69e11fbc5_sk \
    --userCert msp/signcerts/Admin\@org1.example.com-cert.pem \
    --MSP Org1MSP \
    --tlsCert tls/client.crt \
    --tlsKey tls/client.key \
    endorsers \
    --server peer0.org1.example.com:7051 \
    --channel businesschannel \
    --chaincode marblesp \
    --collection=marblesp:collectionMarbles

[
	{
		"Chaincode": "marblesp",
		"EndorsersByGroups": {
			"G0": [
				{
					"MSPID": "Org1MSP",
					"LedgerHeight": 10,
					"Endpoint": "peer1.org1.example.com:7051",
					"Identity": "-----BEGIN CERTIFICATE-----\nMIICKD...UO+g==\n-----END CERTIFICATE-----\n"
				},
				{
					"MSPID": "Org1MSP",
					"LedgerHeight": 10,
					"Endpoint": "peer0.org1.example.com:7051",
					"Identity": "-----BEGIN CERTIFICATE-----\nMIICKD...mgaA==\n-----END CERTIFICATE-----\n"
				}
			],
			"G1": [
				{
					"MSPID": "Org2MSP",
					"LedgerHeight": 10,
					"Endpoint": "peer0.org2.example.com:7051",
					"Identity": "-----BEGIN CERTIFICATE-----\nMIICKT...cGaA=\n-----END CERTIFICATE-----\n"
				},
				{
					"MSPID": "Org2MSP",
					"LedgerHeight": 10,
					"Endpoint": "peer1.org2.example.com:7051",
					"Identity": "-----BEGIN CERTIFICATE-----\nMIICKD...pVTw==\n-----END CERTIFICATE-----\n"
				}
			]
		},
		"Layouts": [
			{
				"quantities_by_group": {
					"G0": 1
				}
			},
			{
				"quantities_by_group": {
					"G1": 1
				}
			}
		]
	}
]
```

结果将按组展示符合要求的背书节点的信息。

#### saveConfig 子命令

该命令并不与 Peer 节点打交道，它将通过参数指定的变量信息保存为本地文件。这样用户在执行后续命令时候可以指定该文件，而无需再指定各个参数值。

需要通过 `--configFile=CONFIGFILE` 来指定所存放的参数信息文件路径。

例如，保存指定的参数信息到本地的 discover_config.yaml 文件，可以执行如下命令：

```bash
$ discover \
    --peerTLSCA tls/ca.crt \
    --userKey msp/keystore/f76cf3c92dac81103c82d5490c417ac0123c279f93213f65947d8cc69e11fbc5_sk \
    --userCert msp/signcerts/Admin\@org1.example.com-cert.pem \
    --MSP Org1MSP \
    --tlsCert tls/client.crt \
    --tlsKey tls/client.key \
    --configFile discover_config.yaml \
    saveConfig
```

命令执行完成后，查看本地的 `discover_config.yaml` 文件内容如下：

```yaml
version: 0
tlsconfig:
  certpath: /etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt
  keypath: /etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key
  peercacertpath: /etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/ca.crt
  timeout: 0s
signerconfig:
  mspid: Org1MSP
  identitypath: /etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem
  keypath: /etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/f76cf3c92dac81103c82d5490c417ac0123c279f93213f65947d8cc69e11fbc5_sk
```

有了这个参数文件，当再使用同样的参数时就无需手动指定，直接使用 `--configFile discover_config.yaml` 即可。

当然，用户也可以手动编写参数文件，但直接使用 saveConfig 命令自动生成将更加方便、高效。
