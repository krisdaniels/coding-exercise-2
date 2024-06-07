package orderedmap

import "sync"

type OrderedMap interface {
	Set(key, value string)
	Get(key string) string
	Delete(key string)
	Iterate(iter func(key, value string))
}

type orderedMapItem struct {
	Prev  *orderedMapItem
	Next  *orderedMapItem
	Key   string
	Value string
}

var _ OrderedMap = (*orderedMap)(nil)

type orderedMap struct {
	kv    map[string]*orderedMapItem
	lock  sync.RWMutex
	first *orderedMapItem
	last  *orderedMapItem
}

func NewOrdererMap() OrderedMap {
	return &orderedMap{
		kv: map[string]*orderedMapItem{},
	}
}

func (m *orderedMap) Set(key, value string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	v, ok := m.kv[key]
	if ok {
		v.Value = value
		return
	}

	item := &orderedMapItem{
		Key:   key,
		Value: value,
	}

	if m.first == nil {
		m.first = item
	} else {
		item.Prev = m.last
	}

	if m.last != nil {
		m.last.Next = item
	}

	m.last = item
	m.kv[key] = item
}

func (m *orderedMap) Get(key string) string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	item, ok := m.kv[key]
	if ok {
		return item.Value
	}
	return ""
}

func (m *orderedMap) Delete(key string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.kv[key]
	if ok {
		if v == m.first {
			m.first = v.Next
			m.first.Prev = nil
		} else if v == m.last {
			m.last = v.Prev
			m.last.Next = nil
		} else {
			v.Prev.Next = v.Next
			v.Next.Prev = v.Prev
		}

		delete(m.kv, key)
	}
}

func (m *orderedMap) Iterate(iter func(key, value string)) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if m.first == nil {
		return
	}

	for n := m.first; n != nil; n = n.Next {
		iter(n.Key, n.Value)
	}
}
