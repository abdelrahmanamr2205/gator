package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AbdelrahmanAmr2205/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error: too few arguments\nUsage: gator %s <username>", cmd.name)
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("error: too many arguments\nUsage: gator %s <username>", cmd.name)
	}

	t := time.Now()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		Name:      cmd.args[0],
	}
	_, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("User created successfully")

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("error: too few arguments\nUsage: gator %s <username>", cmd.name)
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("error: too many arguments\nUsage: gator %s <username>", cmd.name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("User set successfully")

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("error: too many arguments\nUsage: gator %s", cmd.name)
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("error: too many arguments\nUsage: gator %s", cmd.name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("No users registered yet")
	}

	for _, user := range users {
		if user == s.conf.CurrentUserName {
			fmt.Println("*", user, "(current)")
		} else {
			fmt.Println("*", user)
		}
	}

	return nil
}
