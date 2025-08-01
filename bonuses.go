// TODO
//
// Trap case of scoring variables changed during rally, may need to propagate changes.
//
// Rest/matching group maintenance.
//
// Image selector / image upload.
package main

import (
	"fmt"
	"net/http"
	"os"
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

var BonusDisplayScreen = `
	<div class="topline">
		{{if ne .B.BonusID ""}}
		<fieldset>
			<button title="Delete this Bonus?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" style="display:none;" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
		{{end}}
		<fieldset>
			<button id="updatedb" class="hideuntil" title="Delete Bonus" disabled onclick="deleteBonus(this)"></button>
		</fieldset>
		<fieldset>
			<button title="show next" onclick="window.location.href='/bonus?b={{.B.BonusID}}&prev'">⋘</button>
			<button title="show previous" onclick="window.location.href='/bonus?b={{.B.BonusID}}&next'">⋙</button>
			<button title="back to list" onclick="window.location.href='/bonuses'">↥☰↥</button>
		</fieldset>

	</div>
	<article class="bonus">
		<fieldset>
			<label for="BonusID">Code</label>
			<input type="text" id="BonusID" name="BonusID" {{if ne "" .B.BonusID}}readonly{{else}}autofocus{{end}} class="BonusID" value="{{.B.BonusID}}" onchange="addBonus(this)">
		</fieldset>
		<fieldset>
			<label for="BriefDesc">Description</label>
			<input type="text" id="BriefDesc" data-b="{{.B.BonusID}}" {{if ne .B.BonusID ""}}autofocus{{end}} data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" name="BriefDesc" class="BriefDesc" value="{{.B.BriefDesc}}">
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
			<span class="flag {{if.FlagA}}selected{{end}}" title="Alert!">
				<img class="icon" src="/img?i=alert" alt="!" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagA"  data-b="{{.B.BonusID}}"  onclick="saveBonus(this)" name="ScoringFlag" {{if .FlagB}}checked{{end}} value="A">
			</span>
			<span class="flag {{if .FlagB}}selected{{end}}" title="Bike in photo">
				<img class="icon" src="/img?i=bike" alt="B" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagB"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagB}}checked{{end}} value="B">
			</span>
			<span class="flag {{if .FlagD}}selected{{end}}" title="Daylight only">
				<img class="icon" src="/img?i=daylight" alt="D" onclick="toggleButton(this)">
				<input type="checkbox"  class="quietcheck" id="ScoringFlagD"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagD}}checked{{end}} value="D">
			</span>
			<span class="flag {{if .FlagF}}selected{{end}}" title="Face in photo">
				<img class="icon" src="/img?i=face" alt="F" onclick="toggleButton(this)">
				<input type="checkbox"  class="quietcheck" id="ScoringFlagF"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagF}}checked{{end}} value="F">
			</span>
			<span class="flag {{if .FlagN}}selected{{end}}" title="Nighttime only">
				<img class="icon" src="/img?i=night" alt="N" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagN"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagN}}checked{{end}} value="N">
			</span>
			<span class="flag {{if .FlagR}}selected{{end}}" title="Restricted access/hours">
				<img class="icon" src="/img?i=restricted" alt="R" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagR"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagR}}checked{{end}} value="R">
			</span>
			<span class="flag {{if .FlagT}}selected{{end}}" title="Receipt/ticket needed">
				<img class="icon" src="/img?i=receipt" alt="T" onclick="toggleButton(this)">
				<input type="checkbox" class="quietcheck" id="ScoringFlagT"  data-b="{{.B.BonusID}}"  onchange="saveBonus(this)" name="ScoringFlag" {{if .FlagT}}checked{{end}} value="T">
			</span>
		</fieldset>
		<fieldset>
			<label for="Image">Image</label>
			<!--
			<input type="text" id="Image"  data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" name="Image" class="Image" value="{{.B.Image}}">
			-->
			<select id="Image" data-b="{{.B.BonusID}}" name="Image" class="Image" onchange="saveBonus(this)">
			%v
			</select>
		</fieldset><fieldset>
			<img alt="*" id="imgImage" data-bimg-folder="{{.BonusImgFldr}}"class="thumbnail toggle" src="{{.BonusImgFldr}}/{{.B.Image}}" onclick="this.classList.toggle('thumbnail')">
		</fieldset>
		<fieldset>
			<label for="Compulsory">Compulsory?</label> 
			<select id="Compulsory"  data-b="{{.B.BonusID}}"  name="Compulsory" onchange="saveBonus(this)">
				<option value="0" {{if eq .B.Compulsory 0}}selected{{end}}>Optional</option>
				<option value="1" {{if eq .B.Compulsory 1}}selected{{end}}>Compulsory</option>
			</select>
		</fieldset>
		</article>
`

var BonusQAFields = `
<article class="bonus">
<fieldset>
	<label for="Question">Question</label>
	<input type="text" id="Question" name="Question" data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" value="{{.B.Question}}">
</fieldset>
<fieldset>
	<label for="Answer">Answer</label>
	<input type="text" id="Answer" name="Answer" data-b="{{.B.BonusID}}" data-save="saveBonus" oninput="oi(this)" onchange="saveBonus(this)" value="{{.B.Answer}}">
</fieldset>
</article>
`

type bonuscat struct {
	Set     int
	SetName string
	BonusID string
	CatNow  int
	Cats    []CatDefinition
}

var BonusCatSelector = `
<article class="bonus">
	<fieldset>
		<label for="{{.Set}}cat">{{.SetName}}</label>
		<select id="{{.Set}}cat" name="Cat{{.Set}}" data-b="{{.BonusID}}" onchange="saveBonus(this)">
		<option value="0" {{if eq .CatNow 0}}selected{{end}}>{no selection}</option>
		{{$cat := .CatNow}}
		{{range $el := .Cats}}
			<option value="{{$el.Cat}}" {{if eq $el.Cat $cat}}selected{{end}}>{{$el.CatName}}</option>
		{{end}}
		</select>
	</fieldset>
</article>
`

var OtherBonusData = `
		<hr>
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

func buildBonusImageList() []string {

	file, err := os.Open(CS.ImgBonusFolder)
	checkerr(err)
	defer file.Close()
	names, err := file.Readdirnames(0)
	checkerr(err)
	return names
}

func buildImageSelectOptions(imgname string) string {

	il := buildBonusImageList()
	opts := `<option value="">{no selection}</option>`
	currentfound := false
	for _, img := range il {
		opts += `<option value="` + img + `" `
		if img == imgname {
			opts += "selected"
			currentfound = true
		}
		opts += `>` + img + `</option>`
	}
	if !currentfound {
		opts += `<option selected value="` + imgname + `">` + imgname + `</option>`
	}

	return opts

}
func createBonus(w http.ResponseWriter, r *http.Request) {

	bonus := strings.ToUpper(r.PathValue("b"))
	if bonus == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"Blank BonusID"}`)
		return
	}

	sqlx := "INSERT INTO bonuses (BonusID,BriefDesc) VALUES(?,?)"
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
		fmt.Fprint(w, `{"ok":false,"msg":"Duplicate BonusID"}`)
	} else {
		fmt.Fprint(w, `{"ok":true,"msg":"`+bonus+`"}`)
	}
}

func deleteBonus(w http.ResponseWriter, r *http.Request) {

	bonus := strings.ToUpper(r.PathValue("b"))
	if bonus == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"Blank BonusID"}`)
		return
	}

	sqlx := "DELETE FROM bonuses WHERE BonusID=?"
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

func fetchBonusDetails(w http.ResponseWriter, r *http.Request) {

	b := r.FormValue("b")
	if b == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no b specified"}`)
		return
	}
	bd := fetchBonusVars(b)
	fmt.Fprintf(w, `{"ok":true,"msg":"ok","name":"%v"`, bd.BriefDesc)
	fmt.Fprintf(w, `,"flags":"%v","img":"%v"`, bd.Flags, bd.Image)
	fmt.Fprintf(w, `,"notes":"%v"`, bd.Notes)
	fmt.Fprintf(w, `,"askpoints":%v`, jsonBool(bd.AskPoints))
	pm := "p"
	if bd.PointsAreMults {
		pm = "m"
	}
	fmt.Fprintf(w, `,"pointsaremults":"%v"`, pm)
	fmt.Fprintf(w, `,"askmins":%v`, jsonBool(bd.AskMins))
	fmt.Fprintf(w, `,"points":%v`, bd.Points)
	fmt.Fprintf(w, `,"question":"%v"`, bd.Question)
	fmt.Fprintf(w, `,"answer":"%v"`, bd.Answer)
	fmt.Fprintf(w, `,"restmins":%v`, bd.RestMins)
	fmt.Fprint(w, `}`)
}

func list_bonuses(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Bonuses")

	sets := build_axisLabels()

	fmt.Fprint(w, `<p class="intro">Ordinary bonuses generally represent physical locations that entrants must visit and complete some task, typically take a photo.  Descriptions may include limited HTML to affect formatting on score explanations.</p>`)

	fmt.Fprint(w, `<div class="intro bonuslist">`)
	fmt.Fprint(w, `<button class="plus" autofocus title="Add new bonus" onclick="window.location.href='/bonus?back=bonuses'">+</button>`)
	fmt.Fprint(w, ` <input type="text" onchange="showBonus(this.value)" onblur="showBonus(this.value)"  placeholder="Code to show">`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="bonuslist hdr">`)
	fmt.Fprint(w, `<span>Code</span><span>Description</span><span>Points</span><span>Claims</span>`)
	fmt.Fprint(w, `<span class="cats">`)
	for i := range sets {
		if sets[i] == "" {
			continue
		}
		fmt.Fprintf(w, `<span>%v</span>`, sets[i])
	}
	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</div><hr>`)
	fmt.Fprint(w, `</header>`)

	sqlx := "SELECT BonusID,BriefDesc,Points"
	for i := 1; i <= NumCategoryAxes; i++ {
		sqlx += ", Cat" + strconv.Itoa(i)
	}

	sqlx += " FROM bonuses ORDER BY BonusID"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	oe := true
	for rows.Next() {
		var bonus string
		var descr string
		var points int
		var cats [NumCategoryAxes]int
		err := rows.Scan(&bonus, &descr, &points, &cats[0], &cats[1], &cats[2], &cats[3], &cats[4], &cats[5], &cats[6], &cats[7], &cats[8])
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

		fmt.Fprint(w, `<span class="cats">`)
		for i := range sets {
			if sets[i] == "" {
				continue
			}
			sqlx = fmt.Sprintf("SELECT BriefDesc FROM categories WHERE Axis=%v AND Cat=%v", i+1, cats[i])
			catx := getStringFromDB(sqlx, "")
			fmt.Fprintf(w, `<span>%v</span>`, catx)
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
	newrec := false
	if !rows.Next() {
		if rel != "=" {
			list_bonuses(w, r)
			return
		}
		if bonus == "" {
			newrec = true
		} else {
			list_bonuses(w, r)
			return
		}
	}
	var b BonusRec

	if !newrec {
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
	}

	var br BonusDisplayRec
	br.B = b
	bonus = br.B.BonusID
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

	imgopts := buildImageSelectOptions(br.B.Image)
	t, err := template.New("BonusDetail").Parse(fmt.Sprintf(BonusDisplayScreen, imgopts))
	checkerr(err)
	err = t.Execute(w, br)
	checkerr(err)
	if CS.RallyUseQA {
		t, err = template.New("BonusQA").Parse(BonusQAFields)
		checkerr(err)
		err = t.Execute(w, br)
		checkerr(err)
	}

	sets := build_axisLabels()
	for i := range sets {
		if sets[i] == "" {
			continue
		}
		var set bonuscat
		set.Set = i + 1
		set.SetName = sets[i]
		set.BonusID = bonus
		set.CatNow = br.B.Cat[i]
		set.Cats = fetchSetCats(set.Set, true)
		fmt.Printf("%v\n", set)
		t, err := template.New("BonusCat").Parse(BonusCatSelector)
		checkerr(err)
		err = t.Execute(w, set)
		checkerr(err)

	}

}
