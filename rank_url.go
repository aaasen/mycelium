package mycelium

import (
	"log"
	"math"
	"net/url"
)

func RankLength(url string) float64 {
	maxLength := 128.0

	return math.Min(float64(len(url))/maxLength, 1.0)
}

func RankProtocol(rawURL string) float64 {
	parsedURL, err := url.Parse(rawURL)

	if err != nil {
		log.Fatal(err)
	}

	switch parsedURL.Scheme {
	case "http":
		return 1
	case "https":
		return 0.8
	default:
		return 0.0
	}
}

func Rank(url string) float64 {
	return RankLength(url) + RankProtocol(url)
}
