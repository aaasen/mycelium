package crawl

import (
	"log"
)

type DataStore interface {
	Save(*Page) error
}

type DebugDataStore struct {
}

func NewDebugDataStore() *DebugDataStore {
	return &DebugDataStore{}
}

func (self *DebugDataStore) Save(page *Page) error {
	log.Printf("recieved: %v", page.URL)
	return nil
}
