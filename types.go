package main

import (
	"context"
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) HandlerFunc {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error fetching user: %w", err)
		}

		return handler(s, cmd, user)
	}
}
