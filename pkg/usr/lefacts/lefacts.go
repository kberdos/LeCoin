package lefacts

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"lecoin/pkg/repl"
)

type LeBronFacts struct {
	facts []string
}

func NewLebronFacts() *LeBronFacts {
	facts := make([]string, 0)
	facts = append(facts, "LeBron is the GOAT")
	facts = append(facts, "LeBron > MJ")
	facts = append(facts, "It wasn't LeBron's fault the Lakers lost in the 2025 playoffs")
	return &LeBronFacts{
		facts: facts,
	}
}

func (lf *LeBronFacts) Run() {
}

func (lf *LeBronFacts) MakeREPL() *repl.Repl {
	handlers := map[string]repl.ReplHandler{
		"lefact": HandleLeFact,
	}

	return repl.NewRepl(handlers, lf, repl.LF_REPL_NUM)
}

func HandleLeFact(prog repl.Prog, args []string) error {
	if len(args) != 0 {
		return errors.New("this takes no args... you can't argue with legoat")
	}
	lf := prog.(*LeBronFacts)
	idx := rand.IntN(len(lf.facts))
	fmt.Println(lf.facts[idx])
	return nil
}
