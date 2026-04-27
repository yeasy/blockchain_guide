## 监听网络事件

Fabric 的交易提交是异步完成的。客户端提交交易后，需要通过提交状态或事件来确认交易是否已经被 Peer 验证并写入账本。

### 推荐方式：Fabric Gateway 事件

Fabric 2.4 之后，Peer 默认提供 Fabric Gateway 服务。Fabric 2.5 LTS 和 Fabric 3.x 的应用程序应优先使用 Gateway 客户端 API（Go、Node.js、Java）接收事件，而不是使用早期 SDK 的 EventHub 或示例工具 `eventsclient`。

Gateway API 可以处理以下常见场景：

* 提交交易后等待 commit status，获取交易验证码；
* 订阅链码事件，响应链码在成功提交交易中发出的事件；
* 订阅区块事件或过滤区块事件，处理账本新增区块。

Gateway 层内部会使用 Peer 的通道事件服务，但对应用隐藏了底层 Deliver 协议细节。应用侧只需要连接目标 Peer 的 Gateway gRPC 地址（通常与 Peer 服务地址相同，例如 `peer0.org1.example.com:7051`）。

### 底层 Peer 事件服务

如果需要直接接入底层事件流，可以使用 Peer 的通道事件服务。该服务按通道授权和投递事件，默认由通道 Readers 策略控制访问权限。

Fabric 2.5 LTS 和 Fabric 3.x 主要提供三类 Deliver 服务：

* `Deliver`：返回完整区块，区块中包含交易和链码事件信息；
* `DeliverWithPrivateData`：返回完整区块，并附带客户端组织有权限访问的私有数据；
* `DeliverFiltered`：返回过滤后的区块，只包含区块号、交易 ID、交易类型、验证码和过滤后的链码事件；过滤区块不会包含链码事件 payload。

客户端注册事件时，会向 Peer 发送包含 `SeekInfo` 的签名信封，指定起始和结束区块位置。需要持续监听时，结束位置通常设置为最大区块号。若 Peer 启用了双向 TLS，事件请求的通道头还需要包含 TLS 证书哈希。

### 旧 EventHub 和 eventsclient

Fabric 早期版本曾提供 Peer 级别的 EventHub，并常用 `eventsclient` 示例工具连接 `7053` 端口监听事件。这是旧事件模型：事件按 Peer 而非通道注册，访问控制和断点恢复能力都有限。

Fabric 1.1 起事件模型改为通道级 Deliver 服务；Fabric 2.5 LTS 和 Fabric 3.x 文档中不应再把 `eventsclient`、`eventHub` 或 `localhost:7053` 作为推荐做法。历史网络或旧 SDK 维护时，可以把这些内容作为迁移背景；新应用应使用 Gateway 事件 API，或直接实现 Deliver/DeliverFiltered 客户端。
