package mycelium

import ()

type TaskQueue interface {

	// pushes a task onto the queue
	// in this case, the task is a url to be crawled
	Push(string) error

	// pops tasks from the queue
	Pop() ([]string, error)

	Listen(incoming <-chan string, outgoing chan<- string, wantMore <-chan bool)
}
