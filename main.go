package main

import (
	"flag"

	"bsipiczki.com/rss-feed/table"
	"github.com/mmcdole/gofeed"
)

const (
	newInfix  = "new/"
	rssSuffix = ".rss"
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
	feedInput := flag.String("f", "https://www.reddit.com/r/hungary/", "add the feed you want to see")
	flag.Parse()

	var rssURL string
	if *sortNew {
		rssURL = *feedInput + newInfix + rssSuffix
	} else {
		rssURL = *feedInput + rssSuffix
	}

	feed, _ := fetchFeed(rssURL)
	printFeed(feed.Items)
}
