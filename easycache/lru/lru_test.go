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
	lru := Constructor(1024, nil)
	lru.Put("123", String("456"))
	if v, ok := lru.Get("123"); !ok || string(v.(String)) != "456" {
		t.Fatalf("cache hit key = 123 failed")
	}
	if _, ok := lru.Get("456"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}
func TestPut(t *testing.T) {
	lru := Constructor(2, nil)
	lru.Put("123", String("456"))
	lru.Put("456", String("789"))
	if v, ok := lru.Get("123"); ok && v.(String) != "456" || lru.Len() != 2 {
		t.Fatalf("cache 123 put failed")
	}
}

func TestOnCall(t *testing.T) {
	keys := []string{}
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := Constructor(2, callback)
	lru.Put("1", String("1"))
	lru.Put("2", String("2"))
	lru.Put("3", String("3"))
	lru.Put("4", String("4"))
	expect := []string{"1", "2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
