package server

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// rateLimiter is a simple in-memory IP-based rate limiter using a sliding window.
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string][]time.Time
	max      int           // max requests per window
	window   time.Duration // sliding window duration
}

var authRateLimiter = &rateLimiter{
	visitors: make(map[string][]time.Time),
	max:      10,              // 10 attempts per window
	window:   1 * time.Minute, // per minute
}

// allow checks if the IP is within the rate limit. Returns true if allowed.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Filter to requests within the window
	var recent []time.Time
	for _, t := range rl.visitors[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= rl.max {
		rl.visitors[ip] = recent
		return false
	}

	recent = append(recent, now)
	rl.visitors[ip] = recent

	// Periodic cleanup of stale IPs
	if len(rl.visitors) > 1000 {
		for k, v := range rl.visitors {
			if len(v) == 0 || v[len(v)-1].Before(cutoff) {
				delete(rl.visitors, k)
			}
		}
	}

	return true
}

// clientIP extracts the client IP from the request, respecting X-Forwarded-For.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// rateLimitMiddleware wraps an http.HandlerFunc with IP-based rate limiting.
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !authRateLimiter.allow(clientIP(r)) {
			writeJSON(w, 429, jsonResp{Error: "Too many requests. Please try again later."})
			return
		}
		next(w, r)
	}
}
