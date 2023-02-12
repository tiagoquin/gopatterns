// Package example for a confinement
package main

import (
	"fmt"
	"gopatterns"
)

func main() {
	add := func(a, b int) int { return a + b }

	results := gopatterns.Confine(
		func() int { return add(2, 3) },
		func() int { return add(4, 5) },
		func() int { return add(6, 7) },
	)

	for val := range results {
		fmt.Println("Value: ", val)
	}
}
