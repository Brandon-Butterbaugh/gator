package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Brandon-Butterbaugh/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	// Check amount of arguments
	if len(cmd.Args) != 1 {
		return errors.New("invalid amount of arguments for following")
	}

	// Get feed
	feed, err := s.db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatalf("error finding feed with url")
	}

	// Create feed follow
	follow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		log.Fatalf("error creating follow")
	}

	fmt.Println(follow.FeedName)
	fmt.Println(follow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	// Get follows
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatalf("error finding follows")
	}

	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	// Check amount of arguments
	if len(cmd.Args) != 1 {
		return errors.New("invalid amount of arguments for unfollowing")
	}

	// Get feed
	feed, err := s.db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatalf("error finding feed with url")
	}

	// Delete follow
	err = s.db.DeleteFollow(
		context.Background(),
		database.DeleteFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		log.Fatalf("error deleting follow")
	}

	return nil
}
