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
type ynot struct {
	ComboID     string
	ComboDesc   string
	EntrantID   int
	EntrantName string
	ClaimedOK   bool
	Bonuses     []ynotbonus
}

const whynotTemplate = `
<article class="whynot">
<h2>Individual Combo Analysis</h2>
<div>
	Combo <strong>{{.ComboID}}</strong> <em>{{.ComboDesc}}</em> - {{if .ClaimedOK}}Claimed successfully{{else}}Combo not scored{{end}}
</div>
<div>
	Entrant <strong>{{.EntrantID}}</strong> <em>{{.EntrantName}}</em>
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
</div>

</article>
`

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
	y.EntrantName = getStringFromDB(fmt.Sprintf("SELECT "+RiderNameSQL+" FROM entrants WHERE EntrantID=%v", entrant), strconv.Itoa(entrant))
	startHTML(w, "WhyNot?")

	sqlx = "SELECT Decision FROM claims WHERE BonusID=? AND EntrantID=? ORDER BY ClaimTime DESC,OdoReading DESC"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	y.ClaimedOK = true

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
			} else if d > 0 {
				bd.Status = BadClaimIcon
			} else {
				bd.Status = UndecidedClaimIcon
			}
		} else {
			bd.Status = NoClaimIcon
		}
		rows.Close()
		y.ClaimedOK = y.ClaimedOK && bd.Status == GoodClaimIcon
		y.Bonuses = append(y.Bonuses, bd)
	}
	stmt.Close()

	t, err := template.New("ynot").Parse(whynotTemplate)
	checkerr(err)
	err = t.Execute(w, y)
	checkerr(err)
}
