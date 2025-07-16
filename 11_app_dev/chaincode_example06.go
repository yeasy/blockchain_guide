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
var ExpressId int = 0

type ExpressOrder struct{
	Id int
	SenderLocation string
	ReceiverLocation string
	SenderAddress string
	ReceiverAddress string
	SenderPhone string
	ReceiverPhone string
	ExpressMoney int
	ExpressMoneyType string
	ExpressMoneySenderPay int
	ExpressPointAddress []string
	PayingMoney int
	ExpressOrderSign string //0:收货方未签名 1：收货方签名
}

type User struct{
	Name string
	Location string
	Address string
	PriKey string
	PubKey  string
	Phone string
	Money int
}

type Express struct{
	Name string
	Location string
	Phone string
	Money int
	ExpressPointerAddress []string
	Address string
	PriKey string
	PubKey string
}

type ExpressPointer struct{
	Name string
	Location string
	Phone string
	PriKey string
	PubKey string
	ExpressAddress string
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "createUser"{
		return t.createUser(stub,args)
	}else if function == "craeteExpressPointer"{
		return t.createExpressPointer(stub,args)
	}else if function == "createExpress"{
		return t.createExpress(stub,args)
	}else if function == "createExpressOrder"{
		return t.createExpressOrder(stub,args)
	}
	return nil,nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "finishExpressOrder"{
		return t.finishExpressOrder(stub,args)
	}else if function == "addExpressPointer"{
		return t.addExpressPointer(stub,args[0])
	}else if function == "updateExpressOrder"{
		return t.updateExpressOrder(stub,args[0])
	}
	return nil,nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "getExpressOrderById"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 3")
		}
		_,eoBytes,err := getExpressOrderById(stub,args[0])
		if err != nil {
			fmt.Println("Error get data")
			return nil, err
		}
		return eoBytes,nil
	}else if function == "getExpress"{
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		_,exBytes,err := getExpress(stub)
		if err != nil {
			fmt.Println("Error get data")
			return nil, err
		}
		return exBytes,nil
	}else if function == "getUserByAddress"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,userBytes,err := getUserByAddress(stub,args[0])
		if err != nil {
			fmt.Println("Error get data")
			return nil, err
		}
		return userBytes,nil
	}else if function == "getExpressPointerByAddress"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,expressCpBytes,err := getExpressPointerByAddress(stub,args[0])
		if err != nil {
			fmt.Println("Error get data")
			return nil, err
		}
		return expressCpBytes,nil
	}
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

func (t *SimpleChaincode) createUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 4{
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	var user User
	var userBytes []byte
	var number int
	var error err
	var address,priKey,pubKey string
	address,priKey,pubKey = GetAddress()

	number,err= strconv.Atoi(args[3])
	if err != nil{
		return nil,errors.New("Want integer number")
	}

    user = User{Name:args[0],Location:args[1],Address:address,PriKey:priKey,PubKey:pubKey,Phone:args[2],Money:number}
	err = writeUser(stub,user)
	if err!= nil{
		return nil,errors.New("write error")
	}

	userBytes,err = json.Marshal(&user)
	if err!= nil{
		return nil,errors.New("Error retrieve")
	}
	return userBytes, nil
}

func(t *SimpleChaincode) createExpressPointer(stub *shim.ChaincodeStub,args[] string)([]byte,error){
	if len(args) != 3{
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var expressPointer ExpressPointer
	var epBytes []byte
	var address,priKey,pubKey string
	address,priKey,pubKey = GetAddress()

	expressPointer = ExpressPointer{Name:args[0],Location:args[1],ExpressAddress:address,PriKey:priKey,PubKey:pubKey,Phone:phone}
	err = writeUser(stub,expressPointer)
	if err!= nil{
		return nil,errors.New("write error")
	}

	epBytes,err = json.Marshal(&expressPointer)
	if err!= nil{
		return nil,errors.New("Error retrieve")
	}
	return epBytes, nil
}


func(t *SimpleChaincode) createExpress(stub *shim.ChaincodeStub,args[] string)([]byte,error){
	if len(args)!= 3{
		return nil,errors.New("Incorrect number of arguments.Expecting 4")
	}

	if expressId != 0{
		return nil,errors.New("can not create two express company")
	}

	var express Express
	var pointerAddress []string
	var expressBytes []byte
	var money int
	var err error
	var address,priKey,pubKey string

	money,err= strconv.Atoi(args[3])
	if err != nil{
		return nil,errors.New("Want integer number")
	}

	express = Express{Name:args[0],Location:args[1],Phone:args[2],Money:money,Address:address,PriKey:priKey,PubKey:pubKey,ExpressPointerAddress:pointerAddress}
	err = writeExpress(stub,express)

	if err!= nil{
		return nil,errors.New("write error")
	}

	expressBytes,err = json.Marshal(&express)
	if err!= nil{
		return nil,errors.New("Error retrieve")
	}
	return expressBytes, nil
}

func(t *SimpleChaincode) createExpressOrder(stub *shim.ChaincodeStub,args[] string)([]byte,error){
	if len(args) !=8{
		return nil,errors.New("Incorrect number of arguments.Expecting 9")
	}

	var expressOrder ExpressOrder
	var expressPointAddress []string
	var eoBytes []byte
	var err error
	var money int
	var payingMoney int

	money,err = strconv.Atoi(args[7])
	if err != nil{
		return nil,errors.New("want integer number")
	}

	payingMoney,err = strconv.Atoi(args[8])
	if err != nil{
		return nil,errors.New("want integer number")
	}

	expressOrder = ExpressOrder{Id:ExpressOrderId,SenderLocation:args[0],ReceiverLocation:args[1],SenderAddress:args[2],ReceiverAddress:args[3],SenderPhone:args[4],ReceiverPhone:args[5],ExpressMoneyType:args[6],ExpressMoney:args[7],ExpressPointAddress:expressPointAddress,PayingMoney:payingMoney}
	err = writeExpressOrder(expressOrder)

	if err!= nil{
		return err
	}

	ExpressOrderId = ExpressOrderId + 1

	eoBytes,err = json.Marshal(&expressOrder)
	if err!= nil{
		return nil,errors.New("Error retrieve data")
	}
	return eoBytes,nil
}


func(t *SimpleChaincode) addExpressPointer(stub *shim.ChaincodeStub,address string)([]byte,error){
	express,expressBytes,err := getExpress(stub)
	if err != nil{
		return nil,errors.New("Error get data")
	}

	expressOrder,eoBytes,error := getExpressPointerByAddress(stub,address)
	if error != nil{
		return nil,errors.New("Error get data")
	}

	express.ExpressPointerAddress = append(express.ExpressPointerAddress,expressOrder.ExpressAddress)
	err = writeExpress(express)
	if err != nil{
		return nil,errors.New("Error write data")
	}

	expressBytes,err := json.Marshal(&express)
	if err != nil{
		return nil
	}
	return expressBytes,nil
}

func(t *SimpleChaincode) updateExpressOrder(stub *shim.ChaincodeStub,args[] string)([]byte,error){
	if len(arg) != 2{
		return nil,errors.New("Incorrect number of arguments.Expecting 2")
	}

	expressOrder,epBytes,err := getExpressOrderById(args[0])
	if err != nil{
		return nil,errors.New("get error")
	}

	expressOrder.ExpressPointAddress = append(expressOrder.ExpressPointAddress,args[1])
	err = writeExpressOrder(expressOrder)
	if err != nil{
		return nil,nil
	}

	eoBytes,err := json.Marshal(&expressOrder)
	if err != nil{
		return nil
	}

	return eoBytes,nil
}

func(t *SimpleChaincode) finishExpressOrder(stub *shim.ChaincodeStub,args[] string)([]byte,error){
	if len(arg)!= 3{
		return nil,errors.New("Incorrect number of arguments.Expecting 3")
	}

	user,err := getUserByAddress(stub,args[0])
	if err != nil{
		return err
	}

	expressOrder,epBytes,error := getExpressOrderById(args[1])
	if error != nil{
		return error
	}

	express,epBytes,errorone := getExpress(stub)
	errorone != nil{
		return errorone
	}

	if expressOrder.ExpressMoney!=0{
		express.Money = express.Money + expressOrder.ExpressMoney
		expressOrder.ExpressMoney = 0
	}else{
		express.Money = express.Money + expressOrder.PayingMoney
		user.Money = user.Money -expressOrder.PayingMoney
	}

	expressOrder.ExpressOrderSign = args[3]

	err = writeExpress(express)
	if err != nil{
		return err
	}

	err = writeUser(user)
	if err != nil{
		return err
	}

	err = writeExpressOrder(expressOrder)
	if err != nil{
		return err
	}

	eoBytes,err := json.Marshal(&expressOrder)
	if err != nil{
		return nil
	}

	return eoBytes,nil
}

func(t *SimpleChaincode) getExpressOrderById(stub *shim.ChaincodeStub,id string)(ExpressOrder,[]byte,error){
	var expressOrder ExpressOrder
	eoBytes,err := stub.GetState("expressOrder"+id)
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(eoBygtes,&ExpessOrder)
	if err != nil{
		fmt.Println("Error unmarshalling data")
	}
	return expressOrder,eoByes,nil
} 

func(t *SimpleChaincode) getExpress(stub *shim.ChaincodeStub)(Express,[]byte,error){
	var express Express
	exBytes,err := stub.GetState("Express")
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(exBytes,&express)
	if err != nil {
		fmt.Println("Error unmarshalling data")
	}

	return express,exBytes,nil
}

func(t *SimpleChaincode) getUserByAddress(stub *shim.ChaincodeStub,address string)(User,[]byte,error){
	var user User
	userBytes,err := stub.GetState(address)
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(userBytes,&user)
	if err != nil {
		fmt.Println("Error unmarshalling data")
	}

	return user,userBytes,nil
}

func(t *SimpleChaincode) getExpressPointerByAddress(stub *shim.ChaincodeStub,address string)(ExpressPointer,[]byte,error){
	var expressPointer ExpressPointer
	epBytes,err := stub.GetState(address)
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(epBytes,&expressPointer)
	if err != nil {
		fmt.Println("Error unmarshalling data")
	}

	return expressPointer,epBytes,nil
}

func(t *SimpleChaincode) writeExpress(stub *shim.ChaincodeStub,Express express)(error){
	exBytes,err := json.Marshal(&express)
	if err != nil {
		return err
	}
	err = stub.PutState("Express", exBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func(t *SimpleChaincode) writeExpressOrder(stub *shim.ChaincodeStub,ExpressOrder expressOrder)(error){
	eoBytes,err := json.Marshal(&expressOrder)
	if err != nil {
		return err
	}

	id ,_:= strconv.Itoa(expressOrder.Id)
	err = stub.PutState("ExpressOrder"+id,eoBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func(t *SimpleChaincode) writeUser(stub *shim.chaincodeStub,User user)(error){
	userBytes,err := json.Marshal(&user)
	if err != nil {
		return err
	}
	err = stub.PutState(user.Address, userBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func(t *SimpleChaincdoe) writeExpressPointer(stub *shim.chaincodeStub,ExpressPointer expressPointer)(error){
	epBytes,err := json.Marshal(&expressPointer)
	if err != nil {
		return err
	}

	err = stub.PutState(user.Address, userBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}