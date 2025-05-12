package wallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"

	"lecoin/pkg/protocol"
)

type Wallet struct {
	pubkey *ecdsa.PublicKey
}

func NewWallet(pubkey *ecdsa.PublicKey) *Wallet {
	return &Wallet{
		pubkey,
	}
}

func FromString(s string) (*Wallet, error) {
	pkbytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	pkint, err := x509.ParsePKIXPublicKey(pkbytes)
	if err != nil {
		return nil, err
	}

	pubkey, ok := pkint.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	return NewWallet(pubkey), nil
}

func (w *Wallet) GetPubKey() *ecdsa.PublicKey {
	return w.pubkey
}

func PrettyFromString(s string) string {
	w, err := FromString(s)
	if err != nil {
		return "INVALID"
	}
	return w.PrettyString()
}

func (w *Wallet) PrettyString() string {
	// print a prettier and shorter version
	h := w.Hash()
	hb := h[:]
	return hex.EncodeToString(hb)[:PRETTY_STR_LEN]
}

func (w *Wallet) String() string {
	pkBytes, err := x509.MarshalPKIXPublicKey(w.pubkey)
	if err != nil {
		return "INVALID"
	}
	return hex.EncodeToString(pkBytes)
}

func (w *Wallet) Hash() protocol.HashType {
	b := w.Marshal()
	return sha256.Sum256(b)
}

func (w *Wallet) Marshal() []byte {
	buf := make([]byte, 0)
	pkBytes, err := x509.MarshalPKIXPublicKey(w.pubkey)
	if err != nil {
		panic(err)
	}
	buf = append(buf, pkBytes...)
	if len(buf) != WALLET_SZ_BYTES {
		panic("the wallet should marshal to wallet sz bytes")
	}
	return buf
}

func (w *Wallet) Unmarshal(b []byte) error {
	pubkey, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return err
	}
	w.pubkey = pubkey.(*ecdsa.PublicKey)
	return nil
}
