package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"utils"

	"github.com/fatih/color"
)

var Config utils.BaseConfig

// recall func
func checkRecall() {

	fmt.Println("What do you remember?")
	start := time.Now()
	recallString := utils.ReadAndFormat()
	var durations = time.Since(start).Seconds()
	fmt.Printf("\nYou have entered: %s\n", recallString)

	correctStr := utils.FindContentInDB(Config, recallString)
	recallResult := utils.CompareHint(Config, recallString, correctStr, true)

	red := color.New(color.FgHiRed, color.Bold).PrintFunc()
	if durations < 50 && durations > 10 && recallResult {
		red("You have remembered all these numbers and should try some new ones!\n")
		// disable random numbers' record and no more overfitting!
		utils.SetRandomNumbersDisabled(Config, correctStr)
	}

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")
	utils.AddRecalls(Config.DBFilename, []string{datetimeFormatted}, []string{recallString}, []string{fmt.Sprintf("%t", recallResult)})

}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	utils.InitConfig(&Config)
	utils.InitDB(Config)

	if Config.Recall {

		checkRecall()
		return
	}

	datetime := time.Now()
	datetimeFormatted := datetime.Format("2006-01-02 15:04:05")

	uniqueRandomNumbers := utils.GenerateRandomNumbers(Config.Count, Config.Maxium, Config.IsUnique)

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
	if Config.IsStored {
		numbersStr := fmt.Sprintf("%v", uniqueRandomNumbers)
		numbersStr = strings.TrimPrefix(numbersStr, "[")
		numbersStr = strings.TrimSuffix(numbersStr, "]")

		utils.AddNumbers(Config.DBFilename, []string{datetimeFormatted}, []string{numbersStr})
	}
}
