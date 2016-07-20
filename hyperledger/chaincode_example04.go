/*
	author:swb
	emial:swbsin@163.com
	MIT License
*/

package main

import (
	"errors"
	"fmt"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SimpleChaincode struct {
}

var BackGroundNo int = 0
var RecordNo int = 0

type School struct{
	Name string
	Location string
	Address string
	PriKey string
	PubKey string
	StudentAddress []string
}

type Student struct{
	Name string
	Address string
	BackgroundId []int
}

//当离开学校才能记入
type Background struct{
	Id int
	ExitTime int64
	Status string //0:毕业 1：退学 
}

type Record struct{
	Id int
	SchoolAddress string
	StudentAddress string
	SchoolSign string
	ModifyTime int64
	ModifyOperation string
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "createSchool"{
		return t.createSchool(stub,args)
	}else if function == "createStudent"{
		return t.createStudent(stub,args)
	}
}

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "enrollStudent"{
		if len(args)!= 3{
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		return t.enrollStudent(stub,args)
	}else if function == "updateDiploma"{
		if len(args)!= 4{
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		return t.updateDiploma(stub,args)
	}
	return nil,nil
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "getStudentByAddress"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,stuBytes, err := getStudentByAddress(stub,args[0])
		if err != nil {
			fmt.Println("Error get centerBank")
			return nil, err
		}
		return stuBytes, nil
	}else if function == "getRecordById"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,recBytes, err := getRecordById(stub,args[0])
		if err != nil {
			fmt.Println("Error get centerBank")
			return nil, err
		}
		return recBytes, nil
	}else if function == "getRecords"{
		if len(args) != 0 {
			return nil, errors.New("Incorrect number of arguments. Expecting 0")
		}
		records, err := getRecords(stub)
		if err != nil {
			fmt.Println("Error unmarshalling")
			return nil, err
		}
		recBytes, err1 := json.Marshal(&records)
		if err1 != nil {
			fmt.Println("Error marshalling banks")
		}	
		return recBytes, nil
	}else if function == "getSchoolByAddress"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,schBytes, err := getSchoolByAddress(stub,args[0])
		if err != nil {
			fmt.Println("Error get centerBank")
			return nil, err
		}
		return schBytes, nil
	}else if function == "getBackgroundById"{
		if len(args) != 1 {
			return nil, errors.New("Incorrect number of arguments. Expecting 1")
		}
		_,backBytes, err := getBackgroundById(stub,args[0])
		if err != nil {
			fmt.Println("Error get centerBank")
			return nil, err
		}
		return backBytes, nil
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
	var address,
	return address,priKey,pubKey
}


func (t *SimpleChaincode) createSchool(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2{
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	var school School
	var schoolBytes []byte
	var stuAddress []string
	var address,priKey,pubKey string
	address,priKey,pubKey = GetAddress()

	school = School {Name:args[0],Location:args[1],Address:address,PriKey:priKey,PubKey:pubKey,StudentAddress:studentAddress}
	err := writeSchool(stub,school)
	if err != nil{
		return nil, errors.New("write Error" + err.Error())
	}

	schoolBytes ,err = json.Marshal(&school)
	if err!= nil{
		return nil,errors.New("Error retrieving schoolBytes")
	}


	return schoolBytes,nil
}

func (t *SimpleChaincode) createStudent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1{
		return nil,errors.New("Incorrect number of arguments. Expecting 1")
	}

	var student Student
	var studentBytes []byte
	var stuAddress string 
	var bgId []int
	stuAddress,_,_ = GetAddress()

	student = Student{Name:args[0],Address:stuAddress,BackgroundId,bgId}
	err := writeStudent(stub,student)
	if err != nil{
		return nil,errors.New("write Error" + err.Error())
	}

	studentBytes,err = json.Marshal(&student)
	if err!= nil{
		return nil,errors.New("Error retrieving schoolBytes")
	}

	return studentBytes,nil
}

func (t *SimpleChaincode) enrollStudent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	school,schBytes,error:=getSchoolByAddress(stub,args[0])

	if error != nil{
		return nil,errors.New("Error get data")
	}
	student,stuBytes,err:= getStudentByAddress(stub,args[2])
	if err != nil{
		return nil,errors.New("Error get data")
	}
	schoolSign := args[1]
	var record Record
	record = Record{Id:RecordNo,SchoolAddress:args[0],StudentAddress:args[2],SchoolSign:schoolSign,ModifyTime:time.Now().Unix(),ModifyOperation:"2"}

	err = writeRecord(stub,record)
	if err != nil{
		return nil,errors.New("Error write data")
	}

	schoolAddress.StudentAddress = append(schoolAddress.StudentAddress,student.Address)
	err = writeSchool(stub,School)
	if err != nil{
		return nil,errors.New("Error write data")
	}

	err = writeStudent(stub,School)
	if err!= nil{
		return nil,errors.New("Error write data")
	}

	RecordNo = RecordNo + 1
	recordBytes,err = json.Marshal(&record)
	
	if err!= nil{
		return nil,errors.New("Error retrieving schoolBytes")
	}

	return recordBytes,nil
}

func (t *SimpleChaincode) updateDiploma(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var recordBytes []byte
	school,schBytes,error:=getSchoolByAddress(stub,args[0])
	if error != nil{
		return nil,errors.New("Error get data")
	}

	student,stuBytes,err:= getStudentByAddress(stub,args[2])
	if err != nil{
		return nil,errors.New("Error get data")
	}
	schoolSign := args[1]
	var record Record
	record = Record{Id:RecordNo,SchoolAddress:args[0],StudentAddress:args[2],SchoolSign:schoolSign,ModifyTime:time.Now().Unix(),ModifyOperation:args[3]}

	err = writeRecord(stub,record)
	if err != nil{
		return nil,errors.New("Error write data")
	}

	var background Background
	background = Background{Id:BackGroundNo,ExitTime:time.Now().Unix(),Status:args[3]}

	err = writeBackground(stub,background)
	if err != nil{
		return nil,errors.New("Error write data")
	}

	BackGroundNo = BackGroundNo + 1
	recordBytes ,err = json.Marshal(&record)
	
	if err!= nil{
		return nil,errors.New("Error retrieving schoolBytes")
	}
	err = writeStudent(stub,student)
	if err != nil{
		return nil,errors.New("Error write data")
	}
	err = writeSchool(stub,school)
	if err != nil{
		return nil,errors.New("Error write data")
	}

	return recordBytes,nil
}

func getStudentByAddress(stub *shim.ChaincodeStub, address string) (Student,[]byte, error) {
	var student Student
	stuBytes,err := stub.GetState("address")
	if err != nil {
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(stuBytes, &student)
	if err != nil {
		fmt.Println("Error unmarshalling data")
	}
	return student,stuBytes, nil
}

func getSchoolByAddress(stub *shim.ChaincodeStub,address string)(School,[]byte,error){
	var school School
	schBytes,err := stub.GetState("address")
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(schBytes,&school)
	if err != nil{
		fmt.Println("Error unmarshalling data")
	}
	return school,schBytes,nil
}


func getRecordById(stub *shim.ChaincodeStub, id string) (Record,[]byte, error) {
	var record Record
	recBytes,err := stub.GetState("Record"+id)
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(recBytes,&Record)
	if err != nil{
		fmt.Println("Error unmarshalling data")
	}
	return record,recBytes,nil
}

func getRecords(stub *shim.ChaincodeStub) ([]Record, error) {
	var records []Record
	var number string 
	var err error
	var record Record
	if RecordNo<=10 {
		i:=0
		for i<RecordNo {
			number= strconv.Itoa(i)
			record,_, err = getRecordById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			records = append(records,record)
			i = i+1
		}
	} else{
		i:=0
		for i<10{
			number=strconv.Itoa(i)
			record,_, err = getRecordById(stub, number)
			if err != nil {
				return nil, errors.New("Error get detail")
			}
			records = append(records,record)
			i = i+1
		}
		return records, nil
	}
	return nil,nil
}

func getBackgroundById(stub * shim.ChaincodeStub,id string)(Background,[]byte,error){
	var background Background
	backBytes,err := stub.GetState("BackGround"+id)
	if err != nil{
		fmt.Println("Error retrieving data")
	}

	err = json.Unmarshal(backBytes,&BackGround)
	if err != nil{
		fmt.Println("Error unmarshalling data")
	}
	return background,backBytes,nil
}

func writeRecord(stub *shim.ChaincodeStub,record Record) (error) {
	var recId string
	recordBytes,err :=json.Marshal(&record)
	if err != nil{
		return err
	}

	recId,_ = strconv.Itoa(record.Id)
	err = stub.PutState("Record"+recId, recordBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeSchool(stub *shim.ChaincodeStub,school School)(error){
	schBytes ,err := json.Marshal(&school)
	if err != nil{
	    return err
	}

	err = stub.PutState(school.Address,schBytes)
	if err !=nil{
		return errors.New("PutState Error" + err.Error())
	}
	return nil
}

func writeStudent(stub *shim.ChaincodeStub,student Student)(error){
	stuBytes,err :=  json.Marshal(&student)
	if err != nil{
		return err
	}

	err = stub.PutState(student.Address,stuBytes)
	if err != nil{
		return errors.New("PutState Error" + err.Error())
	}	
	return nil
}

func writeBackground(stub *shim.ChaincodeStub,background Background)(error){
	var backId string
	backBytes,err :=json.Marshal(&background)
	if err != nil{
		return err
	}

	backId,_ = strconv.Itoa(background.Id)
	err = stub.PutState("BackGround"+backId, backBytes)
	if err != nil {
		return errors.New("PutState Error" + err.Error())
	}
	return nil	
}