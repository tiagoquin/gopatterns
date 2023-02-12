// Package main example for a basic pipeline
package main

import (
	"context"
	"fmt"
	"gopatterns"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for val := range gopatterns.Take(ctx, gopatterns.Repeat(ctx, 1, 2, 3), 42) {
		fmt.Print(val, " ")
	}
	fmt.Println("Done")
}
