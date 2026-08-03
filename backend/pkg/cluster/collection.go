package cluster

import (
	"sync"

	"github.com/gosimple/slug"
	"go.uber.org/zap"

	"github.com/kyverno/policy-reporter-ui/pkg/utils"
)

type Collection struct {
	mx       *sync.RWMutex
	clusters map[string]*Cluster
}

func (c *Collection) Add(key string, t Config) {
	c.mx.Lock()
	cluster, err := MapClusterToModel(key, t)
	if err == nil {
		zap.L().Debug("added cluster to collection", zap.String("cluster", key))
		c.clusters[key] = cluster
	} else {
		zap.L().Error("failed to add cluster to collection", zap.String("cluster", key), zap.Error(err))
	}
	c.mx.Unlock()
}

func (c *Collection) Remove(key string) {
	c.mx.Lock()
	zap.L().Debug("removed cluster from collection", zap.String("cluster", key))
	delete(c.clusters, key)
	c.mx.Unlock()
}

func (c *Collection) List() []*Cluster {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return utils.ToList(c.clusters)
}

func (c *Collection) All() map[string]*Cluster {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.clusters
}

func (c *Collection) Length() int {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return len(c.clusters)
}

func (c *Collection) Cluster(name string) *Cluster {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.clusters[name]
}

// NewCollection creates a new target Collection.
func NewCollection(clusters ...Config) *Collection {
	collection := &Collection{
		clusters: make(map[string]*Cluster),
		mx:       new(sync.RWMutex),
	}

	for _, t := range clusters {
		collection.Add(slug.Make(t.Name), t)
	}

	return collection
}
