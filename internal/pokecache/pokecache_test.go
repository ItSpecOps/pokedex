package pokecache

import (
    "testing"
    "time"
)

func TestCacheAddGet(t *testing.T) {
    cache := NewCache(2 * time.Second)
    key := "testkey"
    val := []byte("testval")
    cache.Add(key, val)
    got, ok := cache.Get(key)
    if !ok {
        t.Fatalf("expected key to be present")
    }
    if string(got) != string(val) {
        t.Fatalf("expected %s, got %s", val, got)
    }
    // Wait for entry to expire
    time.Sleep(3 * time.Second)
    _, ok = cache.Get(key)
    if ok {
        t.Fatalf("expected key to be expired and removed")
    }
}
