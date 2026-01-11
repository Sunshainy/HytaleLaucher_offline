// Package eventgroup provides utilities for executing groups of tasks
// with progress reporting. It supports both sequential execution with
// progress callbacks and parallel execution similar to errgroup.
package eventgroup

import (
	"sync"
)

// Task represents a single task that can be executed.
type Task func() error

// ProgressCallback is a function that receives progress updates.
// eventName is the name of the event being executed.
// metadata contains additional context about the event.
// progress is a value between 0.0 and 1.0 indicating completion.
type ProgressCallback func(eventName string, metadata map[string]any, progress float64)

// Group holds a collection of tasks to be executed with progress reporting
// or in parallel. It supports two modes of operation:
//
// 1. Sequential with progress: Use Add() to add tasks, then Exec() to run them
//    sequentially with progress reporting.
//
// 2. Parallel execution: Use Go() to spawn goroutines that run concurrently,
//    then Wait() to wait for all of them to complete.
type Group struct {
	tasks []Task

	// For parallel execution (Go/Wait pattern)
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	mu      sync.Mutex
}

// New creates a new Group with the given tasks.
func New(tasks ...Task) *Group {
	return &Group{tasks: tasks}
}

// Add appends a task to the group for sequential execution.
func (g *Group) Add(task Task) {
	g.tasks = append(g.tasks, task)
}

// Len returns the number of tasks in the group.
func (g *Group) Len() int {
	return len(g.tasks)
}

// Go spawns a new goroutine that executes the given function.
// The first call to return a non-nil error will be recorded and
// returned by Wait. Go is safe to call from multiple goroutines.
//
// This is similar to golang.org/x/sync/errgroup but does not
// cancel other goroutines when one returns an error.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.mu.Lock()
				g.err = err
				g.mu.Unlock()
			})
		}
	}()
}

// Wait blocks until all goroutines spawned with Go have completed.
// It returns the first non-nil error encountered by any goroutine,
// or nil if all goroutines completed successfully.
func (g *Group) Wait() error {
	g.wg.Wait()
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.err
}

// Exec executes all tasks in the group sequentially, reporting progress
// through the callback. The progress is scaled and offset according to
// the provided parameters.
//
// Parameters:
//   - callback: Function to receive progress updates
//   - eventName: Name of the event for progress reporting
//   - metadata: Additional context for progress updates
//   - scale: Multiplier for progress values (e.g., 0.5 means tasks represent 50% of total)
//   - offset: Starting offset for progress (e.g., 0.5 means start at 50%)
//
// Returns the first error encountered, or nil if all tasks complete successfully.
func (g *Group) Exec(callback ProgressCallback, eventName string, metadata map[string]any, scale, offset float64) error {
	if len(g.tasks) == 0 {
		return nil
	}

	// Calculate progress per task
	progressPerTask := scale / float64(len(g.tasks))

	// Report initial progress
	callback(eventName, metadata, offset)

	// Execute each task
	for i, task := range g.tasks {
		if err := task(); err != nil {
			return err
		}

		// Report progress after each task
		progress := offset + (float64(i+1) * progressPerTask)
		callback(eventName, metadata, progress)
	}

	// Report final progress
	callback(eventName, metadata, offset+scale)

	return nil
}

// ExecSimple executes all tasks sequentially without progress scaling.
// Progress is reported as a value from 0.0 to 1.0.
func (g *Group) ExecSimple(callback ProgressCallback, eventName string, metadata map[string]any) error {
	return g.Exec(callback, eventName, metadata, 1.0, 0.0)
}
