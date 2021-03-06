package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// PROGRAMVERSION is used a literal string in various places
const PROGRAMVERSION = "Chasm v0.0.0"

func main() {

	fmt.Printf("Serving %v\n", *HTTPPort)

	if !DBInitialised {
		fmt.Printf("Running basic wizard\n")
	}
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)
	http.HandleFunc("/about", aboutChasm)
	http.ListenAndServe(":"+*HTTPPort, nil)
}

func aboutChasm(w http.ResponseWriter, r *http.Request) {
	var cfg struct {
		Version string
		DBPath  string
		Schema  int
		Event   string
	}
	cfg.Version = PROGRAMVERSION
	cfg.DBPath = *DBNAME
	cfg.Schema = 1
	cfg.Event = "IBA Rally"

	t := template.Must(template.New("about.html").Option("missingkey=error").ParseFiles("about.html"))
	err := t.Execute(w, cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
