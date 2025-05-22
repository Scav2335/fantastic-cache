package fan_cache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
// 接口型函数，方便使用者在调用时既能够传入函数作为参数，
// 也能够传入实现了该接口的结构体作为参数。
type GetterFunc func(key string) ([]byte, error)

// Get 定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。
// 这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。
func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int, getter Getter) *Group {
	if getter == nil {
		panic("getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (v ByteView, err error) {
	if key == "" {
		err = fmt.Errorf("key is required")
		return
	}

	var ok bool
	if v, ok = g.mainCache.get(key); ok {
		log.Println("get from main cache")
		return
	}

	return g.load(key)
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

func (g *Group) getLocally(key string) (v ByteView, err error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return
	}

	v = ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, v)

	return
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}
