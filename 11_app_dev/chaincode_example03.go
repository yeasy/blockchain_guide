package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const (
	centerBankKey      = "centerBank"
	bankCounterKey     = "counter:bank"
	companyCounterKey  = "counter:company"
	transferCounterKey = "counter:transaction"
)

// SmartContract models a small central-bank digital currency ledger.
type SmartContract struct {
	contractapi.Contract
}

type CenterBank struct {
	Name        string `json:"name"`
	TotalNumber int    `json:"totalNumber"`
	RestNumber  int    `json:"restNumber"`
}

type Bank struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TotalNumber int    `json:"totalNumber"`
	RestNumber  int    `json:"restNumber"`
}

type Company struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type Transaction struct {
	ID       int   `json:"id"`
	FromType int   `json:"fromType"`
	FromID   int   `json:"fromId"`
	ToType   int   `json:"toType"`
	ToID     int   `json:"toId"`
	Time     int64 `json:"time"`
	Number   int   `json:"number"`
}

// InitLedger creates the central bank and its initial issued balance.
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, name string, total string) (*CenterBank, error) {
	amount, err := parsePositiveAmount(total)
	if err != nil {
		return nil, err
	}

	centerBank := &CenterBank{Name: name, TotalNumber: amount, RestNumber: amount}
	if err := putJSON(ctx, centerBankKey, centerBank); err != nil {
		return nil, err
	}
	return centerBank, nil
}

// CreateBank registers a commercial bank.
func (s *SmartContract) CreateBank(ctx contractapi.TransactionContextInterface, name string) (*Bank, error) {
	id, err := nextID(ctx, bankCounterKey)
	if err != nil {
		return nil, err
	}
	bank := &Bank{ID: id, Name: name}
	if err := putJSON(ctx, bankKey(id), bank); err != nil {
		return nil, err
	}
	return bank, nil
}

// CreateCompany registers a company.
func (s *SmartContract) CreateCompany(ctx contractapi.TransactionContextInterface, name string) (*Company, error) {
	id, err := nextID(ctx, companyCounterKey)
	if err != nil {
		return nil, err
	}
	company := &Company{ID: id, Name: name}
	if err := putJSON(ctx, companyKey(id), company); err != nil {
		return nil, err
	}
	return company, nil
}

// IssueCoin increases the central-bank supply.
func (s *SmartContract) IssueCoin(ctx contractapi.TransactionContextInterface, amount string) (*Transaction, error) {
	value, err := parsePositiveAmount(amount)
	if err != nil {
		return nil, err
	}

	centerBank, err := readCenterBank(ctx)
	if err != nil {
		return nil, err
	}
	centerBank.TotalNumber += value
	centerBank.RestNumber += value
	if err := putJSON(ctx, centerBankKey, centerBank); err != nil {
		return nil, err
	}
	return recordTransaction(ctx, 0, 0, 0, 0, value)
}

// IssueCoinToBank transfers issued currency from the central bank to a commercial bank.
func (s *SmartContract) IssueCoinToBank(ctx contractapi.TransactionContextInterface, bankID string, amount string) (*Transaction, error) {
	id, value, err := parseIDAndAmount(bankID, amount)
	if err != nil {
		return nil, err
	}

	centerBank, err := readCenterBank(ctx)
	if err != nil {
		return nil, err
	}
	if centerBank.RestNumber < value {
		return nil, fmt.Errorf("central bank balance is insufficient")
	}
	bank, err := readBank(ctx, id)
	if err != nil {
		return nil, err
	}

	centerBank.RestNumber -= value
	bank.TotalNumber += value
	bank.RestNumber += value

	if err := putJSON(ctx, centerBankKey, centerBank); err != nil {
		return nil, err
	}
	if err := putJSON(ctx, bankKey(id), bank); err != nil {
		return nil, err
	}
	return recordTransaction(ctx, 0, 0, 1, id, value)
}

// IssueCoinToCompany transfers currency from a bank to a company.
func (s *SmartContract) IssueCoinToCompany(ctx contractapi.TransactionContextInterface, bankID string, companyID string, amount string) (*Transaction, error) {
	fromID, err := strconv.Atoi(bankID)
	if err != nil {
		return nil, fmt.Errorf("bank id must be an integer: %w", err)
	}
	toID, value, err := parseIDAndAmount(companyID, amount)
	if err != nil {
		return nil, err
	}

	bank, err := readBank(ctx, fromID)
	if err != nil {
		return nil, err
	}
	if bank.RestNumber < value {
		return nil, fmt.Errorf("bank %d balance is insufficient", fromID)
	}
	company, err := readCompany(ctx, toID)
	if err != nil {
		return nil, err
	}

	bank.RestNumber -= value
	company.Number += value

	if err := putJSON(ctx, bankKey(fromID), bank); err != nil {
		return nil, err
	}
	if err := putJSON(ctx, companyKey(toID), company); err != nil {
		return nil, err
	}
	return recordTransaction(ctx, 1, fromID, 2, toID, value)
}

// Transfer moves currency between companies.
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, fromCompanyID string, toCompanyID string, amount string) (*Transaction, error) {
	fromID, err := strconv.Atoi(fromCompanyID)
	if err != nil {
		return nil, fmt.Errorf("from company id must be an integer: %w", err)
	}
	toID, value, err := parseIDAndAmount(toCompanyID, amount)
	if err != nil {
		return nil, err
	}

	from, err := readCompany(ctx, fromID)
	if err != nil {
		return nil, err
	}
	if from.Number < value {
		return nil, fmt.Errorf("company %d balance is insufficient", fromID)
	}
	to, err := readCompany(ctx, toID)
	if err != nil {
		return nil, err
	}

	from.Number -= value
	to.Number += value

	if err := putJSON(ctx, companyKey(fromID), from); err != nil {
		return nil, err
	}
	if err := putJSON(ctx, companyKey(toID), to); err != nil {
		return nil, err
	}
	return recordTransaction(ctx, 2, fromID, 2, toID, value)
}

func (s *SmartContract) GetCenterBank(ctx contractapi.TransactionContextInterface) (*CenterBank, error) {
	return readCenterBank(ctx)
}

func (s *SmartContract) GetBankByID(ctx contractapi.TransactionContextInterface, id string) (*Bank, error) {
	bankID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return readBank(ctx, bankID)
}

func (s *SmartContract) GetCompanyByID(ctx contractapi.TransactionContextInterface, id string) (*Company, error) {
	companyID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return readCompany(ctx, companyID)
}

func (s *SmartContract) GetTransactionByID(ctx contractapi.TransactionContextInterface, id string) (*Transaction, error) {
	transactionID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return readTransaction(ctx, transactionID)
}

func (s *SmartContract) GetBanks(ctx contractapi.TransactionContextInterface) ([]Bank, error) {
	count, err := currentID(ctx, bankCounterKey)
	if err != nil {
		return nil, err
	}
	banks := make([]Bank, 0, count)
	for id := 0; id < count; id++ {
		bank, err := readBank(ctx, id)
		if err != nil {
			return nil, err
		}
		banks = append(banks, *bank)
	}
	return banks, nil
}

func (s *SmartContract) GetCompanies(ctx contractapi.TransactionContextInterface) ([]Company, error) {
	count, err := currentID(ctx, companyCounterKey)
	if err != nil {
		return nil, err
	}
	companies := make([]Company, 0, count)
	for id := 0; id < count; id++ {
		company, err := readCompany(ctx, id)
		if err != nil {
			return nil, err
		}
		companies = append(companies, *company)
	}
	return companies, nil
}

// GetCompanys keeps the historical example name available.
func (s *SmartContract) GetCompanys(ctx contractapi.TransactionContextInterface) ([]Company, error) {
	return s.GetCompanies(ctx)
}

func (s *SmartContract) GetTransactions(ctx contractapi.TransactionContextInterface) ([]Transaction, error) {
	count, err := currentID(ctx, transferCounterKey)
	if err != nil {
		return nil, err
	}
	transactions := make([]Transaction, 0, count)
	for id := 0; id < count; id++ {
		transaction, err := readTransaction(ctx, id)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}

func parseIDAndAmount(id string, amount string) (int, int, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return 0, 0, fmt.Errorf("id must be an integer: %w", err)
	}
	value, err := parsePositiveAmount(amount)
	if err != nil {
		return 0, 0, err
	}
	return parsedID, value, nil
}

func parsePositiveAmount(value string) (int, error) {
	amount, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("amount must be an integer: %w", err)
	}
	if amount <= 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}
	return amount, nil
}

func recordTransaction(ctx contractapi.TransactionContextInterface, fromType int, fromID int, toType int, toID int, amount int) (*Transaction, error) {
	id, err := nextID(ctx, transferCounterKey)
	if err != nil {
		return nil, err
	}
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction timestamp: %w", err)
	}

	transaction := &Transaction{ID: id, FromType: fromType, FromID: fromID, ToType: toType, ToID: toID, Time: timestamp.Seconds, Number: amount}
	if err := putJSON(ctx, transactionKey(id), transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func readCenterBank(ctx contractapi.TransactionContextInterface) (*CenterBank, error) {
	var centerBank CenterBank
	if err := readJSON(ctx, centerBankKey, &centerBank); err != nil {
		return nil, err
	}
	return &centerBank, nil
}

func readBank(ctx contractapi.TransactionContextInterface, id int) (*Bank, error) {
	var bank Bank
	if err := readJSON(ctx, bankKey(id), &bank); err != nil {
		return nil, err
	}
	return &bank, nil
}

func readCompany(ctx contractapi.TransactionContextInterface, id int) (*Company, error) {
	var company Company
	if err := readJSON(ctx, companyKey(id), &company); err != nil {
		return nil, err
	}
	return &company, nil
}

func readTransaction(ctx contractapi.TransactionContextInterface, id int) (*Transaction, error) {
	var transaction Transaction
	if err := readJSON(ctx, transactionKey(id), &transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}

func nextID(ctx contractapi.TransactionContextInterface, counterKey string) (int, error) {
	id, err := currentID(ctx, counterKey)
	if err != nil {
		return 0, err
	}
	if err := ctx.GetStub().PutState(counterKey, []byte(strconv.Itoa(id+1))); err != nil {
		return 0, fmt.Errorf("failed to update counter %s: %w", counterKey, err)
	}
	return id, nil
}

func currentID(ctx contractapi.TransactionContextInterface, counterKey string) (int, error) {
	data, err := ctx.GetStub().GetState(counterKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read counter %s: %w", counterKey, err)
	}
	if data == nil {
		return 0, nil
	}
	return strconv.Atoi(string(data))
}

func readJSON(ctx contractapi.TransactionContextInterface, key string, target interface{}) error {
	data, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", key, err)
	}
	if data == nil {
		return fmt.Errorf("state %s does not exist", key)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to decode %s: %w", key, err)
	}
	return nil
}

func putJSON(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to encode %s: %w", key, err)
	}
	if err := ctx.GetStub().PutState(key, data); err != nil {
		return fmt.Errorf("failed to write %s: %w", key, err)
	}
	return nil
}

func bankKey(id int) string {
	return fmt.Sprintf("bank:%d", id)
}

func companyKey(id int) string {
	return fmt.Sprintf("company:%d", id)
}

func transactionKey(id int) string {
	return fmt.Sprintf("transaction:%d", id)
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
