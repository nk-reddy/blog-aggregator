package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nk-reddy/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	myState := state{
		cfg: &cfg,
	}
	myCommands := commands{
		comms: map[string]func(*state, command) error{},
	}

	myCommands.register("login", handlerLogin)

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
