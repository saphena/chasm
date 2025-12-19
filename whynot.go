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
window.location.href='/ynot/'+c.value+'/'+e.value;
}
}
</script>
`
const whynotTemplate = `
<article class="whynot">
<h2>Individual Combo Analysis</h2>
<div>
	<label for="ComboID">Combo</label> 
	<select id="ComboID" onchange="reloadYnot()">
	{{range .Combos}}
		<option value="{{.ComboID}}" {{if eq .ComboID $.ComboID}}selected{{end}}>[ {{.ComboID}} ] {{.ComboDesc}}</option>
	{{end}}
	</select>
	 - {{if .ClaimedOK}}Claimed successfully{{else}}Combo not scored{{end}}
</div>
<div>
	<label for="EntrantID">Entrant</label> 
	<select id="EntrantID" onchange="reloadYnot()">
	{{range .Entrants}}
		<option value="{{.EntrantID}}" {{if eq .EntrantID $.EntrantID}}selected{{end}}>[ {{.EntrantID}} ] {{.EntrantName}}</option>
	{{end}}
	</select>
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

	combo := r.PathValue("combo")
	entrant := intval(r.PathValue("entrant"))

	sqlx := "SELECT Bonuses FROM combos WHERE ComboID='" + combo + "'"
	bonuses := strings.Split(getStringFromDB(sqlx, ""), ",")

	var y ynot
	y.ComboID = combo
	y.EntrantID = entrant
	y.ComboDesc = getStringFromDB("SELECT BriefDesc FROM combos WHERE ComboID='"+combo+"'", combo)
	y.MinTicks = getIntegerFromDB("SELECT MinimumTicks FROM combos WHERE ComboID='"+combo+"'", 0)
	if y.MinTicks < 1 {
		y.MinTicks = len(bonuses)
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
