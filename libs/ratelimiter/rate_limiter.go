package ratelimiter

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrRequestLimitExceeded = errors.New("request limit exceeded")

type RateLimiter struct {
	rps    int
	ticker *time.Ticker
	active atomic.Int32
	mx     sync.Mutex
}

func NewRateLimiter(rps int) *RateLimiter {
	rateLim := RateLimiter{
		rps:    rps,
		ticker: time.NewTicker(time.Second),
	}
	rateLim.active.Store(int32(rps))

	go func() {
		for range rateLim.ticker.C {
			rateLim.active.Store(int32(rps))
		}
	}()

	return &rateLim
}

func (r *RateLimiter) Aquire() {
	r.mx.Lock()
	defer r.mx.Unlock()

	for r.active.Load() <= 0 {
	}

	r.active.Add(-1)
}
