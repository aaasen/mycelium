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

	go crawl.TaskQueue(links_out, links_in)

	for i := 0; i < 2; i++ {
		go crawl.Crawl(links_in, links_out, pages, stop)
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
