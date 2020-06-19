## 数字货币发行与管理
### 简介
该智能合约实现一个简单的商业应用案例，即数字货币的发行与转账。在这之中一共分为三种角色：中央银行，商业银行，企业。其中中央银行可以发行一定数量的货币，企业之间可以进行相互的转账。主要实现如下的功能：

* 初始化中央银行及其发行的货币数量
* 新增商业银行，同时央行并向其发行一定数量的货币
* 新增企业
* 商业银行向企业转给一定数量的数字货币
* 企业之间进行相互的转账
* 查询企业、银行、交易信息

### 主要函数
* `init`：初始化中央银行，并发行一定数量的货币；
* `invoke`：调用合约内部的函数；
* `query`：查询相关的信息；
* `createBank`：新增商业银行，同时央行向其发行一定数量的货币；
* `createCompany`：新增企业；
* `issueCoin`：央行再次发行一定数量的货币（归于交易）；
* `issueCoinToBank`：央行向商业银行转一定数量的数字货币（归于交易）；
* `issueCoinToCp`：商业银行向企业转一定数量的数字货币（归于交易行为）；
* `transfer`：企业之间进行相互转账（归于交易行为）；
* `getCompanys`：获取所有的公司信息，如果企业个数大于10，先访问前10个；
* `getBanks`：获取所有的商业银行信息，如果商业银行个数大于10，先访问前 10 个
* `getTransactions`：获取所有的交易记录 如果交易个数大于10，先访问前 10 个；
* `getCompanyById`：获取某家公司信息；
* `getBankById`：获取某家银行信息；
* `getTransactionBy`：获取某笔交易记录；
* `writeCenterBank`：修改央行信息；
* `writeBank`：修改商业银行信息；
* `writeCompany`：修改企业信息；
* `writeTransaction`：写入交易信息。

### 数据结构设计
* centerBank 中央银行
  * Name：名称
  * TotalNumber：发行货币总数额
  * RestNumber：账户余额
  * ID：ID固定为 0
* bank  商业银行
  * Name：名称
  * TotalNumber：收到货币总数额
  * RestNumber：账户余额
  * ID：银行 ID
* company 企业
  * Name：名称
  * Number：账户余额
  * ID：企业 ID
* transaction 交易内容
  * FromType：发送方角色 //centerBank:0,Bank:1,Company:2
  * FromID：发送方 ID
  * ToType：接收方角色 //Bank:1,Company:2
  * ToID：接收方 ID
  * Time：交易时间
  * Number：交易数额
  * ID：交易 ID
 
### 接口设计
#### `init`
request 参数:

```
args[0] 银行名称
args[1] 初始化发布金额
```

response 参数:

```json
{"Name":"XXX","TotalNumber":"0","RestNumber":"0","ID":"XX"}
```

#### `createBank`

request 参数:
```
args[0] 银行名称
```

response 参数:

```json
{"Name":"XXX","TotalNumber":"0","RestNumber":"0","ID":"XX"}
```

#### `createCompany`

request 参数:

```
args[0] 公司名称
```

response 参数:

```json
{"Name":"XXX","Number":"0","ID":"XX"}
```

#### `issueCoin`

request 参数:

```
args[0] 再次发行货币数额

```
response 参数:

```json
{"FromType":"0","FromID":"0","ToType":"0","ToID":"0","Time":"XX","Number":"XX","ID":"XX"}
```

#### `issueCoinToBank`

request 参数:

```
args[0] 商业银行ID
args[1] 转账数额
```

response 参数:

```json
{"FromType":"0","FromID":"0","ToType":"1","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```

#### `issueCoinToCp`

request 参数:

```
args[0] 商业银行ID
args[1] 企业ID
args[2] 转账数额

```
response 参数:

```json
{"FromType":"1","FromID":"XX","ToType":"2","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```

#### `transfer`

request 参数:
```
args[0] 转账用户ID
args[1] 被转账用户ID
args[2] 转账余额
```

response 参数:

```json
{"FromType":"2","FromID":"XX","ToType":"2","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```

#### `getBanks`

response 参数

```json
[{"Name":"XXX","Number":"XX","ID":"XX"},{"Name":"XXX","Number":"XX","ID":"XX"},...]
```

#### `getCompanys`

response 参数

```json
[{"Name":"XXX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"},{"Name":"XXX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"},...]
```

#### `getTransactions`

response 参数

```json
[{"FromType":"XX","FromID":"XX","ToType":"XX","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"},{"FromType":"XX","FromID":"XX","ToType":"XX","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"},...]
```

#### `getCenterBank`

response 参数

```json
[{"Name":"XX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"}]
```

#### `getBankById`

request 参数

```
args[0] 商业银行ID
```

response 参数

```json
[{"Name":"XX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"}]
```

#### `getCompanyById`

request 参数

```
args[0] 企业ID
```

response 参数

```json
[{"Name":"XXX","Number":"XX","ID":"XX"}]
```

#### `getTransactionById`

request 参数
```
args[0] 交易ID
```

response 参数

```json
{"FromType":"XX","FromID":"XX","ToType":"XX","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```

#### `writeCenterBank`

request 参数

```
CenterBank
```

response 参数

```
err  nil 为成功
```

#### `writeBank`

request 参数

```
Bank
```

response 参数

```
err  nil 为成功
```

#### `writeCompany`

request 参数
```
Company
```

response 参数

```
err  nil 为成功
```

#### `writeTransaction`

request 参数
```
Transaction
```

response 参数

```
err  nil 为成功
···

#### 其它
查询时为了兼顾读速率，将一些信息备份存放在非区块链数据库上也是一个较好的选择。
