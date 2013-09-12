package main

import (
	"github.com/aaasen/crawl"
	"log"
)

func main() {
	stop := make(chan bool)
	links := make(chan string, 1000000)

	links <- "https://news.ycombinator.com/"

	pages := make(chan crawl.Page, 1024)

	for i := 0; i < 100; i++ {
		go crawl.Crawl(stop, links, pages)
	}

	for {
		select {
		case page := <-pages:
			log.Printf("received: %v", page.URL)
		}
	}
}
