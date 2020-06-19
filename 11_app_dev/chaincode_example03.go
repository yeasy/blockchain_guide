/*
	author:swb
	time:16/7/05
	MIT License
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var bankNo int = 0
var cpNo int = 0
var transactionNo int = 0

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type CenterBank struct {
	Name        string
	TotalNumber int
	RestNumber  int
}

type Bank struct {
	Name        string
	TotalNumber int
	RestNumber  int
	ID          int
}

type Company struct {
	Name   string
	Number int
	ID     int
}

type Transaction struct {
	FromType int //CenterBank 0 Bank 1  Company 1
	FromID   int
	ToType   int //Bank 1 Company 2
	ToID     int
	Time     int64
	Number   int
	ID       int
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var totalNumber int
	var centerBank CenterBank
	var cbBytes []byte
	totalNumber, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	centerBank = CenterBank{Name: args[0], TotalNumber: totalNumber, RestNumber: 0}
	err = writeCenterBank(stub,centerBank)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}

	cbBytes,err = json.Marshal(&centerBank)
	if err!= nil{
		return nil,errors.New("Error retrieving cbBytes")
	}
	return cbBytes, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "createBank" {
		return t.createBank(stub, args)
	} else if function == "createCompany" {
		return t.createCompany(stub, args)
	} else if function == "issueCoin" {
		return t.issueCoin(stub, args)
	} else if function == "issueCoinToBank" {
		return t.issueCoinToBank(stub, args)
	} else if function == "issueCoinToCp" {
		return t.issueCoinToCp(stub, args)
	} else if function =="transfer"{
		return t.transfer(stub,args)
	}

	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) createBank(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var bank Bank
	var bankBytes []byte
	var centerBank CenterBank

	bank = Bank{Name:args[0],TotalNumber:0,RestNumber:0,ID:bankNo}

	err := writeCenterBank(stub,centerBank)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}

	bankBytes,err = json.Marshal(&bank)
	if err!= nil{
		return nil,errors.New("Error retrieving cbBytes")
	}

	bankNo = bankNo +1
	return bankBytes, nil
}

func (t *SimpleChaincode) createCompany(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var company Company
	company = Company{Name:args[0],Number:0,ID:cpNo}

	err := writeCompany(stub,company)
	if err != nil{
		return nil, errors.New("write Error" + err.Error())
	}

	cpBytes,err := json.Marshal(&company)
	if(err!=nil){
		return nil,err
	}

	cpNo = cpNo +1
	return cpBytes, nil
}

func (t *SimpleChaincode) issueCoin(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	var centerBank CenterBank
	var tsBytes []byte

	issueNumber ,err:= strconv.Atoi(args[0])
	if err!=nil{
		return nil,errors.New("want Integer number")
	}
	centerBank,_,err = getCenterBank(stub)
	if err !=nil{
		return nil,errors.New("get errors")
	}

	centerBank.TotalNumber = centerBank.TotalNumber + issueNumber
	centerBank.RestNumber = centerBank.RestNumber + issueNumber

	err = writeCenterBank(stub,centerBank)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}

	transaction := Transaction{FromType:0,FromID:0,ToType:0,ToID:0,Time:time.Now().Unix(),Number:issueNumber,ID:transactionNo}
	err = writeTransaction(stub,transaction)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}

	tsBytes,err = json.Marshal(&transaction)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}

	transactionNo = transactionNo +1
	return tsBytes, nil
}

func (t *SimpleChaincode) issueCoinToBank(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	var centerBank CenterBank
	var bank Bank
	var bankId string
	var issueNumber int
	var tsBytes []byte
	var err error
	var bankIdInt int

	bankId = args[0]
	bankIdInt,err = strconv.Atoi(args[0])
	if err!=nil{
		return nil,errors.New("want Integer number")
	}
	issueNumber,err = strconv.Atoi(args[1])
	if err!=nil{
		return nil,errors.New("want Integer number")
	}

	centerBank,_,err = getCenterBank(stub)
	if err !=nil{
		return nil,errors.New("get errors")
	}
	if centerBank.RestNumber<issueNumber{
		return nil,errors.New("Not enough money")
	}

	bank,_,err = getBankById(stub,bankId)
	if err != nil {
		return nil,errors.New("get errors")
	}
	bank.RestNumber = bank.RestNumber + issueNumber
	bank.TotalNumber = bank.TotalNumber + issueNumber
	centerBank.RestNumber = centerBank.RestNumber - issueNumber


	err = writeCenterBank(stub,centerBank)
	if err != nil {
		bank.RestNumber = bank.RestNumber - issueNumber
		bank.TotalNumber = bank.TotalNumber - issueNumber
		centerBank.RestNumber = centerBank.RestNumber + issueNumber
		return nil, errors.New("write errors"+err.Error())
	}

	err = writeBank(stub,bank)
	if err != nil {
		bank.RestNumber = bank.RestNumber - issueNumber
		bank.TotalNumber = bank.TotalNumber - issueNumber
		centerBank.RestNumber = centerBank.RestNumber + issueNumber
		err = writeCenterBank(stub,centerBank)
		if err != nil {
			return nil, errors.New("roll down errors"+err.Error())
		}
		return nil, err
	}

	transaction := Transaction{FromType:0,FromID:0,ToType:1,ToID:bankIdInt,Time:time.Now().Unix(),Number:issueNumber,ID:transactionNo}
	err = writeTransaction(stub,transaction)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}

	tsBytes,err = json.Marshal(&transaction)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}

	transactionNo = transactionNo +1
	return tsBytes, nil
}

func (t *SimpleChaincode) issueCoinToCp(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var company Company
	var bank Bank
	var bankId string
	var bankIdInt int
	var companyId string
	var companyIdInt int
	var issueNumber int
	var tsBytes []byte
	var err error

	bankId = args[0]
	bankIdInt,err = strconv.Atoi(args[0])
	if err!=nil{
		return nil,errors.New("want integer")
	}
	companyId = args[1]
	companyIdInt,err = strconv.Atoi(args[1])
	if err!=nil{
		return nil,errors.New("want integer")
	}
	issueNumber,err= strconv.Atoi(args[2])
	if err!=nil{
		return nil,errors.New("want integer")
	}

	bank,_,err = getBankById(stub,bankId)
	if err != nil {
		return nil,errors.New("get errors")
	}
	if bank.RestNumber<issueNumber{
		return nil,errors.New("Not enough money")	
	}

	company,_,err = getCompanyById(stub,companyId)
	if err != nil {
		return nil,errors.New("get errors")
	}
	bank.RestNumber = bank.RestNumber - issueNumber
	company.Number = company.Number + issueNumber

	err = writeBank(stub,bank)
	if err != nil {
		bank.RestNumber = bank.RestNumber + issueNumber
		company.Number = company.Number - issueNumber
		return nil, err
	}

	err = writeCompany(stub,company)
	if err != nil {
		bank.RestNumber = bank.RestNumber + issueNumber
		company.Number = company.Number - issueNumber
		err = writeBank(stub,bank)
		if err != nil {
			return nil, errors.New("roll down errors"+err.Error())
		}
		return nil, err
	}

	transaction := Transaction{FromType:1,FromID:bankIdInt,ToType:1,ToID:companyIdInt,Time:time.Now().Unix(),Number:issueNumber,ID:transactionNo}
	err = writeTransaction(stub,transaction)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}

	tsBytes,err = json.Marshal(&transaction)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}

	transactionNo = transactionNo +1
	return tsBytes, nil
}

func (t *SimpleChaincode) transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var cpFrom Company
	var cpTo Company
	var cpFromId string
	var cpFromIdInt int
	var cpToId string
	var cpToIdInt int
	var issueNumber int
	var tsBytes [] byte
	var err error

	cpFromId = args[0]
	cpFromIdInt,err = strconv.Atoi(args[0])
	if err!=nil{
		return nil,errors.New("want integer")
	}
	cpToId = args[1]
	cpToIdInt,err = strconv.Atoi(args[1])
	if err!=nil{
		return nil,errors.New("want integer")
	}
	issueNumber,err = strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	cpFrom,_,err = getCompanyById(stub,cpFromId)
	if err != nil {
		return nil,errors.New("get errors")
	}
	if cpFrom.Number<issueNumber{
		return nil,errors.New("Not enough money")	
	}

	cpTo,_,err = getCompanyById(stub,cpToId)
	if err != nil {
		return nil,errors.New("get errors")
	}

	cpFrom.Number = cpFrom.Number - issueNumber
	cpTo.Number = cpTo.Number + issueNumber

	err = writeCompany(stub,cpFrom)
	if err != nil {
		cpFrom.Number = cpFrom.Number - issueNumber
		cpTo.Number = cpTo.Number + issueNumber
		return nil, errors.New("write Error" + err.Error())
	}

	err = writeCompany(stub,cpTo)
	if err != nil {
		cpFrom.Number = cpFrom.Number - issueNumber
		cpTo.Number = cpTo.Number + issueNumber
		err = writeCompany(stub,cpFrom)
		if err !=nil{
			return nil,errors.New("roll down error")
		}
		return nil, errors.New("write Error" + err.Error())
	}

	transaction := Transaction{FromType:2,FromID:cpFromIdInt,ToType:2,ToID:cpToIdInt,Time:time.Now().Unix(),Number:issueNumber,ID:transactionNo}
	err = writeTransaction(stub,transaction)
	if err != nil {
		return nil, errors.New("write Error" + err.Error())
	}

	tsBytes,err = json.Marshal(&transaction)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}

	transactionNo = transactionNo +1
	return tsBytes, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	if function == "getCenterBank" {
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		_,cbBytes, err := getCenterBank(stub)
		if err != nil {
			fmt.Println("Error get centerBank")
			return nil, err
		}
		return cbBytes, nil
	} else if function == "getBankById" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		_,bankBytes, err := getBankById(stub, args[0])
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return bankBytes, nil
	} else if function == "getCompanyById" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		_,cpBytes, err := getCompanyById(stub, args[0])
		if err != nil {
			fmt.Println("Error unmarshalling centerBank")
			return nil, err
		}
		return cpBytes, nil
	} else if function == "getTransactionById" {
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		_,tsBytes, err := getTransactionById(stub, args[0])
		if err != nil {
			fmt.Println("Error unmarshalling")
			return nil, err
		}
		return tsBytes, nil
	} else if function == "getBanks" {
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		banks, err := getBanks(stub)
		if err != nil {
			fmt.Println("Error unmarshalling")
			return nil, err
		}
		bankBytes, err1 := json.Marshal(&banks)
		if err1 != nil {
			fmt.Println("Error marshalling banks")
		}	
		return bankBytes, nil
	} else if function == "getCompanys" {
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		cps, err := getCompanys(stub)
		if err != nil {
			fmt.Println("Error unmarshalling")
			return nil, err
		}
		cpBytes, err1 := json.Marshal(&cps)
		if err1 != nil {
			fmt.Println("Error marshalling banks")
		}	
		return cpBytes, nil
	} else if function == "getTransactions" {
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		tss, err := getTransactions(stub)
		if err != nil {
			fmt.Println("Error unmarshalling")
			return nil, err
		}
		tsBytes, err1 := json.Marshal(&tss)
		if err1 != nil {
			fmt.Println("Error marshalling banks")
		}	
		return tsBytes, nil
	}
	return nil,nil
}

func getCenterBank(stub *shim.ChaincodeStub) (CenterBank, []byte,error) {
	var centerBank CenterBank
	cbBytes, err := stub.GetState("centerBank")
	if err != nil {
		fmt.Println("Error retrieving cbBytes")
	}
	err = json.Unmarshal(cbBytes, &centerBank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return centerBank,cbBytes, nil
}

func getCompanyById(stub *shim.ChaincodeStub, id string) (Company,[]byte, error) {
	var company Company
	cpBytes,err := stub.GetState("company"+id)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(cpBytes, &company)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return company,cpBytes, nil
}

func getBankById(stub *shim.ChaincodeStub, id string) (Bank, []byte,error) {
	var bank Bank
	cbBytes,err := stub.GetState("bank"+id)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(cbBytes, &bank)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return bank,cbBytes, nil
}

func getTransactionById(stub *shim.ChaincodeStub, id string) (Transaction,[]byte, error) {
	var transaction Transaction
	tsBytes,err := stub.GetState("transaction"+id)
	if err != nil {
		fmt.Println("Error retrieving cpBytes")
	}
	err = json.Unmarshal(tsBytes, &transaction)
	if err != nil {
		fmt.Println("Error unmarshalling centerBank")
	}
	return transaction,tsBytes, nil
}

func getBanks(stub *shim.ChaincodeStub) ([]Bank, error) {
	var banks []Bank
	var number string 
	var err error
	var bank Bank
	if bankNo<=10 {
		i:=0
		for i<bankNo {
			number= strconv.Itoa(i)
			bank,_, err = getBankById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			banks = append(banks,bank)
			i = i+1
		}
	} else{
		i:=0
		for i<10{
			number=strconv.Itoa(i)
			bank,_, err = getBankById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			banks = append(banks,bank)
			i = i+1
		}
		return banks, nil
	}
	return nil,nil
}

func getCompanys(stub *shim.ChaincodeStub) ([]Company, error) {
	var companys []Company
	var number string 
	var err error
	var company Company
	if cpNo<=10 {
		i:=0
		for i<bankNo {
			number=strconv.Itoa(i)
			company,_ ,err = getCompanyById(stub,number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			companys = append(companys,company)
			i = i+1
		}
	} else{
		i:=0
		for i<10{
			number=strconv.Itoa(i)
			company,_ ,err = getCompanyById(stub,number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			companys = append(companys,company)
			i = i+1
		}
		return companys, nil
	}
	return nil,nil
}

func getTransactions(stub *shim.ChaincodeStub) ([]Transaction, error) {
	var transactions []Transaction
	var number string 
	var err error
	var transaction Transaction
	if transactionNo<=10 {
		i:=0
		for i<transactionNo {
			number=strconv.Itoa(i)
			transaction,_ ,err = getTransactionById(stub,number)
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
			transaction,_ ,err = getTransactionById(stub,number)
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

func writeCenterBank(stub *shim.ChaincodeStub,centerBank CenterBank) (error) {
	cbBytes, err := json.Marshal(&centerBank)
	if err != nil {
		return err
	}
	err = stub.PutState("centerBank", cbBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeBank(stub *shim.ChaincodeStub,bank Bank) (error) {
	var bankId string
	bankBytes, err := json.Marshal(&bank)
	if err != nil {
		return err
	}
	bankId= strconv.Itoa(bank.ID)
	if err!= nil{
		return errors.New("want Integer number")
	}
	err = stub.PutState("bank"+bankId, bankBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeCompany(stub *shim.ChaincodeStub,company Company) (error) {
	var companyId string
	cpBytes, err := json.Marshal(&company)
	if err != nil {
		return err
	}
	companyId= strconv.Itoa(company.ID)
	if err!= nil{
		return errors.New("want Integer number")
	}
	err = stub.PutState("company"+companyId, cpBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeTransaction(stub *shim.ChaincodeStub,transaction Transaction) (error) {
	var tsId string
	tsBytes, err := json.Marshal(&transaction)
	if err != nil {
		return err
	}
	tsId= strconv.Itoa(transaction.ID)
	if err!= nil{
		return errors.New("want Integer number")
	}
	err = stub.PutState("transaction"+tsId, tsBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}