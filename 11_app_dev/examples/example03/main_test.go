package main

import (
	"testing"

	"github.com/yeasy/blockchain_guide/11_app_dev/internal/testledger"
)

func TestSupplyConservationIncrementalIDsAndRejections(t *testing.T) {
	ctx, stub := testledger.New("cbdc-1", 100)
	contract := new(SmartContract)
	if _, err := contract.InitLedger(ctx, "central", "100"); err != nil {
		t.Fatal(err)
	}
	bank0, _ := contract.CreateBank(ctx, "bank-a")
	bank1, _ := contract.CreateBank(ctx, "bank-b")
	company0, _ := contract.CreateCompany(ctx, "company-a")
	company1, _ := contract.CreateCompany(ctx, "company-b")
	if bank0.ID != 0 || bank1.ID != 1 || company0.ID != 0 || company1.ID != 1 {
		t.Fatalf("IDs are not incremental: banks=%d,%d companies=%d,%d", bank0.ID, bank1.ID, company0.ID, company1.ID)
	}
	stub.SetTransaction("cbdc-2", 101)
	tx0, err := contract.IssueCoinToBank(ctx, "0", "40")
	if err != nil {
		t.Fatal(err)
	}
	stub.SetTransaction("cbdc-3", 102)
	tx1, err := contract.IssueCoinToCompany(ctx, "0", "0", "25")
	if err != nil {
		t.Fatal(err)
	}
	stub.SetTransaction("cbdc-4", 103)
	tx2, err := contract.Transfer(ctx, "0", "1", "10")
	if err != nil {
		t.Fatal(err)
	}
	if tx0.ID != 0 || tx1.ID != 1 || tx2.ID != 2 {
		t.Fatalf("transaction IDs are not incremental: %d,%d,%d", tx0.ID, tx1.ID, tx2.ID)
	}
	central, _ := contract.GetCenterBank(ctx)
	bank, _ := contract.GetBankByID(ctx, "0")
	left, _ := contract.GetCompanyByID(ctx, "0")
	right, _ := contract.GetCompanyByID(ctx, "1")
	if central.RestNumber+bank.RestNumber+left.Number+right.Number != central.TotalNumber {
		t.Fatalf("supply is not conserved: %+v %+v %+v %+v", central, bank, left, right)
	}
	if _, err := contract.Transfer(ctx, "0", "1", "100"); err == nil {
		t.Fatal("insufficient company balance must be rejected")
	}
	if _, err := contract.Transfer(ctx, "0", "0", "1"); err == nil {
		t.Fatal("self-transfer must be rejected")
	}
}
