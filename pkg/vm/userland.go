package vm

import (
	"lecoin/pkg/repl"
)

type UserProg interface {
	MakeREPL() *repl.Repl
	Run()
}
