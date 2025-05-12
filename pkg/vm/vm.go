package vm

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"lecoin/pkg/repl"
	"lecoin/pkg/vm/vfs"
	"lecoin/pkg/vm/vsocket"
	"lecoin/pkg/vswitch"
)

// welcome to LeOS

type VM struct {
	hostname  string
	vfs       *vfs.VFS
	vskt      *vsocket.VSocket
	userprogs []UserProg
	shell     *repl.Repl
}

type VMCfg struct {
	Name string
	Port int
}

func NewVM(cfg VMCfg) *VM {
	var vm VM
	vm.hostname = cfg.Name

	// --- mount vfs ---
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	root := filepath.Join(cwd, fmt.Sprintf("hostdirs/%s", cfg.Name))
	err = os.MkdirAll(root, fs.ModePerm)
	if err != nil {
		panic(err)
	}
	vm.vfs = vfs.NewVFS(root)

	// --- create vsocket ---
	vm.vskt, err = vsocket.NewVSocket(uint16(cfg.Port), vswitch.DEFAULT_PORT)
	if err != nil {
		panic(err)
	}

	// --- create shell ---
	vm.shell = vm.MakeREPL()
	vm.userprogs = make([]UserProg, 0)

	return &vm
}

func (vm *VM) RegisterProg(prog UserProg) {
	vm.userprogs = append(vm.userprogs, prog)
	vm.shell.CombineRepl(prog.MakeREPL())
}

func (vm *VM) Boot() {
	fmt.Println("Booting LeOS...")
	// --- start socket ---
	go vm.vskt.Run()

	// --- run user progs ---
	for _, userprog := range vm.userprogs {
		go userprog.Run()
	}

	// --- run shell ---
	vm.shell.Run()
}
