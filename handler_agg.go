package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/AbdelrahmanAmr2205/gator/internal/database"
	rssfeed "github.com/AbdelrahmanAmr2205/gator/internal/rss_feed"
	"github.com/google/uuid"
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
		t := time.Now()
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   t,
			UpdatedAt:   t,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parseDate(item.PubDate),
			FeedID:      nextfeed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				continue
			} else {
				log.Printf("Couldn't save post %s: %v", post.Title, err)
			}
		}
	}
}

func parseDate(publishDate string) sql.NullTime {
	formats := []string{time.RFC1123Z, time.RFC1123, time.RFC3339, time.RFC822Z, "2006-01-02 15:04:05"}
	for _, format := range formats {
		if t, err := time.Parse(format, publishDate); err == nil {
			// success!
			return sql.NullTime{Time: t.UTC(), Valid: true}
		}
	}
	return sql.NullTime{Time: time.Time{}, Valid: false}
}
