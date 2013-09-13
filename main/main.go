package main

import (
	"github.com/aaasen/crawl"
)

func main() {
	stop := make(chan bool)

	links_in := make(chan string, 1)
	links_in <- "https://news.ycombinator.com/"

	links_out := make(chan string, 100000)

	pages := make(chan crawl.Page, 1024)

	wantMore := make(chan bool)

	taskQueue := crawl.NewDefaultRedisTaskQueue(links_out, links_in, wantMore)
	go taskQueue.Run()

	dataStore := crawl.NewDebugDataStore(pages)

	for i := 0; i < 100; i++ {
		go crawl.Crawl(links_in, links_out, wantMore, pages, stop)
	}

	dataStore.Run()
}
