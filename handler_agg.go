package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nk-reddy/blog-aggregator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("need a refresh duration")
	}

	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", time_between_reqs)
	ticker := time.NewTicker(time_between_reqs)
	for {
		scrapeFeeds(s)
		<-ticker.C
	}
}

func handlerUserPosts(s *state, cmd command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid limit %q: %w", cmd.args[0], err)
		}

		if parsedLimit <= 0 {
			return fmt.Errorf("limit must be greater than zero")
		}

		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\nURL: %s\n Published: %v\n\n", post.Title, post.Url, post.PublishedAt)
	}
	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, item := range fetchedFeed.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}
		published := parsePubDate(item.PubDate)

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: published,
			FeedID:      feed.ID,
		})

		if err == nil {
			continue
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
			continue
		} else {
			log.Printf("could not save post %q: %v", item.Title, err)
		}
	}
	return nil
}

func parsePubDate(pub_date string) sql.NullTime {
	if pub_date == "" {
		return sql.NullTime{Valid: false}
	}

	parsed_pub_date, err := mail.ParseDate(pub_date)
	if err != nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: parsed_pub_date, Valid: true}
}
