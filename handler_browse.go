package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AbdelrahmanAmr2205/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var err error
	if len(cmd.args) > 1 {
		return fmt.Errorf("too many arguments\nUsage: gator browse [number of posts]")
	}
	limit := 2
	if len(cmd.args) == 1 {
		limit, err = strconv.Atoi(cmd.args[0])
	}
	if err != nil || limit < 1 {
		return fmt.Errorf("invalid argument. browse command's argument must be a positive integer")
	}

	posts, err := s.db.GatePostsForUser(context.Background(), database.GatePostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't fetch posts: %w", err)
	}

	for _, post := range posts {
		printPost(post)
		fmt.Println()
		fmt.Println()
	}

	return nil
}

func printPost(post database.Post) {
	fmt.Println("Title:", post.Title)
	if post.PublishedAt.Valid {
		fmt.Println("Published At:", post.PublishedAt.Time.Format("Jan 02, 2006"))
	}
	fmt.Println("Link:", post.Url)
	fmt.Println("Description:", post.Description)
}
