package main

import (
	"bsipiczki.com/rss-feed/table"
	"github.com/mmcdole/gofeed"
)

func fetchFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	return feed, err
}

func printFeed(items []*gofeed.Item) {
	table.Render(items)
}

func main() {
	rssURL := "https://www.reddit.com/r/hungary/.rss?sort=new"
	feed, _ := fetchFeed(rssURL)
	printFeed(feed.Items)
}
