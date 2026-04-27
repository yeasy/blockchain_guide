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
	expressStateKey    = "express"
	orderCounterKey    = "counter:expressOrder"
	senderPaysCode     = "0"
	receiverPaysCode   = "1"
	signedOrderState   = "signed"
	unsignedOrderState = "unsigned"
)

// SmartContract models a small logistics workflow.
type SmartContract struct {
	contractapi.Contract
}

type ExpressOrder struct {
	ID                    int      `json:"id"`
	SenderLocation        string   `json:"senderLocation"`
	ReceiverLocation      string   `json:"receiverLocation"`
	SenderAddress         string   `json:"senderAddress"`
	ReceiverAddress       string   `json:"receiverAddress"`
	SenderPhone           string   `json:"senderPhone"`
	ReceiverPhone         string   `json:"receiverPhone"`
	ExpressMoney          int      `json:"expressMoney"`
	ExpressMoneyType      string   `json:"expressMoneyType"`
	ExpressMoneySenderPay int      `json:"expressMoneySenderPay"`
	ExpressPointAddresses []string `json:"expressPointAddresses"`
	PayingMoney           int      `json:"payingMoney"`
	ExpressOrderSign      string   `json:"expressOrderSign"`
}

type User struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Address  string `json:"address"`
	PriKey   string `json:"priKey"`
	PubKey   string `json:"pubKey"`
	Phone    string `json:"phone"`
	Money    int    `json:"money"`
}

type Express struct {
	Name                  string   `json:"name"`
	Location              string   `json:"location"`
	Phone                 string   `json:"phone"`
	Money                 int      `json:"money"`
	ExpressPointAddresses []string `json:"expressPointAddresses"`
	Address               string   `json:"address"`
	PriKey                string   `json:"priKey"`
	PubKey                string   `json:"pubKey"`
}

type ExpressPoint struct {
	Name           string `json:"name"`
	Location       string `json:"location"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	PriKey         string `json:"priKey"`
	PubKey         string `json:"pubKey"`
	ExpressAddress string `json:"expressAddress"`
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, name string, location string, phone string, money string) (*User, error) {
	moneyValue, err := parseNonNegativeAmount(money)
	if err != nil {
		return nil, err
	}
	address := newAddress(ctx, "user")
	user := &User{
		Name:     name,
		Location: location,
		Address:  address,
		PriKey:   address + "1",
		PubKey:   address + "2",
		Phone:    phone,
		Money:    moneyValue,
	}
	if err := putJSON(ctx, userKey(address), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *SmartContract) CreateExpress(ctx contractapi.TransactionContextInterface, name string, location string, phone string, money string) (*Express, error) {
	moneyValue, err := parseNonNegativeAmount(money)
	if err != nil {
		return nil, err
	}
	address := newAddress(ctx, "express")
	express := &Express{
		Name:                  name,
		Location:              location,
		Phone:                 phone,
		Money:                 moneyValue,
		ExpressPointAddresses: []string{},
		Address:               address,
		PriKey:                address + "1",
		PubKey:                address + "2",
	}
	if err := putJSON(ctx, expressStateKey, express); err != nil {
		return nil, err
	}
	return express, nil
}

func (s *SmartContract) CreateExpressPoint(ctx contractapi.TransactionContextInterface, name string, location string, phone string) (*ExpressPoint, error) {
	address := newAddress(ctx, "expressPoint")
	point := &ExpressPoint{
		Name:     name,
		Location: location,
		Phone:    phone,
		Address:  address,
		PriKey:   address + "1",
		PubKey:   address + "2",
	}
	if err := putJSON(ctx, expressPointKey(address), point); err != nil {
		return nil, err
	}
	return point, nil
}

func (s *SmartContract) AddExpressPoint(ctx contractapi.TransactionContextInterface, pointAddress string) (*Express, error) {
	express, err := readExpress(ctx)
	if err != nil {
		return nil, err
	}
	point, err := readExpressPoint(ctx, pointAddress)
	if err != nil {
		return nil, err
	}

	point.ExpressAddress = express.Address
	express.ExpressPointAddresses = appendUnique(express.ExpressPointAddresses, pointAddress)

	if err := putJSON(ctx, expressPointKey(pointAddress), point); err != nil {
		return nil, err
	}
	if err := putJSON(ctx, expressStateKey, express); err != nil {
		return nil, err
	}
	return express, nil
}

func (s *SmartContract) CreateExpressOrder(ctx contractapi.TransactionContextInterface, senderLocation string, receiverLocation string, senderAddress string, receiverAddress string, senderPhone string, receiverPhone string, payType string, senderPay string, expressMoney string) (*ExpressOrder, error) {
	fee, err := parsePositiveAmount(expressMoney)
	if err != nil {
		return nil, err
	}
	senderPayValue, err := parseNonNegativeAmount(senderPay)
	if err != nil {
		return nil, err
	}
	if payType != senderPaysCode && payType != receiverPaysCode {
		return nil, fmt.Errorf("pay type must be %q for sender or %q for receiver", senderPaysCode, receiverPaysCode)
	}
	if payType == senderPaysCode && senderPayValue < fee {
		return nil, fmt.Errorf("sender prepayment must cover the express fee")
	}

	sender, err := readUser(ctx, senderAddress)
	if err != nil {
		return nil, err
	}
	if payType == senderPaysCode {
		if sender.Money < senderPayValue {
			return nil, fmt.Errorf("sender has insufficient money")
		}
		sender.Money -= senderPayValue
		if err := putJSON(ctx, userKey(senderAddress), sender); err != nil {
			return nil, err
		}
	}
	if _, err := readUser(ctx, receiverAddress); err != nil {
		return nil, err
	}

	id, err := nextID(ctx, orderCounterKey)
	if err != nil {
		return nil, err
	}
	order := &ExpressOrder{
		ID:                    id,
		SenderLocation:        senderLocation,
		ReceiverLocation:      receiverLocation,
		SenderAddress:         senderAddress,
		ReceiverAddress:       receiverAddress,
		SenderPhone:           senderPhone,
		ReceiverPhone:         receiverPhone,
		ExpressMoney:          fee,
		ExpressMoneyType:      payType,
		ExpressMoneySenderPay: senderPayValue,
		ExpressPointAddresses: []string{},
		PayingMoney:           fee,
		ExpressOrderSign:      unsignedOrderState,
	}
	if err := putJSON(ctx, expressOrderKey(id), order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *SmartContract) UpdateExpressOrder(ctx contractapi.TransactionContextInterface, orderID string, pointAddress string) (*ExpressOrder, error) {
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return nil, err
	}
	if _, err := readExpressPoint(ctx, pointAddress); err != nil {
		return nil, err
	}
	order, err := readExpressOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	order.ExpressPointAddresses = appendUnique(order.ExpressPointAddresses, pointAddress)
	if err := putJSON(ctx, expressOrderKey(id), order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *SmartContract) FinishExpressOrder(ctx contractapi.TransactionContextInterface, receiverAddress string, orderID string, receiverSignature string) (*ExpressOrder, error) {
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return nil, err
	}
	order, err := readExpressOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	if order.ReceiverAddress != receiverAddress {
		return nil, fmt.Errorf("receiver does not match order")
	}
	if !validSignature(receiverAddress, receiverSignature) {
		return nil, fmt.Errorf("invalid receiver signature")
	}
	if order.ExpressOrderSign == signedOrderState {
		return nil, fmt.Errorf("order already signed")
	}

	express, err := readExpress(ctx)
	if err != nil {
		return nil, err
	}
	receiver, err := readUser(ctx, receiverAddress)
	if err != nil {
		return nil, err
	}

	if order.ExpressMoneyType == receiverPaysCode {
		if receiver.Money < order.PayingMoney {
			return nil, fmt.Errorf("receiver has insufficient money")
		}
		receiver.Money -= order.PayingMoney
		express.Money += order.PayingMoney
	} else {
		express.Money += order.ExpressMoneySenderPay
	}
	order.ExpressOrderSign = signedOrderState
	order.PayingMoney = 0
	order.ExpressMoneySenderPay = 0

	if err := putJSON(ctx, userKey(receiverAddress), receiver); err != nil {
		return nil, err
	}
	if err := putJSON(ctx, expressStateKey, express); err != nil {
		return nil, err
	}
	if err := putJSON(ctx, expressOrderKey(id), order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *SmartContract) GetExpressOrderByID(ctx contractapi.TransactionContextInterface, id string) (*ExpressOrder, error) {
	orderID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return readExpressOrder(ctx, orderID)
}

func (s *SmartContract) GetExpress(ctx contractapi.TransactionContextInterface) (*Express, error) {
	return readExpress(ctx)
}

func (s *SmartContract) GetUserByAddress(ctx contractapi.TransactionContextInterface, address string) (*User, error) {
	return readUser(ctx, address)
}

func (s *SmartContract) GetExpressPointByAddress(ctx contractapi.TransactionContextInterface, address string) (*ExpressPoint, error) {
	return readExpressPoint(ctx, address)
}

// GetExpressPointerByAddress keeps the historical method spelling available.
func (s *SmartContract) GetExpressPointerByAddress(ctx contractapi.TransactionContextInterface, address string) (*ExpressPoint, error) {
	return s.GetExpressPointByAddress(ctx, address)
}

func readExpress(ctx contractapi.TransactionContextInterface) (*Express, error) {
	var express Express
	if err := readJSON(ctx, expressStateKey, &express); err != nil {
		return nil, err
	}
	return &express, nil
}

func readUser(ctx contractapi.TransactionContextInterface, address string) (*User, error) {
	var user User
	if err := readJSON(ctx, userKey(address), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func readExpressPoint(ctx contractapi.TransactionContextInterface, address string) (*ExpressPoint, error) {
	var point ExpressPoint
	if err := readJSON(ctx, expressPointKey(address), &point); err != nil {
		return nil, err
	}
	return &point, nil
}

func readExpressOrder(ctx contractapi.TransactionContextInterface, id int) (*ExpressOrder, error) {
	var order ExpressOrder
	if err := readJSON(ctx, expressOrderKey(id), &order); err != nil {
		return nil, err
	}
	return &order, nil
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

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

func newAddress(ctx contractapi.TransactionContextInterface, prefix string) string {
	sum := sha256.Sum256([]byte(prefix + ":" + ctx.GetStub().GetTxID()))
	return hex.EncodeToString(sum[:])
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

func userKey(address string) string {
	return "user:" + address
}

func expressPointKey(address string) string {
	return "expressPoint:" + address
}

func expressOrderKey(id int) string {
	return fmt.Sprintf("expressOrder:%d", id)
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
