package ratelimiter

import (
	"sync"
	"time"
)

type Entry struct {
	count   int
	resetAt time.Time
}
type RateLimiter struct {
	mu     sync.Mutex
	limits map[string]*Entry
	window time.Duration
	max    int
}

func New(max int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limits: make(map[string]*Entry),
		max:    max,
		window: window,
	}
}

func (r *RateLimiter) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for ip, entry := range r.limits {
		if time.Now().After(entry.resetAt) {
			delete(r.limits, ip)
		}
	}
}
func (r *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, exists := r.limits[ip]
	if !exists || now.After(entry.resetAt) {
		r.limits[ip] = &Entry{
			count:   1,
			resetAt: now.Add(r.window),
		}
		return true
	}
	if entry.count >= r.max {
		return false
	}

	entry.count++
	return true
}
