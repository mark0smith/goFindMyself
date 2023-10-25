package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// https://stackoverflow.com/a/45766707
func Timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %.2fs!\n", name, time.Since(start).Seconds())
	}
}

func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		// path/to/whatever exists
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does *not* exist
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false
	}
}

// difference returns the elements in `a` that aren't in `b`.
// https://stackoverflow.com/a/45428032
func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func WriteInfo(filename, info string) {
	fil, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	defer fil.Close()
	fil.WriteString(info)
}

// generate `max` random numbers
// if `unique` is ture, all numbers are unique
func GenerateRandomNumbers(n, max int, unique bool) []int {
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

// Insert a value into a slice at index
// https://stackoverflow.com/a/61822301
func Insert(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
