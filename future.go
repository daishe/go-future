// Package future implements futures - simple value wrappers that integrate seamlessly with other Go's concurrency primitives, allowing easy retrieval of the results of an asynchronous operation when it becomes available.
package future

import (
	"context"
	"sync/atomic"
)

// Future is a wrapper that allows return a result of an asynchronous operation at some point in the future.
//
// Futures are similar to channels with capacity of 1, with a notable difference that futures cannot be closed (unlike channels) and they store value, making them easier to use for single value broadcasts.
type Future[T any] struct {
	vp atomic.Pointer[T]             // value pointer
	dp atomic.Pointer[chan struct{}] // done pointer
}

// Resolved creates a new future that is already resolved with the provided value.
func Resolved[T any](v T) *Future[T] {
	f := &Future[T]{}
	f.vp.Store(&v)
	d := make(chan struct{})
	close(d)
	f.dp.Store(&d)
	return f
}

func (f *Future[T]) done() chan struct{} {
	if dp := f.dp.Load(); dp != nil {
		return *dp
	}
	d := make(chan struct{})
	f.dp.CompareAndSwap(nil, &d)
	return *f.dp.Load()
}

// Resolve resolves the future with the provided value. It panics if the future was already resolved.
func (f *Future[T]) Resolve(v T) {
	if !f.TryResolve(v) {
		panic("future: already resolved")
	}
}

// TryResolve attempts to resolve the given future with the provided value. It returns false if the future was already resolved, otherwise it resolves it with the provided value and returns true.
func (f *Future[T]) TryResolve(v T) bool {
	if !f.vp.CompareAndSwap(nil, &v) {
		return false
	}
	close(f.done())
	return true
}

// Get awaits for the resolvement of the given future and returns its value.
func (f *Future[T]) Get() T {
	if vp := f.vp.Load(); vp != nil {
		return *vp
	}
	f.Wait()
	return *f.vp.Load()
}

// Wait awaits for the resolvement of the given future.
func (f *Future[T]) Wait() {
	<-f.done()
}

// Done returns channel that will be closed when the given future is resolved.
func (f *Future[T]) Done() <-chan struct{} {
	return f.done()
}

// Await waits for either the given context to be cancelled - in which case the function returns false - or for all of the supplied struct channels to be closed - in which case the function returns true.
func Await(ctx context.Context, chs ...<-chan struct{}) bool {
	if len(chs) == 0 {
		return ctx.Err() == nil
	}
	select {
	case <-ctx.Done():
		return false
	case <-chs[0]:
		return Await(ctx, chs[1:]...)
	}
}
