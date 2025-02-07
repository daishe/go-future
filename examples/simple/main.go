package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/daishe/go-future"
)

func asyncOperation() *future.Future[int] {
	f := &future.Future[int]{}
	go func() {
		time.Sleep(500 * time.Millisecond) // simulate work
		f.Resolve(42)
	}()
	return f
}

func main() {
	f := asyncOperation()
	wg := &sync.WaitGroup{}

	// (1) Client that block until the value is ready
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Result:", f.Get())
	}()

	// (2) Client that react to completion via Done()
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-f.Done():
			fmt.Println("Future resolved to", f.Get())
		}
	}()

	// (3) Client that uses cancelable waiting with a context
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		if future.Await(ctx, f.Done()) {
			fmt.Println("Got result before timeout:", f.Get())
		} else {
			fmt.Println("Timed out!")
		}
	}()

	wg.Wait()
}
