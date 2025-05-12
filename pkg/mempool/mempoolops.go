package mempool

import (
	"errors"

	"lecoin/pkg/tx"
)

// insert a single transaction, validating on its way in
func (mp *Mempool) InsertTx(t tx.Tx) error {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	_, contains := mp.txs[t.Id()]
	if t.Verify() == nil {
		mp.txs[t.Id()] = t
		if !contains { // notify waiters if new transaction hits
			mp.cond.Broadcast()
		}
		return nil
	}
	return t.Verify()
}

// insert a batch of transactions, validating on their way in
func (mp *Mempool) InsertTxs(txs []tx.Tx) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	// TODO: extension - another condition for write blocking (capacity-based)
	oldlen := len(mp.txs)
	for _, t := range txs {
		if t.Verify() == nil {
			mp.txs[t.Id()] = t
		}
	}
	// notify waiters if any new transactions hit
	if len(mp.txs) > oldlen {
		mp.cond.Broadcast()
	}
}

// pull off some transactions (no priority at this point, random)
func (mp *Mempool) PopTxs(n int) ([]tx.Tx, error) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	// wait until we have enough txs to pop
	for len(mp.txs) < n {
		mp.cond.Wait()
	}
	txs := make([]tx.Tx, 0)

	// pop n transactions
	c := 0
	for txid, t := range mp.txs {
		txs = append(txs, t)
		delete(mp.txs, txid)
		c += 1
		if c >= n {
			break
		}
	}
	if len(txs) != n {
		return nil, errors.New("did not pop n transactions")
	}
	return txs, nil
}
