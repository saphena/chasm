package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DBNAME names the database file
var DBNAME *string = flag.String("db", "chasm.db", "database file")

// HTTPPort is the web port to serve
var HTTPPort *string = flag.String("port", "8080", "Web port")

var runOnline *bool = flag.Bool("online", false, "act as webserver")

// DBH provides access to the database
var DBH *sql.DB

var RallyTimezone *time.Location

func getIntegerFromDB(sqlx string, defval int) int {

	str := getStringFromDB(sqlx, strconv.Itoa(defval))
	res, err := strconv.Atoi(str)
	if err == nil {
		return res
	}
	return defval
}
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

	err = json.Unmarshal([]byte(defaultCS), &CS)
	checkerr(err)
	err = json.Unmarshal([]byte(secondDefault), &CS)
	checkerr(err)
	fmt.Printf("%v\n", CS)

	RallyTimezone, err = time.LoadLocation(CS.RallyTimezone)
	checkerr(err)

	recalc_all()
	//recalc_scorecard(2)

	if !*runOnline {
		return
	}

	fileserver := http.FileServer(http.Dir("."))
	http.Handle("/images/", fileserver)
	http.HandleFunc("/about", showAboutChasm)
	http.HandleFunc("/claims", list_claims)
	http.HandleFunc("/combo", show_combo)
	http.HandleFunc("/combos", show_combos)
	http.HandleFunc("/ebc", showEBC)
	http.HandleFunc("/ebclist", list_EBC_claims)
	http.HandleFunc("/recalc", recalc_handler)
	http.HandleFunc("/rule", show_rule)
	http.HandleFunc("/rules", show_rules)
	http.HandleFunc("/saveebc", saveEBC)
	http.HandleFunc("/x", json_requests)
	http.HandleFunc("/updtcrule", update_rule)
	http.HandleFunc("/", central_dispatch)
	http.ListenAndServe(":"+*HTTPPort, nil)

}

func json_requests(w http.ResponseWriter, r *http.Request) {

	var res struct {
		OK  bool
		Msg string
	}
	r.ParseForm()
	f := r.FormValue("f")
	switch f {
	case "axiscats":
		a := r.FormValue("a")
		s := r.FormValue("s")
		if a == "" || s == "" {
			log.Printf("a=%v, s=%v\n", a, s)
			fmt.Fprint(w, `{"ok":false,"msg":"missing args"}`)
			return
		}
		aa, err := strconv.Atoi(a)
		if err != nil {
			aa = 1
		}
		ss, err := strconv.Atoi(s)
		if err != nil {
			ss = 0
		}
		fmt.Printf("json: %v %v\n", aa, ss)
		res.Msg = strings.Join(optsSingleAxisCats(aa, ss), "")
		res.OK = true

		b, err := json.Marshal(res)
		checkerr(err)
		log.Println(string(b))
		fmt.Fprint(w, string(b))
		return
	case "saveebc":
		saveEBC(w, r)
		fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
		return
	}

	fmt.Fprint(w, `{"ok":true,"msg":"<option>one</option><option>two</option>"}`)
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

func show_combo(w http.ResponseWriter, r *http.Request) {

	comboid := r.FormValue("c")
	if comboid == "" {
		fmt.Fprint(w, "no comboid!")
		return
	}
	cr := loadCombos(comboid)
	if len(cr) < 1 {
		fmt.Fprint(w, "no such comboid")
		return
	}
	showSingleCombo(w, cr[0])
}
func show_rule(w http.ResponseWriter, r *http.Request) {

	const Leg = 0

	n, err := strconv.Atoi(r.FormValue("r"))
	if err != nil {
		n = 1
	}
	CompoundRules = build_compoundRuleArray(Leg)
	for _, cr := range CompoundRules {
		if cr.Ruleid == n {
			showSingleRule(w, cr)
			return
		}
	}
	fmt.Fprint(w, `OMG`)
}

func central_dispatch(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "")

	showMenu(w, "main")
}
