package main

import (
	"fmt"
	"math/rand"
)

func generateUniqueRandomNumbers(n, max int) []int {
	set := make(map[int]bool)
	var result []int
	for len(set) < n {
		value := rand.Intn(max)
		if !set[value] {
			set[value] = true
			result = append(result, value)
		}
	}
	return result
}

func main() {
	uniqueRandomNumbers := generateUniqueRandomNumbers(10, 100)
	fmt.Printf("Random Number is %d\n", uniqueRandomNumbers)
}
