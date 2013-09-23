package mycelium

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	robots "github.com/temoto/robotstxt-go"
)

type RobotFilter struct {
	UserAgent string

	robotData map[string]*robots.RobotsData
}

func NewRobotFilter() *RobotFilter {
	return &RobotFilter{
		UserAgent: "Mycelium",
		robotData: make(map[string]*robots.RobotsData),
	}
}

// Checks if the given url is allowed to be crawled by robots and
// retrieves the page using http.Get() if it is.
func (self *RobotFilter) PoliteGet(url string) (*http.Response, error) {
	if !self.Allowed(url) {
		return nil, errors.New(fmt.Sprintf("robots not allowed for: %v", url))
	}

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Checks if the given url is allowed to be crawled using
// github.com/temoto/robotstxt-go
func (self *RobotFilter) Allowed(rawUrl string) bool {
	url, err := url.Parse(rawUrl)

	if err != nil {
		log.Printf("error parsing %s: %v", rawUrl, err)
		return false
	}

	if robotData, ok := self.robotData[url.Host]; ok {
		return robotData.TestAgent(url.Path, self.UserAgent)
	} else {
		robotUrl := fmt.Sprintf("%s://%s/robots.txt", url.Scheme, url.Host)
		resp, err := http.Get(robotUrl)

		if err != nil {
			log.Printf("error getting robots.txt for %s: %v", robotUrl, err)
			return false
		}

		newRobotData, err := robots.FromResponse(resp)
		resp.Body.Close()

		if err != nil {
			log.Printf("error parsing robots.txt at %s: %v", robotUrl, err)
			return false
		}

		self.robotData[url.Host] = newRobotData
		return self.robotData[url.Host].TestAgent(url.Path, self.UserAgent)
	}
}
