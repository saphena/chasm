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
	Speed    string
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
	fmt.Fprintf(w, `<input type="hidden" name="ok" value="%v">`, r.FormValue("ok"))
	fmt.Fprintf(w, `<span class="pushright"><label for="plusok">+ok</label> <input type="checkbox"  id="plusok" %v onchange="showQOkChanged(this)"></span>`, checked)
	checked = ""
	showspeed := r.FormValue("speed") != ""
	if showspeed {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="hidden" name="speed" value="%v">`, r.FormValue("speed"))
	fmt.Fprintf(w, `<span class="pushright"><label for="showspeed">+speed</label> <input type="checkbox" id="showspeed" %v onchange="showQSpeedChanged(this)"></span>`, checked)

	checked = ""
	if showHot {
		checked = "checked"
	}
	fmt.Fprintf(w, `<span class="pushright"><label for="showhot">%v</label> <input type="checkbox" id="showhot" %v onchange="showQHotChanged(this)"></span>`, hot_icon, checked)
	fmt.Fprint(w, `</form></div>`)

	fmt.Fprint(w, `<div class="rankhead">`)

	showReloadTicker(w, r.URL.String())

	fmt.Fprint(w, `</div></header>`)
	fmt.Fprint(w, `<div class="rankings">`)

	fmt.Fprint(w, `<form id="rankingsfrm">`)
	fmt.Fprintf(w, `<input type="hidden" name="seq" value="%v">`, seq)
	fmt.Fprintf(w, `<input type="hidden" name="desc" value="%v"`, desc)
	fmt.Fprint(w, `</form>`)

	fmt.Fprint(w, `<fieldset class="row hdr rankings">`)
	fmt.Fprint(w, `<fieldset class="col hdr mid link" onclick="reloadRankings('seq','FinishPosition')">Rank</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col hdr link" onclick="reloadRankings('seq','RiderName')">Name</fieldset>`)
	mk := "Miles"
	mph := "MPH"
	if CS.Basics.RallyUnitKms {
		mk = "Kms"
		mph = "km/h"
	}
	fmt.Fprintf(w, `<fieldset class="col hdr mid link" onclick="reloadRankings('seq','CorrectedMiles')">%v</fieldset>`, mk)
	fmt.Fprint(w, `<fieldset class="col hdr right link" onclick="reloadRankings('seq','TotalPoints')">Points</fieldset>`)
	fmt.Fprintf(w, `<fieldset class="col hdr right link" onclick="reloadRankings('seq','PPM')">P&divide;%v</fieldset>`, string(mk[0]))
	if showspeed {
		fmt.Fprintf(w, `<fieldset class="col hdr right link" onclick="reloadRankings('seq','AvgSpeed')">%v</fieldset>`, mph)
	}
	fmt.Fprint(w, `</fieldset>`)

	sqlx := "SELECT ifnull(FinishPosition,0),RiderName,ifnull(PillionName,''),ifnull(CorrectedMiles,0),ifnull(TotalPoints,0),EntrantStatus"
	sqlx += ",IfNull((TotalPoints*1.0) / CorrectedMiles,0) As PPM,ifnull(AvgSpeed,''),EntrantID"
	sqlx += " FROM entrants"
	sqlx += " WHERE EntrantStatus IN ("
	sqlx += strconv.Itoa(EntrantFinisher) + "," + strconv.Itoa(EntrantDNF)
	if showok {
		sqlx += "," + strconv.Itoa(EntrantOK)
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

		err = rows.Scan(&rr.Rank, &rr.Name, &pn, &rr.Distance, &rr.Points, &rr.Status, &rr.Eff, &rr.Speed, &rr.Entrant)
		checkerr(err)
		if pn != "" {
			rr.Name += " &amp; " + pn
		}
		fmt.Fprintf(w, `<fieldset class="row rankings link" onclick="window.location.href='/score?e=%v'">`, rr.Entrant)
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
			eff = fmt.Sprintf("%.1f", rr.Eff)
		}
		fmt.Fprintf(w, `<fieldset class="col right">%v</fieldset>`, eff)
		if showspeed {
			fmt.Fprintf(w, `<fieldset class="col right">%v</fieldset>`, rr.Speed)
		}
		fmt.Fprint(w, `</fieldset>`)
	}
	fmt.Fprint(w, `</div>`)
}
