package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
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
	num := flag.Int("n", 30, "Number of Random Numbers")
	maxium := flag.Int("m", 100, "Generated number should not bigger than ?")
	unique := flag.Bool("u", true, "Should all generated number be unique?")
	remember := flag.Bool("r", false, "Should I remember all generated numbers?")
	flag.Parse()

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")

	uniqueRandomNumbers := generateRandomNumbers(*num, *maxium, *unique)
	fmt.Printf("%s Random Number is %d\n", datetimeFormatted, uniqueRandomNumbers)

	// type of remember is a pointer, add `*` prefix to get its value
	if *remember {

		fil, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		defer fil.Close()
		info := fmt.Sprintf("%s %d\n", datetimeFormatted, uniqueRandomNumbers)
		fil.WriteString(info)
	}
}
