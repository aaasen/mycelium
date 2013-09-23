package mycelium

import (
	"log"
)

type Crawler struct {
	roboFilter *RobotFilter
}

func NewCrawler() *Crawler {
	return &Crawler{
		roboFilter: NewRobotFilter(),
	}
}

func (self *Crawler) Listen(linksIn <-chan string, linksOut chan<- string, wantMore chan<- bool, pages chan<- Page) {
	for {
		select {
		case link := <-linksIn:
			log.Printf("fetching: %v\n", link)
			resp, getErr := self.roboFilter.PoliteGet(link)

			if getErr != nil {
				log.Printf("when getting page: %v", getErr)
				break
			}

			page := NewPageFromResponse(resp)

			for _, aLink := range page.GetLinks() {
				linksOut <- aLink
			}

			pages <- *page
		default:
			wantMore <- true
		}
	}
}
