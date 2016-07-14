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

var schoolId int = 0
var RecordId int = 0

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
	BackgroundId []string
}

//当离开学校才能记入
type Background struct{
	SchoolAddress string
	StartTime string
	EndTime string
	Status //0:退学 1：毕业 
}

type Record struct{
	Id int
	SchoolAddress string
	StudentAddress string
	ModifyTime string
	ModifyOperation string
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


func (t *SimpleChaincode) createSchool(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil,nil
}

func (t *SimpleChaincode) createStudent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil,nil
}

func (t *SimpleChaincode) enrollStudent(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil,nil
}


func (t *SimpleChaincode) updateDiploma(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil,nil
}

func getStudentByAddress(stub *shim.ChaincodeStub, id string) ([]Background,[]byte, error) {
	return nil,nil,nil
}


func getRecordById(stub *shim.ChaincodeStub, id string) (Record,[]byte, error) {
	return nil,nil,nil
}

func getRecords(stub *shim.ChaincodeStub) ([]Record, error) {
	return nil,nil
}

func writeRecord(stub *shim.ChaincodeStub,centerBank CenterBank) (error) {
	return nil
}