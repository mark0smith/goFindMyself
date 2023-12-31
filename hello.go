package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"utils"

	"github.com/fatih/color"
)

// recall func
func checkRecall(dbfile string, showhint int) {

	fmt.Println("What do you remember?")
	recallString := utils.ReadAndFormat()
	fmt.Printf("\nYou have entered: %s\n", recallString)

	correctStr := utils.FindContentInDB(dbfile, recallString)
	recallResult := utils.CompareHint(dbfile, recallString, correctStr, showhint, true)

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")
	utils.AddRecalls(dbfile, []string{datetimeFormatted}, []string{recallString}, []string{fmt.Sprintf("%t", recallResult)})

}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	num := flag.Int("n", 30, "Number of Random Numbers.")
	maxium := flag.Int("m", 100, "Generated number wont't be bigger than this number.")
	unique := flag.Bool("u", true, "If set, all generated numbers will be unique.")
	remember := flag.Bool("r", false, "If set, generated numbers will be logged into `rememberLogfile`.\nYou should set this if you want to do recall test later!")
	recall := flag.Bool("recall", false, "If set, run a recall test, instead of generating random numbers.")
	recallShowHint := flag.Int("hint", 0, "If set, when recall test failes, hint will be given.\n0 for no hint, 1 for diff hint, 2 for full hint")
	dbFilename := flag.String("db", "data.db", "Filename of database")

	flag.Parse()

	if *recall {
		defer utils.Timer("checkRecall")()
		checkRecall(*dbFilename, *recallShowHint)
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
		numbersStr := fmt.Sprintf("%v", uniqueRandomNumbers)
		numbersStr = strings.TrimPrefix(numbersStr, "[")
		numbersStr = strings.TrimSuffix(numbersStr, "]")
		utils.AddNumbers(*dbFilename, []string{datetimeFormatted}, []string{numbersStr})
	}
}
