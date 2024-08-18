package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// DBNAME names the database file
var DBNAME *string = flag.String("db", "chasm.db", "database file")

// HTTPPort is the web port to serve
var HTTPPort *string = flag.String("port", "8080", "Web port")

var runOnline *bool = flag.Bool("online", false, "act as webserver")

// DBH provides access to the database
var DBH *sql.DB

func getStringFromDB(sqlx string, defval string) string {

	rows, err := DBH.Query(sqlx)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if rows.Next() {
		var val string
		rows.Scan(&val)
		return val
	}
	return defval
}

func main() {

	log.Println("Hello sailor")
	flag.Parse()

	dbx, _ := filepath.Abs(*DBNAME)
	fmt.Printf("Using %v\n\n", dbx)

	var err error
	DBH, err = sql.Open("sqlite3", dbx)
	if err != nil {
		panic(err)
	}

	sqlx := "SELECT DBInitialised FROM config"
	dbi := getStringFromDB(sqlx, "0")
	if dbi != "1" {
		fmt.Println("Duff database")
		return
	}
	recalc_all()
	//recalc_scorecard(2)

	if !*runOnline {
		return
	}

	http.HandleFunc("/", central_dispatch)
	http.HandleFunc("/about", about_chasm)
	http.HandleFunc("/recalc", recalc_handler)
	http.HandleFunc("/rule", show_rule)
	http.ListenAndServe(":"+*HTTPPort, nil)

}

func recalc_handler(w http.ResponseWriter, r *http.Request) {

	e := r.FormValue("e")
	if e == "" {
		recalc_all()
	} else {
		n, err := strconv.Atoi(e)
		if err != nil {
			w.Write([]byte(`{ok:false,msg:"e not numeric"}`))
			return
		}
		recalc_scorecard(n)
	}
	w.Write([]byte(`{ok:true,msg:"ok"}`))
}
func show_rule(w http.ResponseWriter, r *http.Request) {

	const Leg = 0

	n, err := strconv.Atoi(r.FormValue("r"))
	if err != nil {
		n = 1
	}
	CompoundRules = build_compoundRuleArray(Leg)
	if n < len(CompoundRules) {
		showSingleRule(w, CompoundRules[n])
	} else {
		w.Write([]byte("OMG"))
	}

}
func about_chasm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello there, I say, I say")
}

const mockupFrontPage = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>CHASM</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
body {
	margin: 0;
	font-size: 14pt;
	font-family				: Verdana, Arial, Helvetica, sans-serif; 

}
.topbar {
	background-color: lightgray;
	border none none solid 2px none;
	width: 100%;
	margin: 0;
	padding: 5px;
}
.about {
	float: right;
	padding-right: 1em;
	font-size: 10pt;
	vertical-align: middle;
	display: table-cell;
}
</style>
</head>
<body>
`
const homeIcon = `
<input title="Return to main menu" style="padding:1px;" type="button" value=" 🏠 " onclick="window.location='admin.php'">`

func central_dispatch(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, mockupFrontPage)
	fmt.Fprint(w, `<div class="topbar">`+homeIcon+` 12 Days Euro Rally<span class="about">About CHASM</span></div>`)
}
