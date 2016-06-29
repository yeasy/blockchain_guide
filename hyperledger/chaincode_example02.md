## 智能合约示例二

### 实现功能

[chaincode_example02.go](chaincode_example02.go) 主要实现如下的功能：

* 初始化 A、B 两个账户，并为两个账户赋初始值
* 在 A、B 两个账户之间进行转账
* 查询A、B两个账户上的余额
* 删除账户。

### 主要函数

* `init`：初始化 A、B 两个账户；
* `invoke`：实现 A、B 账户间的转账；
* `query`：查询 A、B 账户上的余额；
* `delete`：删除账户。

### 依赖的包
```golang
import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)
```
`strconv` 实现 int 与 string 类型之间的转换。

在invoke 函数中，存在：
```golang
X, err = strconv.Atoi(args[2])
	Aval = Aval - X
	Bval = Bval + X
```

当 `args[2]<0` 时，A 账户余额增加，否则 B 账户余额减少。

### 可扩展功能
实例中未包含新增账户并初始化的功能。开发者可以根据自己的业务模型进行添加。
