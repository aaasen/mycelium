package mycelium

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type RedisTaskQueue struct {
	conn redis.Conn

	Incoming <-chan string
	Outgoing chan<- string
	WantMore <-chan bool
}

func NewDefaultRedisTaskQueue(incoming <-chan string, outgoing chan<- string, wantMore <-chan bool) *RedisTaskQueue {
	conn, err := redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatal(err)
	}

	return NewRedisTaskQueue(incoming, outgoing, wantMore, conn)
}

func NewRedisTaskQueue(incoming <-chan string, outgoing chan<- string, wantMore <-chan bool, conn redis.Conn) *RedisTaskQueue {
	return &RedisTaskQueue{
		Incoming: incoming,
		Outgoing: outgoing,
		WantMore: wantMore,
		conn:     conn,
	}
}

func (self *RedisTaskQueue) Run() {
	defer self.conn.Close()

	for {
		select {
		case newLink := <-self.Incoming:
			err := self.Push(newLink)

			if err != nil {
				log.Panicf("error pushing to redis queue: %v", err)
			}
		case <-self.WantMore:
			links, err := self.Pop()

			if err != nil {
				log.Panicf("error popping from redis queue: %v", err)
				break
			}

			for _, link := range links {
				self.Outgoing <- link
			}
		}
	}
}

func (self *RedisTaskQueue) Push(link string) error {
	_, err := self.conn.Do("ZADD", "uncrawled", time.Now().Unix(), link)

	return err
}

func (self *RedisTaskQueue) Pop() ([]string, error) {
	links, getErr := redis.Strings(self.conn.Do("ZRANGE", "uncrawled", 0, 1))

	if getErr != nil {
		return nil, getErr
	}

	_, remErr := self.conn.Do("ZREMRANGEBYRANK", "uncrawled", 0, 1)

	if remErr != nil {
		return nil, remErr
	}

	return links, nil
}
