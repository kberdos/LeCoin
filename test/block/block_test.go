package block_test

import (
	"bytes"
	"fmt"
	"testing"

	"lecoin/pkg/lechain"
	"lecoin/pkg/leminer"
	"lecoin/pkg/lkm"
	"lecoin/pkg/tx"
)

func TestBlock(t *testing.T) {
	t.Run("marshal and unmarshal a block", test_marshal_block)
}

func test_marshal_block(t *testing.T) {
	k1, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	k2, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	chain := lechain.NewLechain()
	miner := leminer.NewLeMiner(chain, k1)
	blockchan := make(chan *lechain.Block)
	stopchan := make(chan bool)

	// send yourself 0 btc a couple times to make a valid block
	twt := tx.NewTwoWayTx(k1.GetWallet(), k1.GetWallet(), 0)
	if err = k1.SignTransaction(twt); err != nil {
		t.Fatal(err)
	}
	miner.InsertTx(twt)
	twt = tx.NewTwoWayTx(k2.GetWallet(), k2.GetWallet(), 0)
	if err = k2.SignTransaction(twt); err != nil {
		t.Fatal(err)
	}
	miner.InsertTx(twt)

	// mine the txs
	fmt.Println("starting mining...")
	go miner.Mine(make(chan bool), stopchan, blockchan, 2)
	b := <-blockchan
	stopchan <- true

	fmt.Println("got past the mining part")
	newblock := lechain.BlankBlock()
	newblock.Unmarshal(b.FullMarshal())
	if !bytes.Equal(newblock.FullMarshal(), b.FullMarshal()) {
		t.Fatal("marshaling is cooked")
	}
}

