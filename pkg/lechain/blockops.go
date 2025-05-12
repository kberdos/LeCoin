package lechain

import (
	"encoding/binary"
	"errors"
	"fmt"

	"lecoin/pkg/protocol"
	"lecoin/pkg/tx"
)

func (b *Block) InsertTransactions(txs []tx.Tx) error {
	if len(b.transactions)+len(txs) > MAX_TXS {
		return errors.New("cannot add: block is full")
	}
	for _, t := range txs {
		b.transactions[t.Id()] = t
	}
	return nil
}

func (b *Block) TxSlice() []tx.Tx {
	txs := make([]tx.Tx, 0)
	for _, t := range b.transactions {
		txs = append(txs, t)
	}
	return txs
}

func (block *Block) Timestamp() error {
	block.timestamp = protocol.Lepoch()
	return nil
}

func (block *Block) GetTimestamp() int64 {
	return block.timestamp
}

func (b *Block) verify_txs() ([]tx.Tx, error) {
	// verify tx fields make sense
	if len(b.transactions) > MAX_TXS {
		return nil, errors.New("too many txs in the block")
	}
	if b.MinerTx == nil {
		return nil, errors.New("no miner tx")
	}
	if _, ok := b.transactions[b.MinerTx.Id()]; !ok {
		return nil, errors.New("miner tx not in tx map")
	}

	// verify all transactions and make sure there's only one minertx
	invalidtxs := make([]tx.Tx, 0)
	var err error
	for _, t := range b.transactions {
		switch v := t.(type) {
		case *tx.MinerTx:
			if v.Id() != b.MinerTx.Id() {
				err = errors.New("more than one minertx in block")
				invalidtxs = append(invalidtxs, t)
			}
		case *tx.TwoWayTx:
		default:
			panic(fmt.Sprintf("unexpected tx.Tx: %#v", v))
		}
		if e := t.Verify(); err != nil {
			err = e
			invalidtxs = append(invalidtxs, t)
		}
	}
	if err != nil {
		return invalidtxs, err
	}
	return nil, nil
}

// verifies everything about the block, including the nonce
func (b *Block) Verify() error {
	// do not repeat previous verifications
	if b.valid {
		return nil
	}

	// you must have won the challenge
	if !b.TryNonce(b.Nonce) {
		return fmt.Errorf("could not verify nonce %d", b.Nonce)
	}
	if _, err := b.verify_txs(); err != nil {
		return err
	}

	b.valid = true
	return nil
}

func (b *Block) TryNonce(nonce uint32) bool {
	hash := b.HashWithNonce(nonce)
	shift := 64 - b.target
	bits := binary.BigEndian.Uint64(hash[:8]) >> shift
	if bits&(1<<64-1) != 0 {
		return false
	}
	b.BlockHash = hash
	// fmt.Printf("win with nonce %d\n", nonce)
	return true
}
