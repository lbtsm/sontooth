package distributedCache

import (
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(100))
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemove(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := String("1234"), String("1234"), String("1234")
	max := len(k1) + len(k2) + v1.Len() + v2.Len()
	lru := New(int64(max))
	lru.Add(k1, v1)
	lru.Add(k2, v2)
	lru.Add(k3, v3)
	if lru.Len() != 2 {
		t.Fatalf("cache twoList length is not except")
	}
	lru.BianLi(lru.twoList.Front(), t)
}
