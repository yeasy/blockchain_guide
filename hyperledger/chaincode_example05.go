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
	"encoding/json"
	"encoding/base64"
	"encoding/hex"
	"io"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var homeNo int = 0
var transactionId int = 0

type Home struct{
	Address string
	Energy int
	Money int
	Id int
	Status int
	PriKey string
	PubKey string
}

type Transaction struct{
	BuyerAddress string
	BuyerAddressSign string
	SellerAddress string
	SellerAddressSign string
	Energy int
	Money int
	Id int
	Time string
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	return nil,nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "createUser"{
		return t.createUser(stub,args)
	}
	return nil,errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "getHomeByAddress"{
		if len(args) != 1{
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,homeBytes,err := getHomeByAddress(stub,args[0])
		if err != nil {
				fmt.Println("Error get home")
				return nil, err
			}
		return homeBytes,nil
	}
	return nil,errors.New("Received unknown function invocation")
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
	var energy , money int
	var err error
	var homeBytes []byte
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	address,priKey,pubKey := GetAddress()
	energy ,err = strconv.Atoi(args[0])
	if err !=nil{
		return nil,errors.New("want Integer number")
	}
	money , err = strconv.Atoi(args[1])
	if err != nil{
		return nil,errors.New("want Integer number")
	}
	home :=Home{Address:address,Energy:energy,Money:money,Id:homeNo,Status:1,PriKey:priKey,PubKey:pubKey}
	err = writeHome(stub,home)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}
	homeBytes,err = json.Marshal(&home)
	if err!= nil{
		return nil,errors.New("Error retrieve")
	}
	return homeBytes, nil
}

func writeHome(stub *shim.ChaincodeStub,home Home) (error) {
	homeBytes, err := json.Marshal(&home)
	if err != nil {
		return err
	}
	err = stub.PutState(home.Address, homeBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func getHomeByAddress(stub *shim.ChaincodeStub, address string) (Home,[]byte, error) {
	var home Home
	homeBytes,err := stub.GetState(address)
	if err != nil {
		fmt.Println("Error retrieving home")
	}
	err = json.Unmarshal(homeBytes, &home)
	if err != nil {
		fmt.Println("Error unmarshalling home")
	}
	return home,homeBytes, nil
}

func getHomesByEnergyline(stub *shim.ChaincodeStub, address string)([]Home,[]byte,error){
	return nil,nil,nil	
}

func buyByAddress(stub *shim.ChaincodeStub,args[] string)（[]byte,error）{

}