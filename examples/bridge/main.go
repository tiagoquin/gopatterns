// Package main example for a bridge
package main

import (
	"context"
	"fmt"
	"gopatterns"
)

func main() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range gopatterns.Bridge(context.Background(), genVals()) {
		fmt.Printf("%v ", v)
	}
}
