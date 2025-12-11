package main

import (
	"errors"
	"log"
	"os"

	"github.com/Brandon-Butterbaugh/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {

	cfg := config.Read()
	s := state{
		cfg: &cfg,
	}

	cmds := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

	osArgs := os.Args
	if len(osArgs) < 2 {
		log.Fatalf("Not enough arguments")
	}
	cmd := command{
		name: osArgs[1],
		args: osArgs[2:],
	}

	err := cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}

type command struct {
	name string
	args []string
}

type commands struct {
	cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	// Check if command exists
	function, ok := c.cmds[cmd.name]

	// If it does run, else it's not a registered command
	if ok {
		err := function(s, cmd)
		if err != nil {
			return err
		}
		return err
	} else {
		return errors.New("unknown command")
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}
