package main

import (
	"github.com/aaasen/crawl"
	"log"
)

func main() {
	stop := make(chan bool)

	links_in := make(chan string, 1)
	links_in <- "https://news.ycombinator.com/"

	links_out := make(chan string, 1024)

	pages := make(chan crawl.Page, 1024)

	wantMore := make(chan bool)

	taskQueue := crawl.NewDefaultTaskQueue(links_out, links_in, wantMore)
	go taskQueue.Run()

	for i := 0; i < 2; i++ {
		go crawl.Crawl(links_in, links_out, wantMore, pages, stop)
	}

	numPages := 0

	for {
		select {
		case page := <-pages:
			numPages++
			log.Printf("received %v: %v", numPages, page.URL)
		}
	}
}
