package mycelium

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type RedisDataStore struct {
	conn redis.Conn
}

func NewRedisDataStore(conn redis.Conn) *RedisDataStore {
	return &RedisDataStore{
		conn: conn,
	}
}

func NewDefaultRedisDataStore() *RedisDataStore {
	conn, err := redis.Dial("tcp", ":6379")

	if err != nil {
		log.Fatalf("in NewDefaultRedisDataStore(): %v", err)
	}

	return NewRedisDataStore(conn)
}

func (self *RedisDataStore) Listen(pages <-chan Page) {
	for {
		select {
		case page := <-pages:
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

func (self *RedisDataStore) Stop() {
	self.conn.Close()
}
