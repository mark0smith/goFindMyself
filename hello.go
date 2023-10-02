package main

import (
	"flag"
	"fmt"
	"math/rand"
)

// generate `max` random numbers
// if `unique` is ture, all numbers are unique
func generateRandomNumbers(n, max int, unique bool) []int {
	set := make(map[int]bool)
	var result []int
	for len(set) < n {
		value := rand.Intn(max)
		// fmt.Printf("[+] Run in %v Mode, Current Result is %d\n", unique, result)
		if unique {
			if !set[value] {
				set[value] = true
				result = append(result, value)
			}
		} else {
			set[value] = true
			result = append(result, value)

		}
		if len(result) >= n {
			break
		}

	}
	return result
}

func main() {
	num := flag.Int("n", 10, "Number of Random Numbers")
	maxium := flag.Int("m", 100, "Generated number should not bigger than ?")
	unique := flag.Bool("u", false, "Should all generated number be unique?")
	flag.Parse()

	uniqueRandomNumbers := generateRandomNumbers(*num, *maxium, *unique)
	fmt.Printf("Random Number is %d\n", uniqueRandomNumbers)
}
