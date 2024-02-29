package fifo

import (
	"container/list"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func Test_fifoCahce_Get(t *testing.T) {
	cache := &Cache{
		maxBytes:  15,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: nil,
	}

	cache.Add("key1", String("1234"))
	if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := cache.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
