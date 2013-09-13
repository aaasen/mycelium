package mycelium

import (
	"log"
	"net/http"
)

type Crawler struct {
}

func NewCrawler() *Crawler {
	return &Crawler{}
}

func Get(url string) (*Page, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return NewPage(resp), nil
}

func (self *Crawler) Listen(linksIn <-chan string, linksOut chan<- string, wantMore chan<- bool, pages chan<- Page, stop <-chan bool) {
	for {
		select {
		case stopSignal := <-stop:
			if stopSignal {
				return
			}
		case link := <-linksIn:
			log.Printf("fetching: %v\n", link)
			page, err := Get(link)

			if err != nil {
				log.Printf("Get(): %v", err)
				break
			}

			for _, aLink := range page.Links {
				linksOut <- aLink
			}

			pages <- *page
		default:
			wantMore <- true
		}
	}
}
