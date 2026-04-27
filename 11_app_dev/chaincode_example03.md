## 数字货币发行与管理

[chaincode_example03.go](chaincode_example03.go) 演示一个简化的数字货币业务：中央银行发行货币，商业银行接收发行额度，企业接收和转移余额。

### 数据结构

* `CenterBank`：中央银行名称、发行总量、可用余额；
* `Bank`：商业银行 ID、名称、收到总量、可用余额；
* `Company`：企业 ID、名称、余额；
* `Transaction`：转出方、转入方、金额和交易时间。

### 主要交易函数

* `InitLedger`：初始化中央银行及初始发行量；
* `CreateBank`：注册商业银行；
* `CreateCompany`：注册企业；
* `IssueCoin`：中央银行增发；
* `IssueCoinToBank`：中央银行向商业银行划拨；
* `IssueCoinToCompany`：商业银行向企业划拨；
* `Transfer`：企业之间转账；
* `GetCenterBank`、`GetBankByID`、`GetCompanyByID`、`GetTransactionByID`、`GetBanks`、`GetCompanies`、`GetTransactions`：读取状态。

示例使用账本状态维护递增 ID，并使用 Fabric 交易时间戳记录交易时间，避免在链码中直接调用本地系统时间。
