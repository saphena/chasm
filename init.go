package main

import (
	"database/sql"
	"flag"
	"fmt"
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

	fmt.Printf("%v\nCopyright (c) Bob Stammers 2021\n\n", PROGRAMVERSION)

	flag.Parse()
	dbx, _ := filepath.Abs(*DBNAME)
	fmt.Printf("Using %v\n", dbx)

	DBH, err = sql.Open("sqlite3", dbx)
	if err != nil {
		panic(err)
	}
	sql := "SELECT eventname FROM config"
	rows, err := DBH.Query(sql)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		rows.Scan((&EventName))
		DBInitialised = EventName != ""
	} else {
		DBInitialised = false
		EventName = ""
	}

}
