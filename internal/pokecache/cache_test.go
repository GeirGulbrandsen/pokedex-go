package pokecache

import (
	"bytes"
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	c := NewCache(1 * time.Minute)
	want := []byte("pikachu")

	c.Add("pokemon", want)
	got, ok := c.Get("pokemon")
	if !ok {
		t.Fatal("expected key to be present in cache")
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestCacheReapLoopRemovesExpiredEntries(t *testing.T) {
	interval := 20 * time.Millisecond
	c := NewCache(interval)
	c.Add("stale", []byte("value"))

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if _, ok := c.Get("stale"); !ok {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatal("expected cache entry to be reaped but it was still present")
}
