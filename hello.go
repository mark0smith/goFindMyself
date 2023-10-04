package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
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

// recall func
func checkRecall() {

	fmt.Println("What do you remember?")

	reader := bufio.NewReader(os.Stdin)
	recallString, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// replace all whitespace characters with a single space
	reg := regexp.MustCompile(`\s+`)
	recallString = reg.ReplaceAllString(recallString, " ")

	recallString = strings.TrimSuffix(recallString, "\n")
	recallString = strings.TrimSuffix(recallString, " ")
	fmt.Printf("\nYou have entered: %s\n", recallString)

	toCheck := fmt.Sprintf("[%s]", recallString)
	content, err := os.ReadFile("./log.txt")
	if err != nil {
		panic(err)
	}
	if strings.Contains(string(content), toCheck) {
		fmt.Println("You have a correct memory!")
	} else {
		fmt.Println("Are you sure you remember it right?")
		// fmt.Printf("%s not in %s", toCheck, string(content))
	}

}

func writeInfo(filename, info string) {
	fil, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	defer fil.Close()
	fil.WriteString(info)
}

func main() {
	num := flag.Int("n", 30, "Number of Random Numbers")
	maxium := flag.Int("m", 100, "Generated number should not bigger than ?")
	unique := flag.Bool("u", true, "Should all generated number be unique?")
	remember := flag.Bool("r", false, "Should I remember all generated numbers?")
	recall := flag.Bool("recall", false, "Do you find me?")
	flag.Parse()

	if *recall {
		checkRecall()
		return
	}

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")

	uniqueRandomNumbers := generateRandomNumbers(*num, *maxium, *unique)
	fmt.Printf("%s Random Number is %d\n", datetimeFormatted, uniqueRandomNumbers)

	// type of remember is a pointer, add `*` prefix to get its value
	if *remember {
		info := fmt.Sprintf("%s %d\n", datetimeFormatted, uniqueRandomNumbers)
		filename := "./log.txt"
		writeInfo(filename, info)
	}
}
