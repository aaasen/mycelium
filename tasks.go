package crawl

import (
	"github.com/garyburd/redigo/redis"
	"log"
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
	for {
		select {
		case newTask := <-self.Incoming:
			self.Outgoing <- newTask

			// add to redis
		case wantMore := <-self.WantMore:
			// fetch from redis, send to outgoing

			log.Println(wantMore)
		}
	}
}
