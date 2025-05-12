package protocol

// signatures, transactions, blocks, blockchains —  they are all verifiable
type Verifiable interface {
	Verify() error
}
