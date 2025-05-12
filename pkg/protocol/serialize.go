package protocol

const (
	ASN1_PUB_KEY_LEN = 91
)

type Serializable interface {
	// serialization without 'excluded' fields (such as keys, block hashes, timestamps, etc.)
	// primarily used for hashing or verification
	Marshal() []byte
	// serialization of the entire struct (exluding INTERNAL fields) for sending on the wire
	FullMarshal() []byte
	// populate yourself from a serialization
	Unmarshal([]byte) error
	// spit out your message type
	MsgType() MessageType
}
