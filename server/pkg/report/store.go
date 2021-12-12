package report

import "sync"

type ResultStore struct {
	store []Result
	size  int
	rwm   *sync.RWMutex
}

func (s *ResultStore) Add(r Result) {
	s.rwm.Lock()
	last := len(s.store)

	if len(s.store) == s.size {
		last = len(s.store) - 1
	}

	s.store = append([]Result{r}, s.store[:last]...)
	s.rwm.Unlock()
}

func (s *ResultStore) List() []Result {
	defer s.rwm.RUnlock()

	s.rwm.RLock()
	return s.store
}

func NewResultStore(size int) *ResultStore {
	return &ResultStore{
		size:  size,
		rwm:   new(sync.RWMutex),
		store: make([]Result, 0),
	}
}
