package report

import "sync"

type Store interface {
	Add(r Result) error
	List() ([]Result, error)
}

type InMemoryStore struct {
	store []Result
	size  int
	rwm   *sync.RWMutex
}

func (s *InMemoryStore) Add(r Result) error {
	s.rwm.Lock()
	last := len(s.store)

	if len(s.store) == s.size {
		last = len(s.store) - 1
	}

	s.store = append([]Result{r}, s.store[:last]...)
	s.rwm.Unlock()

	return nil
}

func (s *InMemoryStore) List() ([]Result, error) {
	defer s.rwm.RUnlock()

	s.rwm.RLock()
	return s.store, nil
}

func NewResultStore(size int) *InMemoryStore {
	return &InMemoryStore{
		size:  size,
		rwm:   new(sync.RWMutex),
		store: make([]Result, 0),
	}
}
