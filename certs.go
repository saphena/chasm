package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

//go:embed certs.css
var certscss string

//go:embed certedit.js
var certeditjs string

//go:embed jodit.css
var joditcss string

//go:embed jodit.js
var joditjs string

type CertFields struct {
	Bike       string
	CrewName   string
	Distance   int
	EntrantID  int
	Place      int
	Points     int
	Rank       string
	MKlit      string
	RallyTitle string
	TeamName   string
}

type TeamRec struct {
	TeamID   int
	TeamName string
}
type TeamMap map[int]string

var certsHTMLheader = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>chasm</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="/css?file=certscss">
</head>
<body>
`

var certeditHTMLheader = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>chasm</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="/css?file=maincss">
<link rel="stylesheet" href="/css?file=joditcss">
<link rel="stylesheet" href="/css?file=certscss">
<script src="/js?file=mainscript"></script>
<script src="/js?file=joditjs"></script>
<script src="/js?file=certeditjs" defer="defer"></script>
</head>
<body>
`

func edit_certificate(w http.ResponseWriter, r *http.Request) {

	entrant := intval(r.FormValue("e"))
	class := intval(r.FormValue("class"))
	sqlx := "SELECT ifnull(html,'') FROM certificates"
	if r.FormValue("class") != "" {
		sqlx += " WHERE Class=" + strconv.Itoa(class)
	}
	html := getStringFromDB(sqlx, "")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Fprint(w, certeditHTMLheader)
	fmt.Fprint(w, `<header>`)
	showTopbar(w, "Finisher certificate")

	fmt.Fprint(w, `<form action="/savecert">`)
	fmt.Fprintf(w, `<input type="hidden" name="class" value="%v">`, class)
	fmt.Fprintf(w, `<input type="hidden" name="e" value="%v">`, entrant)
	fmt.Fprint(w, `<button id="savecert" disabled="disabled">Save changes</button>`)
	fmt.Fprint(w, `</header>`)
	fmt.Fprintf(w, `<article class="certificate"><textarea id="editor" name="html">%v</textarea></article>`, html)
	fmt.Fprint(w, `</form>`)

}
func load_teams() TeamMap {

	sqlx := "SELECT TeamID,ifnull(BriefDesc,TeamID) FROM teams "

	ta := make(TeamMap)

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var tr TeamRec
		err = rows.Scan(&tr.TeamID, &tr.TeamName)
		checkerr(err)
		ta[tr.TeamID] = tr.TeamName
	}
	rows.Close()
	return ta
}

func ordinal_ranklit(n int) string {

	suffix := ""
	switch n {
	case 11:
	case 12:
	case 13:
		suffix = "th"
	default:
		mod := n % 10
		switch mod {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		default:
			suffix = "th"
		}
	}
	return fmt.Sprintf("%v<sup>%v</sup>", n, suffix)

}
func print_certs(w http.ResponseWriter, r *http.Request) {

	sqlx := "SELECT ifnull(Bike,''), ifnull(RiderFirst,''),ifnull(RiderLast,''),ifnull(PillionFirst,''),ifnull(PillionLast,''),ifnull(CorrectedMiles,0),EntrantID"
	sqlx += ",ifnull(FinishPosition,0),ifnull(TeamID,0),TotalPoints"
	sqlx += " FROM entrants "
	if r.FormValue("all") == "" {
		sqlx += " WHERE EntrantStatus=" + strconv.Itoa(EntrantFinisher)
	}
	sqlx += " ORDER BY FinishPosition DESC"

	mklit := CS.UnitMilesLit
	if CS.Basics.RallyUnitKms {
		mklit = CS.UnitKmsLit
	}
	rallyTitle := CS.Basics.RallyTitle

	teams := load_teams()

	html := getStringFromDB("SELECT ifnull(html,'') FROM certificates", "")
	htmltmplt := fmt.Sprintf(`<div class="certframe"><div class="certificate">%v</div></div>`, html)
	tmplt, err := template.New("certs").Parse(htmltmplt)
	checkerr(err)

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Fprint(w, certsHTMLheader)

	for rows.Next() {
		var cf CertFields
		var pillion string
		var team int
		var rf string
		var rl string
		var pf string
		var pl string
		err = rows.Scan(&cf.Bike, &rf, &rl, &pf, &pl, &cf.Distance, &cf.EntrantID, &cf.Place, &team, &cf.Points)
		checkerr(err)
		cf.CrewName = strings.TrimSpace(rf + " " + rl)
		pillion = strings.TrimSpace(pf + " " + pl)
		cf.TeamName = teams[team]
		if pillion != "" {
			cf.CrewName += " &amp; " + pillion
		}
		cf.Rank = ordinal_ranklit(cf.Place)
		cf.MKlit = mklit
		cf.RallyTitle = rallyTitle

		err = tmplt.Execute(w, cf)
		checkerr(err)
	}

}

func save_certificate(w http.ResponseWriter, r *http.Request) {

	fmt.Printf(`%v`, r)

	html := r.FormValue(("html"))
	entrant := intval(r.FormValue("e"))
	class := intval(r.FormValue("class"))

	sqlx := "UPDATE certificates SET html=? WHERE EntrantID=? AND Class=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(html, entrant, class)
	checkerr(err)
	startHTML(w, "certificate saved")
	showMenu(w, "setup")
}
