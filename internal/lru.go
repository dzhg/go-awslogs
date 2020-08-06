package internal

import (
	"container/list"
	"errors"
)

// LRU is a simplified version of golang-lru from hashicorp
// see: https://github.com/hashicorp/golang-lru
type LRU struct {
	size      int
	evictList *list.List
	elements  map[interface{}]*list.Element
}

type entry struct {
	key   interface{}
	value interface{}
}

// NewLRU creates a new LRU instance
func NewLRU(size int) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("size must be positive")
	}

	return &LRU{
		size:      size,
		evictList: list.New(),
		elements:  make(map[interface{}]*list.Element),
	}, nil
}

// Add adds a new entry to the LRU instance
func (u *LRU) Add(key, value interface{}) (evicted bool) {

	if e, ok := u.elements[key]; ok {
		// exists
		u.evictList.MoveToFront(e)
		e.Value.(*entry).value = value
		return false
	}

	newEntry := &entry{key, value}
	element := u.evictList.PushFront(newEntry)
	u.elements[key] = element

	evict := u.evictList.Len() > u.size
	if evict {
		u.removeOldest()
	}
	return evict
}

// Get returns the entry by the given key
func (u *LRU) Get(key interface{}) (value interface{}, ok bool) {
	if element, ok := u.elements[key]; ok {
		u.evictList.MoveToFront(element)
		if element.Value.(*entry) == nil {
			return nil, false
		}
		return element.Value.(*entry).value, true
	}

	return nil, false
}

// Remove deletes the key from the LRU instance
func (u *LRU) Remove(key interface{}) bool {
	if element, ok := u.elements[key]; ok {
		u.evictList.Remove(element)
		kv := element.Value.(*entry)
		delete(u.elements, kv.key)
		return ok
	}
	return false
}

// Contains check if the key is in the LRU instance
func (u *LRU) Contains(key interface{}) bool {
	_, ok := u.elements[key]
	return ok
}

func (u *LRU) removeOldest() {
	element := u.evictList.Back()
	if element != nil {
		u.evictList.Remove(element)
		delete(u.elements, element.Value.(*entry).key)
	}
}
