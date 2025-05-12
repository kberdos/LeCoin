package lechain

import "lecoin/pkg/tx"

func verify_block(b *Block) bool {
	return b.Verify() == nil
}

// for now, all we care about is the block having the default target # of hashes
func verify_target(prevblock, block *Block) bool {
	// TODO: EXTENSION: verify the challenge matches the time stamps
	return block.target == START_TARGET
}

// verify the transactions of a (potential) new block with the existing chain.
// mainly used for mining, and does NOT verify the block nonce
func (lc *LeChain) VerifyNewBlock(b *Block) ([]tx.Tx, error) {
	lc.rwlock.RLock()
	defer lc.rwlock.RUnlock()
	if txs, err := b.verify_txs(); err != nil {
		return txs, err
	}
	// get the old head
	var oldhead *HeadBlock
	var err error
	for _, hb := range lc.headblocks {
		oldhead, err = lc.traverseBack(hb, b.PrevHash, verify_block, verify_target)
		if err == nil {
			break
		}
	}
	if err != nil {
		return make([]tx.Tx, 0), err
	}
	// process the transactions with INSERT set
	_, txs, err := oldhead.lbm.ProcessTransactions(b.transactions, true)
	if err != nil {
		return txs, err
	}
	return nil, nil
}

// verify a branch of the chain through the headblock (read locking)
// XXX: not sure where this is useful
func (lc *LeChain) verify_chain(hb *HeadBlock) error {
	lc.rwlock.RLock()
	defer lc.rwlock.RUnlock()
	// do not repeat verifications
	if hb.valid {
		return nil
	}
	// traverse back to the genesis block, verifying every block
	// and challenge. If no errors, then the lbm is valid and we
	// can bless this chain
	_, err := lc.traverseBack(hb, GENESIS_HASH, verify_block, verify_target)
	hb.valid = (err != nil)
	return err
}

// compute the longest (and thus solely valid) chain.
// the chain that gets to us first wins
// assumed the lc is write locked
func (lc *LeChain) set_head() {
	// go through the headblocks and find the longest / oldest
	cur_head := lc.head
	changed := false
	for _, hb := range lc.headblocks {
		// longer wins
		if hb.length > cur_head.length {
			cur_head = hb
			changed = true
		}
	}
	// head changed, so tell all listeners
	lc.head = cur_head
	if changed {
		// fmt.Printf("head changed! length %d\n", lc.head.length)
		lc.Broadcast()
	}
}
