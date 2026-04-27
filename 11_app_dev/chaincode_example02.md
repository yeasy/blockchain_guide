## 链码示例二：交易资产

[chaincode_example02.go](chaincode_example02.go) 使用 Go Contract API 实现两个账户之间的整数余额转账。

### 主要交易函数

* `InitAccounts`：初始化两个账户及其余额；
* `Transfer`：从一个账户向另一个账户转账；
* `ReadAccount`：读取账户余额；
* `Delete`：删除账户。

示例依赖当前 Go Contract API：

```go
import "github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
```

转账逻辑会检查金额必须为正数、账户必须存在、付款账户余额必须充足。客户端应用应通过 Gateway API 的 submit 类调用提交 `Transfer`，通过 evaluate 类调用读取 `ReadAccount`。
