package vm

import (
	"errors"
	"fmt"

	"lecoin/pkg/repl"
)

func (vm *VM) MakeREPL() *repl.Repl {
	handlers := map[string]repl.ReplHandler{
		"cd":    HandleCd,
		"pwd":   HandlePwd,
		"mkdir": HandleMkdir,
		"ls":    HandleListdir,
	}
	return repl.NewRepl(handlers, vm, repl.VM_REPL_NUM)
}

func HandleCd(stack repl.Prog, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: cd <tgtdir>")
	}
	vm := stack.(*VM)
	return vm.Cd(args[0])
}

func HandlePwd(stack repl.Prog, args []string) error {
	if len(args) != 0 {
		return errors.New("usage: pwd")
	}
	vm := stack.(*VM)
	fmt.Println(vm.GetCwd())
	return nil
}

func HandleMkdir(stack repl.Prog, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: mkdir <newdir>")
	}
	vm := stack.(*VM)
	return vm.Mkdir(args[0])
}

func HandleListdir(stack repl.Prog, args []string) error {
	var dir string
	if len(args) == 1 {
		dir = args[0]
	} else if len(args) == 0 {
		dir = "."
	} else {
		return errors.New("usage: ls <opt:dir>")
	}

	vm := stack.(*VM)
	entries, err := vm.Listdir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}
	return nil
}
