package main

import (
	"github.com/aaasen/crawl"
	"log"
	"time"
)

func main() {

	stop := make(chan bool)
	links := make(chan string, 2048)

	links <- "https://news.ycombinator.com/"

	log.Println("hey")

	pages := make(chan crawl.Page, 1024)

	go crawl.Crawl(stop, links, pages)

	time.Sleep(time.Second * 2)

	stop <- true
}
