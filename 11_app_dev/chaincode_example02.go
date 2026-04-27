package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract transfers integer balances between named accounts.
type SmartContract struct {
	contractapi.Contract
}

type Account struct {
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

// InitAccounts creates two accounts with their opening balances.
func (s *SmartContract) InitAccounts(ctx contractapi.TransactionContextInterface, a string, aBalance string, b string, bBalance string) error {
	aValue, err := parseAmount(aBalance)
	if err != nil {
		return fmt.Errorf("invalid balance for %s: %w", a, err)
	}
	bValue, err := parseAmount(bBalance)
	if err != nil {
		return fmt.Errorf("invalid balance for %s: %w", b, err)
	}
	if err := putBalance(ctx, a, aValue); err != nil {
		return err
	}
	return putBalance(ctx, b, bValue)
}

// Transfer moves amount from one account to another.
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, from string, to string, amount string) (*Account, error) {
	value, err := parseAmount(amount)
	if err != nil {
		return nil, err
	}
	if value <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	fromBalance, err := getBalance(ctx, from)
	if err != nil {
		return nil, err
	}
	toBalance, err := getBalance(ctx, to)
	if err != nil {
		return nil, err
	}
	if fromBalance < value {
		return nil, fmt.Errorf("account %s has insufficient balance", from)
	}

	fromBalance -= value
	toBalance += value

	if err := putBalance(ctx, from, fromBalance); err != nil {
		return nil, err
	}
	if err := putBalance(ctx, to, toBalance); err != nil {
		return nil, err
	}

	return &Account{Name: from, Balance: fromBalance}, nil
}

// ReadAccount returns the current account balance.
func (s *SmartContract) ReadAccount(ctx contractapi.TransactionContextInterface, name string) (*Account, error) {
	balance, err := getBalance(ctx, name)
	if err != nil {
		return nil, err
	}
	return &Account{Name: name, Balance: balance}, nil
}

// Delete removes an account.
func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface, name string) error {
	return ctx.GetStub().DelState(name)
}

func parseAmount(value string) (int, error) {
	amount, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("expected integer amount: %w", err)
	}
	return amount, nil
}

func getBalance(ctx contractapi.TransactionContextInterface, name string) (int, error) {
	data, err := ctx.GetStub().GetState(name)
	if err != nil {
		return 0, fmt.Errorf("failed to read account %s: %w", name, err)
	}
	if data == nil {
		return 0, fmt.Errorf("account %s does not exist", name)
	}
	return parseAmount(string(data))
}

func putBalance(ctx contractapi.TransactionContextInterface, name string, balance int) error {
	if name == "" {
		return fmt.Errorf("account name must not be empty")
	}
	return ctx.GetStub().PutState(name, []byte(strconv.Itoa(balance)))
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
