package lechain

import (
	"lecoin/pkg/lbm"
	"lecoin/pkg/protocol"
	"lecoin/pkg/tx"
	"lecoin/pkg/wallet"
)

// the header for a block, which has the pervious and next block hashes, number of transactions, etc.

type BlockHeader struct {
	PrevHash  protocol.HashType // 0 for genesis block
	BlockHash protocol.HashType // hash of the whole block (not including this field)
	target    uint32            // number of leading 0s we're looking for
	timestamp int64             // used in hash
	valid     bool              // INTERNAL: store our verification of the block
}

// an individual block, which is both Serializable and Hashable
type Block struct {
	Nonce uint32
	BlockHeader
	MinerTx      *tx.MinerTx // payout to the miner
	transactions tx.TxMap    // map from hash of transaction to transaction
}

func BlankBlock() *Block {
	b := &Block{
		transactions: make(tx.TxMap),
	}
	b.PrevHash = protocol.HashType{}
	return b
}

// the genesis block is valid and has blockhash equal to GENESIS_HASH
func GenesisBlock() *Block {
	b := &Block{}
	b.BlockHash = GENESIS_HASH
	b.PrevHash = GENESIS_HASH
	b.valid = true
	return b
}

// NOTE: exported func for testing, but shouldn't be
func NewBlock(prevhash protocol.HashType, miner *wallet.Wallet, reward tx.TxAmount, target uint32) *Block {
	b := &Block{
		MinerTx:      tx.NewMinerTx(miner, reward),
		transactions: make(tx.TxMap),
		Nonce:        0,
	}
	// insert the miner transaction to the transaction list
	b.transactions[b.MinerTx.Id()] = b.MinerTx
	b.PrevHash = prevhash
	b.timestamp = protocol.NO_TIMESTAMP
	b.target = target
	b.valid = false // do not trust a new block by default
	return b
}

// specifies the head of the blockchain
type HeadBlock struct {
	block_hash protocol.HashType // hash of the head block
	lbm        *lbm.LeBalanceMap // the lbm for this branch
	length     uint32            // length (in blocks) of this branch
	valid      bool              // INTERNAL: store our verification of the head block
	// TODO: concurrency
}
