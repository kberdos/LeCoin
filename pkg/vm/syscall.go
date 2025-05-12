package vm

import (
	"os"

	"lecoin/pkg/protocol"
	"lecoin/pkg/vm/vsocket"
)

// syscalls... essentially just wrappers around component functions

// ----- VFS calls -----

func (vm *VM) WriteFile(path string, data []byte) error {
	return vm.vfs.WriteFile(path, data)
}

func (vm *VM) ReadFile(path string) ([]byte, error) {
	return vm.vfs.ReadFile(path)
}

func (vm *VM) Mkdir(path string) error {
	return vm.vfs.Mkdir(path)
}

func (vm *VM) Cd(path string) error {
	return vm.vfs.Cd(path)
}

func (vm *VM) GetCwd() string {
	return vm.vfs.GetCwd()
}

func (vm *VM) Exists(path string) bool {
	return vm.vfs.Exists(path)
}

func (vm *VM) Listdir(dir string) ([]os.DirEntry, error) {
	return vm.vfs.Listdir(dir)
}

// ----- VSocket calls -----

func (vm *VM) NetSend(msg protocol.Serializable, rport uint16) error {
	return vm.vskt.Send(msg, rport)
}

func (vm *VM) NetBroadcast(msg protocol.Serializable) error {
	return vm.vskt.Broadcast(msg)
}

func (vm *VM) RegisterNetChan(msgType protocol.MessageType, constructor protocol.MessageConstructor) vsocket.NetChan {
	return vm.vskt.RegisterChannel(msgType, constructor)
}
