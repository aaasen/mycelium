package main

import (
	"github.com/aaasen/mycelium"
)

func main() {
	stop := make(chan bool)

	links_in := make(chan string, 1)
	links_in <- "https://news.ycombinator.com/"

	links_out := make(chan string, 100000)

	pages := make(chan mycelium.Page, 1024)

	wantMore := make(chan bool)

	taskQueue := mycelium.NewDefaultRedisTaskQueue()
	go taskQueue.Listen(links_out, links_in, wantMore)

	dataStore := mycelium.NewDefaultRedisDataStore()

	crawler := mycelium.NewCrawler()

	for i := 0; i < 100; i++ {
		go crawler.Listen(links_in, links_out, wantMore, pages, stop)
	}

	dataStore.Listen(pages)

	defer dataStore.Stop()
	defer taskQueue.Stop()
}
