package lechain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"sort"

	"lecoin/pkg/protocol"
	"lecoin/pkg/tx"
)

/* operations for marshaling and hashing blocks */

func (b *Block) MsgType() protocol.MessageType {
	return protocol.BLOCK_MSG_TYPE
}

func (b *Block) HashWithNonce(nonce uint32) protocol.HashType {
	b_bytes := b.Marshal()
	binary.BigEndian.PutUint32(b_bytes[:], nonce)
	return sha256.Sum256(b_bytes)
}

func (b *Block) Hash() protocol.HashType {
	// serialize
	b_bytes := b.Marshal()
	// sha256 it lol
	return sha256.Sum256(b_bytes)
}

func (b *Block) Marshal() []byte {
	buf := make([]byte, 0)
	buf = binary.BigEndian.AppendUint32(buf, b.Nonce)
	buf = append(buf, b.PrevHash[:]...)
	buf = binary.BigEndian.AppendUint64(buf, uint64(b.timestamp))
	buf = binary.BigEndian.AppendUint32(buf, b.target)
	buf = append(buf, b.MinerTx.FullMarshal()...)

	// ensure consistent order of transactions in map
	keys := make([]tx.Txid, 0, len(b.transactions))
	for k := range b.transactions {
		if k != b.MinerTx.Id() {
			keys = append(keys, k)
		}
	}
	if len(keys) != len(b.transactions)-1 {
		panic("incorrect # of txs marshaled")
	}
	sort.Slice(keys, func(i, j int) bool { return bytes.Compare(keys[i][:], keys[j][:]) < 0 })

	// include the number of txs
	buf = binary.BigEndian.AppendUint16(buf, uint16(len(keys)))
	for _, k := range keys {
		t := b.transactions[k]
		buf = append(buf, t.FullMarshal()...)
	}

	return buf
}

func (b *Block) FullMarshal() []byte {
	buf := b.Marshal()
	return buf
}

func (b *Block) Unmarshal(data []byte) error {
	// TODO: (extension) data length check
	pos := 0
	b.Nonce = binary.BigEndian.Uint32(data[pos : pos+4])
	pos += 4
	copy(b.PrevHash[:], data[pos:pos+32])
	pos += 32
	b.timestamp = int64(binary.BigEndian.Uint64(data[pos : pos+8]))
	pos += 8
	b.target = binary.BigEndian.Uint32(data[pos : pos+4])
	pos += 4
	// minertx
	minertx := tx.BlankMinerTx()
	minertx.Unmarshal(data[pos:])
	b.MinerTx = minertx
	b.transactions[minertx.Id()] = minertx
	pos += tx.MINER_TX_BASE_SZ + int(minertx.Siglen)
	// now, all the other transactions
	numtxs := binary.BigEndian.Uint16(data[pos : pos+2])
	pos += 2
	for range numtxs {
		t := tx.BlankTwoWayTx()
		t.Unmarshal(data[pos:])
		b.transactions[t.Id()] = t
		// FIXME: make this a method of serializable
		pos += tx.TWO_WAY_TX_BASE_SZ + int(t.Siglen)
	}
	// lastly, set the block hash
	b.BlockHash = b.Hash()
	return nil
}
