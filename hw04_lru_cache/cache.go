package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	item, has := c.items[key]
	if !has {
		if len(c.items) == c.capacity {
			i := c.queue.Back()
			c.queue.Remove(i)
			delete(c.items, i.Key)
		}
		item = c.queue.PushFront(value)
		item.Key = key
		c.items[key] = item
	} else {
		item.Value = value
		c.queue.MoveToFront(item)
	}

	c.mu.Unlock()
	return has
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, has := c.items[key]
	if has {
		c.queue.MoveToFront(item)
		return item.Value, has
	}

	return nil, has
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.mu.Unlock()
}
