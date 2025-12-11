package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	// Check amount of arguments
	if len(cmd.args) != 1 {
		return errors.New("invalid amount of arguments for login")
	}

	// Set config user to username in args
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.args[0])
	return nil
}
