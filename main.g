package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SlyMarbo/rss"
)

func printFeed(feed *rss.Feed) {
	// Print the feed title
	fmt.Printf("Feed Title: %s\n\n", feed.Title)

	// Iterate through the items and print details
	for _, item := range feed.Items {
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n", item.Summary)
		fmt.Printf("Publication Date: %s\n\n", item.Date)
	}
}

func main() {
	rssURL := "https://www.reddit.com/r/hungary/.rss?sort=new"

	// Fetch the RSS feed initially
	feed, err := rss.Fetch(rssURL)
	if err != nil {
		log.Fatalf("Failed to fetch the RSS feed: %v", err)
	}

	// Print the initial feed
	printFeed(feed)

	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Updating the RSS feed...")
		err := feed.Update()
		if err != nil {
			log.Printf("Failed to update the RSS feed: %v", err)
		} else {
			printFeed(feed)
		}
	}
}

