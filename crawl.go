package crawl

import (
	"log"
	"net/http"
)

func Fetch(url string, pages chan Page) {
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("Fetch(): %v", err)
		return
	}

	page := NewPage(resp)

	pages <- *page
}

func Get(url string) (*Page, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return NewPage(resp), nil
}

func Crawl(stop chan bool, links chan string, pages chan Page) {
	for {
		select {
		case stopSignal := <-stop:
			log.Println("stopping")
			if stopSignal {
				return
			}
		case link := <-links:
			log.Printf("fetching: %v\n", link)
			page, err := Get(link)

			if err != nil {
				log.Printf("Get(): %v", err)
			}

			pages <- *page

			for _, link := range page.Links {
				links <- link
			}

		default:
			log.Println("no links")
		}
	}
}
