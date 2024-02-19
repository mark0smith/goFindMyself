package utils

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func createDB(Config BaseConfig) {
	db, err := sql.Open("sqlite3", Config.DBFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE "RandomNumbers" (
		"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"datetime"	TEXT NOT NULL,
		"numbers"	TEXT NOT NULL
	);
	CREATE TABLE "Recall" (
		"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"datetime"	TEXT NOT NULL,
		"recall"	TEXT NOT NULL,
		"result"	INTEGER NOT NULL
	);
	CREATE TABLE "Number" (
		"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"number"	INTEGER NOT NULL UNIQUE,
		"missingCount"	INTEGER NOT NULL,
		"correctCount"	INTEGER NOT NULL,
		"wrongCount"	INTEGER NOT NULL
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func writeDB(logFile string, recallFile string, Config BaseConfig) {

	content, err := os.ReadFile(logFile)
	if err != nil {
		panic(err)
	}

	var datatimeSlice []string
	var numbersSlice []string

	var re = regexp.MustCompile(`(?m)^(?P<datetime>[\d\- :]+) \[(?P<numbers>[\d ]+)\]$`)
	for _, line := range bytes.Split(content, []byte("\n")) {
		if re.MatchString(string(line)) {
			matches := re.FindStringSubmatch(string(line))
			datetime := matches[re.SubexpIndex("datetime")]
			numbers := matches[re.SubexpIndex("numbers")]
			datatimeSlice = append(datatimeSlice, datetime)
			numbersSlice = append(numbersSlice, numbers)
		}
	}

	AddNumbers(Config.DBFilename, datatimeSlice, numbersSlice)

	content, err = os.ReadFile(recallFile)
	if err != nil {
		panic(err)
	}
	datatimeSlice = nil
	numbersSlice = nil
	var resultSlice []string

	re = regexp.MustCompile(`(?m)^(?P<datetime>[\d\- :]+) Recall: (?P<numbers>[\d ]+), Result: (?P<result>\w+)$`)
	for _, line := range bytes.Split(content, []byte("\n")) {
		if re.MatchString(string(line)) {
			matches := re.FindStringSubmatch(string(line))
			result := matches[re.SubexpIndex("result")]
			datetime := matches[re.SubexpIndex("datetime")]
			numbers := matches[re.SubexpIndex("numbers")]
			datatimeSlice = append(datatimeSlice, datetime)
			numbersSlice = append(numbersSlice, numbers)
			resultSlice = append(resultSlice, result)
		}
	}
	AddRecalls(Config.DBFilename, datatimeSlice, numbersSlice, resultSlice)

	for _, recallString := range numbersSlice {
		correctStr := FindContentInDB(Config, recallString)
		CompareHint(Config, recallString, correctStr, false)
	}

}

func InitDB(Config BaseConfig) {
	if FileExist(Config.DBFilename) {
		return
	} else {
		createDB(Config)
		logFile := "log.txt"
		recallFile := "recall_log.txt"
		if FileExist(logFile) && FileExist(recallFile) {
			writeDB(logFile, recallFile, Config)
		}

	}
}

func AddNumbers(dbFile string, datetime []string, numbers []string) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into RandomNumbers(datetime, numbers) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if err != nil {
		panic(err)
	}

	if len(datetime) != len(numbers) {
		panic("length of datetime and numbers mismatch while insert into db")
	}
	for i, v := range datetime {
		_, err = stmt.Exec(v, numbers[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

func AddRecalls(dbFile string, datetime []string, numbers []string, result []string) {

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into Recall(datetime, recall,result) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if len(datetime) != len(numbers) || len(datetime) != len(result) || len(numbers) != len(result) {
		panic("length of datetime and numbers and results mismatch while insert into db")
	}
	for i, v := range datetime {
		_, err = stmt.Exec(v, numbers[i], result[i])
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

// add wrong numbers to db
// type 1 for wrong numbers, type 2 for missing numbers
func AddWrongNumbers(Config BaseConfig, wrongType int, numbers []string) {
	db, err := sql.Open("sqlite3", Config.DBFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, number := range numbers {
		// check if input is number and not bigger than Maxium
		if len(number) < 1 {
			continue
		}
		numberInt, err := strconv.Atoi(number)
		if err != nil {
			// fmt.Printf("%v is not a number. Ignoring ...\n", number)
			continue
		}
		if numberInt > Config.Maxium {
			// fmt.Printf("Number %v is bigger than maxium `%v`. Don't add it into wrong logs.\n", numberInt, Config.Maxium)
			continue
		}

		stmt, err := tx.Prepare("select missingCount,wrongCount from Number where number = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		var missingCount, wrongCount string
		err = stmt.QueryRow(number).Scan(&missingCount, &wrongCount)

		if err != nil {
			stmt, err = tx.Prepare("insert into Number(number,missingCount,wrongCount,correctCount) values(?,?, ?,0)")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			missingCount := 0
			wrongCount := 0
			if wrongType == 1 {
				wrongCount += 1
			} else {
				missingCount += 1
			}
			_, err = stmt.Exec(number, missingCount, wrongCount)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			stmt, err = tx.Prepare("update Number set missingCount = ? , wrongCount = ?  where number = ?")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			missingCount, _ := strconv.Atoi(missingCount)
			wrongCount, _ := strconv.Atoi(wrongCount)
			if wrongType == 1 {
				wrongCount += 1
			} else {
				missingCount += 1
			}
			_, err = stmt.Exec(missingCount, wrongCount, number)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

// add correct numbers to db
func AddCorrectNumbers(Config BaseConfig, numbers []string) {
	db, err := sql.Open("sqlite3", Config.DBFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for _, number := range numbers {
		stmt, err := tx.Prepare("select correctCount from Number where number = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		var correctCount string
		err = stmt.QueryRow(number).Scan(&correctCount)

		if err != nil {
			// if not found then insert a new record
			stmt, err = tx.Prepare("insert into Number(number,correctCount,missingCount,wrongCount) values(?,?,0,0)")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			correctCount := 1

			_, err = stmt.Exec(number, correctCount)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			stmt, err = tx.Prepare("update Number set correctCount = ?  where number = ?")
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			correctCount, _ := strconv.Atoi(correctCount)
			correctCount += 1

			_, err = stmt.Exec(correctCount, number)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}

// find content in filename from user input
func FindContentInDB(Config BaseConfig, recallString string) string {
	recallString = FormatUserInput(recallString)

	recallSlice := strings.Split(recallString, " ")
	firstNumCount := min(len(recallSlice), 2)
	queryContent := strings.Join(recallSlice[:firstNumCount], " ")
	queryContent += "%"

	db, err := sql.Open("sqlite3", Config.DBFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("select numbers from RandomNumbers where numbers like ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var correctStr string
	err = stmt.QueryRow(queryContent).Scan(&correctStr)

	result := ""
	if err != nil {

	} else {
		result = correctStr
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return result
}
