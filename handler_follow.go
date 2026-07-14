package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nk-reddy/blog-aggregator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("need a URL to follow")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("record created - feed name: %s; user: %s\n", feed_follow.FeedName, feed_follow.UserName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("need a URL to unfollow")
	}

	feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("successfully unfollowed %s", feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(follows) == 0 {
		fmt.Println("not followng any feeds")
		return nil
	}

	fmt.Println("Following:")
	for _, follow := range follows {
		fmt.Printf(" * %s\n", follow.FeedName)
	}
	return nil
}
