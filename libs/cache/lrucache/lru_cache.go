package lrucache

import (
	"container/list"
	"sync"
)

type Item[K comparable, V any] struct {
	Key   K
	Value V
}

type LRUCache[K comparable, V any] struct {
	capacity int
	queue    *list.List
	mutex    *sync.RWMutex
	items    map[K]*list.Element
}

func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		queue:    list.New(),
		mutex:    new(sync.RWMutex),
		items:    make(map[K]*list.Element),
	}
}

func (c *LRUCache[K, V]) Add(key K, value V) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		element.Value.(*Item[K, V]).Value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		c.clear()
	}

	item := &Item[K, V]{
		Key:   key,
		Value: value,
	}

	element := c.queue.PushFront(item)
	c.items[item.Key] = element

	return true
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	element, exists := c.items[key]
	if !exists {
		var value V
		return value, false
	}

	c.queue.MoveToFront(element)
	return element.Value.(*Item[K, V]).Value, true
}

func (c *LRUCache[K, V]) Remove(key K) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.deleteItem(element)
		return true
	}
	return false
}

func (c *LRUCache[K, V]) Contains(key K) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, exists := c.items[key]
	return exists
}

func (c *LRUCache[K, V]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.items)
}

func (c *LRUCache[K, V]) clear() {
	if element := c.queue.Back(); element != nil {
		c.deleteItem(element)
	}
}

func (c *LRUCache[K, V]) deleteItem(element *list.Element) {
	item := c.queue.Remove(element).(*Item[K, V])
	delete(c.items, item.Key)
}
