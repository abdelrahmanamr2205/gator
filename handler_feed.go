package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AbdelrahmanAmr2205/gator/internal/database"
	rssfeed "github.com/AbdelrahmanAmr2205/gator/internal/rss_feed"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Too many arguments\nUsage: gator %s", cmd.name)
	}

	feed, err := rssfeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Incorrect number of arguments\nUsage: gator %s <feed_name> <feed_url>", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}

	t := time.Now()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	}

	fmt.Println("Feed created successfully")

	fmt.Println(feed)

	return nil
}
