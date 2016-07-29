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
var transactionNo int = 0

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
	Energy int
	Money int
	Id int
	Time int64
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "createUser"{
		return t.createUser(stub,args)
	}
	return nil,nil
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "changeStatus"{
		if len(args)!= 3{
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		return t.changeStatus(stub,args)
	}else if function == "buyByAddress"{
		if len(args)!= 4{
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		return t.buyByAddress(stub,args)
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
	}else if function == "getHomes"{
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		homes, err := getHomes(stub)
		if err != nil {
			fmt.Println("Error unmarshalling")
			return nil, err
		}
		homeBytes, err1 := json.Marshal(&homes)
		if err1 != nil {
			fmt.Println("Error marshalling banks")
		}	
		return homeBytes, nil
	}else if function == "getTransactions"{
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		transactions, err := getTransactions(stub)
		if err != nil {
			fmt.Println("Error unmarshalling")
			return nil, err
		}
		txBytes, err1 := json.Marshal(&transaction)
		if err1 != nil {
			fmt.Println("Error marshalling data")
		}	
		return txBytes, nil
	}else if function == "getTransactionById"{
		if len(args) != 1{
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,txBytes,err := getTransactionById(stub,args[0])
		if err != nil{
			return nil, err
		}
		return txBytes,nil
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

func buyByAddress(stub *shim.ChaincodeStub,args[] string)（[]byte,error）{
	if len(args)!= 4{
		return nil,errors.New("Incorrect number of arguments. Expecting 4")
	}
	homeSeller,homeSellerBytes,err := getHomeByAddress(stub,args[0])
	homeBuyer,homeBuyerBytes,error := getHomeByAddress(stub,args[2])

	if args[1] != args[2]+"11"{
		buyValue ,erro:= strconv.Atoi(args[3])
		if erro != nil{
			return nil, Errors.New("want integer number")
		}
		if homeSeller.Energy < buyValue && homeBuyer.Money < buyValue{
			return nil,Errors.New("not enough money or energy")
		}

		homeSeller.Energy = homeSeller.Energy - buyValue
		homeSeller.Money = homeSeller.Money + buyValue
		homeBuyer.Energy = homeBuyer.Energy + buyValue
		homeBuyer.Money = homeBuyer.Money - buyValue

		err = writeHome(stub,homeSeller)
		if err == nil{
			return nil,err
		}

		err = writeHome(stub,homeBuyer)
		if err == nil{
			return nil,err
		}

		transaction := Transaction{BuyAddress:args[2],BuyerAddressSign:args[1],SellerAddress:args[0],Energy:buyValue,Money:buyValue,Id:transactionNo,Time:time.Now().Unix()}
		err = writeTransaction(stub,transaction)
		if err != nil{
			return nil,err
		}
		transactionNo = transactionNo + 1
		txBytes ,err = json.Marshal(&transaction)

		if err!= nil{
			return nil,errors.New("Error retrieving schoolBytes")
		}

		return txBytes,nil 
	}
	return nil,err
}

func changeStatus(stub *shim.ChaincodeStub,args[] string)([]byte,error){
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	home,homeBytes,err := getHomeByAddress(stub,args[0])
	if err != nil{
		return nil
	}

	if args[1] == args[0]+"11"{
		status,_ := strconv.Atoi(args[2])
		home.Status = status
		err = writeHome(stub,home)
		if err != nil{
			return homeBytes,nil
		}
	}
	return nil,err
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

func getHomes(stub *shim.ChaincodeStub, address string)([]Home,error){
	var homes []Home
	var number string 
	var err error
	var home Home
	if homeNo<=10 {
		i:=0
		for i<homeNo {
			number= strconv.Itoa(i)
			home,_, err = getHomeById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			homes = append(homes,home)
			i = i+1
		}
	} else{
		i:=0
		for i<10{
			number=strconv.Itoa(i)
			home,_, err = getHomeById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			homes = append(homes,home)
			i = i+1
		}
		return homes, nil
	}
	return nil,nil	
}

func getTransactionById(stub *shim.ChaincodeStub,id string)(Transaction,[]byte,error){
	var transaction Transaction
	txBytes,err := stub.GetState("transaction"+id)
	if err != nil{
		fmt.Println("Error retrieving home")
	}

	err = json.Unmarshal(txBytes,&transaction)
	if err != nil {
		fmt.Println("Error unmarshalling home")
	}

	return transaction,txBytes,nil
}

func getTransactions(stub *shim.ChaincodeStub)([]Transaction,error){
	var transactions []Transaction
	var number string 
	var err error
	var transaction Transaction
	if transactionNo<=10 {
		i:=0
		for i<transactionNo {
			number= strconv.Itoa(i)
			transaction,_, err = getTransactioinById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			transactions = append(transactions,transaction)
			i = i+1
		}
	} else{
		i:=0
		for i<10{
			number=strconv.Itoa(i)
			transaction,_, err = getTransactionById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			transactions = append(transactions,transaction)
			i = i+1
		}
		return transactions, nil
	}
	return nil,nil	
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

func writeTransaction(stub *shim.ChaincodeStub,transaction Transaction)(error){
	txBytes,err := json.Marshal(&transaction)
	if err !- nil{
		return nil
	}

	id ,_:= strconv.Itoa(transaction.Id)
	err = stub.PutState("transaction"+id,txBytes)
	if err != nil{
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}