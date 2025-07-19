package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

const ComboScoreMethodMults = 1
const ComboScoreMethodPoints = 0

var tmpltSingleCombo = `
<div class="topline">
	{{if ne .Comboid ""}}
		<fieldset>
			<button title="Delete this Combo?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" style="display:none;" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
	{{end}}
	<fieldset>
		<button id="updatedb" class="hideuntil" title="Delete this record" disabled onclick="updateComboDB(this)"></button>
	</fieldset>
</div>
<p class="intro">Combos can be set to score different values depending on the number of underlying bonuses scored. 
	By default all underlying bonuses must be scored. 
	Descriptions may include limited HTML to affect formatting on score explanations.
</p>
<article class="combo">
	<form>
		<fieldset class="field">
			<label for="ComboID">ComboID</label>
			<input id="ComboID" class="ComboID" {{if ne "" .Comboid}}readonly{{else}}autofocus{{end}} name="Combo code" value="{{.Comboid}}" onchange="addCombo(this)">
		</fieldset>
		<fieldset class="field">
			<label for="BriefDesc">Description</label>
			<input id="BriefDesc" class="BriefDesc" name="BriefDesc" data-c="{{.Comboid}}" value="{{.BriefDesc}}" data-save="saveCombo" oninput="oi(this)" onchange="saveCombo(this)" >
		</fieldset>
		<fieldset>
			<label for="BonusList">Underlying bonuses</label>
			<input type="text" id="BonusList" class="BonusList"  data-c="{{.Comboid}}"  data-save="saveCombo" oninput="oi(this)" onchange="saveCombo(this)" name="BonusList" value="{{.BonusList}}">
		</fieldset>

		<fieldset class="field">
			<label for="MinimumTicks">Minimum bonuses to score</label>
			<input id="MinimumTicks" type="number" class="MinimumTicks"  data-c="{{.Comboid}}" data-save="saveCombo" oninput="oi(this)" onchange="saveCombo(this)" name="MinimumTicks" value="{{.MinTicks}}">
		</fieldset>
		<fieldset>
			<!-- <label for="PointsList">Points or Multipliers</label> -->

			<select id="ScoreMethod" name="ScoreMethod" data-save="saveCombo" oninput="oi(this)" onchange="saveCombo(this)" data-c="{{.Comboid}}" >
				<option value="0" {{if eq .ScoreMethod 0}}selected{{end}}>Points</option>
				<option value="1" {{if ne .ScoreMethod 0}}selected{{end}}>Multipliers</option>
			</select>

			<input type="text" id="PointsList" class="hide"  data-save="saveCombo"  data-c="{{.Comboid}}" oninput="oi(this)" onchange="saveCombo(this)" name="ScorePoints" value="{{.PointsList}}">
			<fieldset id="PointsListArray">
				<fieldset id="PointsListArrayHdrs">
				</fieldset>
				<fieldset id="PointsListArrayVals">
					<input type="number" id="Points0" class="Points" name="PointsVal" value="{{.PointsList}}">
				</fieldset>
			</fieldset>
		</fieldset>

		<fieldset>
			<label for="Compulsory">Compulsory?</label>
			<select id="Compulsory" name="Compulsory" data-save="saveCombo"  data-c="{{.Comboid}}" oninput="oi(this)" onchange="saveCombo(this)">
				<option value="0" {{if not .Compulsory}}selected{{end}}>optional</option>
				<option value="1" {{if .Compulsory}}selected{{end}}>COMPULSORY</option>
			</select>
		</fieldset>
	</form>
</article>
<script>extractComboPointsArray()</script>
`

func createCombo(w http.ResponseWriter, r *http.Request) {
	bonus := r.FormValue("b")
	if bonus == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"Blank ComboID"}`)
		return
	}

	sqlx := "INSERT INTO combos (ComboID,BriefDesc) VALUES(?,?)"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	res, err := stmt.Exec(bonus, bonus)
	if err != nil {
		fmt.Fprint(w, `{"ok":false,"msg":"`+err.Error()+`"}`)
		return
	}
	//checkerr(err)
	ra, err := res.RowsAffected()
	checkerr(err)
	if ra != 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"Duplicate ComboID"}`)
	} else {
		fmt.Fprint(w, `{"ok":true,"msg":"`+bonus+`"}`)
	}
}

func deleteCombo(w http.ResponseWriter, r *http.Request) {

	bonus := strings.ToUpper(r.FormValue("c"))
	if bonus == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"Blank ComboID"}`)
		return
	}

	sqlx := "DELETE FROM combos WHERE ComboID=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(bonus)
	if err != nil {
		fmt.Fprint(w, `{"ok":false,"msg":"`+err.Error()+`"}`)
		return
	}
	fmt.Fprint(w, `{"ok":true,"msg":"`+bonus+`"}`)
}

func saveCombo(w http.ResponseWriter, r *http.Request) {

	bonus := r.FormValue("c")
	if bonus == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no ComboID supplied"}`)
		return
	}
	fld := r.FormValue("ff")
	if fld == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no fieldname supplied"}`)
		return
	}
	val := r.FormValue(fld)
	sqlx := "UPDATE combos SET " + fld + "=? WHERE ComboID=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(val, bonus)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}

func show_combos(w http.ResponseWriter, r *http.Request) {

	const Combox = `
	Combos are scored automatically when their underlying ordinary or combo bonuses are scored. 
	`

	startHTML(w, "Combos")

	ComboBonuses = loadCombos("")

	fmt.Fprintf(w, `<p class="intro">%v</p>`, Combox)

	fmt.Fprint(w, `<div class="intro bonuslist">`)
	fmt.Fprint(w, `<button class="plus" autofocus title="Add new combo" onclick="window.location.href='/combo?back=combos'">+</button>`)
	fmt.Fprint(w, ` <input type="text" onchange="showCombo(this.value)" onblur="showCombo(this.value)"  placeholder="Code to show">`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="bonuslist hdr">`)
	fmt.Fprint(w, `<span>Code</span><span>Description</span><span>Points</span><span>Claims</span>`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `</header>`)

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

	t, err := template.New("combo").Parse(tmpltSingleCombo)
	checkerr(err)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	startHTMLBL(w, "Combo detail", bl)
	fmt.Fprint(w, `</header>`)

	err = t.Execute(w, c)
	checkerr(err)
}
