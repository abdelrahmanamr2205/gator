package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AbdelrahmanAmr2205/gator/internal/database"
	rssfeed "github.com/AbdelrahmanAmr2205/gator/internal/rss_feed"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Incorrect number of arguments\nUsage: gator %s <duration_between_requests>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("error parsing duration: %w", err)
	}

	for range time.Tick(timeBetweenRequests) {
		scrapeFeeds(s)
	}

	return nil
}

func scrapeFeeds(s *state) {
	nextfeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}

	feed, err := rssfeed.FetchFeed(context.Background(), nextfeed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", nextfeed.Name, err)
		return
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:        nextfeed.ID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", nextfeed.Name, err)
		return
	}
	for _, item := range feed.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", nextfeed.Name, len(feed.Channel.Item))
}

func printFeed(feed rssfeed.RSSFeed) {
	fmt.Printf("Feed: %s\n", feed.Channel.Title)
	fmt.Printf("Description: %s", feed.Channel.Description)
	fmt.Println()
	for _, item := range feed.Channel.Item {
		fmt.Printf("Article Title: %s\n", item.Title)
		fmt.Printf("Description:\n%s\n", item.Description)
		fmt.Println()
	}
}
