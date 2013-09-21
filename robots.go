package mycelium

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	robots "github.com/temoto/robotstxt-go"
)

type RobotFilter struct {
	robotData map[string]*robots.RobotsData
}

func NewRobotFilter() *RobotFilter {
	return &RobotFilter{
		robotData: make(map[string]*robots.RobotsData),
	}
}

func (self *RobotFilter) allowed(rawUrl string) bool {
	log.Println(rawUrl)

	url, err := url.Parse(rawUrl)

	if err != nil {
		log.Printf("error parsing %s: %v", rawUrl, err)
		return false
	}

	if robotData, ok := self.robotData[url.Host]; ok {
		return robotData.TestAgent(url.Path, "Mycelium")
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
		return self.robotData[url.Host].TestAgent(url.Path, "Mycelium")
	}
}
