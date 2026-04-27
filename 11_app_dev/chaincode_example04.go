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
	recordCounterKey     = "counter:record"
	backgroundCounterKey = "counter:background"
)

// SmartContract records school, student, and diploma changes.
type SmartContract struct {
	contractapi.Contract
}

type School struct {
	Name             string   `json:"name"`
	Location         string   `json:"location"`
	Address          string   `json:"address"`
	PriKey           string   `json:"priKey"`
	PubKey           string   `json:"pubKey"`
	StudentAddresses []string `json:"studentAddresses"`
}

type Student struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	BackgroundIDs []int  `json:"backgroundIds"`
}

type Background struct {
	ID       int    `json:"id"`
	ExitTime int64  `json:"exitTime"`
	Status   string `json:"status"`
}

type Record struct {
	ID              int    `json:"id"`
	SchoolAddress   string `json:"schoolAddress"`
	StudentAddress  string `json:"studentAddress"`
	SchoolSign      string `json:"schoolSign"`
	ModifyTime      int64  `json:"modifyTime"`
	ModifyOperation string `json:"modifyOperation"`
}

func (s *SmartContract) CreateSchool(ctx contractapi.TransactionContextInterface, name string, location string) (*School, error) {
	address := newAddress(ctx, "school")
	school := &School{
		Name:             name,
		Location:         location,
		Address:          address,
		PriKey:           address + "1",
		PubKey:           address + "2",
		StudentAddresses: []string{},
	}
	if err := putJSON(ctx, schoolKey(address), school); err != nil {
		return nil, err
	}
	return school, nil
}

func (s *SmartContract) CreateStudent(ctx contractapi.TransactionContextInterface, name string) (*Student, error) {
	address := newAddress(ctx, "student")
	student := &Student{Name: name, Address: address, BackgroundIDs: []int{}}
	if err := putJSON(ctx, studentKey(address), student); err != nil {
		return nil, err
	}
	return student, nil
}

func (s *SmartContract) EnrollStudent(ctx contractapi.TransactionContextInterface, schoolAddress string, schoolSignature string, studentAddress string) (*Record, error) {
	school, err := readSchool(ctx, schoolAddress)
	if err != nil {
		return nil, err
	}
	if !validSchoolSignature(school, schoolSignature) {
		return nil, fmt.Errorf("invalid school signature")
	}
	if _, err := readStudent(ctx, studentAddress); err != nil {
		return nil, err
	}

	school.StudentAddresses = appendUnique(school.StudentAddresses, studentAddress)
	if err := putJSON(ctx, schoolKey(schoolAddress), school); err != nil {
		return nil, err
	}
	return writeRecord(ctx, schoolAddress, studentAddress, schoolSignature, "enroll")
}

func (s *SmartContract) UpdateDiploma(ctx contractapi.TransactionContextInterface, schoolAddress string, schoolSignature string, studentAddress string, status string) (*Record, error) {
	school, err := readSchool(ctx, schoolAddress)
	if err != nil {
		return nil, err
	}
	if !validSchoolSignature(school, schoolSignature) {
		return nil, fmt.Errorf("invalid school signature")
	}
	student, err := readStudent(ctx, studentAddress)
	if err != nil {
		return nil, err
	}

	backgroundID, err := nextID(ctx, backgroundCounterKey)
	if err != nil {
		return nil, err
	}
	timestamp, err := txUnixTime(ctx)
	if err != nil {
		return nil, err
	}
	background := &Background{ID: backgroundID, ExitTime: timestamp, Status: status}
	student.BackgroundIDs = append(student.BackgroundIDs, backgroundID)

	if err := putJSON(ctx, backgroundKey(backgroundID), background); err != nil {
		return nil, err
	}
	if err := putJSON(ctx, studentKey(studentAddress), student); err != nil {
		return nil, err
	}
	return writeRecord(ctx, schoolAddress, studentAddress, schoolSignature, status)
}

func (s *SmartContract) GetStudentByAddress(ctx contractapi.TransactionContextInterface, address string) (*Student, error) {
	return readStudent(ctx, address)
}

func (s *SmartContract) GetSchoolByAddress(ctx contractapi.TransactionContextInterface, address string) (*School, error) {
	return readSchool(ctx, address)
}

func (s *SmartContract) GetRecordByID(ctx contractapi.TransactionContextInterface, id string) (*Record, error) {
	recordID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return readRecord(ctx, recordID)
}

func (s *SmartContract) GetBackgroundByID(ctx contractapi.TransactionContextInterface, id string) (*Background, error) {
	backgroundID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return readBackground(ctx, backgroundID)
}

func (s *SmartContract) GetRecords(ctx contractapi.TransactionContextInterface) ([]Record, error) {
	count, err := currentID(ctx, recordCounterKey)
	if err != nil {
		return nil, err
	}
	records := make([]Record, 0, count)
	for id := 0; id < count; id++ {
		record, err := readRecord(ctx, id)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}
	return records, nil
}

func writeRecord(ctx contractapi.TransactionContextInterface, schoolAddress string, studentAddress string, schoolSignature string, operation string) (*Record, error) {
	id, err := nextID(ctx, recordCounterKey)
	if err != nil {
		return nil, err
	}
	timestamp, err := txUnixTime(ctx)
	if err != nil {
		return nil, err
	}
	record := &Record{
		ID:              id,
		SchoolAddress:   schoolAddress,
		StudentAddress:  studentAddress,
		SchoolSign:      schoolSignature,
		ModifyTime:      timestamp,
		ModifyOperation: operation,
	}
	if err := putJSON(ctx, recordKey(id), record); err != nil {
		return nil, err
	}
	return record, nil
}

func readSchool(ctx contractapi.TransactionContextInterface, address string) (*School, error) {
	var school School
	if err := readJSON(ctx, schoolKey(address), &school); err != nil {
		return nil, err
	}
	return &school, nil
}

func readStudent(ctx contractapi.TransactionContextInterface, address string) (*Student, error) {
	var student Student
	if err := readJSON(ctx, studentKey(address), &student); err != nil {
		return nil, err
	}
	return &student, nil
}

func readRecord(ctx contractapi.TransactionContextInterface, id int) (*Record, error) {
	var record Record
	if err := readJSON(ctx, recordKey(id), &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func readBackground(ctx contractapi.TransactionContextInterface, id int) (*Background, error) {
	var background Background
	if err := readJSON(ctx, backgroundKey(id), &background); err != nil {
		return nil, err
	}
	return &background, nil
}

func validSchoolSignature(school *School, signature string) bool {
	return signature == school.PriKey || signature == school.PriKey+":signed" || signature == school.Address+"1"
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

func schoolKey(address string) string {
	return "school:" + address
}

func studentKey(address string) string {
	return "student:" + address
}

func recordKey(id int) string {
	return fmt.Sprintf("record:%d", id)
}

func backgroundKey(id int) string {
	return fmt.Sprintf("background:%d", id)
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
