package main

import (
	"fmt"
	"github.com/aaasen/crawl"
)

func main() {
	pages := make(chan crawl.Page)

	go crawl.Fetch("https://news.ycombinator.com/", pages)
	go crawl.Fetch("http://www.google.com/robots.txt", pages)
	go crawl.Fetch("http://www.google.com/", pages)

	for page := range pages {
		fmt.Println(string(page.URL.String()))
	}
}
