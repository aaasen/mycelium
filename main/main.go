package main

import (
	"github.com/aaasen/mycelium"
)

func main() {
	stop := make(chan bool)

	links_in := make(chan string, 2048)

	links_out := make(chan string, 2048)

	pages := make(chan mycelium.Page, 2048)

	wantMore := make(chan bool)

	taskQueue := mycelium.NewDefaultRedisTaskQueue()
	dataStore := mycelium.NewDefaultRedisDataStore()
	crawler := mycelium.NewCrawler()

	for i := 0; i < 1; i++ {
		go taskQueue.Listen(links_out, links_in, wantMore)
	}

	for i := 0; i < 100; i++ {
		go crawler.Listen(links_in, links_out, wantMore, pages, stop)
	}

	// for i := 0; i < 100; i++ {
	// 	go dataStore.Listen(pages)
	// }

	defer dataStore.Stop()
	defer taskQueue.Stop()

	dataStore.Listen(pages)
}
