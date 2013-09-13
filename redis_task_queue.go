package mycelium

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type RedisTaskQueue struct {
	conn redis.Conn
}

func NewDefaultRedisTaskQueue() *RedisTaskQueue {
	conn, err := redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatalf("in NewDefaultRedisTaskQueue():", err)
	}

	return NewRedisTaskQueue(conn)
}

func NewRedisTaskQueue(conn redis.Conn) *RedisTaskQueue {
	return &RedisTaskQueue{
		conn: conn,
	}
}

func (self *RedisTaskQueue) Listen(incoming <-chan string, outgoing chan<- string, wantMore <-chan bool) {
	defer self.conn.Close()

	for {
		select {
		case newLink := <-incoming:
			err := self.Push(newLink)

			if err != nil {
				log.Panicf("error pushing to redis queue: %v", err)
			}
		case <-wantMore:
			links, err := self.Pop()

			if err != nil {
				log.Panicf("error popping from redis queue: %v", err)
				break
			}

			for _, link := range links {
				outgoing <- link
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

func (self *RedisTaskQueue) Stop() {
	self.conn.Close()
}
