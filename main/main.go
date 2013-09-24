package main

import (
	"github.com/aaasen/mycelium"
)

func main() {
	links_in := make(chan string, 10000)
	links_in <- "http://reddit.com/r/all"

	links_out := make(chan string, 2048)

	pages := make(chan *mycelium.Page, 2048)

	wantMore := make(chan bool)

	taskQueue := mycelium.NewDefaultRedisTaskQueue()
	go taskQueue.Listen(links_out, links_in, wantMore)
	defer taskQueue.Stop()

	dataStore := mycelium.NewDefaultRedisDataStore()
	defer dataStore.Stop()
	go dataStore.Listen(pages)

	mycelium.InfiniteCrawl(links_in, links_out, wantMore, pages)
}
