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
}

func show_qlist(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Rankings")

	fmt.Fprint(w, `<div class="rankhead">`)

	showReloadTicker(w, r.URL.String())

	fmt.Fprint(w, `</div></header>`)
	fmt.Fprint(w, `<div class="rankings">`)

	seq := r.FormValue("seq")
	if seq == "" {
		seq = "FinishPosition"
	}
	desc := r.FormValue("desc")

	fmt.Fprint(w, `<form id="rankingsfrm">`)
	fmt.Fprintf(w, `<input type="hidden" name="seq" value="%v">`, seq)
	fmt.Fprintf(w, `<input type="hidden" name="desc" value="%v"`, desc)
	fmt.Fprint(w, `</form>`)

	fmt.Fprint(w, `<fieldset class="row hdr rankings">`)
	fmt.Fprint(w, `<fieldset class="col hdr mid link" onclick="reloadRankings('seq','FinishPosition')">Rank</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col hdr link" onclick="reloadRankings('seq','RiderName')">Name</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col hdr mid link" onclick="reloadRankings('seq','CorrectedMiles')">Distance</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col hdr right link" onclick="reloadRankings('seq','TotalPoints')">Points</fieldset>`)
	fmt.Fprint(w, `</fieldset>`)

	sqlx := "SELECT ifnull(FinishPosition,0),RiderName,ifnull(PillionName,''),ifnull(CorrectedMiles,0),ifnull(TotalPoints,0),EntrantStatus"
	sqlx += " FROM entrants"
	sqlx += " ORDER BY EntrantStatus DESC," + seq + " " + desc

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var rr RankRecord
		var pn string
		err = rows.Scan(&rr.Rank, &rr.Name, &pn, &rr.Distance, &rr.Points, &rr.Status)
		checkerr(err)
		if pn != "" {
			rr.Name += " &amp; " + pn
		}
		fmt.Fprint(w, `<fieldset class="row rankings">`)
		status := EntrantStatusLits[rr.Status]
		if rr.Status == EntrantFinisher {
			status = strconv.Itoa(rr.Rank)
		}
		fmt.Fprintf(w, `<fieldset class="col mid">%v</fieldset>`, status)
		fmt.Fprintf(w, `<fieldset class="col">%v</fieldset>`, rr.Name)
		fmt.Fprintf(w, `<fieldset class="col mid">%v</fieldset>`, rr.Distance)
		fmt.Fprintf(w, `<fieldset class="col right">%v</fieldset>`, rr.Points)
		fmt.Fprint(w, `</fieldset>`)
	}
	fmt.Fprint(w, `</div>`)
}
