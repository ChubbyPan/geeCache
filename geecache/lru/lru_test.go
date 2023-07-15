package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := NewLRUCache(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key: key1, value: %v (should be 1234) failed", v)
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "a", "b", "c"
	v1, v2, v3 := "111111", "22222", "33333"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewLRUCache(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("a"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest a failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := NewLRUCache(int64(10), callback)
	lru.Add("a", String("11"))
	lru.Add("bbb", String("22222"))
	lru.Add("ccc", String("33333"))
	lru.Add("ddd", String("44444"))
	lru.Add("eee", String("55555"))

	expect := []string{"a", "bbb", "ccc", "ddd"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
