package clerk

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_UsagePool_Get_ReturnsTheInstance(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, 5*time.Second)

	assert.NotNil(t, pool.Get())
}

func Test_UsagePool_Get_IncrementsReferenceCounter(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	pool.Get()

	assert.Equal(t, uint(1), pool.counter)
}

func Test_UsagePool_Get_ResetsExpiryTime(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	pool.expiresAt = time.Now()
	pool.Get()
	defer pool.Release()

	assert.False(t, pool.IsUnused())
}

func Test_UsagePool_Release_DecrementsReferenceCounter(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	pool.Get()
	pool.Release()

	assert.Equal(t, uint(0), pool.counter)
}

func Test_UsagePool_Release_NegativeReferenceCountPanics(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	assert.Panics(t, func() {
		pool.Release()
	}, ErrNegativeReferenceCount)
}

func Test_UsagePool_IsUnused_ReturnsTrue_WhenReferenceCounterIsZero(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	time.Sleep(time.Millisecond)

	assert.True(t, pool.IsUnused())
}

func Test_UsagePool_IsUnused_ReturnsFalse_WhenReferenceCounterIsNotZero(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	pool.Get()

	assert.False(t, pool.IsUnused())
}

func Test_UsagePool_IsUnused_ReturnsTrue_WhenExpiryTimeIsInThePast(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	time.Sleep(time.Millisecond)

	assert.True(t, pool.IsUnused())
}

func Test_UsagePool_IsUnused_ReturnsFalse_WhenExpiryTimeIsInTheFuture(t *testing.T) {
	pool := NewUsagePool(&struct{}{}, time.Millisecond)

	assert.False(t, pool.IsUnused())
}
