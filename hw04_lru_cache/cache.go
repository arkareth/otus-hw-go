package hw04lrucache

import (
	"sync"
)

type Key string

type elem struct {
	Key   Key
	Value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	sync.Mutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	defer c.Unlock()
	c.Lock()

	i, ok := c.items[key]
	switch {
	case ok:
		i.Value = &elem{Key: key, Value: value}
		c.queue.MoveToFront(i)

		return true
	default:
		c.items[key] = c.queue.PushFront(&elem{Key: key, Value: value})
		if c.queue.Len() > c.capacity {
			l := c.queue.Back()
			c.queue.Remove(l)
			delete(c.items, l.Value.(*elem).Key)
		}

		return false
	}
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	defer c.Unlock()
	c.Lock()

	i, ok := c.items[key]
	switch {
	case ok:
		c.queue.MoveToFront(i)
		return i.Value.(*elem).Value, true
	default:
		return nil, false
	}
}

func (c *lruCache) Clear() {
	defer c.Unlock()
	c.Lock()

	c.capacity = 0
	c.queue = nil
	c.items = nil
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
