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

	fs := http.FileServer(http.Dir("./files"))
	http.Handle("/", fs)

	http.HandleFunc("/chasm", centralDispatch)
	http.HandleFunc("/setup", showWizard)
	http.HandleFunc("/about", aboutChasm)
	http.HandleFunc("/ajax", fetchData)
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

	t := template.Must(template.New("about.html").Option("missingkey=error").ParseFiles("files/about.html"))
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
	if r.FormValue("c") == "region" {
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

	}

}

func centralDispatch(w http.ResponseWriter, r *http.Request) {

	if !DBInitialised {
		fmt.Printf("Running basic wizard\n")
		//sendTemplate(w, r, "wizard")
		showWizardPage(w, r, 1)
	}

}

/*****
func sendTemplate(w http.ResponseWriter, r *http.Request, tmplt string) {

	fmt.Printf("%v\n", r.URL)

	var cfg struct {
	}

	t := template.Must(template.New(tmplt + ".html").Option("missingkey=error").ParseFiles("files/" + tmplt + ".html"))

	err := t.Execute(w, cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
*****/
