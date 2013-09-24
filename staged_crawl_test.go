package mycelium

import (
	"log"
	"testing"
)

func testStagedCrawl(t *testing.T, level, expectedPages int) {
	pages := StagedCrawl([]string{"http://127.0.0.1:8000/"}, level)

	if len(pages) != expectedPages {
		for _, page := range pages {
			log.Println(page)
		}

		t.Errorf("expected %v page, got %v", expectedPages, len(pages))
	}
}

func TestStagedCrawlLevel1(t *testing.T) {
	testStagedCrawl(t, 1, 1)
}

func TestStagedCrawlLevel2(t *testing.T) {
	testStagedCrawl(t, 2, 3)
}

func TestStagedCrawlLevel3(t *testing.T) {
	testStagedCrawl(t, 3, 7)
}

func TestStagedCrawlLevel4(t *testing.T) {
	testStagedCrawl(t, 4, 10)
}
