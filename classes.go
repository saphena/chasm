package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type classrec struct {
	Class       int
	BriefDesc   string
	AutoAssign  int
	MinPoints   int
	MinBonuses  int
	BonusesReqd string
	LowestRank  int
}

func build_classlist() []classrec {

	res := make([]classrec, 0)
	sqlx := "SELECT Class,ifnull(BriefDesc,''),AutoAssign,MinPoints,MinBonuses,ifnull(BonusesReqd,''),LowestRank FROM classes ORDER BY Class"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var cr classrec
		err = rows.Scan(&cr.Class, &cr.BriefDesc, &cr.AutoAssign, &cr.MinPoints, &cr.MinBonuses, &cr.BonusesReqd, &cr.LowestRank)
		checkerr(err)
		res = append(res, cr)
	}
	return res
}

var classIntro = `
<article class="intro">
<p>Classes may be used to produce different certificates for different groups or 'classes' of entrant.</p>
<p> Class can be assigned manually, as in the RBLR1000 to distinguish route for example, or can be assigned automatically using entrant scores, bonuses visited and/or rank.</p>
<p>Class 0 is the default class for all entrants and may not have any filters applied. Other classes are examined in numeric order starting at 1 until the filter criteria are matched. If no matching class is found, 0 is applied.</p>
</article>
<p><br></p>
`
var classlisthdr = `
<article class="classes">
<button autofocus title="Add new claim" class="plus" onclick="addNewClass(this)"> + </button>
</article>
<article class="classes">
<div class="row hdr">
<span class="class">#</span><span>Class</span><span>Auto?</span><span class="minpoints">Points</span>
<span class="minbonuses">Bonuses</span><span class="rank">Rank</span>
</div>
</article>
<hr>
</header>
`

func show_classes(w http.ResponseWriter, r *http.Request) {

	classes := build_classlist()

	startHTML(w, "Classes")

	fmt.Fprint(w, classIntro)
	fmt.Fprint(w, classlisthdr)

	fmt.Fprint(w, `<article class="classes">`)
	for _, c := range classes {
		fmt.Fprint(w, `<div class="row" `)
		if c.Class != 0 {
			fmt.Fprintf(w, `onclick="window.location.href='/class/%v?back=/classes'"`, c.Class)
		}
		fmt.Fprint(w, `>`)
		fmt.Fprintf(w, `<span class="class">%v</span>`, c.Class)
		fmt.Fprintf(w, `<span>%v</span>`, c.BriefDesc)
		ma := "manual"
		if c.AutoAssign != 0 {
			ma = "automatic"
		}
		fmt.Fprintf(w, `<span>%v</span>`, ma)
		fmt.Fprintf(w, `<span>%v</span>`, c.MinPoints)
		fmt.Fprintf(w, `<span>%v</span>`, c.MinBonuses)
		fmt.Fprintf(w, `<span>%v</span>`, c.LowestRank)
		fmt.Fprint(w, `</div>`)
	}
	fmt.Fprint(w, `</article>`)

}

var classTopline = `
	<div class="topline">
		<fieldset>
			<button title="Delete this Class?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" style="display:none;" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
		<fieldset>
			<button id="updatedb" class="hideuntil" title="Delete Bonus" disabled onclick="deleteClass(this)"></button>
		</fieldset>

	</div>
`

var tmpltShowClass = `
<article class="showclass">
<fieldset>
	<label for="BriefDesc">Class</label>
	<input type="hidden" id="Class" name="Class" readonly value="{{.Class}}">
	<span>{{.Class}} </span>
	<input type="text" id="BriefDesc" name="BriefDesc" data-save="saveClass" oninput="oi(this)" onchange="saveClass(this)" value="{{.BriefDesc}}">
</fieldset>
<fieldset>
	<label for="AutoAssign">How assigned</label>
	<select id="AutoAssign" name="AutoAssign" onchange="saveClass(this)">
		<option value="0" {{if eq .AutoAssign 0}}selected{{end}}>Manual</option>
		<option value="1" {{if eq .AutoAssign 1}}selected{{end}}>Automatic
	</select>
</fieldset>
<fieldset id="aafields" {{if eq .AutoAssign 0}}class="hide"{{end}}>
	<fieldset>
		<label for="MinPoints">Minimum points</label>
		<input type="number" id="MinPoints" name="MinPoints" class="small" data-save="saveClass" oninput="oi(this)" onchange="saveClass(this)" value="{{.MinPoints}}">
	</fieldset>
	<fieldset>
		<label for="MinBonuses">Minimum bonuses</label>
		<input type="number" id="MinBonuses" name="MinBonuses" class="small" data-save="saveClass" oninput="oi(this)" onchange="saveClass(this)" value="{{.MinBonuses}}">
	</fieldset>
	<fieldset>
		<label for="BonusesReqd">Required bonuses</label>
		<input type="text" id="BonusesReqd" name="BonusesReqd" data-save="saveClass" oninput="oi(this)" onchange="saveClass(this)" value="{{.BonusesReqd}}">
	</fieldset>
	<fieldset>
		<label for="LowestRank">Lowest rank</label>
		<input type="number" id="LowestRank" name="LowestRank" class="small" data-save="saveClass" oninput="oi(this)" onchange="saveClass(this)" value="{{.LowestRank}}">
	</fieldset>
</fieldset>
</article>
`

func createNewClass() int {

	sqlx := "SELECT max(Class) FROM classes"
	cls := getIntegerFromDB(sqlx, 0) + 1
	sqlx = fmt.Sprintf("INSERT INTO classes (Class,BriefDesc) VALUES(%v,'%v')", cls, cls)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	return cls
}

func deleteClass(w http.ResponseWriter, r *http.Request) {

	clsid := r.PathValue("class")
	if clsid == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"incomplete request"}`)
		return
	}
	sqlx := "DELETE FROM classes WHERE Class=" + clsid
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}
func saveClass(w http.ResponseWriter, r *http.Request) {

	clsid := r.FormValue("c")
	fld := r.FormValue("ff")
	if clsid == "" || fld == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"incomplete request"}`)
		return
	}
	sqlx := "UPDATE classes SET " + fld + "=? WHERE Class=?"

	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(r.FormValue(fld), clsid)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}
func showClass(w http.ResponseWriter, r *http.Request) {

	class := intval(r.PathValue("class"))
	if class == 0 {
		class = createNewClass()
	}
	classes := build_classlist()

	startHTMLBL(w, "Class", r.FormValue("back"))

	if class > 0 {
		fmt.Fprint(w, classTopline)
	}
	fmt.Fprint(w, classIntro)

	t, err := template.New("class").Parse(tmpltShowClass)
	checkerr(err)
	for _, c := range classes {
		if c.Class == class {
			err = t.Execute(w, c)
			checkerr(err)
			break
		}
	}
}
