package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type ScorecardRec struct {
	EntrantID   int
	RiderName   string
	PillionName string
	Rank        int
	TeamID      int
	Miles       int
	Status      int
	Scorex      string
	Speed       string
	Points      int
}

func showScorecard(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	entrant := intval(r.FormValue("e"))
	if entrant < 1 {
		return
	}
	sqlx := "SELECT " + RiderNameSQL
	sqlx += ",ifnull(PillionName,''),FinishPosition,TeamID,ifnull(CorrectedMiles,0),EntrantStatus,ifnull(Scorex,''),ifnull(AvgSpeed,''),TotalPoints"
	sqlx += " FROM entrants WHERE EntrantID=" + strconv.Itoa(entrant)

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	if !rows.Next() {
		return
	}

	var sr ScorecardRec
	err = rows.Scan(&sr.RiderName, &sr.PillionName, &sr.Rank, &sr.TeamID, &sr.Miles, &sr.Status, &sr.Scorex, &sr.Speed, &sr.Points)
	checkerr(err)
	team := sr.RiderName
	if sr.PillionName != "" {
		team += " &amp; " + sr.PillionName
	}

	startHTML(w, "Scorecard")

	mk := CS.UnitMilesLit
	if CS.Basics.RallyUnitKms {
		mk = CS.UnitKmsLit
	}

	fmt.Fprintf(w, `<div class="scorecard">`)
	fmt.Fprintf(w, `<div class="topline noprint"><span>#%v %v</span><span>%v %v</span><span>%v points</span><span>%v</span></div>`, entrant, team, sr.Miles, mk, sr.Points, EntrantStatusLits[sr.Status])
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `</header>`)
	fmt.Fprintf(w, `<div class="scorex" title="Doubleclick to shows claims" ondblclick="window.location.href='/claims?esel=%v'">%v</div>`, entrant, sr.Scorex)

}

func showScorecards(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	startHTML(w, "Scorecards")

	sqlx := "SELECT ifnull(RiderFirst,''),ifnull(RiderLast,''),ifnull(PillionName,''),EntrantID FROM entrants ORDER BY RiderLast,RiderFirst"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	fmt.Fprint(w, `<form action="/score">`)
	fmt.Fprint(w, `<label for="EntrantID">Entrant</label> `)
	fmt.Fprint(w, `<input type="number" id="EntrantID" class="EntrantID" name="e"> `)
	fmt.Fprint(w, `<input type="submit" value=" show ">`)
	fmt.Fprint(w, `</form>`)
	fmt.Fprint(w, `<article class="reviewlist">`)

	for rows.Next() {
		var e EntrantDetails
		err = rows.Scan(&e.RiderFirst, &e.RiderLast, &e.PillionName, &e.EntrantID)
		checkerr(err)
		fmt.Fprintf(w, `<div class="row link" onclick="window.location.href='/score?e=%v';">`, e.EntrantID)
		fmt.Fprintf(w, `<span class="col">%v</span>`, e.EntrantID)
		x := ""
		if e.PillionName != "" {
			x = " &amp; " + e.PillionName
		}
		fmt.Fprintf(w, `<span class="col"><strong>%v</strong>, %v %v</span>`, e.RiderLast, e.RiderFirst, x)
		sqlx = fmt.Sprintf("SELECT count(DISTINCT BonusID) FROM claims WHERE EntrantID=%v", e.EntrantID)
		nc := getIntegerFromDB(sqlx, 0)
		sqlx += " AND Decision>0"
		nr := getIntegerFromDB(sqlx, 0)
		fmt.Fprintf(w, `<span class="col" title="Bonus claims">%v <strong>%v</strong></span>`, printNZ(nc), printNZ(nr))
		fmt.Fprint(w, `</div>`)
	}
	fmt.Fprint(w, `</article>`)

}
