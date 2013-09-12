package crawl

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
	links    []*url.URL
}

func NewPage(resp *http.Response) *Page {
	return &Page{
		URL:      resp.Request.URL.String(),
		response: resp,
		links:    getLinks(resp.Body),
	}
}

func getLinks(content io.Reader) []*url.URL {
	doc, err := goquery.NewDocumentFromReader(content)

	if err != nil {
		log.Fatal(err)
	}

	urls := []*url.URL{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, hrefExists := s.Attr("href"); hrefExists {
			url, err := url.Parse(href)

			if err != nil {
				log.Fatalf("GetLinks(): %v", err)
			}

			urls = append(urls, url)
		}
	})

	return urls
}
