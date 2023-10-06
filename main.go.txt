package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/pkg/browser"
)

// appversion is used as a literal string in various places
const appversion = "Chasm v0.0.0"

func main() {

	fmt.Printf("Serving %v\n", *HTTPPort)

	fs := http.FileServer(http.Dir(Docroot))
	http.Handle("/", fs)

	http.HandleFunc("/chasm", centralDispatch)
	http.HandleFunc("/setup", showWizard)
	http.HandleFunc("/about", aboutChasm)
	http.HandleFunc("/ajax", fetchData)
	http.HandleFunc("/form", showForm)
	launchUI()
	http.ListenAndServe(":"+*HTTPPort, nil)
}

func launchUI() {

	time.Sleep(5 * time.Second)
	starturl := "http://localhost:" + *HTTPPort + "/chasm"

	browser.OpenURL(starturl)

}
func aboutChasm(w http.ResponseWriter, r *http.Request) {
	var cfg struct {
		Version string
		DBPath  string
		Schema  int
		Event   string
	}
	cfg.Version = appversion
	cfg.DBPath = *DBNAME
	cfg.Schema = 1
	cfg.Event = "IBA Rally"

	t := template.Must(template.New("about.html").Option("missingkey=error").ParseFiles(Docroot + "/" + *Language + "/about.html"))
	err := t.Execute(w, cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}

func fetchData(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Printf("Error parsing fetchData form %v\n", err)
		return
	}
	if r.FormValue("c") == "getregion" {
		sql := "SELECT * FROM regions WHERE region=?"
		stmt, _ := DBH.Prepare(sql)
		rows, _ := stmt.Query(r.FormValue("key"))
		defer rows.Close()
		if rows.Next() {
			type Region struct {
				Region, Localtz, Hostcountry, Locale string
				MilesKms, Decimalcomma               int
			}
			var reg Region
			rows.Scan(&reg.Region, &reg.Localtz, &reg.Hostcountry, &reg.Locale, &reg.MilesKms, &reg.Decimalcomma)
			a, _ := json.Marshal(reg)
			w.Write(a)
		}
		return
	}
	if r.FormValue("c") == "putreason" {
		sql := "INSERT OR REPLACE INTO reasons (Code, Briefdesc, Action, Param) VALUES(?,?,?,?)"

		_, err := DBH.Exec(sql, r.FormValue("Code"), r.FormValue("Briefdesc"), r.FormValue("Action"), r.FormValue("Param"))
		if err != nil {
			fmt.Printf("Exec %v failed %v\n", sql, err)
		}

		w.Write([]byte("ok"))
		return
	}

	if r.FormValue("c") == "delreason" {
		sql := "DELETE FROM reasons WHERE Code=?"

		_, err := DBH.Exec(sql, r.FormValue("Code"))
		if err != nil {
			fmt.Printf("Exec %v failed %v\n", sql, err)
		}

		w.Write([]byte("ok"))
		return
	}
}

func centralDispatch(w http.ResponseWriter, r *http.Request) {

	if !DBInitialised {
		fmt.Printf("Running basic wizard\n")
		//sendTemplate(w, r, "wizard")
		showWizardPage(w, r, 1)
		return
	}
	sql := "SELECT count(*) AS rex FROM entrants"
	rows, _ := DBH.Query(sql)
	var rex int = 0
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&rex)
	}
	fmt.Printf("Entrants = %v\n", rex)
	if rex < 1 {
		showSetupMenu(w, r)
		return
	}
	showMainMenu(w, r)

}

func showForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	frm := r.FormValue("f")
	if frm == "" {
		showSetupMenu(w, r)
		return
	}
	if frm == "editreasons" {
		sendReasons(w, r, frm)
	} else {
		sendTemplate(w, r, frm)
	}

}

func showSetupMenu(w http.ResponseWriter, r *http.Request) {
	sendTemplate(w, r, "setupmenu")
}

func showMainMenu(w http.ResponseWriter, r *http.Request) {
	sendTemplate(w, r, "mainmenu")
}

func sendReasons(w http.ResponseWriter, r *http.Request, tmplt string) {

	cfg := fetchReasons()

	t := template.Must(template.New(tmplt + ".html").Option("missingkey=error").ParseFiles(Docroot + "/" + *Language + "/" + tmplt + ".html"))

	err := t.Execute(w, cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
func sendTemplate(w http.ResponseWriter, r *http.Request, tmplt string) {

	fmt.Printf("%v\n", r.URL)

	var cfg RALLYCONFIG

	fetchConfig(&cfg)

	fmt.Printf("%v\n", cfg)

	t := template.Must(template.New(tmplt + ".html").Option("missingkey=error").ParseFiles(Docroot + "/" + *Language + "/" + tmplt + ".html"))

	err := t.Execute(w, cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
