package crawl

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type TaskQueue struct {
	conn redis.Conn

	Incoming <-chan string
	Outgoing chan<- string
}

func NewDefaultTaskQueue(Incoming <-chan string, Outgoing chan<- string) *TaskQueue {
	conn, err := redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatal(err)
	}

	return NewTaskQueue(Incoming, Outgoing, conn)
}

func NewTaskQueue(Incoming <-chan string, Outgoing chan<- string, conn redis.Conn) *TaskQueue {
	return &TaskQueue{
		Incoming: Incoming,
		Outgoing: Outgoing,
		conn:     conn,
	}
}

func (self *TaskQueue) Run() {
	for {
		select {
		case newTask := <-self.Incoming:
			self.Outgoing <- newTask
		}
	}
}
