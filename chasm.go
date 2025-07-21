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

func ajaxFetchBonusDetails(w http.ResponseWriter, r *http.Request) {

	b := r.FormValue("b")
	if b == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no b specified"}`)
		return
	}
	bd := fetchBonusVars(b)
	fmt.Fprintf(w, `{"ok":true,"msg":"ok","name":"%v"`, bd.BriefDesc)
	fmt.Fprintf(w, `,"flags":"%v","img":"%v"`, bd.Flags, bd.Image)
	fmt.Fprintf(w, `,"notes":"%v"`, bd.Notes)
	fmt.Fprintf(w, `,"askpoints":%v`, jsonBool(bd.AskPoints))
	pm := "p"
	if bd.PointsAreMults {
		pm = "m"
	}
	fmt.Fprintf(w, `,"pointsaremults":"%v"`, pm)
	fmt.Fprintf(w, `,"askmins":%v`, jsonBool(bd.AskMins))
	fmt.Fprintf(w, `,"points":%v`, bd.Points)
	fmt.Fprintf(w, `,"question":"%v"`, bd.Question)
	fmt.Fprintf(w, `,"answer":"%v"`, bd.Answer)
	fmt.Fprintf(w, `,"restmins":%v`, bd.RestMins)
	fmt.Fprint(w, `}`)
}

func ajaxFetchEntrantDetails(w http.ResponseWriter, r *http.Request) {

	e := intval(r.FormValue("e"))
	if e < 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"no e specified"}`)
		return
	}
	ed := fetchEntrantDetails(e)
	if ed.PillionName != "" {
		ed.RiderName += " &amp; " + ed.PillionName
	}
	tr := jsonBool(ed.PillionName != "" || ed.TeamID > 0)

	fmt.Fprintf(w, `{"ok":true,"msg":"ok","name":"%v","team":%v}`, ed.RiderName, tr)

}

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

func printNZ(i int) string {
	if i == 0 {
		return ""
	}
	return strconv.Itoa(i)
}
func jsonBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

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
	//	fmt.Printf("%v\n", CS)
	//	x, _ := json.Marshal(CS)
	//	fmt.Printf("%v\n", string(x))

	RallyTimezone, err = time.LoadLocation(CS.Basics.RallyTimezone)
	checkerr(err)

	if !*runOnline {
		return
	}
	fmt.Printf("Serving on port %v\n", *HTTPPort)
	fmt.Println()

	// DEBUG PURPOSES ONLY
	//recalc_all()

	fileserver := http.FileServer(http.Dir("."))
	http.Handle("/images/", fileserver)
	http.HandleFunc("/about", showAboutChasm)
	http.HandleFunc("/bonus", show_bonus)
	http.HandleFunc("/bonuses", list_bonuses)
	http.HandleFunc("/cards", showScorecards)
	http.HandleFunc("/cats", showCategorySets)
	http.HandleFunc("/certs", print_certs)
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
	http.HandleFunc("/help", show_help)
	http.HandleFunc("/img", builtin_images)
	http.HandleFunc("/import", showImport)
	http.HandleFunc("/js", send_js)
	http.HandleFunc("/menu", show_menu)
	http.HandleFunc("/niy", niy)
	http.HandleFunc("/odos", show_odo_checks)
	http.HandleFunc("/qlist", show_qlist)
	http.HandleFunc("/recalc", recalc_handler)
	http.HandleFunc("/reset", showResetOptions)
	http.HandleFunc("/rule", show_rule)
	http.HandleFunc("/rules", show_rules)
	http.HandleFunc("/savecert", save_certificate)
	http.HandleFunc("/saveebc", saveEBC)
	http.HandleFunc("/score", showScorecard)
	http.HandleFunc("/setup", show_setup)
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
	case "fetche":
		ajaxFetchEntrantDetails(w, r)
		return
	case "fetchb":
		ajaxFetchBonusDetails(w, r)
		return
	case "fetchcats":
		showCategoryCats(w, r)

		return
	case "putodo":
		update_odo(w, r)
		return
	case "putcfg":
		ajaxUpdateSettings(w, r)
		return
	}

	fmt.Fprintf(w, `{"ok":false,"msg":"[%v] not implemented yet"}`, f)
}

func niy(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "NIY")

	fmt.Fprintf(w, `<p class="error">NOT IMPLEMENTED YET</p><p>%v</p>`, r)
}
func recalc_handler(w http.ResponseWriter, r *http.Request) {

	const recalcfrm = `
	<article class="recalc">
	<p>This procedure will recalculate all scorecards. This involves rebuilding them from scratch by reprocessing the claims log. This should only take a few moments but it will need exclusive access to the database.</p>
	<p>It's quite safe to do this during a live rally.</p>
	<form action="/recalc">
		<input type="hidden" name="ok" value="ok">
		<button autofocus>Recalculate scorecards</button>
	</form>
	</article>
	`
	startHTML(w, "Recalc scorecards")

	e := r.FormValue("e")
	ok := r.FormValue("ok")
	if ok == "ok" {
		if e == "" {
			recalc_all()
		} else {
			n, err := strconv.Atoi(e)
			if err != nil {
				fmt.Fprintf(w, `<p class="error">%v is not numeric</p>`, e)
				return
			}
			recalc_scorecard(n)
		}
		fmt.Fprint(w, `</header><p class="thatsall">Scorecards recalculated</p>`)
		return
	}
	fmt.Fprint(w, recalcfrm)

}

func show_combo(w http.ResponseWriter, r *http.Request) {

	comboid := r.FormValue("c")
	var cb ComboBonus
	/*
		if comboid == "" {
			fmt.Fprint(w, "no comboid!")
			return
		}
	*/
	if comboid != "" {
		cr := loadCombos(comboid)
		if len(cr) < 1 {
			fmt.Fprint(w, "no such comboid")
			return
		}
		cb = cr[0]
	}
	showSingleCombo(w, cb, r.FormValue("back"))
}

func show_help(w http.ResponseWriter, r *http.Request) {

	topic := r.FormValue("topic")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprint(w, htmlheader)
	fmt.Fprint(w, `<p>I'm so sorry, there is no help`)
	if topic != "" {
		fmt.Fprintf(w, " for '%v'", topic)
	}
	fmt.Fprint(w, `</p>`)
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

func show_menu(w http.ResponseWriter, r *http.Request) {

	menu := r.FormValue("menu")
	if menu == "" {
		menu = "main"
	}

	startHTML(w, menu)
	showMenu(w, menu)

}

func show_setup(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "setup")

	showMenu(w, "setup")
}

func central_dispatch(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "")

	showMenu(w, "main")
}
