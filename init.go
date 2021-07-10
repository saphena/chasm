package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DBNAME names the database file
var DBNAME *string = flag.String("db", "chasm.db", "database file")

// HTTPPort is the web port to serve
var HTTPPort *string = flag.String("port", "8080", "Web port")

// DBInitialised shows the configuration status of the database
var DBInitialised bool

// EventName is the name of the current event
var EventName string

// DBH provides access to the database
var DBH *sql.DB

func init() {

	var err error

	fmt.Printf("%v\nCopyright (c) Bob Stammers 2021\n\n", appversion)

	flag.Parse()
	dbx, _ := filepath.Abs(*DBNAME)
	fmt.Printf("Using %v\n\n", dbx)

	DBH, err = sql.Open("sqlite3", dbx)
	if err != nil {
		panic(err)
	}

	sql := "SELECT Eventname,DBInitialised FROM config"

	rows, err := DBH.Query(sql)
	if err != nil {
		if getYN("Database is not setup. Establish now? [Y/n] ") {
			createDatabase()
			rows, _ = DBH.Query(sql)
		} else {
			fmt.Printf("\nDatabase is not setup, run terminating\n\n")
			os.Exit(1)
		}
	}
	defer rows.Close()

	if rows.Next() {
		var dbi int
		rows.Scan(&EventName, &dbi)
		DBInitialised = (EventName != "") && (dbi != 0)
	} else {
		DBInitialised = false
		EventName = ""
	}

}
