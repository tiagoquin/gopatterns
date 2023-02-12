// Package gopatterns implements various concurrency patterns
// inspired by https://www.oreilly.com/library/view/concurrency-in-go/9781491941294/ch04.html
// with generics
package gopatterns

import (
	"context"
	"sync"
)

// Confine executes sequentially N functions.
// Closes after executing them.
// Retrieve the returned values in results
func Confine[T any](fns ...func() T) <-chan T {
	results := make(chan T, len(fns))

	go func() {
		defer close(results)

		for _, do := range fns {
			results <- do()
		}
	}()
	return results
}

// Repeat generates indefinitely the values passed into channel
func Repeat[T any](ctx context.Context, values ...T) <-chan T {
	stream := make(chan T)

	go func() {
		defer close(stream)

		for {
			for _, v := range values {
				select {
				case <-ctx.Done():
					return
				case stream <- v:
				}
			}
		}
	}()
	return stream
}

// Take takes num elements from the channel stream.
// Forwards them into the returned channel
func Take[T any](ctx context.Context, stream <-chan T, num int) <-chan T {
	taken := make(chan T)

	go func() {
		defer close(taken)

		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done():
				return
			case taken <- <-stream:
			}
		}
	}()

	return taken
}

// FanIn drains in parallel a number of channels [multiplex].
// Returns the values in a single channel
func FanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	wg.Add(len(channels))

	multiplexed := make(chan T)

	drain := func(c <-chan T) {
		defer wg.Done()
		for val := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexed <- val:
			}
		}
	}

	for _, c := range channels {
		go drain(c)
	}

	go func() {
		wg.Wait()
		close(multiplexed)
	}()

	return multiplexed
}

// OrDone either forwards the values of c, or quits if ctx is done.
func OrDone[T any](ctx context.Context, c <-chan T) <-chan T {
	values := make(chan T)
	go func() {
		defer close(values)
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-c:
				if !ok {
					return
				}

				select {
				case values <- v:
				case <-ctx.Done():
				}
			}
		}
	}()

	return values
}

// Tee receives T values from in and forwards them to left and right channels
func Tee[T any](ctx context.Context, in <-chan T) (_, _ <-chan T) {
	left := make(chan T)
	right := make(chan T)
	go func() {
		defer close(left)
		defer close(right)

		for val := range OrDone(ctx, in) {
			var out1, out2 = left, right

			for i := 0; i < 2; i++ {
				select {
				case <-ctx.Done():
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return left, right
}

// Bridge takes a sequence of channels and bridges them into a single channel
func Bridge[T any](ctx context.Context, streams <-chan <-chan T) <-chan T {
	values := make(chan T)
	go func() {
		defer close(values)
		for {
			var stream <-chan T

			select {
			case maybe, ok := <-streams:
				if !ok {
					return
				}
				stream = maybe
			case <-ctx.Done():
				return
			}

			for val := range OrDone(ctx, stream) {
				select {
				case values <- val:
				case <-ctx.Done():
				}
			}
		}
	}()

	return values
}

// Or will close the returned channel once any of the given channels closes or sends something.
// It allows to combine multiple close channels
func Or[T any](channels ...<-chan T) <-chan T {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan T)
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()

	return orDone
}
