package main

import (
	"encoding/json"
	"fmt"
	"os"

	"lecoin/pkg/lcm"
	"lecoin/pkg/usr/lefacts"
	"lecoin/pkg/vm"
)

type HostCfg struct {
	Name string   `json:"name"`
	Port int      `json:"port"`
	Args []string `json:"args"`
}

func ParseCfg(cfgpath string) (*HostCfg, error) {
	f, err := os.Open(cfgpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var hostcfg HostCfg
	dec := json.NewDecoder(f)
	if err = dec.Decode(&hostcfg); err != nil {
		return nil, err
	}

	return &hostcfg, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: client <path to cfg>")
		os.Exit(1)
	}

	cfg, err := ParseCfg(os.Args[1])
	if err != nil {
		panic(err)
	}

	vmcfg := vm.VMCfg{
		Name: cfg.Name,
		Port: cfg.Port,
	}

	vm := vm.NewVM(vmcfg)

	isminer := cfg.Args[0] == "miner"
	mgr, err := lcm.NewLCM(vm, isminer)
	if err != nil {
		panic("could not initialize lecoin manager")
	}
	vm.RegisterProg(mgr)

	lf := lefacts.NewLebronFacts()
	vm.RegisterProg(lf)

	vm.Boot()
}
