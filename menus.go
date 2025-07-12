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
	{"Current rankings", "Show state of play", "/qlist", ""},
	{"Finisher certificates", "Print Finisher certificates", "/certs", "window.open(#url#,'certs')"},
	{"Check-OUT", "Check-out at start of rally", "/odos?check=out", ""},
	{"Check-IN", "Check-in at end of rally", "/odos?check=in", ""},
	{"Rally setup &amp; config", "Access all components", "/setup", ""},
}

var setupmenu = menu{
	{"Rally Parameters", "Configuration of this rally", "/config", ""},
	{"Edit certificate content", "Maintain certificate templates", "/editcert", ""},
	{"Entrants", "Maintain entrant details", "/menu?menu=entrants", ""},
	{"Bonuses", "Ordinary and combo bonuses", "/menu?menu=bonuses", ""},
	{"Time penalties", "Time penalties", "/timep", ""},
	{"Advanced setup", "Advanced configuration", "/menu?menu=asetup", ""},
}

var entrantmenu = menu{
	{"Full entrant records", "All details of entrants", "/entrants", ""},
	{"Odometer check-OUT", "Check-out at start of rally", "/odos?check=out", ""},
	{"Odometer check-IN", "Check-in at end of rally", "/odos?check=in", ""},
	{"Teams", "Maintain teams and membership", "/teams", ""},
	{"Cohorts", "Split entrants into cohorts", "/cohorts", ""},
	{"Import entrants", "Load entrants from spreadsheet", "/import?type=entrants", ""},
}

var bonusmenu = menu{
	{"Ordinary bonuses", "Ordinary bonuses", "/bonuses", ""},
	{"Combos", "Combination bonuses", "/combos", ""},
	{"Import bonuses", "Load ordinary bonuses from spreadsheet", "/import?type=bonuses", ""},
	{"Import combos", "Load combinations from spreadsheet", "/import?type=combos", ""},
}

var advancedmenu = menu{
	{"Categories", "Categories for use with compound rules", "/niy", ""},
	{"Complex calculation rules", "Scoring rules for use with categories", "/niy", ""},
	{"Classes", "Classes", "/niy", ""},
	{"Recalculate scorecards", "Recalculate scorecards", "/recalc", ""},
}
var menus = map[string]*menu{"main": &mainMenu, "setup": &setupmenu, "entrants": &entrantmenu, "bonuses": &bonusmenu, "asetup": &advancedmenu}

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
