package lesender

import (
	"fmt"
	"sync"

	"lecoin/pkg/lbm"
	"lecoin/pkg/lechain"
	"lecoin/pkg/lkm"
	"lecoin/pkg/tx"
	"lecoin/pkg/wallet"
)

/*
used to inspect, send, and receieve lecoin
*/

type wallet_map map[int]string

type LeSender struct {
	chain   *lechain.LeChain
	lekey   *lkm.LeKeyManager
	wallets wallet_map
	Upchan  chan tx.Tx
	rwlock  sync.RWMutex
	abort   chan bool
}

func NewLeSender(chain *lechain.LeChain, lekey *lkm.LeKeyManager) *LeSender {
	return &LeSender{
		chain:   chain,
		lekey:   lekey,
		wallets: make(wallet_map),
		Upchan:  make(chan tx.Tx),
		abort:   make(chan bool),
	}
}

func (ls *LeSender) Run() {
}

// self send 0 btc (for testing)
func (ls *LeSender) LeSelfSend() error {
	w := ls.lekey.GetWallet()
	t := tx.NewTwoWayTx(w, w, 0)
	if err := ls.lekey.SignTransaction(t); err != nil {
		return err
	}
	ls.Upchan <- t
	return nil
}

// get my lebalance
func (ls *LeSender) LeBalance() tx.TxAmount {
	return ls.chain.GetBalance(ls.lekey.GetWallet())
}

func (ls *LeSender) LeBalances() lbm.LbmMap {
	// update the wallets field for sending
	return ls.chain.GetBalances()
}

func (ls *LeSender) LeListBalances() string {
	return ls.chain.BalanceString()
}

func (ls *LeSender) LeListUsers() string {
	ls.rwlock.Lock()
	defer ls.rwlock.Unlock()
	balances := ls.LeBalances()
	i := 0
	out := ""
	wallets := make(wallet_map)
	for k := range balances {
		wallets[i] = k
		wStr := wallet.PrettyFromString(k)
		out += fmt.Sprintf("%d: \t %s \n", i, wStr)
		i += 1
	}
	ls.wallets = wallets
	return out
}
