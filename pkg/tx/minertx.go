package tx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"

	"lecoin/pkg/protocol"
	"lecoin/pkg/wallet"
)

type MinerTx struct {
	txid      Txid
	miner     *wallet.Wallet
	amount    TxAmount
	Siglen    uint32
	signature []byte
	timestamp int64 // probably optional
}

func NewMinerTx(miner *wallet.Wallet, amount TxAmount) *MinerTx {
	tx := &MinerTx{
		miner:     miner,
		amount:    amount,
		signature: make([]byte, 0),
		timestamp: protocol.NO_TIMESTAMP,
	}
	if err := tx.Timestamp(); err != nil {
		panic("this certainly should not happen ever")
	}
	tx.txid = Txid(tx.Hash())
	return tx
}

func BlankMinerTx() *MinerTx {
	return &MinerTx{
		miner:     &wallet.Wallet{},
		signature: make([]byte, 0),
	}
}

func (tx *MinerTx) MsgType() protocol.MessageType {
	return protocol.TX_MSG_TYPE
}

func (tx *MinerTx) Id() Txid {
	return tx.txid
}

func (tx *MinerTx) IntoRecords() []TxRecord {
	r := make([]TxRecord, 0)
	r = append(r, TxRecord{
		Wallet: tx.miner,
		Amount: tx.amount,
	})
	return r
}

// marshals without the key
func (tx *MinerTx) Marshal() []byte {
	buf := new(bytes.Buffer)
	if _, err := buf.Write(tx.miner.Marshal()); err != nil {
		panic(err)
	}
	if err := binary.Write(buf, binary.BigEndian, tx.amount); err != nil {
		panic(err)
	}
	if err := binary.Write(buf, binary.BigEndian, tx.timestamp); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (tx *MinerTx) Hash() protocol.HashType {
	tx_bytes := tx.Marshal()
	return sha256.Sum256(tx_bytes)
}

// marshal then add on the extra fields (txid and signature)
func (tx *MinerTx) FullMarshal() []byte {
	buf := tx.Marshal()
	buf = append(buf, tx.txid[:]...)
	buf = binary.BigEndian.AppendUint32(buf, tx.Siglen)
	buf = append(buf, tx.signature...)
	return buf
}

// BASE_SIZE: 91 + 4 + 8 + 32 + 4 + S = 139 + S
func (t *MinerTx) Unmarshal(b []byte) error {
	pos := 0
	if err := t.miner.Unmarshal(b[pos : pos+wallet.WALLET_SZ_BYTES]); err != nil {
		return err
	}
	pos += wallet.WALLET_SZ_BYTES
	reader := bytes.NewReader(b[pos : pos+4])
	if err := binary.Read(reader, binary.BigEndian, &t.amount); err != nil {
		return err
	}
	pos += 4
	t.timestamp = int64(binary.BigEndian.Uint64(b[pos : pos+8]))
	pos += 8
	copy(t.txid[:], b[pos:pos+sha256.Size])
	pos += sha256.Size
	t.Siglen = binary.BigEndian.Uint32(b[pos : pos+4])
	pos += 4
	t.signature = append(t.signature, b[pos:pos+int(t.Siglen)]...)
	return nil
}

func (tx *MinerTx) Sign(signature []byte) error {
	if len(tx.signature) != 0 {
		return errors.New("transaction already has a signature")
	}
	tx.signature = signature
	tx.Siglen = uint32(len(signature))
	return nil
}

func (tx *MinerTx) Verify() error {
	nonkeybytes := tx.Marshal()
	if ecdsa.VerifyASN1(tx.miner.GetPubKey(), nonkeybytes, tx.signature) {
		return nil
	}
	return errors.New("could not verify miner tx")
}

func (tx *MinerTx) Timestamp() error {
	if tx.timestamp != protocol.NO_TIMESTAMP {
		return errors.New("miner tx is already timestsamped")
	}
	tx.timestamp = protocol.Lepoch()
	return nil
}

func (tx *MinerTx) GetTimestamp() int64 {
	return tx.timestamp
}
