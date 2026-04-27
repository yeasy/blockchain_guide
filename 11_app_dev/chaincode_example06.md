## 物流供应链简单案例

[chaincode_example06.go](chaincode_example06.go) 演示物流公司、快递点、寄件人、收件人和寄货单的基本流程。

### 数据结构

* `User`：用户姓名、位置、地址、联系方式、余额和模拟密钥；
* `Express`：物流公司名称、地址、联系方式、余额和快递点列表；
* `ExpressPoint`：快递点名称、地址、联系方式和所属物流公司；
* `ExpressOrder`：寄货单、寄收双方、物流费用、途经快递点和签收状态。

### 主要交易函数

* `CreateUser`：创建寄件人或收件人；
* `CreateExpress`：创建物流公司；
* `CreateExpressPoint`：创建快递点；
* `AddExpressPoint`：将快递点加入物流公司；
* `CreateExpressOrder`：创建寄货单；
* `UpdateExpressOrder`：追加订单途经快递点；
* `FinishExpressOrder`：收件人签收并完成付款；
* `GetExpressOrderByID`、`GetExpress`、`GetUserByAddress`、`GetExpressPointByAddress`：读取状态。

示例保留两种付款方式：寄件人预付和收件人到付。真实系统中，签收验签与资金结算应接入正式身份、支付和风控流程。
