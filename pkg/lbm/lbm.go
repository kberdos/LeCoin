package lbm

import (
	"errors"
	"fmt"
	"sync"

	"lecoin/pkg/tx"
	"lecoin/pkg/wallet"
)

type LbmMap map[string]tx.TxAmount

// basically just a map of wallets to balances.
type LeBalanceMap struct {
	balances LbmMap
	rwlock   sync.RWMutex
	// mutex
}

func NewLBM(oldlbm *LeBalanceMap) *LeBalanceMap {
	balances := make(LbmMap)
	if oldlbm != nil {
		// deep copy old balances map
		for k, v := range oldlbm.balances {
			balances[k] = v
		}
	}
	return &LeBalanceMap{
		balances: balances,
	}
}

func (lbm *LeBalanceMap) GetBalance(w *wallet.Wallet) tx.TxAmount {
	lbm.rwlock.RLock()
	defer lbm.rwlock.RUnlock()
	bal, ok := lbm.balances[w.String()]
	if !ok {
		return 0
	}
	return bal
}

func (lbm *LeBalanceMap) GetBalances() LbmMap {
	lbm.rwlock.RLock()
	defer lbm.rwlock.RUnlock()
	out := make(LbmMap)
	for k, v := range lbm.balances {
		out[k] = v
	}
	return out
}

/*
exhibits functional pattern -- produces a new lbm since we
are going to be likely creating a new HeadBlock with this lbm
returns any invalid transactions (so they can be axed)
*/
func (lbm *LeBalanceMap) ProcessTransactions(txs tx.TxMap, insert bool) (*LeBalanceMap, []tx.Tx, error) {
	lbm.rwlock.Lock()
	defer lbm.rwlock.Unlock()
	newlbm := NewLBM(lbm)

	var err error = nil
	invalidtxs := make([]tx.Tx, 0)
	for _, tx := range txs {
	outer:
		for _, txr := range tx.IntoRecords() {
			if e := newlbm.ProcessRecord(txr, insert); e != nil {
				err = e
				invalidtxs = append(invalidtxs, tx)
				break outer
			}
		}
	}
	if err != nil {
		return nil, invalidtxs, err
	}
	return newlbm, nil, nil
}

// errors record makes a balance negative
// if insert SET, insert the record. Otherwise, undo the record
func (lbm *LeBalanceMap) ProcessRecord(txr tx.TxRecord, insert bool) error {
	amt, ok := lbm.balances[txr.Wallet.String()]
	if !ok {
		lbm.balances[txr.Wallet.String()] = 0
	}
	delta := txr.Amount
	if !insert { // if we're undoing, negate the delta
		delta = -delta
	}
	if amt+delta < 0 {
		return errors.New("invalid: would become lebroke")
	}
	lbm.balances[txr.Wallet.String()] += delta
	return nil
}

func (lbm *LeBalanceMap) String() string {
	lbm.rwlock.Lock()
	defer lbm.rwlock.Unlock()
	out := ""
	for w, balance := range lbm.balances {
		wStr := wallet.PrettyFromString(w)
		out += fmt.Sprintf("Wallet: %s \t Balance: %0.2f\n", wStr, balance)
	}
	return out
}

// TODO: hash, serialize, verify
