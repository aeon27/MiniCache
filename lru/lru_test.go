package lru

import (
	"fmt"
	"reflect"
	"testing"
)

type Element string

func (ele Element) Len() int {
	return len(ele)
}

func TestGet(t *testing.T) {
	lru := NewCache(int64(0), nil)
	lru.Add("key1", Element("value1"))

	fmt.Printf("cache information: %+v", lru.cache)

	if v, ok := lru.Get("key1"); !ok || v.(Element) != "value1" {
		t.Fatalf("cache hit [key1 = value1] failed!")
		t.Fatalf("cache: %+v", lru.cache)
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed!")
		t.Fatalf("cache information: %+v", lru.cache)
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap), nil)
	lru.Add(k1, Element(v1))
	lru.Add(k2, Element(v2))
	lru.Add(k3, Element(v3))

	fmt.Printf("cache: %+v", lru.cache)

	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("test RemoveOldest failed!\n")
		t.Fatalf("cache information: %+v", lru.cache)
	}
}

func TestOnEvicted(t *testing.T) {
	evictedKeys := make([]string, 0)
	k1, k2, k3, k4 := "key1", "key2", "key3", "key4"
	v1, v2, v3, v4 := "value1", "value2", "value3", "value4"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap), func(key string, value Value) {
		evictedKeys = append(evictedKeys, key)
		fmt.Printf("func OnEnvicted executed, key = %s, value = %v\n", key, value)
	})
	lru.Add(k1, Element(v1))
	lru.Add(k2, Element(v2))
	lru.Add(k3, Element(v3))
	lru.Update(k2, Element("val2-2"))
	lru.Add(k4, Element(v4))

	expected := []string{"key1", "key3"}

	if !reflect.DeepEqual(expected, evictedKeys) {
		t.Fatalf("Call OnEvicted failed! Expected keys are: %s", expected)
	}
}
