package main

import "sync"

type bounds struct {
	minTemperature int8
	maxTemperature int8
	minHumidity    int8
	maxHumidity    int8
}

type boundsCache struct {
	mu     sync.RWMutex
	bounds bounds
}

func boundsCacheInit() *boundsCache {
	return &boundsCache{
		mu: sync.RWMutex{},
		bounds: bounds{
			minTemperature: 0,
			maxTemperature: 50,
			minHumidity:    0,
			maxHumidity:    100,
		},
	}
}

func (c *boundsCache) set(b bounds) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bounds = b
}

func (c *boundsCache) get() bounds {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bounds
}
