package main

import (
	"fmt"
	"os"
	"strconv"

	"lecoin/pkg/vswitch"
)

func main() {
	var port int
	if len(os.Args) == 2 {
		port, err := strconv.Atoi(os.Args[1])
		if err != nil || port < 0 || port >= 1<<16 {
			panic("bad port")
		}
	} else if len(os.Args) == 1 {
		port = vswitch.DEFAULT_PORT
	} else {
		fmt.Println("usage: vswitch <opt: port>")
		os.Exit(1)
	}

	vs := vswitch.NewVSwitch(port)
	if vs == nil {
		panic("vswitch could not be created")
	}

	fmt.Printf("Running VSwitch on port %d\n", port)
	vs.Run()
}
