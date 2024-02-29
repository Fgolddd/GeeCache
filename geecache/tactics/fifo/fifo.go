package fifo

import "container/list"

type Cache struct {
	maxBytes  int64
	nbytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//更新缓存
		c.nbytes += int64(value.Len()) - int64(ele.Value.(*entry).value.Len())
		ele.Value.(*entry).value = value
	} else {
		kv := &entry{key, value}
		ele := c.ll.PushBack(kv)
		c.cache[key] = ele
		c.nbytes += int64(len(kv.key)) + int64(kv.value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveFront()
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		e := ele.Value.(*entry)
		return e.value, true
	}
	return

}

func (c *Cache) RemoveFront() {
	ele := c.ll.Front()
	if ele != nil {
		kv := c.ll.Remove(ele).(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
