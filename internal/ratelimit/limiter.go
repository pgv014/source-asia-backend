package ratelimit

import (
	"sync"
	"time"
)

type UserStats struct {
	Accepted int `json:"accepted_requests"`
	Rejected int `json:"rejected_requests"`
}

type Limiter struct {
	mu sync.Mutex

	requests map[string][]time.Time
	rejected map[string]int
}

func NewLimiter() *Limiter {
	return &Limiter{
		requests: make(map[string][]time.Time),
		rejected: make(map[string]int),
	}
}

func (l *Limiter) Allow(userID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	windowStart := now.Add(-1 * time.Minute)

	validRequests := []time.Time{}

	for _, t := range l.requests[userID] {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	l.requests[userID] = validRequests

	if len(validRequests) >= 5 {
		l.rejected[userID]++
		return false
	}

	l.requests[userID] = append(l.requests[userID], now)

	return true
}

func (l *Limiter) GetStats() map[string]UserStats {
	l.mu.Lock()
	defer l.mu.Unlock()

	stats := make(map[string]UserStats)

	for userID, requests := range l.requests {
		stats[userID] = UserStats{
			Accepted: len(requests),
			Rejected: l.rejected[userID],
		}
	}

	return stats
}
