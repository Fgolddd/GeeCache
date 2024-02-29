package singleflight

import (
	"log"
	"os"
	"sync"
)

// an in-flight or completed Do call
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mtx sync.Mutex
	m   map[string]*call
}

// function fn only will be called once
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mtx.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//if goroutine is requiring
	if c, ok := g.m[key]; ok {
		//release the lock
		g.mtx.Unlock()

		//waiting return
		euid := os.Geteuid()
		log.Printf("requiring------, waiting euid: %d return", euid)
		c.wg.Wait()

		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mtx.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mtx.Lock()
	delete(g.m, key)
	g.mtx.Unlock()

	return c.val, c.err
}
