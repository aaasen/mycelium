package mycelium

import ()

type DataStore interface {
	Save(*Page) error
	Listen(<-chan *Page)
}
