package util

import (
	"math/rand"
	"net/http"
	"time"
)

func ShouldRetry(status int) bool {
	return status == http.StatusTooManyRequests || status >= 500
}

// added some Jitter to solve thundering herd problem
func Backoff(attempt int) time.Duration {
	base := time.Duration(1<<attempt) * 200 * time.Millisecond
	if base > 5*time.Second {
		base = 5 * time.Second
	}
	jitter := time.Duration(rand.Int63n(int64(base / 2)))
	return base + jitter
}

// func BackOff(attempt int) time.Duration {
// 	backoff := time.Duration(1<<attempt) * 200 * time.Millisecond //backoff : 200ms 400ms 800ms and so on
// 	if backoff > 5*time.Second {
// 		return 5 * time.Second // Cap it at 5 seconds
// 	}
// 	return backoff
// }
