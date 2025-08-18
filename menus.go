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
	{"Judge new claims", "Process incoming claims from email", "/ebclist", ""},
	{"Review scorecards", "", "/cards", ""},
	{"Show current rankings", "Show state of play", "/qlist", ""},
	{"Show all claims", "Full list of processed claims", "/claims", ""},
	{"Check-OUT @ start", "Take odo readings at start of rally", "/odos?check=out", ""},
	{"Check-IN @ finish", "Take odo readings at end of rally", "/odos?check=in", ""},
	{"Reporting", "Certificates, analyses,exports", "/menu?menu=Reports", ""},
	{"Rally setup &amp; config", "Access all components", "/setup", ""},
}

var setupmenu = menu{
	{"Rally Parameters", "Configuration of this rally", "/config", ""},
	{"Edit certificate content", "Maintain certificate templates", "/editcert", ""},
	{"Entrants", "Maintain entrant details", "/menu?menu=Entrants", ""},
	{"Bonuses / Combos", "Ordinary and combo bonuses", "/menu?menu=Bonuses", ""},
	{"Time penalties", "Time penalties", "/timep", ""},
	{"Complex calculation rules", "Scoring rules for use with categories", "/rules", ""},
	{"Here be dragons!", "Rally database utilities", "/menu?menu=Utilities", ""},
}

var entrantmenu = menu{
	{"Full entrant records", "All details of entrants", "/entrants", ""},
	{"Teams", "Maintain teams and membership", "/teams", ""},
	{"Certificate Classes", "Classes", "/classes", ""},
	{"Import entrants", "Load entrants from spreadsheet", "/import?type=entrants", ""},
}

var bonusmenu = menu{
	{"Ordinary bonuses", "Ordinary bonuses", "/bonuses", ""},
	{"Bonus categories", "Categories for use with compound rules", "/cats", ""},
	{"Combos", "Combination bonuses", "/combos", ""},
	{"Import bonuses", "Load ordinary bonuses from spreadsheet", "/import?type=bonuses", ""},
	{"Import combos", "Load combinations from spreadsheet", "/import?type=combos", ""},
}

var advancedmenu = menu{
	{"Recalculate scorecards", "Recalculate scorecards", "/recalc", ""},
	{"Reset Rally", "Reset the rally database", "/reset", ""},
}
var reportsmenu = menu{
	{"Finisher certificates", "Print Finisher certificates", "/certs", "window.open(#url#,'certs')"},
	{"Current rankings", "Show state of play", "/qlist", ""},
	{"Edit certificate content", "Maintain certificate templates", "/editcert", ""},
	{"Bonus analysis", "Exportable spreadsheet", "/report/ba", "window.open(#url#,'reports')"},
	{"Export Finishers CSV", "Download CSV of Finishers", "/report/fincsv", ""},
	{"Export Finishers JSON", "Download JSON of Finishers", "/report/finjson", ""},
}
var menus = map[string]*menu{"main": &mainMenu, "Setup": &setupmenu, "Entrants": &entrantmenu, "Bonuses": &bonusmenu, "Utilities": &advancedmenu, "Reports": &reportsmenu}

func show_menu(w http.ResponseWriter, r *http.Request) {

	menu := r.FormValue("menu")
	if menu == "" {
		menu = "main"
	}

	startHTML(w, menu)
	showMenu(w, menu)

}

func show_setup(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Setup")

	showMenu(w, "Setup")
}

func showMenu(w http.ResponseWriter, menu string) {

	m, ok := menus[menu]
	if !ok {
		return
	}

	fmt.Fprint(w, `</header>`)
	fmt.Fprint(w, `<nav class="menu">`)
	for _, v := range *m {
		onclick := ""
		if v.Click == "" {
			onclick = "window.location.href='" + v.Url + "'"
		} else {
			onclick = strings.ReplaceAll(v.Click, "#url#", "'"+v.Url+"'")
		}
		fmt.Fprintf(w, `<button class="menu" onclick="%v" title="%v">%v</button>`, onclick, v.Hint, v.Tag)
	}
	fmt.Fprint(w, `</nav>`)

}
