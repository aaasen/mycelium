package mycelium

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisDataStore struct {
	Pages <-chan Page

	conn redis.Conn
}

func NewRedisDataStore(pages <-chan Page, conn redis.Conn) *RedisDataStore {
	return &RedisDataStore{
		Pages: pages,
		conn:  conn,
	}
}

func NewDefaultRedisDataStore(pages <-chan Page) *RedisDataStore {
	conn, err := redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatalf("in NewRedisDataStore(): %v", err)
	}

	return NewRedisDataStore(pages, conn)
}

func (self *RedisDataStore) Run() {
	for {
		select {
		case page := <-self.Pages:
			err := self.Save(&page)

			if err != nil {
				log.Panicf("in RedisDataStore.Save(): %v", err)
			}
		}
	}
}

func (self *RedisDataStore) Save(page *Page) error {
	log.Printf("received: %v", page.URL)

	_, err := self.conn.Do("SET", page.URL, "page")

	return err
}
