package singleflight

import "sync"

//  an in-flight or completed Do call
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
	if c, ok := g.m[key]; ok {
		g.mtx.Unlock()
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
