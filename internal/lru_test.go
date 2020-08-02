package internal

import (
	"testing"
)

func TestLRU(t *testing.T) {
	lru, err := NewLRU(3)

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	evicted := lru.Add("a", 1)
	if evicted {
		t.Fatalf("unexpected evicition: a -> 1")
	}

	evicted = lru.Add("b", 2)
	if evicted {
		t.Fatalf("unexpected eviction: b -> 2")
	}

	evicted = lru.Add("c", 3)
	if evicted {
		t.Fatalf("unexpected eviction: c -> 3")
	}

	evicted = lru.Add("d", 4)
	if !evicted {
		t.Fatalf("expect eviction: d -> 4")
	}

	_, ok := lru.Get("a")
	if ok {
		t.Fatalf("unexpected entry: a")
	}

	_, ok = lru.Get("b")
	if !ok {
		t.Fatalf("expect entry: b")
	}
}
