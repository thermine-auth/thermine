package ratelimit

import (
	"testing"
	"time"
)

func TestLimiterAllowsABurstThenRefills(t *testing.T) {
	now := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)
	l := New(3)
	l.now = func() time.Time { return now }

	for i := range 3 {
		if ok, _ := l.Allow("203.0.113.7"); !ok {
			t.Fatalf("request %d refused inside the burst", i+1)
		}
	}

	ok, wait := l.Allow("203.0.113.7")
	if ok || wait <= 0 || wait > 20*time.Second {
		t.Fatalf("fourth request = %v, wait %v; want refused, about 20s", ok, wait)
	}

	if ok, _ := l.Allow("198.51.100.1"); !ok {
		t.Error("another address shares the first one's budget")
	}

	now = now.Add(20 * time.Second)
	if ok, _ := l.Allow("203.0.113.7"); !ok {
		t.Error("no token after a third of a minute at 3 a minute")
	}
}

func TestLimiterOffAtZero(t *testing.T) {
	l := New(0)
	for range 1000 {
		if ok, _ := l.Allow("x"); !ok {
			t.Fatal("a zero limit refused a request")
		}
	}
}

func TestLimiterForgetsIdleAddresses(t *testing.T) {
	now := time.Date(2026, 9, 15, 12, 0, 0, 0, time.UTC)
	l := New(60)
	l.now = func() time.Time { return now }

	l.Allow("a")
	now = now.Add(2 * time.Minute)
	l.Allow("b")

	if _, kept := l.buckets["a"]; kept {
		t.Error("an address idle long enough to refill was kept")
	}
}
