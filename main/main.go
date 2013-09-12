package main

import (
	"github.com/aaasen/crawl"
	"log"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	pages := make(chan crawl.Page)
	crawl.Crawl([]string{"https://news.ycombinator.com/"}, pages, wg)

	wg.Wait()

	for i := 0; i < 5; i++ {
		page := <-pages
		log.Println(page)
	}
}
