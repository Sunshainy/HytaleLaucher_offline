package throttle

import (
	"sync"
	"time"
)

// State is a generic thread-safe container for caching values with error tracking.
// It stores a value of type T along with the time it was last updated and any error
// that occurred during the last update.
//
// The State type is typically used with pointer types (e.g., *SomeStruct) to cache
// the results of expensive operations like network requests.
type State[T any] struct {
	mu        sync.RWMutex
	value     T
	updatedAt time.Time
	err       error
}

// Get returns the currently cached value and any error from the last update.
// It acquires a read lock to ensure thread-safe access.
func (s *State[T]) Get() (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value, s.err
}

// Set stores a new value and clears any previous error.
// It acquires a write lock to ensure thread-safe access.
func (s *State[T]) Set(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = value
	s.updatedAt = time.Now()
	s.err = nil
}

// SetError stores an error from a failed update attempt.
// The previous value is preserved.
// It acquires a write lock to ensure thread-safe access.
func (s *State[T]) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.updatedAt = time.Now()
	s.err = err
}

// Update atomically sets both the value and error.
// This is useful when an operation may succeed or fail.
// It acquires a write lock to ensure thread-safe access.
func (s *State[T]) Update(value T, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value = value
	s.updatedAt = time.Now()
	s.err = err
}

// UpdatedAt returns the time of the last update.
// It acquires a read lock to ensure thread-safe access.
func (s *State[T]) UpdatedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updatedAt
}

// HasValue returns true if there is a cached value (updated at least once without error).
// For pointer types, this also checks if the value is non-nil.
func (s *State[T]) HasValue() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.updatedAt.IsZero() && s.err == nil
}

// Error returns any error from the last update attempt.
// It acquires a read lock to ensure thread-safe access.
func (s *State[T]) Error() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.err
}
