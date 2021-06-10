package synx

import (
	"sync"
)

type locker struct {
	lock *sync.Mutex
}

func newLock() locker {
	return locker{
		lock: &sync.Mutex{},
	}
}

// Lock the mutex
func (l *locker) Lock() {
	l.lock.Lock()
}

// Unlock the mutex
func (l *locker) Unlock() {
	l.lock.Unlock()
}
