package main

import (
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
	http.HandleFunc("/about", aboutChasm)
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

func centralDispatch(w http.ResponseWriter, r *http.Request) {

	if !DBInitialised {
		fmt.Printf("Running basic wizard\n")
		//sendTemplate(w, r, "wizard")
		showWizardPage(w, r, 1)
	}

}

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
