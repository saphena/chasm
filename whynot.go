package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type ynotbonus struct {
	BonusID   string
	BonusDesc string
	Status    string
}
type ynotcombo struct {
	ComboID   string
	ComboDesc string
}
type ynotentrant struct {
	EntrantID   int
	EntrantName string
}
type ynot struct {
	ComboID     string
	ComboDesc   string
	MinTicks    int
	MaxTicks    int
	ScoredTicks int
	EntrantID   int
	EntrantName string
	Combos      []ynotcombo
	Entrants    []ynotentrant
	ClaimedOK   bool
	Bonuses     []ynotbonus
}

const reloadYnot = `
<script>
function reloadYnot() {
let c = document.getElementById('ComboID');
let e = document.getElementById('EntrantID');
if (c && e) {
window.location.href='/ynot?c='+c.value+'&e='+e.value;
}
}
</script>
`

const whynotTemplate = `
<div id="whynothelp" class="popover" popover>
<h1>Why combo not scored</h1>
<p>This helps answer the question "<em>why has this rider not scored this combo?</em>".</p>
<p>This shows, for each underlying bonus of the combo, whether or not this rider has a good claim for it, highlighting what's missing from the score.</p>
<p>Only the effective (most recent) claim is onsidered in each case. To check for earlier claims, examine the claims log for the rider.</p>
</div>
<article class="whynot">
<h2>Individual Combo Analysis</h2>
<div>This is useful for understanding why a particular combo has not yet been scored by an entrant. <input type="button" class="popover" popovertarget="whynothelp" value="[more details here]"></div>
<div>
	<label for="ComboID">Combo</label> 
	<select id="ComboID" onchange="reloadYnot()">
	{{range .Combos}}
		<option value="{{.ComboID}}" {{if eq .ComboID $.ComboID}}selected{{end}}>[ {{.ComboID}} ] {{.ComboDesc}}</option>
	{{end}}
	</select>
	 - needs {{.MinTicks}} of {{.MaxTicks}}
	 - {{if .ClaimedOK}} Claimed successfully {{else}} Combo not scored &#9785; {{end}}
	 <button onclick="window.location.href='/combo?c={{.ComboID}}'">Combo details</button>
</div>
<div>
	<label for="EntrantID">Entrant</label> 
	<select id="EntrantID" onchange="reloadYnot()">
	{{range .Entrants}}
		<option value="{{.EntrantID}}" {{if eq .EntrantID $.EntrantID}}selected{{end}}>[ {{.EntrantID}} ] {{.EntrantName}}</option>
	{{end}}
	</select>
	<button onclick="loadPage('/claims?esel={{.EntrantID}}')">Claims log</button>
</div>
<hr>
<div class="ynotbonuslist">
	{{range .Bonuses}}
	<div>
	<span>{{.BonusID}}</span>
	<span>{{.BonusDesc}}</span>
	<span>{{.Status}}</span>
	</div>
	{{end}}
	<hr>
	<div><span>Bonuses scored: {{.ScoredTicks}}</span> <span>Minimum needed: {{.MinTicks}}</span></div>

</div>

</article>
`

func ynotcombos() []ynotcombo {

	sqlx := "SELECT ComboID,BriefDesc FROM combos ORDER BY ComboID"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	res := make([]ynotcombo, 0)
	for rows.Next() {
		var y ynotcombo
		err = rows.Scan(&y.ComboID, &y.ComboDesc)
		checkerr(err)
		res = append(res, y)
	}
	return res
}

func ynotentrants() []ynotentrant {

	sqlx := "SELECT EntrantID," + RiderNameSQL + " FROM entrants ORDER BY EntrantID"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	res := make([]ynotentrant, 0)
	for rows.Next() {
		var y ynotentrant
		err = rows.Scan(&y.EntrantID, &y.EntrantName)
		checkerr(err)
		res = append(res, y)
	}
	return res
}

func whynotcombo(w http.ResponseWriter, r *http.Request) {

	const GoodClaimIcon = `Good claim &#10004;`
	const BadClaimIcon = `Claim rejected &#10008;`
	const UndecidedClaimIcon = `Claim not yet decided`
	const NoClaimIcon = `no claim found &#9785;`

	combo := r.FormValue("c")
	entrant := intval(r.FormValue("e"))

	if combo == "" {
		combo = getStringFromDB("SELECT ComboID FROM combos ORDER BY ComboID", "") // Get first combo
	}
	if entrant < 1 {
		entrant = getIntegerFromDB("SELECT EntrantID FROM entrants ORDER BY EntrantID", 0) // Get first entrant
	}
	sqlx := "SELECT Bonuses FROM combos WHERE ComboID='" + combo + "'"
	bonuses := strings.Split(getStringFromDB(sqlx, ""), ",")

	var y ynot
	y.ComboID = combo
	y.EntrantID = entrant
	y.ComboDesc = getStringFromDB("SELECT BriefDesc FROM combos WHERE ComboID='"+combo+"'", combo)
	y.MaxTicks = len(bonuses)
	y.MinTicks = getIntegerFromDB("SELECT MinimumTicks FROM combos WHERE ComboID='"+combo+"'", 0)
	if y.MinTicks < 1 {
		y.MinTicks = y.MaxTicks
	}
	y.EntrantName = getStringFromDB(fmt.Sprintf("SELECT "+RiderNameSQL+" FROM entrants WHERE EntrantID=%v", entrant), strconv.Itoa(entrant))
	y.Combos = ynotcombos()
	y.Entrants = ynotentrants()
	startHTML(w, "WhyNot?")

	fmt.Fprint(w, reloadYnot)
	sqlx = "SELECT Decision FROM claims WHERE BonusID=? AND EntrantID=? ORDER BY ClaimTime DESC,OdoReading DESC"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	y.ClaimedOK = false

	for _, b := range bonuses {
		var bd ynotbonus
		bd.BonusID = b
		bd.BonusDesc = getStringFromDB("SELECT BriefDesc FROM bonuses WHERE BonusID='"+b+"'", b)
		rows, err := stmt.Query(b, entrant)
		checkerr(err)
		defer rows.Close()
		if rows.Next() {
			var d int
			err = rows.Scan(&d)
			checkerr(err)
			if d == 0 {
				bd.Status = GoodClaimIcon
				y.ScoredTicks++
			} else if d > 0 {
				bd.Status = BadClaimIcon
			} else {
				bd.Status = UndecidedClaimIcon
			}
		} else {
			bd.Status = NoClaimIcon
		}
		rows.Close()
		y.Bonuses = append(y.Bonuses, bd)
	}
	stmt.Close()

	y.ClaimedOK = y.ScoredTicks >= y.MinTicks

	t, err := template.New("ynot").Parse(whynotTemplate)
	checkerr(err)
	err = t.Execute(w, y)
	checkerr(err)
}
