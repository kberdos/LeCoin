package tx_test

import (
	"bytes"
	"testing"

	"lecoin/pkg/lkm"
	"lecoin/pkg/tx"
)

func TestTx(t *testing.T) {
	t.Run("marshal and unmarshal miner tx", test_marshal_minertx)
	t.Run("marshal and unmarshal two way tx", test_marshal_twt)
}

func test_marshal_twt(t *testing.T) {
	k1, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	k2, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	twt := tx.NewTwoWayTx(k1.GetWallet(), k2.GetWallet(), 2)
	if err = k1.SignTransaction(twt); err != nil {
		t.Fatal(err)
	}
	newtx := tx.BlankTwoWayTx()
	newtx.Unmarshal(twt.FullMarshal())
	if !bytes.Equal(newtx.FullMarshal(), twt.FullMarshal()) {
		t.Fatal("twt marshaling is cooked")
	}
}

func test_marshal_minertx(t *testing.T) {
	k1, err := lkm.NewLKM(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	twt := tx.NewMinerTx(k1.GetWallet(), 1)
	if err = k1.SignTransaction(twt); err != nil {
		t.Fatal(err)
	}
	newtx := tx.BlankMinerTx()
	newtx.Unmarshal(twt.FullMarshal())
	if !bytes.Equal(newtx.FullMarshal(), twt.FullMarshal()) {
		t.Fatal("miner marshaling is cooked")
	}
}
