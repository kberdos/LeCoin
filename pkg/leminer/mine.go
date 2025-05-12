package leminer

import (
	"fmt"
	"math/rand/v2"

	"lecoin/pkg/lechain"
	"lecoin/pkg/tx"
)

// the main mining loop to be run
// num_tx : the number of transactions (NOT INCLUDING THE MINER TX) to include in a block
func (lm *LeMiner) Mine(chainchan, abort chan bool, blockchan chan *lechain.Block, num_tx int) {
	// fmt.Printf("mining for target %d\n", num_tx)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-chainchan:
				continue
			default:
				b := lm.chain.NewChainBlock(lm.lekey.GetWallet())
				if err := lm.build_block(b, num_tx); err != nil {
					fmt.Println("failed to build block")
					fmt.Println(err)
					continue
				}
				fmt.Println("built a block")
				// mine the block
				if lm.mine_block(b, chainchan) {
					// fmt.Println("mined successfully")
					// successfully mined a block! here you go :-)

					// try to push the block onto the chain
					fmt.Printf("pushing block with %d txs\n", num_tx)
					_, err := lm.chain.PushBlock(b, blockchan)
					if err != nil {
						fmt.Println(err)
						fmt.Println("failed to push block ... undoing")
						go lm.mempool.InsertTxs(b.TxSlice())
						continue
					}
				} else {
					// failed :-( try to put the transactions back in the mempool
					fmt.Println("failed to mine block ... undoing")
					go lm.mempool.InsertTxs(b.TxSlice())
				}
			}
		}
	}()
	chainchan <- <-abort // tell mine_block to quit it when you're done
	done <- true         // tell the entire miner to quit it when you're done
}

// given a block, attempt to mine it
// abort when signaled (i.e. when the chain updates)
func (lm *LeMiner) mine_block(b *lechain.Block, abort chan bool) bool {
	counter := 0
	for {
		select {
		case <-abort:
			fmt.Println("aborted")
			return false
		default:
			nonce := rand.Uint32()
			if b.TryNonce(nonce) {
				b.Nonce = nonce
				return true
			}

			counter += 1
			if counter%100 == 0 {
				// fmt.Printf("Tried %d nonces\n", counter)
				if err := b.Timestamp(); err != nil {
					fmt.Println(err)
					return false
				}
			}
		}
	}
}

// wrapper around the mempool's insert
func (lm *LeMiner) InsertTx(t tx.Tx) error {
	return lm.mempool.InsertTx(t)
}

// construct a block -- coordinate with the mempool to pull off transactions
// blocks until enough transactions are available
func (lm *LeMiner) build_block(b *lechain.Block, n int) error {
	// pull off some transactions
	txs, err := lm.mempool.PopTxs(n)
	if err != nil {
		return err
	}
	// add them to the block
	if err := b.InsertTransactions(txs); err != nil {
		lm.mempool.InsertTxs(txs)
		return err
	}
	// sign the miner transaction
	if err := lm.lekey.SignTransaction(b.MinerTx); err != nil {
		lm.mempool.InsertTxs(txs)
		return err
	}
	// verify the block against the chain
	if invalidtxs, err := lm.chain.VerifyNewBlock(b); err != nil {
		// insert back only the valid txs
		to_remove := make(tx.TxMap)
		for _, tx := range invalidtxs {
			to_remove[tx.Id()] = tx
		}
		to_add := make([]tx.Tx, 0)
		for _, tx := range txs {
			if _, ok := to_remove[tx.Id()]; !ok {
				to_add = append(to_add, tx)
			}
		}
		lm.mempool.InsertTxs(to_add)
		return err
	}
	return nil
}
