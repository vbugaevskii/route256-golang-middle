package lrucache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU_Add_existElementWithFullQueueSync_moveToFront(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", "56")
	emptyMap := make(map[string]int)
	lru.Add("someKey3", emptyMap)

	//Act
	lru.Add("someKey1", 10)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])
	assert.Equal(t, "someKey1", frontItem.Key)
	assert.Equal(t, 10, frontItem.Value)
	assert.Equal(t, "someKey2", backItem.Key)
	assert.Equal(t, "56", backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Add_existElementSync_moveToFront(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", "56")

	//Act
	lru.Add("someKey1", 10)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])
	assert.Equal(t, "someKey1", frontItem.Key)
	assert.Equal(t, 10, frontItem.Value)
	assert.Equal(t, "someKey2", backItem.Key)
	assert.Equal(t, "56", backItem.Value)
	assert.Equal(t, 2, lru.queue.Len())
}

func TestLRU_Add_newElementWithFullQueueSync_clearAndPushToFront(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", "56")
	lru.Add("someKey3", Item[string, any]{"key", 7})

	//Act
	lru.Add("someKey4", 5)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])
	assert.Equal(t, "someKey4", frontItem.Key)
	assert.Equal(t, 5, frontItem.Value)
	assert.Equal(t, "someKey2", backItem.Key)
	assert.Equal(t, "56", backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())

	_, exists := lru.Get("someKey1")
	assert.False(t, exists)
}

func TestLRU_Add_newElementSync_pushToFront(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)

	//Act
	lru.Add("someKey2", 5)

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])
	assert.Equal(t, "someKey2", frontItem.Key)
	assert.Equal(t, 5, frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 2, lru.queue.Len())
}

func TestLRU_Add_newElementAsync_allKeysExists(t *testing.T) {
	//Arrange
	wg := sync.WaitGroup{}
	lru := NewLRUCache[string, any](3)
	wg.Add(3)

	//Act
	go func() {
		lru.Add("someKey1", 8)
		wg.Done()
	}()
	go func() {
		lru.Add("someKey2", "56")
		wg.Done()
	}()
	go func() {
		lru.Add("someKey3", Item[string, any]{"key", 7})
		wg.Done()
	}()

	wg.Wait()

	//Assert
	val, _ := lru.Get("someKey1")
	assert.Equal(t, 8, val)

	val, _ = lru.Get("someKey2")
	assert.Equal(t, "56", val)

	val, _ = lru.Get("someKey3")
	assert.Equal(t, Item[string, any]{"key", 7}, val)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Get_hasElement_returnItAndMoveToFront(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", 5)
	lru.Add("someKey3", "90")

	//Act
	item, _ := lru.Get("someKey2")

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])
	assert.Equal(t, 5, item)
	assert.Equal(t, "someKey2", frontItem.Key)
	assert.Equal(t, 5, frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Get_hasNotElement_returnNil(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", 5)
	lru.Add("someKey3", "90")

	//Act
	item, _ := lru.Get("someKey")

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])
	assert.Nil(t, item)
	assert.Equal(t, "someKey3", frontItem.Key)
	assert.Equal(t, "90", frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Remove_hasElement_removeIt(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", 5)
	lru.Add("someKey3", "90")

	//Act
	_ = lru.Remove("someKey2")

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])

	_, exists := lru.Get("someKey2")
	assert.False(t, exists)
	assert.Equal(t, "someKey3", frontItem.Key)
	assert.Equal(t, "90", frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 2, lru.queue.Len())
}

func TestLRU_Remove_hasNotElement_doNothing(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", 5)
	lru.Add("someKey3", "90")

	//Act
	result := lru.Remove("someKey")

	//Assert
	frontItem := lru.queue.Front().Value.(*Item[string, any])
	backItem := lru.queue.Back().Value.(*Item[string, any])
	assert.False(t, result)
	assert.Equal(t, "someKey3", frontItem.Key)
	assert.Equal(t, "90", frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_Contains(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)
	lru.Add("someKey1", 8)
	lru.Add("someKey2", 5)
	lru.Add("someKey3", "90")

	//Assert
	assert.True(t, lru.Contains("someKey1"))
	assert.True(t, lru.Contains("someKey2"))
	assert.True(t, lru.Contains("someKey3"))
	assert.False(t, lru.Contains("someKey4"))
}

func TestLRU_Len(t *testing.T) {
	//Arrange
	lru := NewLRUCache[string, any](3)

	//Assert
	assert.Equal(t, lru.Len(), 0)

	lru.Add("someKey1", 8)
	assert.Equal(t, lru.Len(), 1)

	lru.Add("someKey2", 5)
	assert.Equal(t, lru.Len(), 2)

	lru.Add("someKey3", "90")
	assert.Equal(t, lru.Len(), 3)

	lru.Add("someKey4", "90")
	assert.Equal(t, lru.Len(), 3)

	_ = lru.Remove("someKey3")
	assert.Equal(t, lru.Len(), 2)

	_ = lru.Remove("someKey2")
	assert.Equal(t, lru.Len(), 1)

	_ = lru.Remove("someKey4")
	assert.Equal(t, lru.Len(), 0)
}
