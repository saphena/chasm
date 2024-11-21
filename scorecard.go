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
	sqlx := "SELECT RiderName,ifnull(PillionName,''),FinishPosition,TeamID,CorrectedMiles,EntrantStatus,ifnull(Scorex,''),ifnull(AvgSpeed,''),TotalPoints"
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
	if CS.RallyUnitKms {
		mk = CS.UnitKmsLit
	}

	fmt.Fprintf(w, `<div class="scorecard">`)
	fmt.Fprintf(w, `<div class="topline noprint"><span>#%v %v</span><span>%v %v</span><span>%v points</span><span>%v</span></div>`, entrant, team, sr.Miles, mk, sr.Points, EntrantStatusLits[sr.Status])
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `</header>`)
	fmt.Fprintf(w, `<div class="scorex">%v</div>`, sr.Scorex)

}
