package mycelium

import (
	"log"
)

type DataStore interface {
	Save(*Page) error
}

type DebugDataStore struct {
	Pages <-chan Page
}

func NewDebugDataStore(pages <-chan Page) *DebugDataStore {
	return &DebugDataStore{
		Pages: pages,
	}
}

func (self *DebugDataStore) Run() {
	for {
		select {
		case page := <-self.Pages:
			self.Save(&page)
		}
	}
}

func (self *DebugDataStore) Save(page *Page) error {
	log.Printf("received: %v", page.URL)

	return nil
}
