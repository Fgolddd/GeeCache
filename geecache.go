package geecache

import (
	"fmt"
	"log"
	"sync"
)

// the main data structure of a cache
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

// a Getter loads data for a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// a GetterFunc implements Getter with a function
type GetterFunc func(ket string) ([]byte, error)

var (
	mtx    sync.RWMutex
	groups = make(map[string]*Group)
)

// create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mtx.Lock()
	defer mtx.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mtx.RLock()
	defer mtx.RUnlock()

	g := groups[name]
	return g
}

// get value by a key from cache
// func (g *Group) Get(key string) (ByteView, error) {
// 	if key == "" {
// 		return ByteView{}, fmt.Errorf("key is required")
// 	}

// 	if v, ok := g.mainCache.get(key); ok {
// 		log.Println("[GeeCache] hint")
// 		return v, nil
// 	}

// 	return g.load(key)
// }

// func (g *Group) load(key string) (value ByteView, err error) {
// 	return g.getLocally(key)
// }

// func (g *Group) getLocally(key string) (ByteView, error) {
// 	bytes, err := g.getter.Get(key)
// 	if err != nil {
// 		return ByteView{}, err
// 	}

// 	value := ByteView{b: cloneBytes(bytes)}
// 	g.populateCache(key, value)

// 	return value, nil
// }

// func (g *Group) populateCache(key string, value ByteView) {
// 	g.mainCache.add(key, value)
// }

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
