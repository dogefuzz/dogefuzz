package common

import "sync"

// Reference: https://levelup.gitconnected.com/implementation-of-thread-safe-dictionary-data-structure-in-golang-2bcb235fd9e4
type IKey = interface{}
type IValue = interface{}

type ConcurrentMap struct {
	items map[IKey]IValue
	lock  sync.RWMutex
}

func (d *ConcurrentMap) Add(key IKey, value IValue) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.items == nil {
		d.items = make(map[IKey]IValue)
	}
	d.items[key] = value
}

func (d *ConcurrentMap) Remove(key IKey) bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	_, ok := d.items[key]
	if ok {
		delete(d.items, key)
	}
	return ok
}

func (d *ConcurrentMap) Get(key IKey) IValue {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.items[key]
}

func (d *ConcurrentMap) Exist(key IKey) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
	_, ok := d.items[key]
	return ok
}

func (d *ConcurrentMap) Clear() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.items = make(map[IKey]IValue)
}

func (d *ConcurrentMap) Size() int {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return len(d.items)
}

func (d *ConcurrentMap) GetKeys() []IKey {
	d.lock.RLock()
	defer d.lock.RUnlock()
	keys := []IKey{}
	for i := range d.items {
		keys = append(keys, i)
	}
	return keys
}

func (d *ConcurrentMap) GetValues() []IValue {
	d.lock.RLock()
	defer d.lock.RUnlock()
	values := []IValue{}
	for i := range d.items {
		values = append(values, d.items[i])
	}
	return values
}
