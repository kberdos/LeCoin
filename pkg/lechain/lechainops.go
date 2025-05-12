package lechain

import (
	"bytes"
	"errors"

	"lecoin/pkg/protocol"
	"lecoin/pkg/wallet"
)

/*
specifies operations on the blockchain
*/

/*
given a block hash, find the block associated with that block hash
*/
func (lc *LeChain) getBlock(block_hash protocol.HashType) (*Block, error) {
	block, ok := lc.blocks[block_hash]
	if !ok {
		return &Block{}, errors.New("could not find block")
	}
	return block, nil
}

/*
Traverse backwards until we find the block, constructing a new headblock
- returns the headblock'd version of the block with the updated lbm
- also returns whether or not the block was already a head
- The following hofs are useful for verifying the entire chain
  - bp: BlockPredicate that can be run on every block (errors if eval to false)
  - tbp: TwoBlockPredicate that can be rnu on every pair of blocks (errors if eval to false)

- assumes at least rlocking
*/
// TODO: (extension) prevent loops with a seen list?
func (lc *LeChain) traverseBack(hb *HeadBlock, target_hash protocol.HashType, bp BlockPredicate, tbp TwoBlockPredicate) (*HeadBlock, error) {
	// fmt.Printf("t block: %08b\n", target_hash)
	// fmt.Printf("f block: %08b\n", hb.block_hash)
	lbm := hb.lbm
	block, err := lc.getBlock(hb.block_hash)
	if err != nil {
		return &HeadBlock{}, err
	}
	// fmt.Println("we have this head")
	var steps uint32 = 0
	for !bytes.Equal(block.BlockHash[:], target_hash[:]) {
		if bytes.Equal(block.BlockHash[:], GENESIS_HASH[:]) || steps == hb.length-1 {
			return &HeadBlock{}, errors.New("dead end: back to the genesis block")
		}
		if bp != nil && !bp(block) {
			return &HeadBlock{}, errors.New("block predicate failed")
		}
		// undo all transactions (insert is UNSET)
		if lbm, _, err = lbm.ProcessTransactions(block.transactions, false); err != nil {
			return &HeadBlock{}, err
		}
		prevblock, err := lc.getBlock(block.PrevHash)
		// fmt.Printf("prevblock has hash %08d\n", prevblock.BlockHash)
		if err != nil {
			return &HeadBlock{}, err
		}
		if tbp != nil && !tbp(prevblock, block) {
			return &HeadBlock{}, errors.New("two block predicate failed")
		}
		block = prevblock
		// fmt.Println("set block to prevblock")
		steps += 1
	}
	// and construct a headblock
	return &HeadBlock{
		block_hash: block.BlockHash,
		lbm:        lbm,
		length:     hb.length - steps,
	}, nil
}

/*
Pushes a new block onto the head of the chain
Requires that the previous block is an existing head of the chain
- find a HeadBlock to link to
- change the fields in the headblock struct to point to this block
- if found block is not a headblock, create a new head block and traverse back (hopefully not far!) to the block and unto transactions
- add the block to the hashmap
- update the lbm for this branch with the new transactions
- at every step, verify and error if anything is invalid
- this entire operation should be atomic with write lock
*/
// FIXME: should separate the blockchan from pushblock and do that separately
func (lc *LeChain) PushBlock(b *Block, blockchan chan *Block) (*HeadBlock, error) {
	// fmt.Println("before the locks")
	lc.rwlock.Lock()
	defer lc.rwlock.Unlock()
	// fmt.Println("past the locks")
	// validate this new block
	if err := b.Verify(); err != nil {
		return nil, err
	}
	// fmt.Println("past verification")
	// traverse backwards through each head to find the head of the chain corresponding to this block
	var unc_head *HeadBlock
	var err error
	for _, hb := range lc.headblocks {
		unc_head, err = lc.traverseBack(hb, b.PrevHash, nil, nil)
		if err == nil {
			break
		}
	}
	if err != nil {
		return &HeadBlock{}, err
	}
	// fmt.Printf("found unc head with length %d and hash %08b\n", unc_head.length, unc_head.block_hash)
	// update the lbm of the old head to make a new lbm
	lbm, _, err := unc_head.lbm.ProcessTransactions(b.transactions, true)
	if err != nil {
		return &HeadBlock{}, err
	}
	// fmt.Println("past process transactions")
	// add the new head
	new_hb := &HeadBlock{
		block_hash: b.BlockHash,
		lbm:        lbm,
		length:     unc_head.length + 1,
	}
	// validate the target through the timestamps
	unc_block, err := lc.getBlock(unc_head.block_hash)
	if err != nil {
		return &HeadBlock{}, err
	}
	if !verify_target(unc_block, b) {
		return &HeadBlock{}, errors.New("could not verify target on new block")
	}

	// insert the new headblock and new block and recalculate the valid head
	// fmt.Printf("b len of headblocks is %d\n", unc_head.length)
	// fmt.Printf("INSERTING hash is %08b\n", new_hb.block_hash)
	lc.headblocks[new_hb.block_hash] = new_hb
	// fmt.Printf("a len of headblock is %d\n", new_hb.length)
	lc.blocks[b.BlockHash] = b
	lc.set_head()
	if blockchan != nil {
		blockchan <- b
	}
	return new_hb, nil
}

/*
Generate an empty block from the valid head of the chain
*/
func (lc *LeChain) NewChainBlock(miner *wallet.Wallet) *Block {
	// TODO: (extension) calculate the reward and target
	lc.rwlock.RLock()
	defer lc.rwlock.RUnlock()
	return NewBlock(lc.head.block_hash, miner, START_REWARD, START_TARGET)
}

// orphan block (private)

// broadcast block?

//
