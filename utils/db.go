package utils

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"regexp"

	_ "github.com/mattn/go-sqlite3"
)

func createDB(dbfile string) {
	db, err := sql.Open("sqlite3", dbfile)
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
		"number"	INTEGER NOT NULL,
		"missingCount"	INTEGER,
		"wrongCount"	INTEGER
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

func writeDB(logFile, recallFile, dbFile string) {

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
	AddNumbers(dbFile, datatimeSlice, numbersSlice)

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
	AddRecalls(dbFile, datatimeSlice, numbersSlice, resultSlice)

}

func InitDB(logFile, recallFile, dbFile string) {
	if FileExist(dbFile) {
		// fmt.Printf("Database file already exists and I will not init DB again!\n")
		return
	} else {
		createDB(dbFile)
		writeDB(logFile, recallFile, dbFile)
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
func AddWrongNumbers(dbFile string, wrongType int, numbers []string) {
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

	// if len(datetime) != len(numbers) {
	// 	panic("length of datetime and numbers mismatch while insert into db")
	// }
	// for i, v := range datetime {
	// 	_, err = stmt.Exec(v, numbers[i])
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
