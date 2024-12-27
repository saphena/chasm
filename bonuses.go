package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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

func list_bonuses(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Bonuses")

	fmt.Fprint(w, `<p class="intro">Ordinary bonuses generally represent physical locations that entrants must visit and complete some task, typically take a photo. They are presented on scorecards in code order. Numeric only codes should all have the same number of digits (use leading '0' if necessary). Descriptions may include limited HTML to affect formatting on score explanations.</p>`)

	fmt.Fprint(w, `<div class="intro">`)
	fmt.Fprint(w, `<button class="plus" onclick="addNewBonus()" title="Add new bonus">+</button>`)
	fmt.Fprint(w, ` <input type="text" onchange="showBonus(this.value)" placeholder="Code to show">`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="bonuslist hdr">`)
	fmt.Fprint(w, `<span>Code</span><span>Description</span><span>Points</span><span></span><span class="right">Claims</span>`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `</header>`)

	sqlx := "SELECT BonusID,BriefDesc,Points FROM bonuses ORDER BY BonusID"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var bonus string
		var descr string
		var points int
		err := rows.Scan(&bonus, &descr, &points)
		checkerr(err)
		fmt.Fprint(w, `<div class="bonuslist">`)
		fmt.Fprintf(w, `<span>%v</span><span>%v</span><span>%v</span>`, bonus, descr, points)
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

	startHTML(w, "Bonus detail")
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

	var selected string
	var checked string

	fmt.Fprint(w, `</header>`)
	fmt.Fprint(w, `<article class="bonus">`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="BonusID">Code</label>`)
	fmt.Fprintf(w, `<input type="text" id="BonusID" name="BonusID" class="BonusID" value="%v">`, b.BonusID)
	fmt.Fprint(w, `</fieldset>`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="BriefDesc">Description</label>`)
	fmt.Fprintf(w, `<input type="text" id="BriefDesc" name="BriefDesc" class="BriefDesc" value="%v">`, b.BriefDesc)
	fmt.Fprint(w, `</fieldset>`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="Points">Points</label>`)
	fmt.Fprintf(w, `<input type="number" id="Points" name="Points" class="Points" value="%v">`, b.Points)
	fmt.Fprint(w, ` <select id="AskPoints" name="AskPoints">`)
	selected = ""
	if b.AskPoints == 0 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="0" %v>Fixed</option>`, selected)
	selected = ""
	if b.AskPoints == 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="1" %v>Variable</option>`, selected)
	selected = ""
	if b.AskPoints == 2 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="2" %v>Multiplier</option>`, selected)

	fmt.Fprint(w, `</select>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="Notes">Scoring notes</label>`)
	fmt.Fprintf(w, `<input type="text" id="Notes" name="Notes" class="Notes" value="%v">`, b.Notes)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label>Scoring flags</label>`)
	fmt.Fprint(w, `<span title="Alert!">`)
	fmt.Fprint(w, `<label for="ScoringFlagA" class="short"><img class="icon" src="/img?i=alert" alt="!"></label>`)
	checked = ""
	if strings.Contains(b.Flags, "A") {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" id="ScoringFlagA" name="ScoringFlagA" %v value="A">`, checked)
	fmt.Fprint(w, `</span>`)

	fmt.Fprint(w, `<span title="Bike in photo">`)
	fmt.Fprint(w, `<label for="ScoringFlagB" class="short"><img class="icon" src="/img?i=bike" alt="B"></label>`)
	checked = ""
	if strings.Contains(b.Flags, "B") {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" id="ScoringFlagB" name="ScoringFlagB" %v value="B">`, checked)
	fmt.Fprint(w, `</span>`)

	fmt.Fprint(w, `<span title="Daylight only">`)
	fmt.Fprint(w, `<label for="ScoringFlagD" class="short"><img class="icon" src="/img?i=daylight" alt="D"></label>`)
	checked = ""
	if strings.Contains(b.Flags, "D") {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" id="ScoringFlagD" name="ScoringFlagD" %v value="D">`, checked)
	fmt.Fprint(w, `</span>`)

	fmt.Fprint(w, `<span title="Face in photo">`)
	fmt.Fprint(w, `<label for="ScoringFlagF" class="short"><img class="icon" src="/img?i=face" alt="F"></label>`)
	checked = ""
	if strings.Contains(b.Flags, "F") {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" id="ScoringFlagF" name="ScoringFlagF" %v value="F">`, checked)
	fmt.Fprint(w, `</span>`)

	fmt.Fprint(w, `<span title="Nighttime only">`)
	fmt.Fprint(w, `<label for="ScoringFlagN" class="short"><img class="icon" src="/img?i=night" alt="N"></label>`)
	checked = ""
	if strings.Contains(b.Flags, "N") {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" id="ScoringFlagN" name="ScoringFlagN" %v value="N">`, checked)
	fmt.Fprint(w, `</span>`)

	fmt.Fprint(w, `<span title="Restricted access/hours">`)
	fmt.Fprint(w, `<label for="ScoringFlagR" class="short"><img class="icon" src="/img?i=restricted" alt="R"></label>`)
	checked = ""
	if strings.Contains(b.Flags, "R") {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" id="ScoringFlagR" name="ScoringFlagR" %v value="R">`, checked)
	fmt.Fprint(w, `</span>`)

	fmt.Fprint(w, `<span title="Receipt/ticket needed">`)
	fmt.Fprint(w, `<label for="ScoringFlagT" class="short"><img class="icon" src="/img?i=receipt" alt="T"></label>`)
	checked = ""
	if strings.Contains(b.Flags, "T") {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" id="ScoringFlagT" name="ScoringFlagT" %v value="T">`, checked)
	fmt.Fprint(w, `</span>`)

	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="Image">Image</label>`)
	fmt.Fprintf(w, ` <input type="text" id="Image" name="Image" class="Image" value="%v">`, b.Image)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `</article>`)
}
