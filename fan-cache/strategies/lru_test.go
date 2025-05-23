package strategies

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestLruCache_Get(t *testing.T) {
	lru := NewLruCache(0, nil)
	lru.Set("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || v.(String) != "1234" {
		t.Fatalf("ache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("strategies miss key2 failed")
	}
}

func TestLruCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	lru := NewLruCache(len(k1+k2+v1+v2), nil)
	lru.Set(k1, String(v1))
	lru.Set(k2, String(v2))
	lru.Set(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := NewLruCache(10, callback)
	lru.Set("key1", String("123456"))
	lru.Set("k2", String("k2"))
	lru.Set("k3", String("k3"))
	lru.Set("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
