## 链码示例一：信息公证

[chaincode_example01.go](chaincode_example01.go) 使用 Go Contract API 实现一个最小键值账本：

* `InitLedger`：初始化 `hello_world` 的值；
* `Write`：创建或更新指定 key 的 value；
* `Read`：读取指定 key 的 value。

代码入口使用 `contractapi.NewChaincode` 创建链码实例，并调用 `Start` 由 Fabric peer 托管执行。

```go
func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err)
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}
```

客户端应用在 Fabric v2.4+ 中应通过 Gateway API 调用合约。读取状态时使用 evaluate 类调用 `Read`；修改状态时使用 submit 类调用 `Write`。
