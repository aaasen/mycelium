package mycelium

import (
	"log"
	"net/http"
)

type Crawler struct {
	roboFilter *RobotFilter
}

func NewCrawler() *Crawler {
	return &Crawler{
		roboFilter: NewRobotFilter(),
	}
}

func Get(url string) (*Page, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return NewPage(resp), nil
}

func (self *Crawler) Listen(linksIn <-chan string, linksOut chan<- string, wantMore chan<- bool, pages chan<- Page) {
	for {
		select {
		case link := <-linksIn:
			if !self.roboFilter.allowed(link) {
				log.Println("not allowed: %v", link)
				break
			}

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
