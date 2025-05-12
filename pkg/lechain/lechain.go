package lechain

import (
	"sync"

	"lecoin/pkg/lbm"
	"lecoin/pkg/protocol"
	"lecoin/pkg/tx"
	"lecoin/pkg/wallet"
)

type LeChain struct {
	// add: the headblock that we really trust (valid, longest, and first arrived)
	headblocks map[protocol.HashType]*HeadBlock // map of all head nodes
	blocks     map[protocol.HashType]*Block
	head       *HeadBlock
	handles    []chan bool // listeners
	rwlock     sync.RWMutex
}

// TODO: (extension) load from a filepath
func NewLechain() *LeChain {
	chain := &LeChain{
		headblocks: make(map[protocol.HashType]*HeadBlock),
		blocks:     make(map[protocol.HashType]*Block),
		rwlock:     sync.RWMutex{},
	}
	gblock := GenesisBlock()
	ghead := &HeadBlock{
		block_hash: gblock.BlockHash,
		lbm:        lbm.NewLBM(nil),
		length:     1, // includes the genesis block
		valid:      true,
	}
	chain.headblocks[gblock.BlockHash] = ghead
	chain.blocks[gblock.BlockHash] = gblock
	chain.head = ghead
	return chain
}

func (chain *LeChain) Length() uint32 {
	chain.rwlock.RLock()
	defer chain.rwlock.RUnlock()
	return chain.head.length
}

func (chain *LeChain) GetBalance(w *wallet.Wallet) tx.TxAmount {
	return chain.head.lbm.GetBalance(w)
}

func (chain *LeChain) GetBalances() lbm.LbmMap {
	return chain.head.lbm.GetBalances()
}

func (chain *LeChain) BalanceString() string {
	return chain.head.lbm.String()
}
