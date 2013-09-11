package crawl

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Page struct {
	URL     *url.URL
	Content []byte
}

func NewPage(resp *http.Response) *Page {
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	return &Page{
		URL:     resp.Request.URL,
		Content: content,
	}
}
