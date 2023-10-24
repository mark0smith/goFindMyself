package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"utils"

	"github.com/fatih/color"
)

// recall func
func checkRecall(rememberLogfile string, recallLogfile string, recallLog bool, showhint int, dbfile string) {

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
	recallString = strings.TrimSpace(recallString)
	fmt.Printf("\nYou have entered: %s\n", recallString)

	toCheck := fmt.Sprintf("[%s]", recallString)
	content, err := os.ReadFile(rememberLogfile)
	if err != nil {
		panic(err)
	}

	recallResult := strings.Contains(string(content), toCheck)

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")
	utils.AddRecalls(dbfile, []string{datetimeFormatted}, []string{recallString}, []string{fmt.Sprintf("%t", recallResult)})
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
					missingNumbers := utils.Difference(correctSlice, recallSlice)
					if len(missingNumbers) > 5 {
						info += fmt.Sprintf("You are missing %d numbers, which is too many for hinting. You should remember it again!\n", len(missingNumbers))
					} else {
						if len(missingNumbers) > 0 {
							red := color.New(color.FgRed, color.Bold).SprintFunc()
							info += fmt.Sprintf("You are missing these numbers: %s\n", red(strings.Join(missingNumbers, " ")))
						}
						wrongNumbers := utils.Difference(recallSlice, correctSlice)
						if len(wrongNumbers) > 0 {
							yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
							info += fmt.Sprintf("You add these numbers which should't exist: %s\n", yellow(strings.Join(wrongNumbers, " ")))
						}
						if len(missingNumbers) == 0 && len(wrongNumbers) == 0 {
							cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
							info += fmt.Sprintf("You have remember all the numbers, but the %s are wrong!", cyan("orders"))
						}
					}

				} else if showhint == 2 {
					//info += fmt.Sprintf("The Right: %s\nThe Wrong: %s", correctStr, recallString)

					// colorize missing and wrong numbers
					correctSlice := strings.Split(correctStr, " ")
					missingNumbers := utils.Difference(correctSlice, recallSlice)
					wrongSlice := strings.Split(recallString, " ")
					wrongNumbers := utils.Difference(recallSlice, correctSlice)

					red := color.New(color.FgRed, color.Bold).SprintFunc()
					yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
					green := color.New(color.FgGreen, color.Bold).SprintFunc()

					var correctStrColored []string
					var wrongStrColored []string

					if len(missingNumbers) == 0 && len(wrongNumbers) == 0 {
						for idx := range correctSlice {
							rVal := correctSlice[idx]
							wVal := wrongSlice[idx]
							if rVal == wVal {
								correctStrColored = append(correctStrColored, rVal)
								wrongStrColored = append(wrongStrColored, wVal)
							} else {
								correctStrColored = append(correctStrColored, green(rVal))
								wrongStrColored = append(wrongStrColored, red(wVal))
							}
						}
					} else {
						for _, val := range correctSlice {
							if slices.Contains(missingNumbers, val) {
								correctStrColored = append(correctStrColored, red(val))
							} else {
								correctStrColored = append(correctStrColored, val)
							}
						}

						for _, val := range wrongSlice {
							if slices.Contains(wrongNumbers, val) {
								wrongStrColored = append(wrongStrColored, yellow(val))
							} else {
								wrongStrColored = append(wrongStrColored, val)
							}
						}
					}
					correctStr = strings.Join(correctStrColored, " ")
					wrongStr := strings.Join(wrongStrColored, " ")
					info += fmt.Sprintf("The Right: %s\nThe Wrong: %s\n", correctStr, wrongStr)
				}
			}
			fmt.Printf("%s", info)
		}
	}

	if recallLog {
		datetime := time.Now()
		datetimeFormatted := datetime.Format("2006-01-02 15:04:05")
		info := fmt.Sprintf("%s Recall: %s, Result: %v\n", datetimeFormatted, recallString, recallResult)
		utils.WriteInfo(recallLogfile, info)
	}
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
	dbFilename := flag.String("db", "data.db", "Filename of database")
	migratedb := flag.Bool("migratedb", true, "If set, recall info will be logged into `recallLogfile`.")

	flag.Parse()

	if *migratedb {
		utils.InitDB(*rememberLogfile, *recallLogfile, *dbFilename)
	}

	if *recall {
		defer utils.Timer("checkRecall")()
		checkRecall(*rememberLogfile, *recallLogfile, *recallLog, *recallShowHint, *dbFilename)
		return
	}

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")

	uniqueRandomNumbers := utils.GenerateRandomNumbers(*num, *maxium, *unique)

	// format output string, making it easier to remember
	var outputStr string
	boldBlue := color.New(color.Bold, color.FgHiBlue).SprintFunc()
	for idx, val := range uniqueRandomNumbers {
		info := ""
		if (idx+1)%5 == 0 {
			info = fmt.Sprintf(" %-2s", boldBlue(fmt.Sprintf("%-2d", val)))
			if (idx+1)%10 == 0 {
				info += "\n"
			}
		} else {
			if (idx+1)%10 == 1 {
				info = fmt.Sprintf("%-2d", val)
			} else {
				info = fmt.Sprintf(" %-2d", val)
			}
		}
		outputStr += info
	}

	fmt.Printf("%s Random Number Generated:\n%s\n", datetimeFormatted, outputStr)

	// type of remember is a pointer, add `*` prefix to get its value
	if *remember {
		info := fmt.Sprintf("%s %d\n", datetimeFormatted, uniqueRandomNumbers)
		filename := *rememberLogfile
		utils.WriteInfo(filename, info)
		numbersStr := fmt.Sprintf("%v", uniqueRandomNumbers)
		numbersStr = strings.TrimPrefix(numbersStr, "[")
		numbersStr = strings.TrimSuffix(numbersStr, "]")
		utils.AddNumbers(*dbFilename, []string{datetimeFormatted}, []string{numbersStr})
	}
}
