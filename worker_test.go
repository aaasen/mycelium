package mycelium

import (
	"strings"
	"testing"
	"time"
)

func TestGetPage(t *testing.T) {
	page, err := NewWorker().GetPage("http://127.0.0.1:8000")

	if err != nil {
		t.Error(err.Error())
	}

	cleanBody := strings.Trim(page.Body, "\r\t\n ")
	expectedBody := `<p>/</p>
<a href="/a">a</a>
<a href="/b">b</a>`

	if cleanBody != expectedBody {
		t.Errorf("received:\n%v\n", cleanBody)
	}
}

func TestGetPages(t *testing.T) {
	pages := NewWorker().GetPages(
		[]string{"http://127.0.0.1:8000/", "http://127.0.0.1:8000/a/"},
		time.Second*5)

	expectedBodyIndex := `<p>/</p>
<a href="/a">a</a>
<a href="/b">b</a>`

	expectedBodyA := `<p>/a/</p>
<a href="1">1</a>
<a href="1">1</a>`

	for _, page := range pages {
		page.Body = strings.Trim(page.Body, "\r\t\n ")
	}

	if !(pages[0].Body == expectedBodyIndex ||
		pages[0].Body == expectedBodyA &&
			pages[1].Body == expectedBodyIndex ||
		pages[1].Body == expectedBodyA) {
		t.Fail()
	}
}
