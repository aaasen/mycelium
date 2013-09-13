package mycelium

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL string

	response *http.Response
	Links    []string
}

func NewPage(resp *http.Response) *Page {
	return &Page{
		URL:      resp.Request.URL.String(),
		Links:    getLinks(resp.Body),
		response: resp,
	}
}

func getLinks(content io.Reader) []string {
	doc, err := goquery.NewDocumentFromReader(content)

	if err != nil {
		log.Printf("GetLinks(): %v", err)
		return []string{}
	}

	urls := []string{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, hrefExists := s.Attr("href"); hrefExists {
			url, err := url.Parse(href)

			if err != nil {
				log.Printf("GetLinks(): %v", err)
				return
			}

			// TODO: fix this
			if url.IsAbs() {
				urls = append(urls, url.String())
			}
		}
	})

	return urls
}
