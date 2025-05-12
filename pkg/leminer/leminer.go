package leminer

import (
	"errors"
	"sync"

	"lecoin/pkg/lechain"
	"lecoin/pkg/lkm"
	"lecoin/pkg/mempool"
)

type LeMiner struct {
	chain   *lechain.LeChain
	lekey   *lkm.LeKeyManager
	mempool *mempool.Mempool
	// should have an upchan to the lcm
	Upchan  chan *lechain.Block
	Abort   chan bool
	running bool
	rwlock  sync.RWMutex
}

// NOTE: need to create mempool, use chain instead, and key manager
func NewLeMiner(lc *lechain.LeChain, lekey *lkm.LeKeyManager) *LeMiner {
	return &LeMiner{
		chain:   lc,
		lekey:   lekey,
		mempool: mempool.NewMempool(""),
		Upchan:  make(chan *lechain.Block),
		Abort:   make(chan bool),
	}
}

// TODO:
func (miner *LeMiner) Run() error {
	miner.rwlock.Lock()
	if miner.running {
		miner.rwlock.Unlock()
		return errors.New("already running")
	}
	miner.running = true
	miner.rwlock.Unlock()
	// run the mempool
	go miner.mempool.Run()
	// hook into the lechain to receieve messages from it
	chainchan := miner.chain.RegisterListener()
	// the key manager does not broadcast events (for now)

	// mine max txs per block until aborted
	// fmt.Println("running miner")

	go miner.Mine(chainchan, miner.Abort, miner.Upchan, lechain.MAX_TXS-1)
	return nil
}
