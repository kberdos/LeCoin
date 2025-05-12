package lkm

import (
	"crypto/ecdsa"
	"crypto/rand"

	"lecoin/pkg/tx"
)

// sign a transaction (that we produced)
func (lmk *LeKeyManager) SignTransaction(tx tx.Tx) error {
	// marhsal the transaction
	txhash := tx.Marshal()
	// generate the signature of the transaction
	sig, err := ecdsa.SignASN1(rand.Reader, lmk.privkey, txhash)
	if err != nil {
		return err
	}
	if err = tx.Sign(sig); err != nil {
		return err
	}
	return nil
}
