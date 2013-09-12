package crawl

import (
	"fmt"
	"log"
	"net/http"
)

func Fetch(url string, pages chan Page) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("Fetch(): %v", err)
	}

	page := NewPage(resp)

	pages <- *page
}

func Crawl(urls []string, pages chan Page) {
	for _, url := range urls {
		fmt.Println("url")
		go Fetch(url, pages)
	}

	for page := range pages {
		fmt.Println(page.URL)
		fmt.Println(page.links)

		for _, link := range page.links {
			go Fetch(link.String(), pages)
		}
	}
}
