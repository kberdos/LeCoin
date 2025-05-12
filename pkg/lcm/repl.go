package lcm

import (
	"errors"
	"fmt"
	"strconv"

	"lecoin/pkg/repl"
	"lecoin/pkg/tx"
)

func (lcm *LeCoinManager) MakeREPL() *repl.Repl {
	handlers := map[string]repl.ReplHandler{
		"balance": HandleBalance,
		"list":    HandleList,
		"users":   HandleUsers,
		"ss":      HandleSelfSend,
		"send":    HandleSend,
	}
	return repl.NewRepl(handlers, lcm, repl.LCM_REPL_NUM)
}

func HandleBalance(prog repl.Prog, args []string) error {
	if len(args) != 0 {
		return errors.New("balance takes no args")
	}
	lcm := prog.(*LeCoinManager)
	bal := lcm.ls.LeBalance()
	fmt.Printf("Wallet: %s\n", lcm.lekey.GetWallet().PrettyString())
	fmt.Printf("My LeBalance is: %.2f\n", bal)
	return nil
}

func HandleList(prog repl.Prog, args []string) error {
	if len(args) != 0 {
		return errors.New("list takes no args")
	}
	lcm := prog.(*LeCoinManager)
	fmt.Print(lcm.ls.LeListBalances())
	return nil
}

func HandleUsers(prog repl.Prog, args []string) error {
	if len(args) != 0 {
		return errors.New("users takes no args")
	}
	lcm := prog.(*LeCoinManager)
	fmt.Print(lcm.ls.LeListUsers())
	return nil
}

func HandleSelfSend(prog repl.Prog, args []string) error {
	if len(args) != 0 {
		return errors.New("balance takes no args")
	}
	lcm := prog.(*LeCoinManager)
	return lcm.ls.LeSelfSend()
}

func HandleSend(prog repl.Prog, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: send <wallet index> <amount>")
	}
	idx, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	amt, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return err
	}
	fmt.Println(amt)
	lcm := prog.(*LeCoinManager)
	return lcm.ls.LeSendUser(int(idx), tx.TxAmount(amt))
}
