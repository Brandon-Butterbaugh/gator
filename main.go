package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/Brandon-Butterbaugh/gator/internal/config"
	"github.com/Brandon-Butterbaugh/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		handler(s, cmd, user)
		return nil
	}
}

func main() {

	cfg := config.Read()

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("failed to open database")
	}
	defer db.Close()
	dbQueries := database.New(db)

	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	osArgs := os.Args
	if len(osArgs) < 2 {
		log.Fatalf("Not enough arguments")
	}
	cmd := command{
		Name: osArgs[1],
		Args: osArgs[2:],
	}

	err = cmds.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
