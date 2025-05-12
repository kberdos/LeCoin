package protocol

import "crypto/sha256"

/*
Different entities are hashable, such as blocks
*/
type HashType [sha256.Size]byte

type Hashable interface {
	Hash() HashType
}
