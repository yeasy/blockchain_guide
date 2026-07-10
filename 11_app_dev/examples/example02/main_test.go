package main

import (
	"testing"

	"github.com/yeasy/blockchain_guide/11_app_dev/internal/testledger"
)

func TestTransferConservesBalancesAndRejectsInvalidOperations(t *testing.T) {
	ctx, _ := testledger.New("transfer-1", 100)
	contract := new(SmartContract)
	if err := contract.InitAccounts(ctx, "alice", "100", "bob", "50"); err != nil {
		t.Fatal(err)
	}
	if _, err := contract.Transfer(ctx, "alice", "bob", "30"); err != nil {
		t.Fatal(err)
	}
	alice, _ := contract.ReadAccount(ctx, "alice")
	bob, _ := contract.ReadAccount(ctx, "bob")
	if alice.Balance != 70 || bob.Balance != 80 || alice.Balance+bob.Balance != 150 {
		t.Fatalf("unexpected balances: alice=%d bob=%d", alice.Balance, bob.Balance)
	}
	if _, err := contract.Transfer(ctx, "alice", "bob", "100"); err == nil {
		t.Fatal("insufficient funds must be rejected")
	}
	if _, err := contract.Transfer(ctx, "alice", "alice", "10"); err == nil {
		t.Fatal("self-transfer must be rejected instead of minting value")
	}
	alice, _ = contract.ReadAccount(ctx, "alice")
	bob, _ = contract.ReadAccount(ctx, "bob")
	if alice.Balance != 70 || bob.Balance != 80 {
		t.Fatalf("rejected operations changed state: alice=%d bob=%d", alice.Balance, bob.Balance)
	}
}
