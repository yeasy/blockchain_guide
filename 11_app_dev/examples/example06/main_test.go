package main

import (
	"testing"

	"github.com/yeasy/blockchain_guide/11_app_dev/internal/testledger"
)

func TestOrderLifecycleConservationIncrementalIDsAndDuplicates(t *testing.T) {
	ctx, stub := testledger.New("user-1", 100)
	contract := new(SmartContract)
	sender, _ := contract.CreateUser(ctx, "sender", "a", "1", "100")
	stub.SetTransaction("user-2", 101)
	receiver, _ := contract.CreateUser(ctx, "receiver", "b", "2", "50")
	stub.SetTransaction("express-1", 102)
	express, _ := contract.CreateExpress(ctx, "carrier", "hub", "3", "0")
	stub.SetTransaction("point-1", 103)
	point, _ := contract.CreateExpressPoint(ctx, "point", "hub", "4")
	if _, err := contract.AddExpressPoint(ctx, point.Address); err != nil {
		t.Fatal(err)
	}
	if _, err := contract.AddExpressPoint(ctx, point.Address); err == nil {
		t.Fatal("duplicate carrier point must be rejected")
	}
	stub.SetTransaction("order-1", 104)
	order0, err := contract.CreateExpressOrder(ctx, "a", "b", sender.Address, receiver.Address, "1", "2", senderPaysCode, "10", "10")
	if err != nil || order0.ID != 0 {
		t.Fatalf("order0 = %+v, %v", order0, err)
	}
	if _, err := contract.UpdateExpressOrder(ctx, "0", point.Address); err != nil {
		t.Fatal(err)
	}
	if _, err := contract.UpdateExpressOrder(ctx, "0", point.Address); err == nil {
		t.Fatal("duplicate route point must be rejected")
	}
	finished, err := contract.FinishExpressOrder(ctx, receiver.Address, "0", receiver.Address+"1")
	if err != nil || finished.ExpressOrderSign != signedOrderState {
		t.Fatalf("finish = %+v, %v", finished, err)
	}
	if _, err := contract.FinishExpressOrder(ctx, receiver.Address, "0", receiver.Address+"1"); err == nil {
		t.Fatal("duplicate finish must be rejected")
	}
	if _, err := contract.UpdateExpressOrder(ctx, "0", point.Address); err == nil {
		t.Fatal("signed order must reject route changes")
	}
	sender, _ = contract.GetUserByAddress(ctx, sender.Address)
	receiver, _ = contract.GetUserByAddress(ctx, receiver.Address)
	express, _ = contract.GetExpress(ctx)
	if sender.Money+receiver.Money+express.Money != 150 {
		t.Fatalf("money not conserved: sender=%d receiver=%d express=%d", sender.Money, receiver.Money, express.Money)
	}
	stub.SetTransaction("order-2", 105)
	order1, err := contract.CreateExpressOrder(ctx, "a", "b", sender.Address, receiver.Address, "1", "2", receiverPaysCode, "0", "10")
	if err != nil || order1.ID != 1 {
		t.Fatalf("order1 = %+v, %v", order1, err)
	}
	receiver.Money = 5
	stub.PutState(userKey(receiver.Address), testledger.JSON(receiver))
	if _, err := contract.FinishExpressOrder(ctx, receiver.Address, "1", receiver.Address+"1"); err == nil {
		t.Fatal("receiver with insufficient money must be rejected")
	}
}
