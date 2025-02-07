# Futures (for Golang)

The go-future module provides a tiny, dependency‑free implementation of futures (single‑value promises) that works nicely with Go’s existing concurrency primitives - if you only need to produce a single value once and let any number of goroutines read it later, future is a clear, safer abstraction.

A Future[T] is a lock‑free container that can be resolved exactly once and then read from any goroutine. They behave much like a channel with a capacity of 1: they hold a single value that can be produced once and consumed by any number of goroutines, and they provide a signaling mechanism (the Done channel) that is closed when the value becomes available, mirroring the way a buffered‑1 channel becomes readable once a send occurs.

The key distinction lies in broadcast capability - while a single‑value channel can be received by only one goroutine before it becomes empty again, a future stores the result and allows any number of receivers to obtain the same value after it is resolved, enabling the possibility to effectively broadcast a single result to many listeners. Additionally, futures never need to be closed, so users are freed from implementing any special channel‑closing logic when sharing a value across multiple goroutines.

## Installing

First, use go get to install the latest version of the library.

```sh
go get -u github.com/daishe/go-future
```

Next, include go-future in your application:

```go
import "github.com/daishe/go-future"
```

## Quick start

```go
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
		// Simulate work
		time.Sleep(500 * time.Millisecond)
		f.Resolve(42) // resolve exactly once
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
		fmt.Println("Result:", f.Get()) // → Result: 42
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
```

Possible output:

```out
Timed out!
Result: 42
Future resolved to 42
```

See `example` directory for more usage examples.

## License

The project is released under the **Apache License, Version 2.0**. See the full LICENSE file for the complete terms and conditions.
