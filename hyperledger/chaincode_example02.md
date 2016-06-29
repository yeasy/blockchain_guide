##智能合约案例二代码分析
###`chaincode_example02.go`主要实现如下的功能：
- 初始化A、B两个账户，并为两个账户赋初始值
- 在A、B两个账户之间进行转账
- 查询A、B两个账户上的余额
- 删除账户。

###function及实现的功能
- `init`  初始化A、B两个账户
- `invoke`  实现A、B账户间的转账。
- `query`  查询A、B账户上的余额
- `delete` 删除账户。

###包依赖
```
import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)
```
`strconv`实现int与string之间的转换。

在invoke函数中，有：
```
X, err = strconv.Atoi(args[2])
	Aval = Aval - X
	Bval = Bval + X
```
当args[2]<0时，A账户余额增加，B账户余额减少。

###案例不足之处
没有新增账户并初始化的功能。开发者可以根据自己的业务模型进行添加。
