package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// https://stackoverflow.com/a/45766707
func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %.2fs!\n", name, time.Since(start).Seconds())
	}
}

func fileExist(filename string) bool {
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
		"Numbers"	TEXT NOT NULL
	);
	CREATE TABLE "Recall" (
		"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"datetime"	TEXT NOT NULL,
		"recall"	TEXT NOT NULL,
		"result"	INTEGER NOT NULL
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

// func writeLog2db(logFile,dbFile string){
// }

func initDB(logFile, recallFile, dbFile string) {
	if fileExist(dbFile) {
		fmt.Printf("Database file already exists and I will not init DB again!\n")
		return
	} else {
		createDB(dbFile)
	}
}
