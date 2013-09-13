package crawl

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type TaskQueue struct {
	conn redis.Conn

	Incoming <-chan string
	Outgoing chan<- string
	WantMore <-chan bool
}

func NewDefaultTaskQueue(incoming <-chan string, outgoing chan<- string, wantMore <-chan bool) *TaskQueue {
	conn, err := redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatal(err)
	}

	return NewTaskQueue(incoming, outgoing, wantMore, conn)
}

func NewTaskQueue(incoming <-chan string, outgoing chan<- string, wantMore <-chan bool, conn redis.Conn) *TaskQueue {
	return &TaskQueue{
		Incoming: incoming,
		Outgoing: outgoing,
		WantMore: wantMore,
		conn:     conn,
	}
}

func (self *TaskQueue) Run() {
	defer self.conn.Close()

	for {
		select {
		case newLink := <-self.Incoming:
			_, err := self.conn.Do("ZADD", "uncrawled", time.Now().Unix(), newLink)

			if err != nil {
				log.Panicln(err)
			}
		case <-self.WantMore:
			links, getErr := redis.Strings(self.conn.Do("ZRANGE", "uncrawled", 0, 1))

			if getErr != nil {
				log.Panicln(getErr)
				break
			}

			_, remErr := self.conn.Do("ZREMRANGEBYRANK", "uncrawled", 0, 1)

			if remErr != nil {
				log.Panicln(remErr)
				break
			}

			for _, link := range links {
				log.Printf("task: %v\n", link)

				self.Outgoing <- link
			}
		}
	}
}
