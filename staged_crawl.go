package mycelium

import (
	"time"
)

func StagedCrawl(seedUrls []string, stages int) []*Page {
	if stages > 1 {
		pages := NewWorker().GetPages(seedUrls, time.Second*5)

		links := make([]string, 0)
		for _, page := range pages {
			links = append(links, page.GetLinks()...)
		}

		return append(pages, self.Crawl(links, stages-1)...)
	} else if stages == 1 {
		return NewWorker().GetPages(seedUrls, time.Second*5)
	} else {
		return nil
	}
}
