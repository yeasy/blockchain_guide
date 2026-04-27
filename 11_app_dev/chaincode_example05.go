package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

const (
	homeCounterKey        = "counter:home"
	energyTxCounterKey    = "counter:energyTransaction"
	homeAvailableForTrade = 1
)

// SmartContract manages a small community energy market.
type SmartContract struct {
	contractapi.Contract
}

type Home struct {
	Address string `json:"address"`
	Energy  int    `json:"energy"`
	Money   int    `json:"money"`
	ID      int    `json:"id"`
	Status  int    `json:"status"`
	PriKey  string `json:"priKey"`
	PubKey  string `json:"pubKey"`
}

type EnergyTransaction struct {
	BuyerAddress     string `json:"buyerAddress"`
	BuyerAddressSign string `json:"buyerAddressSign"`
	SellerAddress    string `json:"sellerAddress"`
	Energy           int    `json:"energy"`
	Money            int    `json:"money"`
	ID               int    `json:"id"`
	Time             int64  `json:"time"`
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, energy string, money string) (*Home, error) {
	energyValue, err := parseNonNegativeAmount(energy)
	if err != nil {
		return nil, err
	}
	moneyValue, err := parseNonNegativeAmount(money)
	if err != nil {
		return nil, err
	}
	id, err := nextID(ctx, homeCounterKey)
	if err != nil {
		return nil, err
	}
	address := newAddress(ctx, "home")
	home := &Home{
		Address: address,
		Energy:  energyValue,
		Money:   moneyValue,
		ID:      id,
		Status:  homeAvailableForTrade,
		PriKey:  address + "1",
		PubKey:  address + "2",
	}
	if err := putHome(ctx, home); err != nil {
		return nil, err
	}
	return home, nil
}

func (s *SmartContract) BuyByAddress(ctx contractapi.TransactionContextInterface, sellerAddress string, buyerSignature string, buyerAddress string, energy string) (*EnergyTransaction, error) {
	energyValue, err := parsePositiveAmount(energy)
	if err != nil {
		return nil, err
	}
	if !validSignature(buyerAddress, buyerSignature) {
		return nil, fmt.Errorf("invalid buyer signature")
	}

	seller, err := readHome(ctx, sellerAddress)
	if err != nil {
		return nil, err
	}
	if seller.Status != homeAvailableForTrade {
		return nil, fmt.Errorf("seller %s is not available for trading", sellerAddress)
	}
	buyer, err := readHome(ctx, buyerAddress)
	if err != nil {
		return nil, err
	}
	if seller.Energy < energyValue {
		return nil, fmt.Errorf("seller has insufficient energy")
	}
	if buyer.Money < energyValue {
		return nil, fmt.Errorf("buyer has insufficient money")
	}

	seller.Energy -= energyValue
	seller.Money += energyValue
	buyer.Energy += energyValue
	buyer.Money -= energyValue

	if err := putHome(ctx, seller); err != nil {
		return nil, err
	}
	if err := putHome(ctx, buyer); err != nil {
		return nil, err
	}

	id, err := nextID(ctx, energyTxCounterKey)
	if err != nil {
		return nil, err
	}
	timestamp, err := txUnixTime(ctx)
	if err != nil {
		return nil, err
	}
	transaction := &EnergyTransaction{
		BuyerAddress:     buyerAddress,
		BuyerAddressSign: buyerSignature,
		SellerAddress:    sellerAddress,
		Energy:           energyValue,
		Money:            energyValue,
		ID:               id,
		Time:             timestamp,
	}
	if err := putJSON(ctx, energyTransactionKey(id), transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *SmartContract) ChangeStatus(ctx contractapi.TransactionContextInterface, address string, signature string, status string) (*Home, error) {
	if !validSignature(address, signature) {
		return nil, fmt.Errorf("invalid owner signature")
	}
	statusValue, err := parseNonNegativeAmount(status)
	if err != nil {
		return nil, err
	}
	home, err := readHome(ctx, address)
	if err != nil {
		return nil, err
	}
	home.Status = statusValue
	if err := putHome(ctx, home); err != nil {
		return nil, err
	}
	return home, nil
}

func (s *SmartContract) GetHomeByAddress(ctx contractapi.TransactionContextInterface, address string) (*Home, error) {
	return readHome(ctx, address)
}

func (s *SmartContract) GetHomes(ctx contractapi.TransactionContextInterface) ([]Home, error) {
	count, err := currentID(ctx, homeCounterKey)
	if err != nil {
		return nil, err
	}
	homes := make([]Home, 0, count)
	for id := 0; id < count; id++ {
		home, err := readHomeByID(ctx, id)
		if err != nil {
			return nil, err
		}
		homes = append(homes, *home)
	}
	return homes, nil
}

func (s *SmartContract) GetTransactionByID(ctx contractapi.TransactionContextInterface, id string) (*EnergyTransaction, error) {
	transactionID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	var transaction EnergyTransaction
	if err := readJSON(ctx, energyTransactionKey(transactionID), &transaction); err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *SmartContract) GetTransactions(ctx contractapi.TransactionContextInterface) ([]EnergyTransaction, error) {
	count, err := currentID(ctx, energyTxCounterKey)
	if err != nil {
		return nil, err
	}
	transactions := make([]EnergyTransaction, 0, count)
	for id := 0; id < count; id++ {
		transaction, err := s.GetTransactionByID(ctx, strconv.Itoa(id))
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}
	return transactions, nil
}

func putHome(ctx contractapi.TransactionContextInterface, home *Home) error {
	if err := putJSON(ctx, homeKey(home.Address), home); err != nil {
		return err
	}
	return ctx.GetStub().PutState(homeIDKey(home.ID), []byte(home.Address))
}

func readHome(ctx contractapi.TransactionContextInterface, address string) (*Home, error) {
	var home Home
	if err := readJSON(ctx, homeKey(address), &home); err != nil {
		return nil, err
	}
	return &home, nil
}

func readHomeByID(ctx contractapi.TransactionContextInterface, id int) (*Home, error) {
	addressBytes, err := ctx.GetStub().GetState(homeIDKey(id))
	if err != nil {
		return nil, fmt.Errorf("failed to read home id %d: %w", id, err)
	}
	if addressBytes == nil {
		return nil, fmt.Errorf("home id %d does not exist", id)
	}
	return readHome(ctx, string(addressBytes))
}

func validSignature(address string, signature string) bool {
	return signature == address+"1"
}

func parsePositiveAmount(value string) (int, error) {
	amount, err := parseNonNegativeAmount(value)
	if err != nil {
		return 0, err
	}
	if amount == 0 {
		return 0, fmt.Errorf("amount must be greater than zero")
	}
	return amount, nil
}

func parseNonNegativeAmount(value string) (int, error) {
	amount, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("amount must be an integer: %w", err)
	}
	if amount < 0 {
		return 0, fmt.Errorf("amount must not be negative")
	}
	return amount, nil
}

func newAddress(ctx contractapi.TransactionContextInterface, prefix string) string {
	sum := sha256.Sum256([]byte(prefix + ":" + ctx.GetStub().GetTxID()))
	return hex.EncodeToString(sum[:])
}

func txUnixTime(ctx contractapi.TransactionContextInterface) (int64, error) {
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return 0, fmt.Errorf("failed to read transaction timestamp: %w", err)
	}
	return timestamp.Seconds, nil
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

func homeKey(address string) string {
	return "home:" + address
}

func homeIDKey(id int) string {
	return fmt.Sprintf("home:id:%d", id)
}

func energyTransactionKey(id int) string {
	return fmt.Sprintf("energyTransaction:%d", id)
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
