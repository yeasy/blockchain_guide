### 物流供应链简单案例
#### 功能描述
该 [智能合约](chaincode_example06.go) 实现了一个简单的供应链应用案例，针对物流行业的应用场景。由于将合约的协议公开，并且签收快递时需要签名，可以在很大程度上保证不被冒领，实现了一手交钱，一手交货，同时提高了效率，确保了透明。

该智能合约中三种角色如下：
- 物流公司（本案例中只有1位）
- 寄货方（本案例中有多位）
- 收货方（本案例中有多位）

业务流程如下：

1、寄货方填写寄货单，物流公司根据寄货单寄快递。

2、寄快递过程中物流公司各个快递点对快递进行扫描，描述目前快递进度，并更新货单状态。寄货方和收货方可以根据单号进行查询。

3、快递到达后，收货方检查商品，确认无误后，扫码并使用私钥签名，支付相关费用，更新订单状态。

在实际中，物流费的支付分为两类：
- 1、寄货方支付。收货方签收快递后先预付给物流公司。
- 2、收货方支付。收货方签收快递后支付给物流公司。

在本案例中暂不考虑货物损坏、收货方失联、货物保值等的相关问题。具体实现逻辑如下：

- 创建账户。为每个用户生成唯一的私钥与地址。
- 生成寄货单。寄货方填写纸质寄货单，物流公司根据此生成电子单。
- 更新寄货单。物流公司旗下快递点根据配送信息更新电子寄货单。
- 收货方签收确认。收货方收到货物后，使用自己的私钥进行签收，完成相应的付款。

账户私钥应该由安装在本地的客户端生成，本例中为了简便，使用模拟私钥和公钥。每位用户的私钥为guid+“1”，公钥为guid+“2”。用户签名为私钥+“1”

#### 数据结构设计
- 寄货单
    - 寄货单编号
    - 寄货方地址
    - 收货方地址
    - 寄货方联系方式
    - 收货方联系方式
    - 物流费用
    - 物流费用支付类型  //0：寄货方支付 1：收货方支付
    - 寄货方预支付费用  //模拟实际预支付，寄货方支付物流费下值为物流费，否则为0
    - 快递配送信息    // 快递运送状态，所经过快递分拨中心与快递点的数组
    - 收货方签名 

- 寄货方
    - 姓名
    - 所在地址
    - 账户地址
    - 账户公钥
    - 联系方式
    - 账户余额
- 收货方
    - 姓名
    - 所在地址
    - 账户地址
    - 账户公钥
    - 账户私钥
    - 联系方式
    - 账户余额
- 物流公司
    - 账户公钥
    - 账户私钥
    - 名称
    - 地址
    - 联系方式
    - 账户余额
    - 物流公司旗下分拨中心与快递点
- 快递点
    - 名称
    - 所在地址
    - 联系方式
    - 快递点公钥
    - 快递点私钥
    - 快递点账户地址

#### function及各自实现的功能
- `init`  初始化物流公司及其下相应快递点
- `invoke`   调用合约内部的函数
- `query`   查询相关的信息
- `createUser` 创建用户 init
- `createExpress` 创建物流公司 init
- `createExpressPoint` 创建快递点 init
- `createExpressOrder` 寄货方创建寄货单  init
- `finishExpressOrder` 收货方签收寄货单 invoke
- `addExpressPointer` 物流公司添加新的快递点  invoke
- `updateExpressOrder` 更新物流公司订单,添加快递点的信息 invoke  


- `getExpressOrderById` 查询订单状态  query
- `getExpress`  获取物流公司信息      query
- `getUserByAddress` 获取用户信息   query
- `getExpressPointByAddress`  获取快递点信息  query   

- `writeExpress` 存储物流公司信息 （以物流公司账户地址进行存储）
- `writeExpressOrder` 存储寄货单  （以“express”+id 进行存储）
- `writeUser` 存储用户信息   （以地址进行存储）
- `writeExpressPoint` 存储物流点信息  （以快递点账户地址进行存储）

#### 接口设计
 `createUser`

request参数
```
args[0] 姓名 
args[1] 所在地址
args[2] 联系方式
args[3] 账户余额
```

response参数
```
user信息的json表示
```

 `createExpressPointer` 

request参数
```
args[0] 姓名
args[1] 所在地址
args[2] 联系方式
```

response参数
```
物流点的信息的json表示
```

 `createExpress`

request 参数
```
args[0]  名称
args[1]  地址
args[2]  联系方式
args[3]  账户余额
```
response 参数
```
物流公司信息的json表示
```

 `addExpressPointer`

request参数
```
args[0] 添加快递点
```

response参数
```
物流公司信息的json表示
```

 `createExpressOrder`

request参数
```
args[0] 寄货方地址
args[1] 收货方地址
args[2] 寄货方账户地址
args[3] 收货方账户地址
args[4] 寄货方联系方式
args[5] 收货方联系方式
args[6] 物流费用支付类型
args[7] 寄货方预支付费用 （收货方支付的话值为0）
args[8] 物流费用
```

response 参数
```
订单信息的json表示
```

 `updateExpressOrder`

request参数
```
args[0]  订单id
args[1]  快递点地址
```

response参数
```
订单信息的json表示
```

  `finishExpressOrder`

request参数
```
args[0] 收货方账户地址
args[1] 账户订单编号
args[2] 收货方签名
```

response参数
```
订单信息的json表示
```

`getExpressOrderById`

request参数：
```
args[0] id
```

response参数：
```
快递订单的json表示
```

`getExpress`

response参数：
```
快递信息的json表示
```

`getUserByAddress`

request参数
```
args[0] address
```

response参数
```
用户信息的json表示
```

`getExpressPointerByAddress`

request参数
```
args[0] address
```

response参数
```
快递点的json信息表示
```

#### 测试
