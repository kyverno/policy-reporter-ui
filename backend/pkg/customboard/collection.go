package customboard

import (
	"sync"

	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

type Collection struct {
	mx     *sync.RWMutex
	boards map[string]*CustomBoard
}

func (c *Collection) Add(key string, t *CustomBoard) {
	c.mx.Lock()
	c.boards[key] = t
	c.mx.Unlock()
}

func (c *Collection) Remove(key string) {
	c.mx.Lock()
	delete(c.boards, key)
	c.mx.Unlock()
}

func (c *Collection) Boards() []*CustomBoard {
	c.mx.RLock()
	defer c.mx.Unlock()

	return utils.ToList(c.boards)
}

func (c *Collection) Board(id string) *CustomBoard {
	c.mx.RLock()
	defer c.mx.Unlock()

	return utils.FindInMap(c.boards, func(b *CustomBoard) bool {
		return b.ID == id
	}, nil)
}

func (c *Collection) Length() int {
	return len(c.boards)
}

// NewCollection creates a new target Collection.
func NewCollection(boards ...*CustomBoard) *Collection {
	collection := &Collection{
		boards: make(map[string]*CustomBoard, 0),
		mx:     new(sync.RWMutex),
	}

	for _, t := range boards {
		if t != nil {
			collection.Add(t.ID, t)
		}
	}

	return collection
}
