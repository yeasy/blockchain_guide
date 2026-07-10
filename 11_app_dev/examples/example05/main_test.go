package main

import (
	"testing"

	"github.com/yeasy/blockchain_guide/11_app_dev/internal/testledger"
)

func TestEnergyTradeConservesAssetsAndStatusLifecycle(t *testing.T) {
	ctx, stub := testledger.New("home-1", 100)
	contract := new(SmartContract)
	seller, _ := contract.CreateUser(ctx, "100", "0")
	stub.SetTransaction("home-2", 101)
	buyer, _ := contract.CreateUser(ctx, "0", "100")
	stub.SetTransaction("trade-1", 102)
	tx, err := contract.BuyByAddress(ctx, seller.Address, buyer.Address+"1", buyer.Address, "20")
	if err != nil || tx.ID != 0 {
		t.Fatalf("trade = %+v, %v", tx, err)
	}
	seller, _ = contract.GetHomeByAddress(ctx, seller.Address)
	buyer, _ = contract.GetHomeByAddress(ctx, buyer.Address)
	if seller.Energy+buyer.Energy != 100 || seller.Money+buyer.Money != 100 {
		t.Fatalf("assets not conserved: seller=%+v buyer=%+v", seller, buyer)
	}
	if _, err := contract.BuyByAddress(ctx, seller.Address, buyer.Address+"1", buyer.Address, "100"); err == nil {
		t.Fatal("insufficient energy must be rejected")
	}
	if _, err := contract.BuyByAddress(ctx, buyer.Address, buyer.Address+"1", buyer.Address, "1"); err == nil {
		t.Fatal("self-trade must be rejected")
	}
	if _, err := contract.ChangeStatus(ctx, seller.Address, seller.Address+"1", "0"); err != nil {
		t.Fatal(err)
	}
	if _, err := contract.ChangeStatus(ctx, seller.Address, seller.Address+"1", "0"); err == nil {
		t.Fatal("duplicate status transition must be rejected")
	}
	if _, err := contract.ChangeStatus(ctx, seller.Address, seller.Address+"1", "2"); err == nil {
		t.Fatal("unknown status must be rejected")
	}
}
