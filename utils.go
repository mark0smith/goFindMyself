package main

import (
	"fmt"
	"time"
)

// https://stackoverflow.com/a/45766707
func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %.2fs!\n", name, time.Since(start).Seconds())
	}
}
