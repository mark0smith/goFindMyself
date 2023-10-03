package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
)

// generate `max` random numbers
// if `unique` is ture, all numbers are unique
func generateRandomNumbers(n, max int, unique bool) []int {
	if unique && n > max {
		fmt.Printf("In Unique mode, number of random numbers (%d) should not bigger than max range (%d).", n, max)
		var result []int
		return result
	}
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
	unique := flag.Bool("u", true, "Should all generated number be unique?")
	flag.Parse()

	uniqueRandomNumbers := generateRandomNumbers(*num, *maxium, *unique)
	fmt.Printf("Random Number is %d\n", uniqueRandomNumbers)
	
	fil, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0420)
        if err != nil {
                fmt.Printf("Error: %s\n", err)
                return
        }
        defer fil.Close()
        fil.WriteString("Hello World!\n")
}
