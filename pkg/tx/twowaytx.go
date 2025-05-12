package tx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"lecoin/pkg/protocol"
	"lecoin/pkg/wallet"
)

// implements Tx interface
type TwoWayTx struct {
	txid      Txid
	sender    *wallet.Wallet
	receiver  *wallet.Wallet
	amount    TxAmount
	timestamp int64
	feeamount TxAmount // ignore
	Siglen    uint32   // need this because signature is variable length
	signature []byte   // signature on the rest of the transaction
}

// TODO: (extension) add in the fee
func NewTwoWayTx(sender, receiver *wallet.Wallet, amount TxAmount) *TwoWayTx {
	twt := &TwoWayTx{
		sender:    sender,
		receiver:  receiver,
		amount:    amount,
		timestamp: protocol.NO_TIMESTAMP,
	}
	if err := twt.Timestamp(); err != nil {
		panic("this is surprising lol")
	}
	twt.txid = Txid(twt.Hash())
	return twt
}

func (t *TwoWayTx) String() string {
	return fmt.Sprintf(
		"TwoWayTx{txid: %x, sender: %s, receiver: %s, amount: %f, timestamp: %d, signature: %x}",
		t.txid,
		t.sender.String(),
		t.receiver.String(),
		t.amount,
		t.timestamp,
		t.signature,
	)
}

func (tx *TwoWayTx) MsgType() protocol.MessageType {
	return protocol.TX_MSG_TYPE
}

func (tx *TwoWayTx) Id() Txid {
	return tx.txid
}

func (tx *TwoWayTx) IntoRecords() []TxRecord {
	r := make([]TxRecord, 0)
	// sender
	r = append(r, TxRecord{
		Wallet: tx.sender,
		Amount: -tx.amount,
	})
	// receiver
	r = append(r, TxRecord{
		Wallet: tx.receiver,
		Amount: tx.amount,
	})

	// (optional) transaction fee to the miner
	// if tx.miner != nil {
	// 	r = append(r, TxRecord{
	// 		Wallet: tx.miner,
	// 		Amount: tx.feeamount,
	// 	})
	// }

	// TODO: add in the miner's fee if you desire (might need to be done at the block level)
	return r
}

func (tx *TwoWayTx) Hash() protocol.HashType {
	tx_bytes := tx.Marshal()
	return sha256.Sum256(tx_bytes)
}

func (tx *TwoWayTx) Marshal() []byte {
	buf := new(bytes.Buffer)
	if _, err := buf.Write(tx.sender.Marshal()); err != nil {
		panic(err)
	}
	if _, err := buf.Write(tx.receiver.Marshal()); err != nil {
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

// marshal + txid and signature
func (tx *TwoWayTx) FullMarshal() []byte {
	buf := tx.Marshal()
	buf = append(buf, tx.txid[:]...)
	buf = binary.BigEndian.AppendUint32(buf, tx.Siglen)
	buf = append(buf, tx.signature[:]...)
	return buf
}

func BlankTwoWayTx() *TwoWayTx {
	return &TwoWayTx{
		sender:    &wallet.Wallet{},
		receiver:  &wallet.Wallet{},
		signature: make([]byte, 0),
	}
}

// size: 2 * 91 + 4 + 8 + 32 + 4 + signature length
func (t *TwoWayTx) Unmarshal(b []byte) error {
	pos := 0
	if err := t.sender.Unmarshal(b[pos : pos+wallet.WALLET_SZ_BYTES]); err != nil {
		return err
	}
	pos += wallet.WALLET_SZ_BYTES
	if err := t.receiver.Unmarshal(b[pos : pos+wallet.WALLET_SZ_BYTES]); err != nil {
		return err
	}
	pos += wallet.WALLET_SZ_BYTES
	reader := bytes.NewReader(b[pos : pos+4])
	if err := binary.Read(reader, binary.BigEndian, &t.amount); err != nil {
		return err
	}
	pos += 4
	// timestamp
	t.timestamp = int64(binary.BigEndian.Uint64(b[pos : pos+8]))
	pos += 8
	copy(t.txid[:], b[pos:pos+sha256.Size])
	pos += sha256.Size
	t.Siglen = binary.BigEndian.Uint32(b[pos : pos+4])
	pos += 4
	t.signature = append(t.signature, b[pos:pos+int(t.Siglen)]...)
	return nil
}

func (tx *TwoWayTx) Sign(signature []byte) error {
	if len(tx.signature) != 0 {
		return errors.New("transaction already has a signature")
	}
	tx.signature = signature
	tx.Siglen = uint32(len(signature))
	return nil
}

func (tx *TwoWayTx) Verify() error {
	nonkeybytes := tx.Marshal()
	if ecdsa.VerifyASN1(tx.sender.GetPubKey(), nonkeybytes, tx.signature) {
		return nil
	}
	return errors.New("could not verify two way transaction")
}

func (tx *TwoWayTx) Timestamp() error {
	if tx.timestamp != protocol.NO_TIMESTAMP {
		return errors.New("two way tx is already timestsamped")
	}
	tx.timestamp = protocol.Lepoch()
	return nil
}

func (tx *TwoWayTx) GetTimestamp() int64 {
	return tx.timestamp
}
