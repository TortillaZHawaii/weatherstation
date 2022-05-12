package main

import "sync"

type node struct {
	data measurement
	next *node
	prev *node
}

type doublyLinkedList struct {
	length int
	head   *node
	tail   *node
}

type measurementsCache struct {
	mu   sync.RWMutex
	list doublyLinkedList
}

func measurementsCacheInit() *measurementsCache {
	return &measurementsCache{
		list: doublyLinkedList{
			length: 0,
			head:   nil,
			tail:   nil,
		},
	}
}

func (c *measurementsCache) addFront(m measurement) {
	c.mu.Lock()
	defer c.mu.Unlock()
	newNode := &node{data: m, next: c.list.head, prev: nil}
	if c.list.head != nil {
		c.list.head.prev = newNode
	}
	c.list.head = newNode
	if c.list.tail == nil {
		c.list.tail = newNode
	}
	c.list.length++
}

func (c *measurementsCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.list.head = nil
	c.list.tail = nil
	c.list.length = 0
}

func (c *measurementsCache) getArray() []measurement {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var list []measurement
	for n := c.list.head; n != nil; n = n.next {
		list = append(list, n.data)
	}
	return list
}
