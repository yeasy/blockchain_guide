## 社区能源共享

[chaincode_example05.go](chaincode_example05.go) 演示社区能源微网中的家庭账户和电力交易。

### 数据结构

* `Home`：家庭地址、可售电量、余额、状态和模拟密钥；
* `EnergyTransaction`：买方、卖方、交易电量、金额和时间。

### 主要交易函数

* `CreateUser`：创建家庭账户；
* `BuyByAddress`：买方向卖方购买指定电量；
* `ChangeStatus`：账户所有者修改是否可交易；
* `GetHomeByAddress`、`GetHomes`、`GetTransactionByID`、`GetTransactions`：读取状态。

交易逻辑会检查买方签名、卖方交易状态、卖方电量和买方余额。交易时间使用 Fabric 交易时间戳。
