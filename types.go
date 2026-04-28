package main

import (
	"fmt"

	"github.com/AbdelrahmanAmr2205/gator/internal/config"
	"github.com/AbdelrahmanAmr2205/gator/internal/database"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

type command struct {
	name string
	args []string
}

type HandlerFunc func(*state, command) error

type commands struct {
	handlers map[string]HandlerFunc
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Unkown Command")
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f HandlerFunc) {
	c.handlers[name] = f
}
