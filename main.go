package main

import (
	"flag"

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
	sortNew := flag.Bool("n", false, "sort by new")
	flag.Parse()

	var rssURL string
	if *sortNew {
		rssURL = "https://www.reddit.com/r/hungary/new/.rss"
	} else {
		rssURL = "https://www.reddit.com/r/hungary/.rss"
	}

	feed, _ := fetchFeed(rssURL)
	printFeed(feed.Items)
}
