### 学历认证
#### 功能描述
该 [智能合约](chaincode_example04.go) 实现了一个简单的征信管理的案例。针对于学历认证领域，由于条约公开，在条约外无法随意篡改的特性，天然具备稳定性和中立性。

该智能合约中三种角色如下：
- 学校
- 个人
- 需要学历认证的机构或公司

学校可以根据相关信息在区块链上为某位个人授予学历，相关机构可以查询某人的学历信息，由于使用私钥签名，确保了信息的真实有效。
为了简单，尽量简化相关的业务，另未完成学业的学生因违纪或外出创业退学，学校可以修改其相应的学历信息。

账户私钥应该由安装在本地的客户端生成，本例中为了简便，使用模拟私钥和公钥。

#### 数据结构设计
- 学校
    - 名称
    - 所在位置
    - 账号地址
    - 账号公钥
    - 账户私钥
    - 学校学生
- 个人
    - 姓名
    - 账号地址
    - 过往学历
- 学历信息
    - 学历信息编号
    - 就读学校
    - 就读年份
    - 完成就读年份
    - 就读状态 //0：毕业 1：退学
- 修改记录（入学也相当于一种修改记录）
    - 编号
    - 学校账户地址（一般根据账户地址可以算出公钥地址，然后可以进行校验）
    - 学校签名
    - 个人账户地址
    - 个人公钥地址（个人不需要公钥地址）
    - 修改时间 
    - 修改操作//0:正常毕业  1：退学 2:入学

对学历操作信息所有的操作都归为记录。    
#### function及各自实现的功能
- `init`  初始化函数，并创建一所学校
- `invoke`   调用合约内部的函数
- `query`   查询相关的信息
 
- `updateDiploma` 由学校更新学生学历信息，并签名（返回记录信息） invoke
- `enrollStudent` 学校招生（返回学校信息） invoke
- `createSchool` 添加一名新学校    init
- `createStudent`  添加一名新学生  init
- `getStudentByAddress` 通过学生的地址访问学生的学历信息  query 
- `getRecordById` 通过Id获取记录  query
- `getRecords` 获取全部记录（如果记录数大于10,返回前10个） query
- `getSchoolByAddress` 通过地址获取学校的信息
- `getBackgroundById` 通过地点获取所存储的学历信息

- `writeRecord` 写入记录
- `writeSchool` 写入新创建的学校
- `writeStudent` 写入新创建的学生

#### 接口设计
 `createSchool`

request参数:
```
args[0] 学校名称
args[1] 学校所在位置
```
response参数:
```
学校信息的json表示，当创建一所新学校时，该学校学生账户地址列表为空
```

`createStudent`

request参数：
```
args[0] 学生的姓名
```

response参数：
```
学生信息的json表示，刚创建过往学历信息列表为空
```

`updateDiploma` 

request参数
```
args[0] 学校账户地址
args[1] 学校签名
args[2] 待修改学生的账户地址
args[3] //对该学生的学历进行怎样的修改，0：正常毕业  1：退学  
```

response参数
```
返回修改记录的json表示
```

`enrollStudent`

request参数:
```
args[0] 学校账户地址
args[1] 学校签名
args[2] 学生账户地址
```

response参数
```
返回修改记录的json表示
```

`getStudentByAddress`

request参数
```
args[0] address
```
response参数
```
学生信息的json表示
```

`getRecordById`

request参数
```
args[0] 修改记录的ID
```
response参数
```
修改记录的json表示
```

`getRecords`

response参数
```
获取修改记录数组（如果个数大于10，返回前10个）
```
`getSchoolByAddress`

request参数
```
args[0] address
```
response参数
```
学校信息的json表示
```

`getBackgroundById`

request参数
```
args[0] ID
```

response参数
```
学历信息的json表示
```

#### 测试
