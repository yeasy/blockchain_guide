# 管理链上代码 (Lifecycle)

从 Fabric 2.0 开始，引入了全新的链码生命周期管理机制（Chaincode Lifecycle），旨在支持更去中心化的治理模式。与 1.x 版本中“一旦实例化全网生效”的模式不同，新生命周期要求组织间达成共识才能启动链码。

本章将详细介绍如何使用 `peer lifecycle chaincode` 命令族来管理 Fabric 2.x/3.0 网络的链码。

## 链码生命周期概览

在 Fabric 2.x/3.0 中，部署一个链码需要经历以下四个核心步骤：

1.  **打包 (Package)**: 将链码源码和元数据打包成 `.tar.gz` 文件。
2.  **安装 (Install)**: 将链码包安装到 **每个** 需要运行该链码的 Peer 节点上。
3.  **批准 (Approve)**: 每个组织根据自己的意愿，为特定的链码定义（版本、背书策略等）投票（批准）。
4.  **提交 (Commit)**: 当获得足够多组织（满足生命周期背书策略，默认是大多数）的批准后，链码定义被提交到通道，正式生效。

之后，用户就可以调用或查询链码了。

![Fabric 2.0 链码生命周期](_images/chaincode_lifecycle_v2.png)

## 1. 打包链码 (Package)

打包操作在本地进行，不需要连接网络。

```bash
# 设置环境变量
export CC_NAME=basic
export CC_VERSION=1.0
export CC_SRC_PATH=../asset-transfer-basic/chaincode-go

# 打包
peer lifecycle chaincode package ${CC_NAME}.tar.gz \
    --path ${CC_SRC_PATH} \
    --lang golang \
    --label ${CC_NAME}_${CC_VERSION}
```

这会生成一个 `basic.tar.gz` 文件。

## 2. 安装链码 (Install)

将打包好的文件安装到 Peer 节点上。此操作是针对**节点**的，需要在所有背书节点上执行。

```bash
peer lifecycle chaincode install ${CC_NAME}.tar.gz
```

安装成功后，系统会返回一个**包标识符 (Package ID)**，格式为 `label:hash`。你需要记录下这个 ID，后续步骤会用到。

```bash
# 查询已安装的链码包ID
peer lifecycle chaincode queryinstalled
# 输出示例: basic_1.0:e23a...
```

## 3. 组织批准 (Approve)

这是 2.x 最关键的变更。每个组织都需要使用自己的 MSP 身份批准链码定义。链码定义包括：名称、版本、序列号、背书策略等。

**注意**：所有组织必须批准**完全相同**的参数（包括 Package ID），才能达成共识。

```bash
# 环境变量
export PACKAGE_ID=basic_1.0:e23a...
export CHANNEL_NAME=mychannel

# 批准链码定义
peer lifecycle chaincode approveformyorg \
    --channelID ${CHANNEL_NAME} \
    --name ${CC_NAME} \
    --version ${CC_VERSION} \
    --package-id ${PACKAGE_ID} \
    --sequence 1 \
    --tls --cafile ${ORDERER_CA}
```

*   `--sequence`: 序列号。首次部署为 1，每次升级链码时需递增（如 2, 3）。
*   `--package-id`: 指定要运行的具体代码包 ID。

你可以随时检查当前通道的批准状态：

```bash
peer lifecycle chaincode checkcommitreadiness \
    --channelID ${CHANNEL_NAME} \
    --name ${CC_NAME} \
    --version ${CC_VERSION} \
    --sequence 1 \
    --output json
```

## 4. 提交链码 (Commit)

当 `checkcommitreadiness` 显示已有足够多的组织（默认是大多数，Majority）批准了该定义，任意一个组织的管理员都可以执行提交操作。

```bash
peer lifecycle chaincode commit \
    --channelID ${CHANNEL_NAME} \
    --name ${CC_NAME} \
    --version ${CC_VERSION} \
    --sequence 1 \
    --tls --cafile ${ORDERER_CA} \
    --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles ${ORG1_CA} \
    --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles ${ORG2_CA}
```

**注意**：`commit` 交易需要收集足够多组织的背书，因此通常需要指定多个 `--peerAddresses` 来收集背书签名。

## 5. 调用与查询

链码提交成功后，即可正常使用。如果链码包含 `Init` 方法且需要初始化，可以使用 `--isInit` 标志调用一次（需要在 Approve 和 Commit 时指定 `--init-required`）。大多数现代链码不需要显式 Init。

**调用 (Invoke):**

```bash
peer chaincode invoke \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.example.com \
    --tls --cafile ${ORDERER_CA} \
    -C ${CHANNEL_NAME} \
    -n ${CC_NAME} \
    --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles ${ORG1_CA} \
    --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles ${ORG2_CA} \
    -c '{"function":"InitLedger","Args":[]}'
```

**查询 (Query):**

```bash
peer chaincode query -C ${CHANNEL_NAME} -n ${CC_NAME} -c '{"Args":["GetAllAssets"]}'
```

## 升级链码

在 Fabric 2.x 中，升级链码本质上是更新链码定义。流程如下：

1.  **Package & Install**: 打包并安装新版本的代码（如 v1.1）。这会生成一个新的 Package ID。
2.  **Approve**: 组织批准新的定义。注意：**增加序列号**（如 `--sequence 2`），更新版本号，并指向新的 `package-id`。
3.  **Commit**: 提交新的定义。

不需要像 1.x 那样执行专门的 `upgrade` 命令，只要 Commit 成功，新代码即刻生效。

## 总结

新的生命周期管理（Lifecycle）虽然步骤看起来繁琐（Install -> Approve -> Commit），但它带来了真正的**去中心化治理**能力。组织之间可以协商链码的升级策略，而不再是被动接受某个管理员的单方面操作。这对于企业级联盟链来说至关重要。
