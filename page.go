package mycelium

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Page struct {
	URL  string
	Body string

	response *http.Response
}

// Creates a new Page using an http.Response
func NewPageFromResponse(resp *http.Response) *Page {
	bodyReader := io.MultiReader(resp.Body)
	body, err := ioutil.ReadAll(bodyReader)

	if err != nil {
		log.Println("error reading response body: %v", err)
	}

	return &Page{
		URL:      resp.Request.URL.String(),
		Body:     string(body),
		response: resp,
	}
}

// Extracts all links (<a> tags with href attributes) from a Page
func (self *Page) GetLinks() []string {
	return getLinks(self.response.Request.URL, strings.NewReader(self.Body))
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
