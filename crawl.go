package mycelium

import (
	"log"
	"net/http"
)

func get(url string) (*Page, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return NewPage(resp), nil
}

func Crawl(links_in <-chan string, links_out chan<- string, wantMore chan<- bool, pages chan<- Page, stop <-chan bool) {
	for {
		select {
		case stopSignal := <-stop:
			if stopSignal {
				return
			}
		case link := <-links_in:
			log.Printf("fetching: %v\n", link)
			page, err := get(link)

			if err != nil {
				log.Printf("Get(): %v", err)
				break
			}

			for _, aLink := range page.Links {
				links_out <- aLink
			}

			pages <- *page
		default:
			wantMore <- true
		}
	}
}
