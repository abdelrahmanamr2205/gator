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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Incorrect number of arguments\nUsage: gator %s <feed_name> <feed_url>", cmd.name)
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

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Println("Feed created and followed successfully")

	fmt.Println(feed)
	fmt.Println(feedFollow)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Too many arguments\nUsage: gator %s", cmd.name)
	}

	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds %w", err)
	}

	for _, feed := range feeds {
		fmt.Println("feed name:", feed.FeedName)
		fmt.Println("feed url:", feed.FeedUrl)
		fmt.Println("owning user:", feed.UserName)
		fmt.Println()
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Incorrect number of arguments\nUsage: gator %s <feed_url>", cmd.name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	t := time.Now()
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: t,
		UpdatedAt: t,
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Printf("Feed %s followed successfully\n", feed.Name)
	fmt.Println(feedFollow)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("Too many arguments\nUsage: gator %s", cmd.name)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching follows: %w", err)
	}

	fmt.Println("You are following:")
	for i, follow := range feedFollows {
		fmt.Printf("%d. %s\n", i+1, follow.FeedName)
	}

	return nil
}
