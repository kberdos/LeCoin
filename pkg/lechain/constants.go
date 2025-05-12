package lechain

import (
	"lecoin/pkg/protocol"
)

const (
	START_REWARD = 6 // initial reward for mining a block
	// FIXME: change back to a higher number (23)
	START_TARGET = 15 // initial nonce hash target # of 0s
	MAX_TXS      = 4  // number of transactions that can live on a block
)

var GENESIS_HASH = protocol.HashType{}

type (
	BlockPredicate    func(b *Block) bool
	TwoBlockPredicate func(prevblock, block *Block) bool
	Handle            uint32
)
