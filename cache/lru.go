package cache

import "container/list"

type LruCache struct {
	maxBytes  int
	nBytes    int
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

func New(maxBytes int, onEvicted func(string, Value)) *LruCache {
	return &LruCache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (l *LruCache) Set(key string, value Value) {
	if ele, ok := l.cache[key]; ok {
		// 修改
		l.ll.MoveToBack(ele)
		e := ele.Value.(*entry)
		l.nBytes += value.Len() - e.value.Len()
		e.value = value
	} else {
		// 新增
		ele = l.ll.PushBack(&entry{key, value})
		l.cache[key] = ele
		l.nBytes += len(key) + value.Len()
	}
	for l.maxBytes > 0 && l.maxBytes < l.nBytes {
		l.RemoveOldest()
	}
}

func (l *LruCache) Get(key string) (val Value, ok bool) {
	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToBack(ele)
		return ele.Value.(*entry).value, true
	}

	return
}

// Len the number of cache entries
func (l *LruCache) Len() int {
	return l.ll.Len()
}

// RemoveOldest removes the oldest item
func (l *LruCache) RemoveOldest() {
	ele := l.ll.Front()
	if ele != nil {
		l.ll.Remove(ele)
		e := ele.Value.(*entry)
		delete(l.cache, e.key)
		l.nBytes -= len(e.key) + e.value.Len()
		if l.OnEvicted != nil {
			l.OnEvicted(e.key, e.value)
		}
	}
}

type entry struct {
	key   string // 删除元素时，方便在map中定位对应的key
	value Value
}

type Value interface {
	Len() int
}
