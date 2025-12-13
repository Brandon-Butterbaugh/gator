package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Brandon-Butterbaugh/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func handlerAgg(s *state, cmd command) error {
	// Check amount of arguments
	if len(cmd.Args) != 1 {
		return errors.New("invalid amount of arguments for agg")
	}

	// Parse arg into time.Duration
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		log.Fatalf("error parsing time argument")
	}
	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	// Get next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("error getting next feed")
	}

	// Mark feed as fetched
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Fatalf("error marking feed as fetched")
	}

	// Fetch the feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Fatalf("error fetching feed")
	}

	// Add items to posts
	for _, item := range rssFeed.Channel.Item {
		// Make sure not null variables can be put into database
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = s.db.CreatePosts(
			context.Background(),
			database.CreatePostsParams{
				ID:        uuid.New(),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Title:     item.Title,
				Url:       item.Link,
				Description: sql.NullString{
					String: item.Description,
					Valid:  true,
				},
				PublishedAt: publishedAt,
				FeedID:      feed.ID,
			},
		)
		if err != nil {
			// ignore duplicate urls
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code != "23505" {
					fmt.Print(err)
					log.Printf("error creating post: %v", err)
					continue
				}
			}
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// Check amount of arguments
	if len(cmd.Args) != 2 {
		return errors.New("invalid amount of arguments for adding a feed")
	}

	// Create new feed
	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      cmd.Args[0],
			Url:       cmd.Args[1],
			UserID:    user.ID,
		},
	)
	if err != nil {
		return err
	}

	// Create feed follow
	_, err = s.db.CreateFeedFollow(
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

	fmt.Print(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	// Get feeds
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("error getting feeds from database")
	}

	// Print feeds and get user names
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			log.Fatalf("error getting user %v from database", feed.UserID)
		}
		fmt.Println(feed.ID)
		fmt.Println(feed.CreatedAt)
		fmt.Println(feed.UpdatedAt)
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(user.Name)

	}

	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Make client and blank feed
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	feed := &RSSFeed{}

	// Make request
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		log.Fatalf("error making request")
	}

	// Set header to identify program
	req.Header.Set("User-Agent", "gator")

	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error getting response")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status error: %v", resp.StatusCode)
	}

	// Read the response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read body error: %v", err)
	}

	// Unmarshal the response
	err = xml.Unmarshal(bodyBytes, feed)
	if err != nil {
		log.Fatalf("xml unmarhsal error: %v", err)
	}

	// Unescape the feed and the items
	unescapeFeed(feed)
	for i := range feed.Channel.Item {
		unescapeItem(&feed.Channel.Item[i])
	}

	return feed, nil
}

func unescapeFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
}

func unescapeItem(item *RSSItem) {
	item.Title = html.UnescapeString(item.Title)
	item.Description = html.UnescapeString(item.Description)
}
