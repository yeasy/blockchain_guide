package main

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract stores and retrieves notarized key/value data.
type SmartContract struct {
	contractapi.Contract
}

// InitLedger writes the initial value under a well-known key.
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, value string) error {
	return ctx.GetStub().PutState("hello_world", []byte(value))
}

// Write creates or updates a value.
func (s *SmartContract) Write(ctx contractapi.TransactionContextInterface, key string, value string) error {
	if key == "" {
		return fmt.Errorf("key must not be empty")
	}
	return ctx.GetStub().PutState(key, []byte(value))
}

// Read returns a value by key.
func (s *SmartContract) Read(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	value, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", key, err)
	}
	if value == nil {
		return "", fmt.Errorf("key %s does not exist", key)
	}
	return string(value), nil
}

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
