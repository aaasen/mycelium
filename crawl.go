package crawl

import (
	"log"
	"net/http"
	"sync"
)

var iterations = 0

func Fetch(url string, pages chan Page, wg sync.WaitGroup) {
	wg.Add(1)
	resp, err := http.Get(url)

	if err != nil {
		wg.Done()
		log.Printf("Fetch(): %v", err)
		return
	}

	page := NewPage(resp)

	pages <- *page
	wg.Done()
}

func Crawl(urls []string, pages chan Page, wg sync.WaitGroup) {
	log.Println(iterations)

	if iterations > 0 {
		return
	}

	iterations++

	for _, url := range urls {
		log.Printf("fetching: %v", url)
		go Fetch(url, pages, wg)
	}

	for page := range pages {
		go Crawl(page.links, pages, wg)
	}
}
