package utils

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

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

// common returns the elements in `a` that are in `b`.
func Common(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; found {
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

func FormatUserInput(recallString string) string {
	reg := regexp.MustCompile(`\s+`)
	recallString = reg.ReplaceAllString(recallString, " ")

	recallString = strings.TrimSuffix(recallString, "\n")
	recallString = strings.TrimSpace(recallString)
	return recallString
}

// find content in filename from user input
func FindContentInFile(rememberLogfile, recallString string) string {
	recallString = FormatUserInput(recallString)

	content, err := os.ReadFile(rememberLogfile)
	if err != nil {
		panic(err)
	}

	recallSlice := strings.Split(recallString, " ")
	firstNumCount := min(len(recallSlice), 2)
	correctRegStr := fmt.Sprintf(`\[%s [\w ]+\]`, strings.Join(recallSlice[:firstNumCount], " "))
	reg := regexp.MustCompile(correctRegStr)
	correctStr := reg.FindString(string(content))
	correctStr = strings.TrimPrefix(correctStr, "[")
	correctStr = strings.TrimSuffix(correctStr, "]")
	return correctStr
}

// read user input and format it
func ReadAndFormat() string {
	reader := bufio.NewReader(os.Stdin)
	recallString, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// replace all whitespace characters with a single space
	reg := regexp.MustCompile(`\s+`)
	recallString = reg.ReplaceAllString(recallString, " ")

	recallString = strings.TrimSuffix(recallString, "\n")
	recallString = strings.TrimSpace(recallString)
	return recallString
}

// compare string
func CompareHint(Config BaseConfig, recallString, correctStr string, showOutput bool) bool {
	recallResult := recallString == correctStr

	recallSlice := strings.Split(recallString, " ")
	correctSlice := strings.Split(correctStr, " ")

	missingNumbers := Difference(correctSlice, recallSlice)
	AddWrongNumbers(Config, 2, missingNumbers)

	wrongNumbers := Difference(recallSlice, correctSlice)
	AddWrongNumbers(Config, 1, wrongNumbers)

	correctNumbers := Common(recallSlice, correctSlice)
	AddCorrectNumbers(Config, correctNumbers)

	if !showOutput {
		return recallResult
	}

	if recallResult {
		fmt.Println("You have a correct memory!")
	} else {
		fmt.Println("Are you sure you remember it right?")
		if Config.HintLevel > 0 {
			// compare two strings and assume first few numbers (min of 2 and slice lenth) is correct

			info := "\nHint Part:\n"
			if len(correctStr) < 1 {
				info += "Totally Wrong! Don't you even remember the first 2 number(s)?\n"
			} else {
				if Config.HintLevel == 1 {

					if len(missingNumbers) > 5 {
						info += fmt.Sprintf("You are missing %d numbers, which is too many for hinting. You should remember it again!\n", len(missingNumbers))
					} else {
						if len(missingNumbers) > 0 {
							red := color.New(color.FgRed, color.Bold).SprintFunc()
							info += fmt.Sprintf("You are missing these numbers: %s\n", red(strings.Join(missingNumbers, " ")))
						}
						wrongNumbers := Difference(recallSlice, correctSlice)
						if len(wrongNumbers) > 0 {
							yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
							info += fmt.Sprintf("You add these numbers which should't exist: %s\n", yellow(strings.Join(wrongNumbers, " ")))
						}
						if len(missingNumbers) == 0 && len(wrongNumbers) == 0 {
							cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
							info += fmt.Sprintf("You have remember all the numbers, but the %s are wrong!", cyan("orders"))
						}
					}

				} else if Config.HintLevel == 2 {
					//info += fmt.Sprintf("The Right: %s\nThe Wrong: %s", correctStr, recallString)

					// colorize missing and wrong numbers
					correctSlice := strings.Split(correctStr, " ")

					red := color.New(color.FgRed, color.Bold).SprintFunc()
					yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
					green := color.New(color.FgGreen, color.Bold).SprintFunc()

					var correctStrColored []string
					var wrongStrColored []string
					var extraStrColored []string

					// only numbers' orders are wrong
					if len(missingNumbers) == 0 && len(wrongNumbers) == 0 {
						for idx := range correctSlice {
							rVal := correctSlice[idx]
							wVal := recallSlice[idx]
							if rVal == wVal {
								correctStrColored = append(correctStrColored, rVal)
								wrongStrColored = append(wrongStrColored, wVal)
							} else {
								correctStrColored = append(correctStrColored, green(rVal))
								wrongStrColored = append(wrongStrColored, red(wVal))
							}
						}
					} else {

						recallSliceFix := make([]string, len(recallSlice))
						copy(recallSliceFix, recallSlice)
						// fmt.Println(recallSliceFix)

						for idx, val := range correctSlice {
							if idx < len(recallSliceFix) && val != recallSliceFix[idx] {
								correctStrColored = append(correctStrColored, green(val))
								recallSlice[idx] = strings.Repeat(" ", len(val))
								extraStrColored = append(extraStrColored, red(recallSliceFix[idx]))

							} else {
								extraStrColored = append(extraStrColored, strings.Repeat(" ", len(val)))
								correctStrColored = append(correctStrColored, val)

							}
						}

						for _, val := range recallSlice {
							if slices.Contains(wrongNumbers, val) {
								wrongStrColored = append(wrongStrColored, yellow(val))
							} else {
								wrongStrColored = append(wrongStrColored, val)
							}
						}

					}
					correctStr = strings.Join(correctStrColored, " ")
					wrongStr := strings.Join(wrongStrColored, " ")
					extraStr := strings.Join(extraStrColored, " ")
					info += fmt.Sprintf("The Right: %s\nThe Wrong: %s\nThe Extra: %s\n", correctStr, wrongStr, extraStr)
				}
			}
			fmt.Printf("%s", info)
		}
	}

	return recallResult
}
