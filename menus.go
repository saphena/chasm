package main

import (
	"fmt"
	"net/http"
	"strings"
)

type menuItem struct {
	Tag   string
	Hint  string
	Url   string
	Click string
}

type menu []menuItem

var mainMenu = menu{
	{"Claims log", "Full list of processed claims", "/claims", ""},
	{"EBC claims judging", "Process incoming claims from email", "/ebclist", ""},
	{"Review scorecards", "", "/cards", ""},
	{"Current rankings", "Show state of play", "/qlist", "window.open(#url#,'qlist')"},
	{"Finisher certificates", "Print Finisher certificates", "/fcerts", "window.open(#url#,'fcerts')"},
	{"Rally setup &amp; config", "Acess all components", "/setup", ""},
}

var menus = map[string]*menu{"main": &mainMenu}

func showMenu(w http.ResponseWriter, menu string) {

	m, ok := menus[menu]
	if !ok {
		return
	}

	fmt.Fprint(w, `<nav class="menu">`)
	for _, v := range *m {
		onclick := ""
		if v.Click == "" {
			onclick = "window.location.href='" + v.Url + "'"
		} else {
			onclick = strings.ReplaceAll(v.Click, "#url#", "'"+v.Url+"'")
		}
		fmt.Fprintf(w, `<button onclick="%v" title="%v">%v</button>`, onclick, v.Hint, v.Tag)
	}
	fmt.Fprint(w, `</nav>`)

}
