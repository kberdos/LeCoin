package tx

import (
	"lecoin/pkg/protocol"
	"lecoin/pkg/wallet"
)

const (
	MINER_TX_BASE_SZ   = 139
	TWO_WAY_TX_BASE_SZ = 230
)

type (
	Txid     protocol.HashType
	TxMap    map[Txid]Tx
	TxAmount float32
)

// individual transaction
type Tx interface {
	protocol.Hashable
	protocol.Serializable
	protocol.Verifiable // verify the transaction is signed properly
	protocol.Stampable  // time stamp transactions
	// miner wallet
	// break into TxRecords
	IntoRecords() []TxRecord
	Sign([]byte) error // add a signature to this transaction (errors if signature already exists)
	Id() Txid
}

// wallet X gets amount Y (can be negative)
type TxRecord struct {
	Wallet *wallet.Wallet
	Amount TxAmount
}
