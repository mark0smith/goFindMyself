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
func checkRecall(rememberLogfile string, recallLogfile string, recallLog bool, showhint int) {

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
	content, err := os.ReadFile(rememberLogfile)
	if err != nil {
		panic(err)
	}

	recallResult := strings.Contains(string(content), toCheck)
	if recallResult {
		fmt.Println("You have a correct memory!")
	} else {
		fmt.Println("Are you sure you remember it right?")
		// fmt.Printf("%s not in %s", toCheck, string(content))
		if showhint > 0 {
			// compare two strings and assume first few numbers (min of 2 and slice lenth) is correct
			recallSlice := strings.Split(recallString, " ")
			firstNumCount := min(len(recallSlice), 2)
			correctRegStr := fmt.Sprintf(`\[%s [\w ]+\]`, strings.Join(recallSlice[:firstNumCount], " "))
			reg := regexp.MustCompile(correctRegStr)
			correctStr := reg.FindString(string(content))
			correctStr = strings.TrimPrefix(correctStr, "[")
			correctStr = strings.TrimSuffix(correctStr, "]")

			info := "\nHint Part:\n"
			if len(correctStr) < 1 {
				info += fmt.Sprintf("Totally Wrong! Don't you even remember the first %d number(s)?", firstNumCount)
			} else {
				if showhint == 1 {
					correctSlice := strings.Split(correctStr, " ")
					missingNumbers := difference(correctSlice, recallSlice)
					if len(missingNumbers) > 5 {
						info += fmt.Sprintf("You are missing %d numbers, which is too many for hinting. You should remember it again!\n", len(missingNumbers))
					} else {
						info += fmt.Sprintf("You are missing these numbers: %s\n", strings.Join(missingNumbers, " "))
						wrongNumbers := difference(recallSlice, correctSlice)
						info += fmt.Sprintf("You add these numbers which should't exist: %s\n", strings.Join(wrongNumbers, " "))
					}

				} else if showhint == 2 {
					info += fmt.Sprintf("The Right: %s\nThe Wrong: %s", correctStr, recallString)
				}
			}
			fmt.Printf("%s", info)
		}
	}

	if recallLog {
		datetime := time.Now()
		datetimeFormatted := datetime.Format("2006-01-02 15:04:05")
		info := fmt.Sprintf("%s Recall: %s, Result: %v\n", datetimeFormatted, recallString, recallResult)
		writeInfo(recallLogfile, info)
	}
}

// difference returns the elements in `a` that aren't in `b`.
// https://stackoverflow.com/a/45428032
func difference(a, b []string) []string {
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

func writeInfo(filename, info string) {
	fil, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	defer fil.Close()
	fil.WriteString(info)
}

func main() {
	num := flag.Int("n", 30, "Number of Random Numbers.")
	maxium := flag.Int("m", 100, "Generated number wont't be bigger than this number.")
	unique := flag.Bool("u", true, "If set, all generated numbers will be unique.")
	remember := flag.Bool("r", false, "If set, generated numbers will be logged into `rememberLogfile`.\nYou should set this if you want to do recall test later!")
	rememberLogfile := flag.String("r_file", "log.txt", "Filename of remember log")
	recall := flag.Bool("recall", false, "If set, run a recall test, instead of generating random numbers.")
	recallLog := flag.Bool("recallLog", true, "If set, recall info will be logged into `recallLogfile`.")
	recallLogfile := flag.String("recall_file", "recall_log.txt", "Filename of recall log")
	recallShowHint := flag.Int("hint", 0, "If set, when recall test failes, hint will be given.\n0 for no hint, 1 for diff hint, 2 for full hint")
	flag.Parse()

	if *recall {
		checkRecall(*rememberLogfile, *recallLogfile, *recallLog, *recallShowHint)
		return
	}

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")

	uniqueRandomNumbers := generateRandomNumbers(*num, *maxium, *unique)
	fmt.Printf("%s Random Number is %d\n", datetimeFormatted, uniqueRandomNumbers)

	// type of remember is a pointer, add `*` prefix to get its value
	if *remember {
		info := fmt.Sprintf("%s %d\n", datetimeFormatted, uniqueRandomNumbers)
		filename := *rememberLogfile
		writeInfo(filename, info)
	}
}
