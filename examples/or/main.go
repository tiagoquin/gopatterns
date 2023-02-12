// Package main example of or
package main

import (
	"fmt"
	"gopatterns"
	"time"
)

func main() {
	emitAfter := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	<-gopatterns.Or(
		emitAfter(time.Minute),
		emitAfter(time.Second),
		emitAfter(time.Hour),
		emitAfter(12*time.Second),
	)

	fmt.Println("Done after ", time.Since(start))
}
