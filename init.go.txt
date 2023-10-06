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

// Language identifies the i18n folder to use
var Language *string = flag.String("lang", "", "Language folder override")

// DBInitialised shows the configuration status of the database
var DBInitialised bool

// EventName is the name of the current event
var EventName string

// DBH provides access to the database
var DBH *sql.DB

// Docroot holds the path to the web assets folder
var Docroot string

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

	sqlx := "SELECT Eventname,DBInitialised, Langcode FROM config"

	rows, err := DBH.Query(sqlx)
	if err != nil {
		if getYN("Database is not setup. Establish now? [Y/n] ") {
			createDatabase(*Language)
			fmt.Printf("Closing new database, please rerun me\n")
			os.Exit(0)
			//rows, _ = DBH.Query(sqlx)
		} else {
			fmt.Printf("\nDatabase is not setup, run terminating\n\n")
			os.Exit(1)
		}
	}
	defer rows.Close()

	if rows.Next() {
		var dbi int
		rows.Scan(&EventName, &dbi, Language)
		DBInitialised = (EventName != "") && (dbi != 0)
	} else {
		DBInitialised = false
		EventName = ""
		*Language = "en"
	}
	Docroot = "files/"

}
