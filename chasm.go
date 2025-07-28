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

// Set to true for production builds to suppress debug traps
const ProductionBuild = false

// DBNAME names the database file
var DBNAME *string = flag.String("db", "chasm.db", "database file")

// HTTPPort is the web port to serve
var HTTPPort *string = flag.String("port", "8080", "Web port")

var runOnline *bool = flag.Bool("online", true, "act as webserver")

// DBH provides access to the database
var DBH *sql.DB

var RallyTimezone *time.Location

func main() {

	fmt.Printf("Chasm v%v  Copyright (c) %v %v\n", ChasmVersion, CopyriteYear, CopyriteHolder)
	flag.Parse()

	dbx, _ := filepath.Abs(*DBNAME)
	fmt.Printf("Using %v\n", dbx)

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
	err = json.Unmarshal([]byte(debugDefaults), &CS)
	checkerr(err)
	err = json.Unmarshal([]byte(getStringFromDB("SELECT ifnull(Settings,'{}') FROM config", "{}")), &CS)
	checkerr(err)
	loadRallyBasics(&CS.Basics)

	RallyTimezone, err = time.LoadLocation(CS.Basics.RallyTimezone)
	checkerr(err)

	if !*runOnline {
		return
	}
	fmt.Printf("Serving on port %v\n", *HTTPPort)
	fmt.Println()

	fileserver := http.FileServer(http.Dir("."))
	http.Handle("/images/", fileserver)
	http.HandleFunc("/about", showAboutChasm)
	http.HandleFunc("DELETE /bonus/{b}", deleteBonus)
	http.HandleFunc("/bonus", show_bonus)
	http.HandleFunc("POST /bonus/{b}", createBonus)
	http.HandleFunc("/bonuses", list_bonuses)
	http.HandleFunc("/cards", showScorecards)
	http.HandleFunc("/cats", showCategorySets)
	http.HandleFunc("/certs", print_certs)
	http.HandleFunc("DELETE /claim/{claimid}", deleteClaim)
	http.HandleFunc("/claim", showClaim)
	http.HandleFunc("/claims", list_claims)
	http.HandleFunc("/combo", show_combo)
	http.HandleFunc("/combos", show_combos)
	http.HandleFunc("/config", editConfigMain)
	http.HandleFunc("/css", send_css)
	http.HandleFunc("/ebc", showEBC)
	http.HandleFunc("/ebclist", list_EBC_claims)
	http.HandleFunc("/editcert", edit_certificate)
	http.HandleFunc("DELETE /entrant/{e}", deleteEntrant)
	http.HandleFunc("/entrant/{e}", showEntrant)
	http.HandleFunc("/entrants", list_entrants)
	http.HandleFunc("/img", builtin_images)
	http.HandleFunc("/import", showImport)
	http.HandleFunc("/js", send_js)
	http.HandleFunc("/menu", show_menu)
	http.HandleFunc("/niy", niy)
	http.HandleFunc("/odos", show_odo_checks)
	http.HandleFunc("/qlist", show_qlist)
	http.HandleFunc("/recalc", recalc_handler)
	http.HandleFunc("/reset", showResetOptions)
	http.HandleFunc("DELETE /rule/{rule}", deleteRule)
	http.HandleFunc("POST /rule", createRule)
	http.HandleFunc("/rule", show_rule)
	http.HandleFunc("/rules", show_rules)
	http.HandleFunc("/savecert", save_certificate)
	http.HandleFunc("/saveebc", saveEBC)
	http.HandleFunc("/saverule", saveRule)
	http.HandleFunc("/score", showScorecard)
	http.HandleFunc("/setup", show_setup)
	http.HandleFunc("/teams", list_teams)
	http.HandleFunc("/updtcrule", update_rule)
	http.HandleFunc("/upload", uploadImportDatafile)
	http.HandleFunc("/x", json_requests)
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
		//fmt.Printf("json: %v %v\n", aa, ss)
		res.Msg = strings.Join(optsSingleAxisCats(aa, ss), "")
		res.OK = true

		b, err := json.Marshal(res)
		checkerr(err)
		log.Println(string(b))
		fmt.Fprint(w, string(b))
		return
	case "addb":
		createBonus(w, r)
		return
	case "addcat":
		addCatCat(w, r)
		return
	case "addco":
		createCombo(w, r)
		return
	case "adde":
		createEntrant(w, r)
		return
	case "delb":
		deleteBonus(w, r)
		return
	case "delcat":
		delCatCat(w, r)
		return
	case "delco":
		deleteCombo(w, r)
		return
	case "saveb":
		saveBonus(w, r)

		return
	case "savecat":
		updateCatName(w, r)
		return
	case "saveco":
		saveCombo(w, r)
		return
	case "saveclaim":
		saveClaim(r)
		fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
		return
	case "savee":
		saveEntrant(w, r)
		return
	case "saveebc":
		saveEBC(w, r)
		fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
		return
	case "savers":
		updateReviewStatus(w, r)
		return
	case "saveset":
		updateSetName(w, r)
		return
	case "setteam":
		setTeam(w, r)
		return
	case "fetche":
		ajaxFetchEntrantDetails(w, r)
		return
	case "fetchb":
		fetchBonusDetails(w, r)
		return
	case "fetchcats":
		showCategoryCats(w, r)
		return

	case "fetchmembers":
		showTeamMembers(w, r)

		return
	case "putodo":
		update_odo(w, r)
		return
	case "putcfg":
		ajaxUpdateSettings(w, r)
		return
	case "ulist":
		comboBonusList(w, r)
		return
	}

	fmt.Fprintf(w, `{"ok":false,"msg":"[%v] not implemented yet"}`, f)
}

func central_dispatch(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Main menu")

	showMenu(w, "main")
}
