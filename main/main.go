package main

import (
	"github.com/aaasen/crawl"
)

func main() {
	pages := make(chan crawl.Page)

	crawl.Crawl([]string{"https://news.ycombinator.com/"}, pages)
	// go crawl.Fetch("http://www.google.com/robots.txt", pages)
	// go crawl.Fetch("http://www.google.com/", pages)

}
