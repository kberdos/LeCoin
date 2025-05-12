package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Prog interface{} // synonymous with UserProg

type ReplHandler func(Prog, []string) error

// A command has both a handler and the protocol it is over
type ReplCommand struct {
	protocol int // ip, tcp, etc
	handler  ReplHandler
}

type ReplCommands map[string]ReplCommand

type ReplConfig map[int]Prog // maps the protocol / REPL nums to Progs

type Repl struct {
	config   ReplConfig
	commands ReplCommands
}

// Make a new REPL with the given handlers and progs / protocol (repl num).
func NewRepl(handlers map[string]ReplHandler, prog Prog, protocol int) *Repl {
	commands := make(ReplCommands)
	for trigger, handler := range handlers {
		commands[trigger] = ReplCommand{protocol, handler}
	}

	config := ReplConfig{
		protocol: prog,
	}
	return &Repl{config, commands}
}

// add new commands / config from newrepl
func (r *Repl) CombineRepl(newrepl *Repl) {
	for trigger, command := range newrepl.commands {
		_, ok := r.commands[trigger]
		if !ok {
			r.commands[trigger] = command
		}
	}
	for protocol, prog := range newrepl.config {
		_, ok := r.config[protocol]
		if !ok {
			r.config[protocol] = prog
		}
	}
}

// Run the REPL, in a goroutine if you wish
func (r *Repl) Run() {
	Scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")
	for Scanner.Scan() {
		t := strings.Trim(Scanner.Text(), " ")
		args := strings.Fields(t)
		if len(args) < 1 {
			fmt.Println("Invalid Command")
			fmt.Printf("> ")
			continue
		}

		if args[0] == "help" {
			r.HandleHelp()
			fmt.Printf("> ")
			continue
		}

		command, ok := r.commands[args[0]]
		if !ok {
			fmt.Println("Invalid Command")
			fmt.Printf("> ")
			continue
		}
		// pull off the right prog based on the command's protocol
		prog, ok := r.config[command.protocol]
		if !ok {
			fmt.Println("Invalid REPL Protocol Found") // really should not happen
		}
		err := command.handler(prog, args[1:])
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		fmt.Printf("> ")
	}
}

func (r *Repl) HandleHelp() {
	for cmd := range r.commands {
		fmt.Println(cmd)
	}
}
