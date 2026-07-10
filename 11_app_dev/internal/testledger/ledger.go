package testledger

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/v2/shim"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Stub is a deliberately small in-memory ledger for the example contracts.
// Embedding the full interface makes unsupported operations fail loudly if a
// future example starts relying on them.
type Stub struct {
	shim.ChaincodeStubInterface
	state     map[string][]byte
	txID      string
	timestamp *timestamppb.Timestamp
}

// New returns a transaction context backed by an in-memory ledger.
func New(txID string, seconds int64) (*contractapi.TransactionContext, *Stub) {
	stub := &Stub{state: make(map[string][]byte)}
	stub.SetTransaction(txID, seconds)
	ctx := new(contractapi.TransactionContext)
	ctx.SetStub(stub)
	return ctx, stub
}

// SetTransaction changes the transaction metadata while preserving state.
func (s *Stub) SetTransaction(txID string, seconds int64) {
	s.txID = txID
	s.timestamp = timestamppb.New(time.Unix(seconds, 0))
}

func (s *Stub) GetTxID() string {
	return s.txID
}

func (s *Stub) GetTxTimestamp() (*timestamppb.Timestamp, error) {
	return s.timestamp, nil
}

func (s *Stub) GetState(key string) ([]byte, error) {
	value := s.state[key]
	return append([]byte(nil), value...), nil
}

func (s *Stub) PutState(key string, value []byte) error {
	s.state[key] = append([]byte(nil), value...)
	return nil
}

func (s *Stub) DelState(key string) error {
	delete(s.state, key)
	return nil
}

// JSON serializes a value for direct test-ledger setup.
func JSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}
