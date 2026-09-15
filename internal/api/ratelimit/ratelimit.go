// Package ratelimit slows down guessing: it limits how often one address may
// call the endpoints that take a password or send an email.
//
// The per-account lockout already stops guessing one account's password. This
// is the other half: one address trying many accounts, creating accounts in
// bulk, or filling someone's inbox with reset links. It counts in memory, per
// server process, so with several servers each keeps its own count — a proxy
// or load balancer in front can enforce one limit across all of them.
package ratelimit

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Limiter is a token bucket per client address: `perMinute` requests a
// minute, refilled continuously, with the whole minute's worth available at
// once.
type Limiter struct {
	rate  float64 // tokens per second
	burst float64

	mu      sync.Mutex
	buckets map[string]*bucket
	swept   time.Time

	now func() time.Time
}

type bucket struct {
	tokens float64
	seen   time.Time
}

// New returns a limiter allowing `perMinute` requests a minute per address.
// Zero or less allows everything.
func New(perMinute int) *Limiter {
	return &Limiter{
		rate:    float64(perMinute) / 60,
		burst:   float64(perMinute),
		buckets: map[string]*bucket{},
		now:     time.Now,
	}
}

// Allow takes a token for `key`, and says whether there was one and, when not,
// how long until there is.
func (l *Limiter) Allow(key string) (bool, time.Duration) {
	if l.burst <= 0 {
		return true, 0
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.sweep(now)

	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst, seen: now}
		l.buckets[key] = b
	}

	b.tokens = min(l.burst, b.tokens+now.Sub(b.seen).Seconds()*l.rate)
	b.seen = now

	if b.tokens < 1 {
		return false, time.Duration((1 - b.tokens) / l.rate * float64(time.Second))
	}

	b.tokens--
	return true, 0
}

// sweep forgets addresses whose bucket has refilled, so the map does not grow
// with every address that ever called. It runs at most once a minute.
func (l *Limiter) sweep(now time.Time) {
	if now.Sub(l.swept) < time.Minute {
		return
	}
	l.swept = now

	full := time.Duration(l.burst / l.rate * float64(time.Second))
	for key, b := range l.buckets {
		if now.Sub(b.seen) > full {
			delete(l.buckets, key)
		}
	}
}

// Middleware refuses a request over the limit with 429 and Retry-After. The
// address is Gin's ClientIP, so behind a proxy XERMESS_TRUSTED_PROXIES has to
// name it — otherwise every client shares the proxy's one budget.
func (l *Limiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, wait := l.Allow(c.ClientIP())
		if !ok {
			seconds := int(wait/time.Second) + 1
			c.Header("Retry-After", strconv.Itoa(seconds))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many attempts; try again in " + strconv.Itoa(seconds) + " seconds"})
			return
		}

		c.Next()
	}
}
