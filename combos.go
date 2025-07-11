package main

import (
	"fmt"
	"net/http"
	"text/template"
)

const ComboScoreMethodMults = 1
const ComboScoreMethodPoints = 0

var tmpltSingleCombo = `
<article class="combo">
	<form>
		<fieldset class="field">
			<label for="ComboID">ComboID</label>
			<input id="ComboID" class="ComboID" name="Combo code" value="{{.Comboid}}">
		</fieldset>
		<fieldset class="field">
			<label for="BriefDesc">Description</label>
			<input id="BriefDesc" class="BriefDesc" name="BriefDesc" value="{{.BriefDesc}}">
		</fieldset>
		<fieldset class="field">
			<label for="ScoreMethod">Value is</label>
			<select id="ScoreMethod" name="ScoreMethod">
				<option value="0" {{if eq .ScoreMethod 0}}selected{{end}}>points</option>
				<option value="1" {{if ne .ScoreMethod 0}}selected{{end}}>multipliers</option>
			</select>
		</fieldset>
		<fieldset class="field">
			<label for="MinimumTicks">Minimum bonuses to score</label>
			<input id="MinimumTicks" type="number" class="MinimumTicks" name="MinimumTicks" value="{{.MinTicks}}">
		</fieldset>
		<fieldset>
			<label for="BonusList">Underlying bonuses</label>
			<input type="text" id="BonusList" class="BonusList" name="BonusList" value="{{.BonusList}}">
		</fieldset>
		<fieldset>
			<label for="Compulsory">Compulsory?</label>
			<select id="Compulsory" name="Compulsory">
				<option value="0" {{if not .Compulsory}}selected{{end}}>optional</option>
				<option value="1" {{if .Compulsory}}selected{{end}}>COMPULSORY</option>
			</select>
		</fieldset>
	</form>
</article>`

func show_combos(w http.ResponseWriter, r *http.Request) {

	const Combox = `
	Combos are scored automatically when their underlying ordinary or combo bonuses are scored. 
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<style>%s</style>`, maincss)
	fmt.Fprintf(w, `<script>%s</script>`, mainscript)

	ComboBonuses = loadCombos("")

	startHTML(w, "Combos")

	fmt.Fprint(w, `</header>`)

	fmt.Fprintf(w, `<p class="intro">%v</p>`, Combox)
	fmt.Fprint(w, `<div class="combos">`)
	for _, cb := range ComboBonuses {
		var mp string
		fmt.Fprintf(w, `<fieldset class="row combo" onclick="window.location.href='/combo?c=%v&back=/combos'">`, cb.Comboid)
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, cb.Comboid)
		fmt.Fprintf(w, `<fieldset class="col BriefDesc">%s</fieldset>`, cb.BriefDesc)
		mp = ""
		if cb.Compulsory {
			mp = "<strong>!</strong>"
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

func showSingleCombo(w http.ResponseWriter, c ComboBonus, bl string) {

	const Combox = `
	Combos can be set to score different values depending on the number of underlying bonuses scored. 
	By default all underlying bonuses must be scored. 
	Descriptions may include limited HTML to affect formatting on score explanations.
	`

	t, err := template.New("combo").Parse(tmpltSingleCombo)
	checkerr(err)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	startHTMLBL(w, "Combo detail", bl)
	fmt.Fprint(w, `</header>`)

	fmt.Fprintf(w, `<p class="intro">%v</p>`, Combox)
	err = t.Execute(w, c)
	checkerr(err)
}
