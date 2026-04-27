## 学历认证

[chaincode_example04.go](chaincode_example04.go) 演示学校、学生和学历变更记录的链上管理。示例为了便于阅读，使用派生字符串模拟地址、公钥、私钥和签名。

### 数据结构

* `School`：学校名称、位置、地址、公私钥和学生地址列表；
* `Student`：学生姓名、地址和学历记录 ID 列表；
* `Background`：学历记录 ID、离校时间和状态；
* `Record`：学校对学生学历状态的变更记录。

### 主要交易函数

* `CreateSchool`：创建学校；
* `CreateStudent`：创建学生；
* `EnrollStudent`：学校登记学生入学；
* `UpdateDiploma`：学校更新学生学历状态；
* `GetStudentByAddress`、`GetSchoolByAddress`、`GetRecordByID`、`GetBackgroundByID`、`GetRecords`：读取状态。

真实系统中，签名与验签应由客户端身份、证书和密码学库完成，本例只保留业务流程骨架。
