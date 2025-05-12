package lesender

import (
	"errors"

	"lecoin/pkg/tx"
	"lecoin/pkg/wallet"
)

/*
Send to someone else based on wallets map index
*/
func (ls *LeSender) LeSendUser(idx int, amt tx.TxAmount) error {
	ls.rwlock.RLock()
	defer ls.rwlock.RUnlock()
	wStr, ok := ls.wallets[idx]
	if !ok {
		return errors.New("could not findwallet")
	}
	return ls.LeSendWallet(wStr, amt)
}

/*
Send to someone else based on wallet string
*/
func (ls *LeSender) LeSendWallet(wStr string, amt tx.TxAmount) error {
	remoteW, err := wallet.FromString(wStr)
	if err != nil {
		return err
	}
	localW := ls.lekey.GetWallet()
	t := tx.NewTwoWayTx(localW, remoteW, amt)
	if err := ls.lekey.SignTransaction(t); err != nil {
		return err
	}
	ls.Upchan <- t
	return nil
}

func (ls *LeSender) LeLoadGenerator() error {
	return nil
}

func (ls *LeSender) StopLeLoadGenerator() error {
	return nil
}
