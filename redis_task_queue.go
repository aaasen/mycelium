package mycelium

import (
	"github.com/garyburd/redigo/redis"
	"log"
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
				log.Fatalf("error pushing to redis queue: %v", err)
			}
		case <-wantMore:
			links, err := self.Pop(100)

			if err != nil {
				log.Fatalf("error popping from redis queue: %v", err)
				break
			}

			for _, link := range links {
				outgoing <- link
			}
		}
	}
}

func (self *RedisTaskQueue) Push(link string) error {
	crawled, err := self.hasBeenCrawled(link)

	if err != nil {
		return err
	}

	if !crawled {
		_, err := self.conn.Do("ZADD", "uncrawled", Rank(link), link)

		return err
	}

	return nil
}

func (self *RedisTaskQueue) Pop(numTasks int) ([]string, error) {
	links, getErr := redis.Strings(self.conn.Do("ZRANGE", "uncrawled", 0, numTasks))

	if getErr != nil {
		return nil, getErr
	}

	_, remErr := self.conn.Do("ZREMRANGEBYRANK", "uncrawled", 0, 1)

	if remErr != nil {
		return nil, remErr
	}

	for _, link := range links {
		err := self.markAsCrawled(link)

		if err != nil {
			return nil, err
		}
	}

	return links, nil
}

func (self *RedisTaskQueue) Stop() {
	self.conn.Close()
}

func (self *RedisTaskQueue) markAsCrawled(url string) error {
	_, err := self.conn.Do("SADD", "crawled", url)

	return err
}

func (self *RedisTaskQueue) hasBeenCrawled(url string) (bool, error) {
	return redis.Bool(self.conn.Do("SISMEMBER", "crawled", url))
}
