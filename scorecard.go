package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	rs_notreviewed  = 0
	rs_teamhappy    = 1
	rs_teamnothappy = 2
	rs_entranthappy = 3
)

type ScorecardRec struct {
	EntrantID         int
	RiderName         string
	PillionName       string
	Rank              int
	TeamID            int
	Miles             int
	Status            int
	Scorex            string
	Points            int
	ReviewedByTeam    int
	AcceptedByEntrant int
	LastReviewed      string
	ReviewStatus      int
}

func recalc_handler(w http.ResponseWriter, r *http.Request) {

	const recalcfrm = `
	<article class="recalc">
	<p>This procedure will recalculate all scorecards. This involves rebuilding them from scratch by reprocessing the claims log. This should only take a few moments but it will need exclusive access to the database.</p>
	<p>It's quite safe to do this during a live rally.</p>
	<form action="/recalc">
		<input type="hidden" name="ok" value="ok">
		<button autofocus>Recalculate scorecards</button>
	</form>
	</article>
	`
	startHTML(w, "Recalc scorecards")

	e := r.FormValue("e")
	ok := r.FormValue("ok")
	if ok == "ok" {
		if e == "" {
			recalc_all()
		} else {
			n, err := strconv.Atoi(e)
			if err != nil {
				fmt.Fprintf(w, `<p class="error">%v is not numeric</p>`, e)
				return
			}
			recalc_scorecard(n)
		}
		fmt.Fprint(w, `</header><p class="thatsall">Scorecards recalculated</p>`)
		return
	}
	fmt.Fprint(w, recalcfrm)

}

func showScorecard(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	entrant := intval(r.FormValue("e"))
	if entrant < 1 {
		showScorecards(w, r)
		return
	}
	sqlx := "SELECT " + RiderNameSQL
	sqlx += ",ifnull(PillionName,''),FinishPosition,TeamID,ifnull(CorrectedMiles,0),EntrantStatus,ifnull(Scorex,''),TotalPoints"
	sqlx += ",ReviewedByTeam,AcceptedByEntrant,ifnull(LastReviewed,''),ReviewStatus"
	sqlx += " FROM entrants WHERE EntrantID=" + strconv.Itoa(entrant)

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	if !rows.Next() {
		showScorecards(w, r)
		return
	}

	var sr ScorecardRec
	err = rows.Scan(&sr.RiderName, &sr.PillionName, &sr.Rank, &sr.TeamID, &sr.Miles, &sr.Status, &sr.Scorex, &sr.Points, &sr.ReviewedByTeam, &sr.AcceptedByEntrant, &sr.LastReviewed, &sr.ReviewStatus)
	checkerr(err)
	team := sr.RiderName
	if sr.PillionName != "" {
		team += " &amp; " + sr.PillionName
	}

	if r.FormValue("back") != "" {
		startHTMLBL(w, "Scorecard", r.FormValue("back"))
	} else {
		startHTML(w, "Scorecard")
	}

	mk := CS.UnitMilesLit
	if CS.Basics.RallyUnitKms {
		mk = CS.UnitKmsLit
	}

	fmt.Fprintf(w, `<div class="scorecard">`)
	fmt.Fprintf(w, `<div class="topline noprint"><span>#%v %v</span><span>%v %v</span><span>%v points</span><span>%v</span>`, entrant, team, sr.Miles, mk, sr.Points, EntrantStatusLits[sr.Status])
	fmt.Fprintf(w, `<select id="ReviewStatus" name="ReviewStatus" data-e="%v" onchange="saveRS(this);">`, entrant)
	sel := ""
	if sr.ReviewStatus == rs_notreviewed {
		sel = "selected"
	}
	fmt.Fprintf(w, `<option value="%v" %v>not reviewed</option>`, rs_notreviewed, sel)
	sel = ""
	if sr.ReviewStatus == rs_teamhappy {
		sel = "selected"
	}
	fmt.Fprintf(w, `<option value="%v" %v>Team happy</option>`, rs_teamhappy, sel)

	sel = ""
	if sr.ReviewStatus == rs_teamnothappy {
		sel = "selected"
	}
	fmt.Fprintf(w, `<option value="%v" %v>Team UNHAPPY</option>`, rs_teamnothappy, sel)
	sel = ""
	if sr.ReviewStatus == rs_entranthappy {
		sel = " selected"
	}
	fmt.Fprintf(w, `<option value="%v" %v>Rider AGREES</option>`, rs_entranthappy, sel)

	fmt.Fprint(w, `</select>`)
	fmt.Fprint(w, `</div>`) //topline
	fmt.Fprint(w, `</div>`) //scorecard
	fmt.Fprint(w, `</header>`)
	fmt.Fprintf(w, `<div class="scorex" title="Doubleclick to shows claims" ondblclick="window.location.href='/claims?esel=%v'">%v</div>`, entrant, sr.Scorex)

}

func showScorecards(w http.ResponseWriter, r *http.Request) {

	const chk = "&#10003;" //Regular checkmark
	const xxx = "&#10007;"
	const accepted = "<span class='bigtick'> &#10004;</span>" //Heavy checkmark
	const rejected = "&#10008;"

	r.ParseForm()

	startHTML(w, "Scorecards")

	sqlx := "SELECT ifnull(RiderFirst,''),ifnull(RiderLast,''),ifnull(PillionName,''),EntrantID,ReviewedByTeam,AcceptedByEntrant,ifnull(LastReviewed,''),ReviewStatus FROM entrants ORDER BY RiderLast,RiderFirst"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	fmt.Fprint(w, `<article class="reviewhdr"><form action="/score" class="reviewhdr">`)
	fmt.Fprint(w, `<label for="EntrantID">Entrant</label> `)
	fmt.Fprint(w, `<input type="number" autofocus id="EntrantID" class="EntrantID" name="e"> `)
	fmt.Fprint(w, `<input type="hidden" name="back" value="cards">`)
	fmt.Fprint(w, `<button> show </button>`)
	fmt.Fprint(w, `</form></article>`)

	fmt.Fprint(w, `<div class="reviewlist"><div class="row hdr">`)
	fmt.Fprint(w, `<span class="col hdr">Flag</span>`)
	fmt.Fprint(w, `<span class=" ">Name</span>`)
	fmt.Fprintf(w, `<span class="col hdr">Claims <span class="rejects">%v</span></span>`, rejected)
	fmt.Fprint(w, `</div></div><hr>`)
	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<article class="reviewlist">`)

	for rows.Next() {
		var e EntrantDetails
		err = rows.Scan(&e.RiderFirst, &e.RiderLast, &e.PillionName, &e.EntrantID, &e.ReviewedByTeam, &e.AcceptedByEntrant, &e.LastReviewed, &e.ReviewStatus)
		checkerr(err)
		fmt.Fprintf(w, `<div class="row link" onclick="window.location.href='/score?e=%v&back=cards';">`, e.EntrantID)
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
		fmt.Fprintf(w, `<span class="col nums" title="Bonus claims"><span class="numclaims">%v</span> <span class="rejects">%v</span></span>`, printNZ(nc), printNZ(nr))

		fmt.Fprint(w, `<span class="col" title="Review status">`)
		switch e.ReviewStatus {
		case rs_teamhappy:
			fmt.Fprint(w, chk)
		case rs_teamnothappy:
			fmt.Fprint(w, xxx)
		}
		if e.ReviewStatus == rs_entranthappy {
			fmt.Fprintf(w, ` %v`, accepted)
		}
		fmt.Fprint(w, `</span>`)

		fmt.Fprint(w, `</div>`)
	}
	fmt.Fprint(w, `</article>`)

}

// updateReviewStatus updates the entrant review status by ajax
func updateReviewStatus(w http.ResponseWriter, r *http.Request) {

	const ReviewDateFmt = "2006-01-02T15:04"

	rs := intval(r.FormValue("rs"))
	e := intval(r.FormValue("e"))
	lrdt := time.Now().Format(ReviewDateFmt)
	sqlx := fmt.Sprintf("UPDATE entrants SET ReviewStatus=%v,LastReviewed='%v' WHERE EntrantID=%v", rs, lrdt, e)
	//fmt.Println(sqlx)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)

}
