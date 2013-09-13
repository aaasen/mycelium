package crawl

import (
	"log"
	"net/http"
)

func Get(url string) (*Page, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return NewPage(resp), nil
}

func Crawl(links_in <-chan string, links_out chan<- string, pages chan<- Page, stop <-chan bool) {
	for {
		select {
		case stopSignal := <-stop:
			if stopSignal {
				return
			}
		case link := <-links_in:
			log.Printf("fetching: %v\n", link)
			page, err := Get(link)

			if err != nil {
				log.Printf("Get(): %v", err)
				break
			}

			for _, aLink := range page.Links {
				links_out <- aLink
			}

			pages <- *page

		}
	}
}
