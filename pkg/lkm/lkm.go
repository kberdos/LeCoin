package lkm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"

	"lecoin/pkg/vm"
	"lecoin/pkg/wallet"
)

/*
LeKeyManager handles cryptographic signing and verification —  essentially wraps public/private key architecture
*/

type LeKeyManager struct {
	vm      *vm.VM
	privkey *ecdsa.PrivateKey
	pubkey  *ecdsa.PublicKey
	w       *wallet.Wallet
}

func NewLKM(vm *vm.VM, filepath string) (*LeKeyManager, error) {
	var privkey *ecdsa.PrivateKey
	var err error
	if !vm.Exists(filepath) {
		// generate key
		privkey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return &LeKeyManager{}, err
		}

		privkeybytes, err := x509.MarshalECPrivateKey(privkey)
		if err != nil {
			return &LeKeyManager{}, err
		}
		vm.WriteFile(filepath, privkeybytes)
		fmt.Println(len(privkeybytes))

	} else {
		// load key
		privkeybytes, err := vm.ReadFile(filepath)
		if err != nil {
			return &LeKeyManager{}, err
		}
		fmt.Println(len(privkeybytes))

		privkey, err = x509.ParseECPrivateKey(privkeybytes)
		if err != nil {
			return &LeKeyManager{}, err
		}
	}
	pubkey := privkey.PublicKey
	return &LeKeyManager{
		vm:      vm,
		privkey: privkey,
		pubkey:  &pubkey,
		w:       wallet.NewWallet(&pubkey),
	}, nil
}

func (lkm *LeKeyManager) GetWallet() *wallet.Wallet {
	return lkm.w
}
