package utils

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

type CircuitBreaker struct {
	name         string
	maxFailures  int
	resetTimeout time.Duration

	mu            sync.RWMutex
	failures      int
	lastFailTime  time.Time
	state         State
	successCount  int
	halfOpenLimit int
}

func NewCircuitBreaker(name string, maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		name:          name,
		maxFailures:   maxFailures,
		resetTimeout:  resetTimeout,
		state:         StateClosed,
		halfOpenLimit: 3, // Need 3 successful calls to close from half-open
	}
}

func (cb *CircuitBreaker) Call(fn func() error) error {
	cb.mu.Lock()

	// Check if circuit should transition from open to half-open
	if cb.state == StateOpen && time.Since(cb.lastFailTime) > cb.resetTimeout {
		cb.state = StateHalfOpen
		cb.successCount = 0
	}

	state := cb.state
	cb.mu.Unlock()

	if state == StateOpen {
		return ErrCircuitOpen
	}

	// Execute the function
	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.lastFailTime = time.Now()
		cb.successCount = 0

		if cb.failures >= cb.maxFailures {
			cb.state = StateOpen
		}
		return err
	}

	// Success handling
	if cb.state == StateHalfOpen {
		cb.successCount++
		if cb.successCount >= cb.halfOpenLimit {
			cb.state = StateClosed
			cb.failures = 0
		}
	} else if cb.state == StateClosed {
		cb.failures = 0
	}

	return nil
}

func (cb *CircuitBreaker) GetState() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.successCount = 0
}
