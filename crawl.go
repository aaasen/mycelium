package crawl

import (
	"log"
	"net/http"
)

func Fetch(url string, pages chan Page) {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	page := NewPage(resp)

	pages <- *page
}
