package mycelium

import (
	"log"
)

func InfiniteCrawl(linksIn <-chan string, linksOut chan<- string, wantMore chan<- bool, pages chan<- *Page) {
	for {
		select {
		case link := <-linksIn:
			go func(url string) {
				page, err := NewWorker().GetPage(url)

				if err != nil {
					log.Println(err)
					return
				}

				pages <- page

				links := page.GetLinks()

				for _, link := range links {
					linksOut <- link
				}
			}(link)
		default:
			wantMore <- true
		}
	}
}
