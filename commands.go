package main

import (
	"fmt"

	"github.com/nk-reddy/blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	command string
	args    []string
}

type commands struct {
	comms map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	comm, exists := c.comms[cmd.command]
	if !exists {
		return fmt.Errorf("command does not exist")
	}
	err := comm(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.comms[name] = f
}
