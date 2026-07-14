package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nk-reddy/blog-aggregator/internal/config"
	"github.com/nk-reddy/blog-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	myState := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	myCommands := commands{
		comms: map[string]func(*state, command) error{},
	}

	myCommands.register("login", handlerLogin)
	myCommands.register("register", handlerRegister)
	myCommands.register("reset", handlerReset)
	myCommands.register("users", handlerUsers)
	myCommands.register("agg", handlerAgg)
	myCommands.register("feeds", handlerFeeds)
	myCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	myCommands.register("follow", middlewareLoggedIn(handlerFollow))
	myCommands.register("following", middlewareLoggedIn(handlerFollowing))
	myCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	userArgs := os.Args
	if len(userArgs) < 2 {
		fmt.Println("not enough arguments passed in")
		os.Exit(1)
	}

	userCommand := command{
		command: userArgs[1],
		args:    userArgs[2:],
	}

	err = myCommands.run(&myState, userCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
