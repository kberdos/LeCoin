package leminer_test

import (
	"fmt"
	"testing"

	"lecoin/pkg/lechain"
	"lecoin/pkg/leminer"
	"lecoin/pkg/lkm"
	"lecoin/pkg/tx"
)

func TestMine(t *testing.T) {
	t.Run("mine block with zero tx", test_mine_zero_tx)
	t.Run("mine block with one tx", test_mine_one_tx)
	t.Run("mine block with two txs", test_mine_two_tx)
}

func test_mine_zero_tx(t *testing.T) {
	k1, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}

	chain := lechain.NewLechain()
	if chain.Length() != 1 {
		t.Fatalf("incorrect chain length: expected 1, got %d\n", chain.Length())
	}
	miner := leminer.NewLeMiner(chain, k1)
	blockchan := make(chan *lechain.Block)
	stopchan := make(chan bool)
	// mine 0 txs
	go miner.Mine(make(chan bool), stopchan, blockchan, 0)
	b := <-blockchan
	stopchan <- true
	fmt.Printf("Successfuly hashed to %08b with nonce %d\n", b.BlockHash, b.Nonce)
	if chain.Length() != 2 {
		t.Fatalf("incorrect chain length: expected 2, got %d\n", chain.Length())
	}
}

func test_mine_one_tx(t *testing.T) {
	k1, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	k2, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}

	chain := lechain.NewLechain()
	if chain.Length() != 1 {
		t.Fatalf("incorrect chain length: expected 1, got %d\n", chain.Length())
	}
	miner := leminer.NewLeMiner(chain, k1)
	blockchan := make(chan *lechain.Block)
	stopchan := make(chan bool)
	// mine 0 txs
	go miner.Mine(make(chan bool), stopchan, blockchan, 0)
	b := <-blockchan
	stopchan <- true
	fmt.Printf("Successfuly hashed to %08b with nonce %d\n", b.BlockHash, b.Nonce)
	if chain.Length() != 2 {
		t.Fatalf("incorrect chain length: expected 2, got %d\n", chain.Length())
	}
	// wait a couple seconds for the chain to calm down
	// <-time.NewTimer(2 * time.Second).C

	// insert a transaction then make another block
	twt := tx.NewTwoWayTx(k1.GetWallet(), k2.GetWallet(), 2)
	if err = k1.SignTransaction(twt); err != nil {
		t.Fatal(err)
	}
	miner.InsertTx(twt)

	// mine 1 tx
	go miner.Mine(make(chan bool), stopchan, blockchan, 1)
	b = <-blockchan
	fmt.Printf("Successfuly hashed to %08b with nonce %d\n", b.BlockHash, b.Nonce)
	stopchan <- true
	if chain.Length() != 3 {
		t.Fatalf("incorrect chain length: expected 3, got %d\n", chain.Length())
	}
}

func test_mine_two_tx(t *testing.T) {
	k1, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	k2, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}

	chain := lechain.NewLechain()
	if chain.Length() != 1 {
		t.Fatalf("incorrect chain length: expected 1, got %d\n", chain.Length())
	}
	miner := leminer.NewLeMiner(chain, k1)
	blockchan := make(chan *lechain.Block)
	stopchan := make(chan bool)
	// mine 0 txs
	go miner.Mine(make(chan bool), stopchan, blockchan, 0)
	b := <-blockchan
	stopchan <- true
	fmt.Printf("Successfuly hashed to %08b with nonce %d\n", b.BlockHash, b.Nonce)
	if chain.Length() != 2 {
		t.Fatalf("incorrect chain length: expected 2, got %d\n", chain.Length())
	}

	// insert a transaction then make another block
	twt := tx.NewTwoWayTx(k1.GetWallet(), k2.GetWallet(), 2)
	if err = k1.SignTransaction(twt); err != nil {
		t.Fatal(err)
	}
	miner.InsertTx(twt)

	// mine 1 tx
	go miner.Mine(make(chan bool), stopchan, blockchan, 1)
	b = <-blockchan
	fmt.Printf("Successfuly hashed to %08b with nonce %d\n", b.BlockHash, b.Nonce)
	stopchan <- true
	if chain.Length() != 3 {
		t.Fatalf("incorrect chain length: expected 3, got %d\n", chain.Length())
	}

	// insert another transaction then make another block
	twt = tx.NewTwoWayTx(k2.GetWallet(), k1.GetWallet(), 1)
	if err = k2.SignTransaction(twt); err != nil {
		t.Fatal(err)
	}
	miner.InsertTx(twt)

	// mine 1 tx
	go miner.Mine(make(chan bool), stopchan, blockchan, 1)
	b = <-blockchan
	fmt.Printf("Successfuly hashed to %08b with nonce %d\n", b.BlockHash, b.Nonce)
	stopchan <- true
	if chain.Length() != 4 {
		t.Fatalf("incorrect chain length: expected 3, got %d\n", chain.Length())
	}
}
