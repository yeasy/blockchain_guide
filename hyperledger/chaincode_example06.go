/*
	author:swb
	emial:swbsin@163.com
	MIT License
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var ExpressOrderId int = 0

type ExpressOrder struct{
	Id int
	SenderLocation string
	ReceiverLocation string
	SenderAddress string
	ReceiverAddress string
	SenderPhone string
	ReceiverPhone string
	ExpressMoney int
	ExpressMoneyType int
	ExpressMoneySenderPay int
	ExpressPointAddress []string
	ExpressOrderStatus //0:收货方未签名 1：收货方签名
}

type Participanter struct{
	Name string
	Location string
	Address string
	PriKey string
	PubKey  string
	Phone string
	Money int
}

type ExpressCp struct{
	Name string
	Location string
	Phone string
	Money int
	ExpressPointAddress []string
}

type ExpressCpPoint struct{
	Name string
	Location string
	Phone string
	PriKey string
	PubKey string
	ExpressAddress string
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	return nil,nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil,nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil,nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//生成Address
func GetAddress() (string,string,string) {
	var address,priKey,pubKey string
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "","",""
	}

	h := md5.New()
	h.Write([]byte(base64.URLEncoding.EncodeToString(b)))

	address = hex.EncodeToString(h.Sum(nil))
	priKey = address+"1"
	pubKey = address+"2"

	return address,priKey,pubKey
}



func (t *SimpleChaincode) createParticipanter(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil,nil	
}