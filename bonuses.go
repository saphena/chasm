package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type BonusRec struct {
	BonusID     string
	BriefDesc   string
	Points      int
	Compulsory  int
	Notes       string
	Flags       string
	AskPoints   int
	RestMinutes int
	AskMinutes  int
	Image       string
	Coords      string
	Waffle      string
	Question    string
	Answer      string
	Leg         int
	Cat         [NumCategoryAxes]int
}

type BonusDisplayRec struct {
	B            BonusRec
	FlagA        bool
	FlagB        bool
	FlagD        bool
	FlagF        bool
	FlagN        bool
	FlagR        bool
	FlagT        bool
	BonusImgFldr string
}

const TrashcanIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
  <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
  <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
</svg>`

const FloppyDiskIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-floppy" viewBox="0 0 16 16">
  <path d="M11 2H9v3h2z"/>
  <path d="M1.5 0h11.586a1.5 1.5 0 0 1 1.06.44l1.415 1.414A1.5 1.5 0 0 1 16 2.914V14.5a1.5 1.5 0 0 1-1.5 1.5h-13A1.5 1.5 0 0 1 0 14.5v-13A1.5 1.5 0 0 1 1.5 0M1 1.5v13a.5.5 0 0 0 .5.5H2v-4.5A1.5 1.5 0 0 1 3.5 9h9a1.5 1.5 0 0 1 1.5 1.5V15h.5a.5.5 0 0 0 .5-.5V2.914a.5.5 0 0 0-.146-.353l-1.415-1.415A.5.5 0 0 0 13.086 1H13v4.5A1.5 1.5 0 0 1 11.5 7h-7A1.5 1.5 0 0 1 3 5.5V1H1.5a.5.5 0 0 0-.5.5m3 4a.5.5 0 0 0 .5.5h7a.5.5 0 0 0 .5-.5V1H4zM3 15h10v-4.5a.5.5 0 0 0-.5-.5h-9a.5.5 0 0 0-.5.5z"/>
</svg>`

var BonusDisplayScreen = `
	<div class="topline">
		{{if ne .B.BonusID ""}}
		<fieldset>
			<button title="Delete this Bonus?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
		{{end}}
		<fieldset>
			<button id="updatedb" title="Update DB" disabled>` + FloppyDiskIcon + `</button>
		</fieldset>
		<fieldset>
			<input type="button" title="show next" value="⋘" onclick="window.location.href='/bonus?b={{.B.BonusID}}&prev'">
			<input type="button" title="show previous" value="⋙" onclick="window.location.href='/bonus?b={{.B.BonusID}}&next'">
			<input type="button" title="back to list" value="↥☰↥" onclick="window.location.href='/bonuses'">
		</fieldset>

	</div>
	<article class="bonus">
		<fieldset>
			<label for="BonusID">Code</label>
			<input type="text" id="BonusID" name="BonusID" {{if ne "" .B.BonusID}}readonly{{else}}autofocus{{end}} class="BonusID" value="{{.B.BonusID}}" onchange="addBonus(this)">
		</fieldset>
		<fieldset>
			<label for="BriefDesc">Description</label>
			<input type="text" id="BriefDesc" data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" name="BriefDesc" class="BriefDesc" value="{{.B.BriefDesc}}">
		</fieldset>
		<fieldset>
			<label for="Points">Points</label>
			<input type="number" id="Points"  data-b="{{.B.BonusID}}" data-save="saveBonus"  oninput="oi(this)" onchange="saveBonus(this)" name="Points" class="Points" value="{{.B.Points}}">
			<select id="AskPoints" name="AskPoints"  data-b="{{.B.BonusID}}" onchange="saveBonus(this)" >
				<option value="0" {{if eq .B.AskPoints 0}}selected{{end}}>Fixed</option>
				<option value="1" {{if eq .B.AskPoints 1}}selected{{end}}>Variable</option>
				<option value="2" {{if eq .B.AskPoints 2}}selected{{end}}>Multiply last</option>
			</select>
		</fieldset>
		<fieldset>
			<label for="Notes">Scoring notes</label>
			<textarea id="Notes" name="Notes"  data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" class="Notes">{{.B.Notes}}</textarea>
		</fieldset>
		<fieldset>
			<span class="button {{if.FlagA}}selected{{end}}" title="Alert!">
				<img class="icon" src="/img?i=alert" alt="!" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagA"  data-b="{{.B.BonusID}}"  onclick="saveBonus(this)" name="ScoringFlag" {{if .FlagB}}checked{{end}} value="A">
			</span>
			<span class="button {{if .FlagB}}selected{{end}}" title="Bike in photo">
				<img class="icon" src="/img?i=bike" alt="B" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagB"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagB}}checked{{end}} value="B">
			</span>
			<span class="button {{if .FlagD}}selected{{end}}" title="Daylight only">
				<img class="icon" src="/img?i=daylight" alt="D" onclick="toggleButton(this)">
				<input type="checkbox"  class="quietcheck" id="ScoringFlagD"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagD}}checked{{end}} value="D">
			</span>
			<span class="button {{if .FlagF}}selected{{end}}" title="Face in photo">
				<img class="icon" src="/img?i=face" alt="F" onclick="toggleButton(this)">
				<input type="checkbox"  class="quietcheck" id="ScoringFlagF"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagF}}checked{{end}} value="F">
			</span>
			<span class="button {{if .FlagN}}selected{{end}}" title="Nighttime only">
				<img class="icon" src="/img?i=night" alt="N" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagN"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagN}}checked{{end}} value="N">
			</span>
			<span class="button {{if .FlagR}}selected{{end}}" title="Restricted access/hours">
				<img class="icon" src="/img?i=restricted" alt="R" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagR"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagR}}checked{{end}} value="R">
			</span>
			<span class="button {{if .FlagT}}selected{{end}}" title="Receipt/ticket needed">
				<img class="icon" src="/img?i=receipt" alt="T" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagT"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagT}}checked{{end}} value="T">
			</span>
		</fieldset>
		<fieldset>
			<label for="Image">Image</label>
			<input type="text" id="Image"  data-b="{{.B.BonusID}}" data-save="saveBonus"oninput="oi(this)" onchange="saveBonus(this)" name="Image" class="Image" value="{{.B.Image}}">
		</fieldset><fieldset>
			<img alt="*" data-bimg-folder="{{.BonusImgFldr}}"class="thumbnail" src="{{.BonusImgFldr}}/{{.B.Image}}" onclick="this.classList.toggle('thumbnail')">
		</fieldset>
		<fieldset>
			<label for="Compulsory">Compulsory?</label> 
			<select id="Compulsory"  data-b="{{.B.BonusID}}"  name="Compulsory" onchange="saveBonus(this)">
				<option value="0" {{if eq .B.Compulsory 0}}selected{{end}}>Optional</option>
				<option value="1" {{if eq .B.Compulsory 1}}selected{{end}}>Compulsory</option>
			</select>
		</fieldset>
		<fieldset>
			<label for="RestMinutes">Rest minutes</label> 
			<input id="RestMinutes" type="number" name="RestMinutes"  data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" class="RestMinutes" value="{{.B.RestMinutes}}">
			<select id="AskMinutes"  data-b="{{.B.BonusID}}" name="AskMinutes"  onchange="saveBonus(this)" >
				<option value="0" {{if eq .B.AskMinutes 0}}selected{{end}}>Fixed</option>
				<option value="1" {{if eq .B.AskMinutes 1}}selected{{end}}>Variable</option>
			</select>
		</fieldset>
		<fieldset>
			<label for="Coords">Coords</label> 
			<input id="Coords" name="Coords"  data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" class="Coords" value="{{.B.Coords}}">
		</fieldset>
		<fieldset>
			<label for="Waffle">Waffle</label> 
			<textarea id="Waffle" name="Waffle"  data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" class="Waffle">{{.B.Waffle}}</textarea>
		</fieldset>
	</article>
`

func createBonus(w http.ResponseWriter, r *http.Request) {

	bonus := strings.ToUpper(r.FormValue("b"))
	if bonus == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"Blank BonusID"}`)
		return
	}

	sqlx := "INSERT INTO bonuses (BonusID,BriefDesc) VALUES(?,?)"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	res, err := stmt.Exec(bonus, bonus)
	checkerr(err)
	ra, err := res.RowsAffected()
	checkerr(err)
	if ra != 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"Duplicate BonusID"}`)
	} else {
		fmt.Fprint(w, `{"ok":true,"msg":"`+bonus+`"}`)
	}
}
func list_bonuses(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Bonuses")

	fmt.Fprint(w, `<p class="intro">Ordinary bonuses generally represent physical locations that entrants must visit and complete some task, typically take a photo.  Descriptions may include limited HTML to affect formatting on score explanations.</p>`)

	fmt.Fprint(w, `<div class="intro bonuslist">`)
	fmt.Fprint(w, `<button class="plus" autofocus title="Add new bonus" onclick="window.location.href='/bonus?back=bonuses'">+</button>`)
	fmt.Fprint(w, ` <input type="text" onchange="showBonus(this.value)" onblur="showBonus(this.value)"  placeholder="Code to show">`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="bonuslist hdr">`)
	fmt.Fprint(w, `<span>Code</span><span>Description</span><span>Points</span><span>Claims</span>`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `</header>`)

	sqlx := "SELECT BonusID,BriefDesc,Points FROM bonuses ORDER BY BonusID"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	oe := true
	for rows.Next() {
		var bonus string
		var descr string
		var points int
		err := rows.Scan(&bonus, &descr, &points)
		checkerr(err)
		oex := "even"
		if oe {
			oex = "odd"
		}
		oe = !oe
		fmt.Fprintf(w, `<div class="bonuslist row link %v" onclick="window.location.href='/bonus?b=%v&back=bonuses'">`, oex, bonus)
		fmt.Fprintf(w, `<span>%v</span><span>%v</span><span>%v</span>`, bonus, descr, points)
		n := getIntegerFromDB("SELECT count(DISTINCT EntrantID) FROM claims WHERE BonusID='"+bonus+"'", 0)
		fmt.Fprint(w, `<span class="claims">`)
		if n > 0 {
			fmt.Fprintf(w, `<a href="/claims?bsel=%v"> %v </a>`, bonus, n)
		}
		fmt.Fprint(w, `</span>`)
		fmt.Fprint(w, `</div>`)
	}
}

func saveBonus(w http.ResponseWriter, r *http.Request) {

	bonus := r.FormValue("b")
	if bonus == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no BonusID supplied"}`)
		return
	}
	fld := r.FormValue("ff")
	if fld == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no fieldname supplied"}`)
		return
	}
	val := r.FormValue(fld)
	sqlx := "UPDATE bonuses SET " + fld + "=? WHERE BonusID=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(val, bonus)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}
func show_bonus(w http.ResponseWriter, r *http.Request) {

	bonus := strings.ToUpper(r.FormValue("b"))
	sqlx := "SELECT Bonusid, BriefDesc, Points"
	sqlx += ",Compulsory, ifnull(Notes,''), ifnull(Flags,''), AskPoints"
	sqlx += ",RestMinutes, AskMinutes, ifnull(Image,''), ifnull(Coords,'')"
	sqlx += ",ifnull(Waffle,''),ifnull(Question,''),ifnull(Answer,''),Leg"
	for i := 1; i <= NumCategoryAxes; i++ {
		sqlx += ", Cat" + strconv.Itoa(i)
	}
	sqlx += " FROM bonuses WHERE BonusID"

	rel := "="
	ord := "BonusID"
	ok := r.Form.Has("next")
	if ok {
		rel = ">"
	} else {
		ok = r.Form.Has("prev")
		if ok {
			ord += " DESC"
			rel = "<"
		}
	}

	sqlx += rel + "'" + bonus + "'"
	sqlx += " ORDER BY " + ord

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if !rows.Next() {
		if rel != "=" {
			list_bonuses(w, r)
			return
		}
		fmt.Fprintf(w, `<p class="error">No such bonus '%v'</p>`, bonus)
		return
	}
	var b BonusRec

	s := reflect.ValueOf(&b).Elem()
	numCols := s.NumField() + NumCategoryAxes - 1
	columns := make([]interface{}, numCols)
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)

		if field.Kind() == reflect.Array {
			for j := 0; j < field.Len(); j++ {
				columns[i+j] = field.Index(j).Addr().Interface()
			}
		} else {
			columns[i] = field.Addr().Interface()
		}
	}

	err = rows.Scan(columns...)
	checkerr(err)

	var br BonusDisplayRec
	br.B = b
	br.BonusImgFldr = CS.ImgBonusFolder
	br.FlagA = strings.Contains(b.Flags, "A")
	br.FlagB = strings.Contains(b.Flags, "B")
	br.FlagD = strings.Contains(b.Flags, "D")
	br.FlagF = strings.Contains(b.Flags, "F")
	br.FlagN = strings.Contains(b.Flags, "N")
	br.FlagR = strings.Contains(b.Flags, "R")
	br.FlagT = strings.Contains(b.Flags, "T")

	if r.FormValue("back") != "" {
		startHTMLBL(w, "Bonus detail", r.FormValue("back"))
	} else {
		startHTML(w, "Bonus detail")
	}
	fmt.Fprint(w, `</header>`)

	t, err := template.New("BonusDetail").Parse(BonusDisplayScreen)
	checkerr(err)
	err = t.Execute(w, br)
	checkerr(err)

}
