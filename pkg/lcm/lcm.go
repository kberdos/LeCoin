package lcm

import (
	"fmt"

	"lecoin/pkg/lechain"
	"lecoin/pkg/leminer"
	"lecoin/pkg/lesender"
	"lecoin/pkg/lkm"
	"lecoin/pkg/protocol"
	"lecoin/pkg/tx"
	"lecoin/pkg/vm"
	"lecoin/pkg/vm/vsocket"
)

type LeCoinManager struct {
	// vm for file managing  NOTE: could potentially just keep in lkm
	vm    *vm.VM
	chain *lechain.LeChain
	miner *leminer.LeMiner
	lekey *lkm.LeKeyManager
	ls    *lesender.LeSender
	p2p   *P2P
	// p2p manager (todo)
}

type P2P struct {
	txchan    vsocket.NetChan
	blockchan vsocket.NetChan
}

// TODO: (extension) add keypath + other saved state loading
func NewLCM(vm *vm.VM, miner bool) (vm.UserProg, error) {
	lcm := &LeCoinManager{
		vm:    vm,
		chain: lechain.NewLechain(),
		p2p:   &P2P{},
	}
	vm.Mkdir(".ssh")
	lekey, err := lkm.NewLKM(vm, ".ssh/ecdkey")
	if err != nil {
		return nil, err
	}
	lcm.lekey = lekey
	if miner {
		miner := leminer.NewLeMiner(lcm.chain, lekey)
		lcm.miner = miner
	}
	lcm.ls = lesender.NewLeSender(lcm.chain, lekey)

	// register p2p transaction and block listening
	lcm.p2p.txchan = vm.RegisterNetChan(protocol.TX_MSG_TYPE, func() protocol.Serializable { return tx.BlankTwoWayTx() })
	lcm.p2p.blockchan = vm.RegisterNetChan(protocol.BLOCK_MSG_TYPE, func() protocol.Serializable { return lechain.BlankBlock() })
	return lcm, nil
}

func (lcm *LeCoinManager) Run() {
	// run all of the subcomponents
	chainchan := lcm.chain.RegisterListener()
	minerchan := make(chan *lechain.Block)
	if lcm.miner != nil {
		minerchan = lcm.miner.Upchan
		go lcm.run_miner() // NOTE: we are automatically starting the miner in this case
	}

	for {
		select {
		case c := <-chainchan:
			fmt.Println("chain got a new block!")
			_ = c
		case b := <-minerchan:
			fmt.Println("miner mined a block!")
			lcm.vm.NetBroadcast(b)
		case t := <-lcm.ls.Upchan:
			fmt.Println("got transaction from myself!")
			// broadcast to everyone
			lcm.vm.NetBroadcast(t)
			if lcm.miner != nil {
				lcm.miner.InsertTx(t)
			}
		case s := <-lcm.p2p.txchan:
			fmt.Println("got transaction from someone else!")
			t := s.(tx.Tx)
			if lcm.miner != nil {
				if err := lcm.miner.InsertTx(t); err != nil {
					fmt.Println(err)
				}
			}
		case m := <-lcm.p2p.blockchan:
			fmt.Println("got block from someone else!")
			b := m.(*lechain.Block)
			if _, err := lcm.chain.PushBlock(b, nil); err != nil {
				fmt.Println(err)
			}
		}
	}
}

// TODO: move somewhere else and place in repl
func (lcm *LeCoinManager) run_miner() {
	lcm.miner.Run()
}

func (lcm *LeCoinManager) stop_miner() {
	lcm.miner.Abort <- true
}
