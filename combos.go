package main

import (
	"fmt"
	"net/http"
)

const ComboScoreMethodMults = 1
const ComboScoreMethodPoints = 0

var tmpltSingleCombo = `
<div class="singlecombo">
	<form>
		<fieldset class="field">
			<label for="ComboID">ComboID</label>
			<input id="ComboID" name="Combo code" value="%s">
		</fieldset>
		<fieldset class="field">
			<label for="BriefDesc">Description</label>
			<input id="BriefDesc" name="BriefDesc" value="%s">
		</fieldset>
		<fieldset class="field">
			<label for="ScoreMethod">Value is</label>
			<select id="ScoreMethod" name="ScoreMethod">%s</select>
		</fieldset>
		<fieldset class="col">
			<label for="MinimumTicks">Minmum bonuses scored</label>
			<input id="MinimumTicks" name="MinimumTicks"
		</fieldset>
	</form>
</div>`

func show_combos(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<style>%s</style>`, maincss)
	fmt.Fprintf(w, `<script>%s</script>`, mainscript)

	ComboBonuses = loadCombos("")

	fmt.Fprint(w, `<div class="combos">`)
	for _, cb := range ComboBonuses {
		var mp string
		fmt.Fprint(w, `<fieldset class="row combo">`)
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, cb.Comboid)
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, cb.BriefDesc)
		mp = ""
		if cb.Compulsory {
			mp = "!"
		}
		fmt.Fprintf(w, `<fieldset class="col">%v</fieldset>`, mp)
		mp = ""
		if cb.ScoreMethod == ScoreMethodMults {
			mp = "x"
		}
		fmt.Fprintf(w, `<fieldset class="col">%s %v</fieldset>`, mp, cb.PointsList)
		mp = ""
		if cb.MinTicks < len(cb.Bonuses) {
			mp = fmt.Sprintf("[%v-%v]", cb.MinTicks, len(cb.Bonuses))
		}
		fmt.Fprintf(w, `<fieldset class="col">%s %v</fieldset>`, cb.BonusList, mp)

		fmt.Fprint(w, `</fieldset>`)
	}
	fmt.Fprint(w, `</div>`)
}

func showSingleCombo(w http.ResponseWriter, r ComboBonus) {

	pm := selectOptionArray([]int{ComboScoreMethodPoints, ComboScoreMethodMults}, []string{"points", "multipliers"}, r.ScoreMethod)
	page := fmt.Sprintf(tmpltSingleCombo, r.Comboid, r.BriefDesc, pm)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, page)
}
