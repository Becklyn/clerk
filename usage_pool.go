package clerk

import (
	"sync"
	"time"
)

type UsagePool[T any] struct {
	sync.Mutex
	ref       T
	counter   uint
	expiresIn time.Duration
	expiresAt time.Time
}

func NewUsagePool[T any](
	ref T,
	expiresIn time.Duration,
) *UsagePool[T] {
	return &UsagePool[T]{
		Mutex:     sync.Mutex{},
		ref:       ref,
		counter:   0,
		expiresIn: expiresIn,
		expiresAt: time.Now().Add(expiresIn),
	}
}

func (r *UsagePool[T]) Get() T {
	r.Lock()
	defer r.Unlock()

	r.counter++
	r.expiresAt = time.Now().Add(r.expiresIn)

	return r.ref
}

func (r *UsagePool[T]) Release() {
	r.Lock()
	defer r.Unlock()

	r.counter--
}

func (r *UsagePool[T]) IsUnused() bool {
	r.Lock()
	defer r.Unlock()

	return r.counter == 0 && time.Now().After(r.expiresAt)
}
