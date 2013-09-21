package main

import (
	"github.com/aaasen/mycelium"
)

func main() {
	links_in := make(chan string, 10000)
	links_in <- "https://news.ycombinator.com/"

	links_robot_filter := make(chan string, 2048)

	links_out := make(chan string, 2048)

	pages := make(chan mycelium.Page, 2048)

	wantMore := make(chan bool)

	taskQueue := mycelium.NewDefaultRedisTaskQueue()
	dataStore := mycelium.NewDefaultRedisDataStore()
	robotFilter := mycelium.NewRobotFilter()
	crawler := mycelium.NewCrawler()

	go robotFilter.Listen(links_out, links_robot_filter)
	go taskQueue.Listen(links_robot_filter, links_in, wantMore)

	for i := 0; i < 100; i++ {
		go crawler.Listen(links_in, links_out, wantMore, pages)
	}

	defer dataStore.Stop()
	defer taskQueue.Stop()

	dataStore.Listen(pages)
}
