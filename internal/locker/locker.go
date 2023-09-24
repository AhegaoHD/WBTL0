package locker

import (
	"sync"
	//uuid "github.com/satori/go.uuid"
)

type Locker struct {
	mu sync.Map
}

func New() *Locker {
	return &Locker{
		mu: sync.Map{},
	}
}

func (m *Locker) Lock(key string) func() {
	value, _ := m.mu.LoadOrStore(key, &sync.Mutex{})
	mtx := value.(*sync.Mutex)
	mtx.Lock()

	return func() { mtx.Unlock() }
}
