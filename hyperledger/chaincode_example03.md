##数字货币简单商业应用案例
#### 功能描述
该智能合约实现一个简单的商业应用案例，即数字货币的发行与转账。在这之中一共分为三种角色：中央银行，商业银行，企业。其中中央银行可以发行一定数量的货币，企业之间可以进行相互的转账。主要实现如下的功能：
- 初始化中央银行及其发行的货币数量
- 新增商业银行，同时央行并向其发行一定数量的货币
- 新增企业
- 商业银行向企业转给一定数量的数字货币
- 企业之间进行相互的转账
- 查询企业、银行、交易信息

#### function及各自实现的功能：
- `init`  初始化中央银行，并发行一定数量的货币
- `invoke`   调用合约内部的函数
- `query`   查询相关的信息
- `createBank`   新增商业银行，同时央行向其发行一定数量的货币
- `createCompany`   新增企业
- `issueCoin` 央行再次发行一定数量的货币 （归于交易）
- `issueCoinToBank` 央行向商业银行转一定数量的数字货币 （归于交易）
- `issueCoinToCp`  商业银行向企业转一定数量的数字货币  （归于交易行为）
- `transfer`   企业之间进行相互转账  （归于交易行为）
- `getCompanys`   获取所有的公司信息  如果企业个数大于10，先访问前10个
- `getBanks`    获取所有的商业银行信息  如果商业银行个数大于10，先访问前10个
- `getTransactions` 获取所有的交易记录 如果交易个数大于10，先访问前10个
- `getCompanyById`   获取某家公司信息
- `getBankById`   获取某家银行信息
- `getTransactionBy` 获取某笔交易记录
- `writeCenterBank`  修改央行信息
- `writeBank` 修改商业银行信息
- `writeCompany` 修改企业信息
- `writeTransaction` 写入交易信息


#### 数据结构设计
- centerBank 中央银行
  - Name 名称
  - TotalNumber 发行货币总数额
  - RestNumber 账户余额
  - ID ID固定为0
- bank  商业银行
  - Name 名称
  - TotalNumber 收到货币总数额
  - RestNumber 账户余额
  - ID 银行id
- company 企业
  - Name 名称
  - Number  账户余额
  - ID 企业id
- transaction 交易内容
  - FromType 发送方角色 //centerBank:0,Bank:1,Company:2
  - FromID 发送方ID
  - ToType  接收方角色 //Bank:1,Company:2
  - ToID 接收方ID
  - Time  交易时间
  - Number 交易数额
  - ID 交易ID
 
#### 接口设计
- `init`
request参数:
```
args[0] 银行名称
args[1] 初始化发布金额
```
response参数:
```
{"Name":"XXX","TotalNumber":"0","RestNumber":"0","ID":"XX"}
```

- `createBank`:

request参数:
```
args[0] 银行名称
```
response参数:
```
{"Name":"XXX","TotalNumber":"0","RestNumber":"0","ID":"XX"}
```

- `createCompany`

request参数:
```
args[0] 公司名称
```
response参数:
```
{"Name":"XXX","Number":"0","ID":"XX"}
```
- `issueCoin`

request参数:
```
args[0] 再次发行货币数额
```
response参数:
```
{"FromType":"0","FromID":"0","ToType":"0","ToID":"0","Time":"XX","Number":"XX","ID":"XX"}
```

- `issueCoinToBank`

request参数:
```
args[0] 商业银行ID
args[1] 转账数额
```
response参数:
```
{"FromType":"0","FromID":"0","ToType":"1","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```

- `issueCoinToCp`

request参数:
```
args[0] 商业银行ID
args[1] 企业ID
args[2] 转账数额
```
response参数:
```
{"FromType":"1","FromID":"XX","ToType":"2","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```

- `transfer`

request参数:
```
args[0] 转账用户ID
args[1] 被转账用户ID
args[2] 转账余额
```
response参数:
```
{"FromType":"2","FromID":"XX","ToType":"2","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```
- `getBanks`

response参数
```
[{"Name":"XXX","Number":"XX","ID":"XX"},{"Name":"XXX","Number":"XX","ID":"XX"},...]
```

- `getCompanys`

response参数
```
[{"Name":"XXX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"},{"Name":"XXX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"},...]
```

- `getTransactions`

response参数
```
[{"FromType":"XX","FromID":"XX","ToType":"XX","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"},{"FromType":"XX","FromID":"XX","ToType":"XX","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"},...]
```
- `getCenterBank`

response参数
```
[{"Name":"XX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"}]
```

- `getBankById`

request参数
```
args[0] 商业银行ID
```
response参数
```
[{"Name":"XX","TotalNumber":"XX","RestNumber":"XX","ID":"XX"}]
```

- `getCompanyById`

request参数
```
args[0] 企业ID
```
response参数
```
[{"Name":"XXX","Number":"XX","ID":"XX"}]
```
- `getTransactionById`

request参数
```
args[0] 交易ID
```
response参数
```
{"FromType":"XX","FromID":"XX","ToType":"XX","ToID":"XX","Time":"XX","Number":"XX","ID":"XX"}
```

- `writeCenterBank`

request参数
```
CenterBank
```
response参数
```
err  nil 为成功

- `writeBank`

request参数
```
Bank
```
response参数
```
err  nil 为成功

- `writeCompany`

request参数
```
Company
```
response参数
```
err  nil 为成功

- `writeTransaction`

request参数
```
Transaction
```
response参数
```
err  nil 为成功


#### 其它
查询时为了兼顾读速率，将一些信息备份存放在非区块链数据库上也是一个较好的选择。
