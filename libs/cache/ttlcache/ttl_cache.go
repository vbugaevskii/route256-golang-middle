package ttlcache

import (
	"container/list"
	"context"
	"sync"
	"time"
)

type Item[K comparable, V any] struct {
	Key       K
	Value     V
	CreatedAt time.Time
}

type TTLCache[K comparable, V any] struct {
	ttl      int
	capacity int
	queue    *list.List
	mutex    *sync.RWMutex
	items    map[K]*list.Element
}

func NewTTLCache[K comparable, V any](capacity int, ttl int) *TTLCache[K, V] {
	return &TTLCache[K, V]{
		ttl:      ttl,
		capacity: capacity,
		queue:    list.New(),
		mutex:    new(sync.RWMutex),
		items:    make(map[K]*list.Element),
	}
}

func (c *TTLCache[K, V]) Add(key K, value V) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		element.Value.(*Item[K, V]).Value = value
		element.Value.(*Item[K, V]).CreatedAt = time.Now()
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

func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	element, exists := c.items[key]
	if !exists {
		var value V
		return value, false
	}

	if c.hasExpired(element) {
		c.deleteItem(element)

		var value V
		return value, false
	}

	c.queue.MoveToFront(element)
	return element.Value.(*Item[K, V]).Value, true
}

func (c *TTLCache[K, V]) Remove(key K) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.deleteItem(element)
		return true
	}
	return false
}

func (c *TTLCache[K, V]) Contains(key K) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, exists := c.items[key]
	return !c.hasExpired(element) && exists
}

func (c *TTLCache[K, V]) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.items)
}

func (c *TTLCache[K, V]) clear() {
	if element := c.queue.Back(); element != nil {
		c.deleteItem(element)
	}
}

func (c *TTLCache[K, V]) deleteItem(element *list.Element) {
	item := c.queue.Remove(element).(*Item[K, V])
	delete(c.items, item.Key)
}

func (c *TTLCache[K, V]) hasExpired(element *list.Element) bool {
	item := element.Value.(*Item[K, V])
	return item.CreatedAt.Add(time.Duration(c.ttl) * time.Second).After(time.Now())
}

func (c *TTLCache[K, V]) RunCleaner(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.ttl) * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for e := c.queue.Front(); e != nil; e = e.Next() {
					if c.hasExpired(e) {
						prev := e.Prev()
						c.deleteItem(e)
						e = prev
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
