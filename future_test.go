package future_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/daishe/go-future"
)

type Result[T any] struct {
	Value   T
	Success bool
	Panic   bool
}

func GetResult[T any](f *future.Future[T], fn func(*future.Future[T]) (T, bool)) (r *Result[T]) {
	defer func() {
		if rec := recover(); rec != nil {
			r = &Result[T]{Panic: true}
			return
		}
	}()
	v, ok := fn(f)
	r = &Result[T]{Value: v, Success: ok, Panic: false}
	return
}

func (r *Result[T]) String() string {
	return fmt.Sprintf("{Value:%v,Success:%v,Panic:%v}", r.Value, r.Success, r.Panic)
}

type StartCond struct {
	start chan struct{}
}

func NewStartCond() StartCond {
	return StartCond{start: make(chan struct{})}
}

func (sc StartCond) Start() {
	close(sc.start)
}

func (sc StartCond) Wait() {
	<-sc.start
}

type Results[T any] struct {
	wg        *sync.WaitGroup
	collected []*Result[T]
}

func NewResults[T any](sc StartCond, f *future.Future[T], fns ...func(*future.Future[T]) (T, bool)) *Results[T] {
	r := &Results[T]{
		wg:        &sync.WaitGroup{},
		collected: make([]*Result[T], len(fns)),
	}

	pre := NewStartCond()
	defer pre.Start()

	for i, fn := range fns {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			pre.Wait()
			sc.Wait()
			r.collected[i] = GetResult(f, fn)
		}()
	}

	return r
}

func (r *Results[T]) Range(yield func(*Result[T]) bool) {
	r.wg.Wait()
	for _, x := range r.collected {
		if !yield(x) {
			return
		}
	}
}

func (r *Results[T]) String() string {
	strs := []string{}
	for x := range r.Range {
		strs = append(strs, x.String())
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

func Count[T any](test func(*Result[T]) bool, results ...*Results[T]) (count, totalCount int) {
	for _, res := range results {
		for r := range res.Range {
			if test(r) {
				count++
			}
			totalCount++
		}
	}
	return count, totalCount
}

func FindOne[T any](t *testing.T, test func(*Result[T]) bool, results ...*Results[T]) *Result[T] {
	t.Helper()
	c, _ := Count(test, results...)
	if c == 0 {
		t.Fatalf("no result that passes test found: %v", results)
		return nil
	}
	if c > 1 {
		t.Fatalf("more than one result that passes test found: %v", results)
		return nil
	}
	for _, res := range results {
		for r := range res.Range {
			if test(r) {
				return r
			}
		}
	}
	t.Fatalf("count returned 1, but search found no matching results: %v", results)
	return nil
}

func All[T any](test func(*Result[T]) bool, results ...*Results[T]) bool {
	c, all := Count(test, results...)
	return c == all
}

func AllMust[T any](t *testing.T, test func(*Result[T]) bool, results ...*Results[T]) {
	t.Helper()
	if !All(test, results...) {
		t.Errorf("all must pass: %v", results)
	}
}

func AllExceptOne[T any](test func(*Result[T]) bool, results ...*Results[T]) bool {
	c, all := Count(test, results...)
	return all-c == 1
}

func AllExceptOneMust[T any](t *testing.T, test func(*Result[T]) bool, results ...*Results[T]) {
	t.Helper()
	if !AllExceptOne(test, results...) {
		t.Errorf("all except one must pass: %v", results)
	}
}

func One[T any](test func(*Result[T]) bool, results ...*Results[T]) bool {
	c, _ := Count(test, results...)
	return c == 1
}

func OneMust[T any](t *testing.T, test func(*Result[T]) bool, results ...*Results[T]) {
	t.Helper()
	if !One(test, results...) {
		t.Errorf("exactly one must pass: %v", results)
	}
}

// func None[T any](test func(*Result[T]) bool, results ...*Results[T]) bool {
// 	c, _ := Count(test, results...)
// 	return c == 0
// }

// func NoneMust[T any](t *testing.T, test func(*Result[T]) bool, results ...*Results[T]) {
// 	t.Helper()
// 	if !One(test, results...) {
// 		t.Errorf("all must fail: %v", results)
// 	}
// }

func TryResolve[T any](v T) func(*future.Future[T]) (T, bool) {
	return func(f *future.Future[T]) (T, bool) {
		return v, f.TryResolve(v)
	}
}

func Resolve[T any](v T) func(*future.Future[T]) (T, bool) {
	return func(f *future.Future[T]) (T, bool) {
		f.Resolve(v)
		return v, true
	}
}

func IsDone[T any](f *future.Future[T]) (T, bool) {
	var z T
	select {
	case <-f.Done():
		return z, true
	default:
		return z, false
	}
}

func Get[T any](f *future.Future[T]) (T, bool) {
	return f.Get(), true
}

func WaitAndGet[T any](f *future.Future[T]) (T, bool) {
	f.Wait()
	return f.Get(), true
}

func IsValueEqual[T comparable](to T) func(*Result[T]) bool {
	return func(r *Result[T]) bool {
		return r.Panic == false && r.Value == to
	}
}

func IsSuccessful[T any](r *Result[T]) bool {
	return r.Panic == false && r.Success == true
}

func IsUnsuccessful[T any](r *Result[T]) bool {
	return r.Panic == false && r.Success == false
}

func IsPanic[T any](r *Result[T]) bool {
	return r.Panic == true
}

func TestFuture(t *testing.T) {
	t.Parallel()

	f := &future.Future[int]{}

	preStart, start, postStart := NewStartCond(), NewStartCond(), NewStartCond()
	preStart.Start()

	resolvers := NewResults(start, f, Resolve(11), Resolve(12), Resolve(13), Resolve(14), Resolve(15))
	tryResolvers := NewResults(start, f, TryResolve(21), TryResolve(22), TryResolve(23), TryResolve(24), TryResolve(25))

	preIsDone := NewResults(preStart, f, IsDone, IsDone, IsDone)
	preGot := NewResults(preStart, f, Get, Get, Get)
	got := NewResults(start, f, Get, Get, Get)
	gotWaited := NewResults(start, f, WaitAndGet, WaitAndGet, WaitAndGet)
	postIsDone := NewResults(postStart, f, IsDone, IsDone, IsDone)

	All(IsUnsuccessful, preIsDone)

	start.Start()

	if One(IsSuccessful, resolvers) {
		AllExceptOneMust(t, IsPanic, resolvers)
		AllMust(t, IsUnsuccessful, tryResolvers)
	} else {
		AllMust(t, IsPanic, resolvers)
		AllExceptOneMust(t, IsUnsuccessful, tryResolvers)
	}
	v := FindOne(t, IsSuccessful, resolvers, tryResolvers).Value

	AllMust(t, IsSuccessful, preGot, got, gotWaited)
	AllMust(t, IsValueEqual(v), preGot, got, gotWaited)

	postStart.Start()

	AllMust(t, IsSuccessful, postIsDone)
}

func TestResolved(t *testing.T) {
	t.Parallel()

	f := future.Resolved(1)

	preStart, start, postStart := NewStartCond(), NewStartCond(), NewStartCond()
	preStart.Start()

	resolvers := NewResults(start, f, Resolve(11), Resolve(12), Resolve(13), Resolve(14), Resolve(15))
	tryResolvers := NewResults(start, f, TryResolve(21), TryResolve(22), TryResolve(23), TryResolve(24), TryResolve(25))

	preIsDone := NewResults(preStart, f, IsDone, IsDone, IsDone)
	preGot := NewResults(preStart, f, Get, Get, Get)
	isDone := NewResults(start, f, IsDone, IsDone, IsDone)
	got := NewResults(start, f, Get, Get, Get)
	gotWaited := NewResults(start, f, WaitAndGet, WaitAndGet, WaitAndGet)
	postIsDone := NewResults(postStart, f, IsDone, IsDone, IsDone)

	AllMust(t, IsSuccessful, preIsDone)
	AllMust(t, IsSuccessful, preGot)
	AllMust(t, IsValueEqual(1), preGot)

	start.Start()

	AllMust(t, IsPanic, resolvers)
	AllMust(t, IsUnsuccessful, tryResolvers)

	AllMust(t, IsSuccessful, isDone)
	AllMust(t, IsSuccessful, preGot, got, gotWaited)
	AllMust(t, IsValueEqual(1), preGot, got, gotWaited)

	postStart.Start()

	AllMust(t, IsSuccessful, postIsDone)
}

func TestAwaitSuccessful(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	doneA, doneB, doneC := make(chan struct{}), make(chan struct{}), make(chan struct{})

	got := &future.Future[bool]{}
	go func() {
		got.Resolve(future.Await(ctx, doneA, doneB, doneC))
	}()

	if IsSuccessful(GetResult(got, IsDone)) {
		t.Errorf("await returned %v before closing all channels", got.Get())
	}

	close(doneB)
	if IsSuccessful(GetResult(got, IsDone)) {
		t.Errorf("await returned %v before closing all channels", got.Get())
	}

	close(doneA)
	if IsSuccessful(GetResult(got, IsDone)) {
		t.Errorf("await returned %v before closing all channels", got.Get())
	}

	close(doneC)
	if !got.Get() {
		t.Errorf("await returned %v after closing all channels", got.Get())
	}
}

func TestAwaitCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	doneA, doneB, doneC := make(chan struct{}), make(chan struct{}), make(chan struct{})

	got := &future.Future[bool]{}
	go func() {
		got.Resolve(future.Await(ctx, doneA, doneB, doneC))
	}()

	if IsSuccessful(GetResult(got, IsDone)) {
		t.Errorf("await returned %v before closing all channels", got.Get())
	}

	close(doneB)
	if IsSuccessful(GetResult(got, IsDone)) {
		t.Errorf("await returned %v before closing all channels", got.Get())
	}

	close(doneA)
	if IsSuccessful(GetResult(got, IsDone)) {
		t.Errorf("await returned %v before closing all channels", got.Get())
	}

	cancel()
	if got.Get() {
		t.Errorf("await returned %v after context cancel", got.Get())
	}
}
