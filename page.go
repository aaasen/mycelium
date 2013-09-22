package mycelium

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL  string
	Body string

	response *http.Response
	Links    []string
}

func NewPage(resp *http.Response) *Page {
	bodyReader := io.MultiReader(resp.Body)
	body, err := ioutil.ReadAll(bodyReader)

	if err != nil {
		log.Println("error reading response body: %v", err)
	}

	return &Page{
		URL:      resp.Request.URL.String(),
		Body:     string(body),
		Links:    getLinks(resp.Request.URL, bytes.NewReader(body)),
		response: resp,
	}
}

func getRootPath(url *url.URL) string {
	return fmt.Sprintf("%s://%s", url.Scheme, url.Host)
}

func getLinks(baseUrl *url.URL, content io.Reader) []string {
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

			urls = append(urls, baseUrl.ResolveReference(url).String())
		}
	})

	return urls
}
