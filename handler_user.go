package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Brandon-Butterbaugh/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	// Check amount of arguments
	if len(cmd.Args) != 1 {
		return errors.New("invalid amount of arguments for login")
	}

	// Check if name exists in database
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatalf("User not found")
		return err
	}

	// Set config user to username in args
	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	// Check amount of arguments
	if len(cmd.Args) != 1 {
		return errors.New("invalid amount of arguments for register")
	}

	// Check if name already exists in database
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err == nil {
		log.Fatalf("name already exists")
		return err
	}

	// Create user in database
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      cmd.Args[0],
		},
	)
	if err != nil {
		return err
	}

	// Set current user
	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s created\n", cmd.Args[0])
	printUser(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("failed to reset database")
		return err
	}

	return err
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("failed to get users")
		return err
	}

	var current string
	for _, user := range users {
		current = ""
		if user.Name == s.cfg.CurrentUserName {
			current = " (current)"
		}
		fmt.Printf(" * Name:    %v%s\n", user.Name, current)
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
