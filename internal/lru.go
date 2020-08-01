package internal

import (
	"container/list"
	"errors"
)

type LRU struct {
	size      int
	evictList *list.List
	elements  map[interface{}]*list.Element
}

type entry struct {
	key   interface{}
	value interface{}
}

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

func (u *LRU) Remove(key interface{}) bool {
	if element, ok := u.elements[key]; ok {
		u.evictList.Remove(element)
		kv := element.Value.(*entry)
		delete(u.elements, kv.key)
		return ok
	}
	return false
}

func (u *LRU) Contains(key interface{}) bool {
	_, ok := u.elements[key]
	return ok
}

func (u *LRU) removeOldest() {
	element := u.evictList.Back()
	if element != nil {
		u.evictList.Remove(element)
	}
}