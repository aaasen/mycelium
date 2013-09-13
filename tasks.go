package crawl

import ()

func TaskQueue(incomingTasks <-chan string, outgoingTasks chan<- string) {
	for {
		select {
		case newTask := <-incomingTasks:
			outgoingTasks <- newTask
		}
	}
}
