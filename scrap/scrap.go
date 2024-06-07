package scrap

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Comment struct {
	Author  string
	Content string
}

func Scrap(url string) []Comment {
	comments := []Comment{}

	c := colly.NewCollector(
		colly.AllowedDomains("old.reddit.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	c.OnHTML(".entry", func(e *colly.HTMLElement) {
		author := e.ChildText(".author")
		content := e.ChildText(".md")

		if author != "" && content != "" {
			comment := Comment{
				Author:  author,
				Content: content,
			}
			comments = append(comments, comment)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// threadURL := "https://old.reddit.com/r/hungary/comments/1da3tei/visszal%C3%A9pett_szentkir%C3%A1lyi_alexandra_a/"
	c.Visit(url)

	c.Wait()

	return comments
}
