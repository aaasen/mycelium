package mycelium

import (
	"strings"
	"testing"
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
