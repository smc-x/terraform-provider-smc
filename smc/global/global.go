package global

import "sync"

var (
	Deferred = make([]func(), 0)

	lock = &sync.Mutex{}
)

func Defer(cb func()) {
	lock.Lock()
	defer lock.Unlock()
	Deferred = append(Deferred, cb)
}

func Run() {
	lock.Lock()
	defer lock.Unlock()
	for i := len(Deferred) - 1; i >= 0; i-- {
		Deferred[i]()
	}
}
