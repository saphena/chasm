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

var BonusDisplayScreen = `
	<article class="bonus">
		<fieldset>
			<label for="BonusID">Code</label>
			<input type="text" id="BonusID" name="BonusID" class="BonusID" value="{{.B.BonusID}}">
		</fieldset>
		<fieldset>
			<label for="BriefDesc">Description</label>
			<input type="text" id="BriefDesc" name="BriefDesc" class="BriefDesc" value="{{.B.BriefDesc}}">
		</fieldset>
		<fieldset>
			<label for="Points">Points</label>
			<input type="number" id="Points" name="Points" class="Points" value="{{.B.Points}}">
			<select id="AskPoints" name="AskPoints">
				<option value="0" {{if eq .B.AskPoints 0}}selected{{end}}>Fixed</option>
				<option value="1" {{if eq .B.AskPoints 1}}selected{{end}}>Variable</option>
				<option value="2" {{if eq .B.AskPoints 2}}selected{{end}}>Multiply last</option>
			</select>
		</fieldset>
		<fieldset>
			<label for="Notes">Scoring notes</label>
			<input type="text" id="Notes" name="Notes" class="Notes" value="{{.B.Notes}}">
		</fieldset>
		<fieldset>
			<label>Scoring flags</label>
			<span title="Alert!">
				<label for="ScoringFlagA" class="short"><img class="icon" src="/img?i=alert" alt="!"></label>
				<input type="checkbox" id="ScoringFlagA" name="ScoringFlagA" {{if .FlagA}}checked{{end}} value="A">
			</span>
			<span title="Bike in photo">
				<label for="ScoringFlagB" class="short"><img class="icon" src="/img?i=bike" alt="B"></label>
				<input type="checkbox" id="ScoringFlagB" name="ScoringFlagB" {{if .FlagB}}checked{{end}} value="B">
			</span>
			<span title="Daylight only">
				<label for="ScoringFlagD" class="short"><img class="icon" src="/img?i=daylight" alt="D"></label>
				<input type="checkbox" id="ScoringFlagD" name="ScoringFlagD" {{if .FlagD}}checked{{end}} value="D">
			</span>
			<span title="Face in photo">
				<label for="ScoringFlagF" class="short"><img class="icon" src="/img?i=face" alt="F"></label>
				<input type="checkbox" id="ScoringFlagF" name="ScoringFlagF" {{if .FlagF}}checked{{end}} value="F">
			</span>
			<span title="Nighttime only">
				<label for="ScoringFlagN" class="short"><img class="icon" src="/img?i=night" alt="N"></label>
				<input type="checkbox" id="ScoringFlagN" name="ScoringFlagN" {{if .FlagN}}checked{{end}} value="N">
			</span>
			<span title="Restricted access/hours">
				<label for="ScoringFlagR" class="short"><img class="icon" src="/img?i=restricted" alt="R"></label>
				<input type="checkbox" id="ScoringFlagR" name="ScoringFlagR" {{if .FlagR}}checked{{end}} value="R">
			</span>
			<span title="Receipt/ticket needed">
				<label for="ScoringFlagT" class="short"><img class="icon" src="/img?i=receipt" alt="T"></label>
				<input type="checkbox" id="ScoringFlagT" name="ScoringFlagT" {{if .FlagT}}checked{{end}} value="T">
			</span>
		</fieldset>
		<fieldset>
			<label for="Image">Image</label>
			<input type="text" id="Image" name="Image" class="Image" value="{{.B.Image}}">
		</fieldset><fieldset>
			<img alt="*" data-bimg-folder="{{.BonusImgFldr}}" src="{{.BonusImgFldr}}/{{.B.Image}}">
		</fieldset>
		<fieldset>
			<label for="Compulsory">Compulsory?</label> 
			<select id="Compulsory" >
				<option value="0" {{if eq .B.Compulsory 0}}selected{{end}}>Optional</option>
				<option value="1" {{if eq .B.Compulsory 1}}selected{{end}}>Compulsory</option>
			</select>
		</fieldset>
		<fieldset>
			<label for="RestMinutes">Rest minutes</label> 
			<input id="RestMinutes" name="RestMinutes" class="RestMinutes" value="{{.B.RestMinutes}}">
			<select name="AskMinutes">
				<option value="0" {{if eq .B.AskMinutes 0}}selected{{end}}>Fixed</option>
				<option value="1" {{if eq .B.AskMinutes 1}}selected{{end}}>Variable</option>
			</select>
		</fieldset>
		<fieldset>
			<label for="Coords">Coords</label> 
			<input id="Coords" name="Coords" class="Coords" value="{{.B.Coords}}">
		</fieldset>
		<fieldset>
			<label for="Waffle">Waffle</label> 
			<input id="Waffle" name="Waffle" class="Waffle" value="{{.B.Waffle}}">
		</fieldset>
	</article>
`

func list_bonuses(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Bonuses")

	fmt.Fprint(w, `<p class="intro">Ordinary bonuses generally represent physical locations that entrants must visit and complete some task, typically take a photo.  Descriptions may include limited HTML to affect formatting on score explanations.</p>`)

	fmt.Fprint(w, `<div class="intro">`)
	fmt.Fprint(w, `<button class="plus" onclick="addNewBonus()" title="Add new bonus">+</button>`)
	fmt.Fprint(w, ` <input type="text" onchange="showBonus(this.value)" placeholder="Code to show">`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="bonuslist hdr">`)
	fmt.Fprint(w, `<span>Code</span><span>Description</span><span>Points</span><span class="right">Claims</span>`)
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
		fmt.Fprint(w, `<span class="claims"></span>`)
		fmt.Fprint(w, `</div>`)
	}
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
	sqlx += " FROM bonuses WHERE BonusID='" + bonus + "'"

	if r.FormValue("back") != "" {
		startHTMLBL(w, "Bonus detail", r.FormValue("back"))
	} else {
		startHTML(w, "Bonus detail")
	}
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if !rows.Next() {
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

	t, err := template.New("BonusDetail").Parse(BonusDisplayScreen)
	checkerr(err)
	err = t.Execute(w, br)
	checkerr(err)

}
