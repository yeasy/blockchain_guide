package main

import (
	"testing"

	"github.com/yeasy/blockchain_guide/11_app_dev/internal/testledger"
)

func TestReadWriteLifecycle(t *testing.T) {
	ctx, _ := testledger.New("write-1", 100)
	contract := new(SmartContract)
	if err := contract.Write(ctx, "asset", "v1"); err != nil {
		t.Fatal(err)
	}
	value, err := contract.Read(ctx, "asset")
	if err != nil || value != "v1" {
		t.Fatalf("first read = %q, %v", value, err)
	}
	if err := contract.Write(ctx, "asset", "v2"); err != nil {
		t.Fatal(err)
	}
	value, err = contract.Read(ctx, "asset")
	if err != nil || value != "v2" {
		t.Fatalf("updated read = %q, %v", value, err)
	}
	if _, err := contract.Read(ctx, "missing"); err == nil {
		t.Fatal("missing key must fail explicitly")
	}
}
