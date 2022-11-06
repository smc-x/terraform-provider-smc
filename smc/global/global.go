package global

import "sync"

var (
	deferred = make([]func(), 0)
	lock     = &sync.Mutex{}
)

func Defer(cb func()) {
	lock.Lock()
	defer lock.Unlock()
	deferred = append(deferred, cb)
}

func Run() {
	lock.Lock()
	defer lock.Unlock()
	for i := len(deferred) - 1; i >= 0; i-- {
		deferred[i]()
	}
	deferred = make([]func(), 0)
}
