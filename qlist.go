package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type RankRecord struct {
	Rank     int
	Status   int
	Name     string
	Distance int
	Points   int
	Eff      float64
	Entrant  int
}

func show_qlist(w http.ResponseWriter, r *http.Request) {

	const hot_icon = "&#9832;"

	startHTML(w, "Rankings")

	fmt.Fprintf(w, `<div class="topline left"><form id="optionsfrm">`)

	hot := r.FormValue("hot")
	showHot := hot != ""
	fmt.Fprintf(w, `<input type="hidden" name="hot" value="%v">`, hot)
	seq := r.FormValue("seq")
	if seq == "" {
		seq = "FinishPosition"
	}
	fmt.Fprintf(w, `<input type="hidden" name="seq" value="%v">`, seq)
	desc := r.FormValue("desc")
	fmt.Fprintf(w, `<input type="hidden" name="desc" value="%v">`, desc)

	checked := ""
	showok := r.FormValue("ok") != ""
	if showok {
		checked = "checked"
	}

	fmt.Fprint(w, `<span class="pushright"><label for="plusfinisher">Finishers</label> `)
	fmt.Fprint(w, `<input type="checkbox" disabled id="plusfinisher" checked> </span> `)

	fmt.Fprintf(w, `<input type="hidden" name="ok" value="%v">`, r.FormValue("ok"))
	fmt.Fprintf(w, `<span class="pushright"><label for="plusok">+ok</label> <input type="checkbox"  id="plusok" %v onchange="showQOkChanged(this)"></span>`, checked)
	checked = ""
	showdnf := r.FormValue("dnf") != ""
	if showdnf {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="hidden" name="dnf" value="%v">`, r.FormValue("dnf"))
	fmt.Fprintf(w, `<span class="pushright"><label for="plusdnf">+DNF</label> <input type="checkbox"  id="plusdnf" %v onchange="showQDnfChanged(this)"></span>`, checked)

	checked = ""
	if showHot {
		checked = "checked"
	}
	fmt.Fprintf(w, `<span class="pushright" title="Show in order of finish time"><label for="showhot">%v</label> <input type="checkbox" id="showhot" %v onchange="showQHotChanged(this)"></span>`, hot_icon, checked)
	fmt.Fprint(w, `</form></div>`)

	fmt.Fprint(w, `<div class="rankhead">`)

	showReloadTicker(w, r.URL.String())

	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="rankings">`)

	fmt.Fprint(w, `<form id="rankingsfrm">`)
	fmt.Fprintf(w, `<input type="hidden" name="seq" value="%v">`, seq)
	fmt.Fprintf(w, `<input type="hidden" name="desc" value="%v"`, desc)
	fmt.Fprint(w, `</form>`)

	fmt.Fprint(w, `<fieldset class="row hdr rankings">`)
	fmt.Fprint(w, `<fieldset class="col hdr mid sort" onclick="reloadRankings('seq','FinishPosition')">Rank</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col hdr sort" onclick="reloadRankings('seq','RiderName')">Name</fieldset>`)

	mk := CS.UnitMilesLit
	if CS.Basics.RallyUnitKms {
		mk = CS.UnitKmsLit

	}
	fmt.Fprintf(w, `<fieldset class="col hdr mid sort" onclick="reloadRankings('seq','CorrectedMiles')">%v</fieldset>`, mk)
	fmt.Fprint(w, `<fieldset class="col hdr right sort" onclick="reloadRankings('seq','TotalPoints')">Points</fieldset>`)
	fmt.Fprintf(w, `<fieldset class="col hdr right sort" onclick="reloadRankings('seq','PPM')">P&divide;%v</fieldset>`, string(mk[0]))
	fmt.Fprint(w, `</fieldset></div><!-- rankings --><hr></header>`)

	sqlx := "SELECT ifnull(FinishPosition,0)," + RiderNameSQL + "," + PillionNameSQL + ",ifnull(CorrectedMiles,0),ifnull(TotalPoints,0),EntrantStatus"
	sqlx += ",IfNull((TotalPoints*1.0) / CorrectedMiles,0) As PPM,EntrantID"
	sqlx += " FROM entrants"
	sqlx += " WHERE EntrantStatus IN ("
	sqlx += strconv.Itoa(EntrantFinisher)
	if showok {
		sqlx += "," + strconv.Itoa(EntrantOK)
	}
	if showdnf {
		sqlx += "," + strconv.Itoa(EntrantDNF)
	}
	sqlx += ")"
	sqlx += " ORDER BY "
	if showHot {
		sqlx += "FinishTime DESC, "
	}
	sqlx += " EntrantStatus DESC," + seq + " " + desc

	//	fmt.Println(sqlx)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var rr RankRecord
		var pn string

		err = rows.Scan(&rr.Rank, &rr.Name, &pn, &rr.Distance, &rr.Points, &rr.Status, &rr.Eff, &rr.Entrant)
		checkerr(err)
		if pn != "" {
			rr.Name += " &amp; " + pn
		}
		fmt.Fprintf(w, `<fieldset class="row rankings link" onclick="window.location.href='/score?e=%v&back=/qlist'">`, rr.Entrant)
		status := EntrantStatusLits[rr.Status]
		if rr.Status == EntrantFinisher {
			status = strconv.Itoa(rr.Rank)
		}
		fmt.Fprintf(w, `<fieldset class="col mid">%v</fieldset>`, status)
		fmt.Fprintf(w, `<fieldset class="col">%v</fieldset>`, rr.Name)
		fmt.Fprintf(w, `<fieldset class="col mid">%v</fieldset>`, rr.Distance)
		fmt.Fprintf(w, `<fieldset class="col right">%v</fieldset>`, rr.Points)
		eff := ""
		if rr.Distance > 0 && rr.Points != 0 {
			eff = fmtDecimal("%.1f", rr.Eff)
		}
		fmt.Fprintf(w, `<fieldset class="col right">%v</fieldset>`, eff)

		fmt.Fprint(w, `</fieldset>`)
	}
	fmt.Fprint(w, `</div>`)
}
