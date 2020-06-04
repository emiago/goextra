package cmap

import (
	"fmt"
	"sync"
)

type MParent struct {
	lock sync.RWMutex
}

type Map struct {
	MParent
	Items map[string]interface{}
}

func NewMap() *Map {
	return &Map{
		Items: make(map[string]interface{}),
	}
}

type CallbackMap func(key string, value interface{}) bool
type GetCallbackMap func(value interface{}, exists bool) int

func (m *MParent) Lock() {
	m.lock.Lock()
}

func (m *MParent) Unlock() {
	m.lock.Unlock()
}

func (m *MParent) RLock() {
	m.lock.RLock()
}

func (m *MParent) RUnlock() {
	m.lock.RUnlock()
}

func (m *Map) Set(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.Items[key] = value
}

func (m *Map) SetIfNotExists(key string, value interface{}) (rval interface{}, exists bool) {
	m.Lock()
	defer m.Unlock()
	rval, exists = m.Items[key]
	if !exists {
		rval = value
		m.Items[key] = rval
	}
	return
}

func (m *Map) Get(key string) (value interface{}, exists bool) {
	m.RLock()
	defer m.RUnlock()
	value, exists = m.Items[key]
	return
}

func (m *Map) Del(key string) (value interface{}, exists bool) {
	m.Lock()
	defer m.Unlock()
	value, exists = m.Items[key]
	delete(m.Items, key)
	return
}

func (m *Map) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.Items)
}

func (m *Map) GetAndUpdate(key string, cb GetCallbackMap) int {
	m.Lock()
	defer m.Unlock()
	value, exists := m.Items[key]
	action := cb(value, exists)
	switch {
	case action > 0:
		m.Items[key] = value
	case action < 0:
		delete(m.Items, key)
	}
	return action
}

func (m *Map) Loop(cb CallbackMap) {
	m.RLock()
	defer m.RUnlock()
	for k, v := range m.Items {
		m.RUnlock()
		br := cb(k, v)
		m.RLock()
		if br {
			return
		}
	}
}

func (m *Map) LoopNoLock(cb CallbackMap) {
	for k, v := range m.Items {
		br := cb(k, v)
		if br {
			return
		}
	}
}

func (m *Map) String() string {
	m.RLock()
	defer m.RUnlock()
	return fmt.Sprintf("%v", m.Items)
}

/*Concurent map with Int keys*/

type MapInt struct {
	MParent
	Items map[int]interface{}
}

func NewMapInt() *MapInt {
	return &MapInt{
		Items: make(map[int]interface{}),
	}
}

type CallbackMapInt func(key int, value interface{}) bool

func (m *MapInt) Set(key int, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.Items[key] = value
}

func (m *MapInt) Get(key int) (value interface{}, exists bool) {
	m.RLock()
	defer m.RUnlock()
	value, exists = m.Items[key]
	return
}

func (m *MapInt) Del(key int) (value interface{}, exists bool) {
	m.Lock()
	defer m.Unlock()
	value, exists = m.Items[key]
	delete(m.Items, key)
	return
}

func (m *MapInt) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.Items)
}

func (m *MapInt) Loop(cb CallbackMapInt) {
	m.RLock()
	defer m.RUnlock()
	for k, v := range m.Items {
		m.RUnlock()
		br := cb(k, v)
		m.RLock()
		if br {
			return
		}
	}
}

func (m *MapInt) String() string {
	m.RLock()
	defer m.RUnlock()
	return fmt.Sprintf("%v", m.Items)
}
