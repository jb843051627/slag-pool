package service

import (
	"sync"

	"github.com/jb843051627/slag-pool/internal/model"
)

type ReadingCache struct {
	mu   sync.RWMutex
	data map[int64][]model.SensorReading
}

func NewReadingCache() *ReadingCache {
	return &ReadingCache{data: make(map[int64][]model.SensorReading)}
}

func (c *ReadingCache) Get(key int64) ([]model.SensorReading, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *ReadingCache) Update(key int64, val []model.SensorReading) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.data[key] = val
}