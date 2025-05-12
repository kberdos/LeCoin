package mempool

import (
	"sync"

	"lecoin/pkg/tx"
)

// Essentially just a store of transactions that miners can use to make blocks
// as an EXTENSION, can self-manage cache clearing
// lol another extension could be concurrent kv-store vibes with buckets
// this is VOLATILE, there's no reason to persist this since a client is going
// to retransmit their transaction if they don't see it in the universe after a short
// amount of time.

type Mempool struct {
	// extension: priority queue instead of map
	txs      tx.TxMap
	mtx      sync.Mutex
	cond     *sync.Cond
	capacity int // INTERNAL + extension: the mempool might become full at some point
}

func NewMempool(filepath string) *Mempool {
	// TODO: parse from filepath to deserialize
	mp := &Mempool{
		txs:      make(tx.TxMap),
		capacity: MEMPOOL_DEFAULT_CAPACITY,
	}
	mp.cond = sync.NewCond(&mp.mtx)
	return mp
}

func (mp *Mempool) Run() {
	// TODO: (extension) caching, eviction, etc.
}
