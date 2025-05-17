package gpool

import "time"

// option pattern
type Option func(*pool)

// capacity of the pool
func WithCapacity(capacity int32) Option {
	return func(gp *pool) {
		gp.capacity = capacity
	}
}

// interval for cleaning up workers
func WithInterval(interval time.Duration) Option {
	return func(gp *pool) {
		gp.interval = interval
	}
}
