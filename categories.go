package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

type CatDefinition struct {
	Set     int
	Cat     int
	CatName string
}

const ordered_list_icon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-list-ol" viewBox="0 0 16 16">
  <path fill-rule="evenodd" d="M5 11.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5m0-4a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5"/>
  <path d="M1.713 11.865v-.474H2c.217 0 .363-.137.363-.317 0-.185-.158-.31-.361-.31-.223 0-.367.152-.373.31h-.59c.016-.467.373-.787.986-.787.588-.002.954.291.957.703a.595.595 0 0 1-.492.594v.033a.615.615 0 0 1 .569.631c.003.533-.502.8-1.051.8-.656 0-1-.37-1.008-.794h.582c.008.178.186.306.422.309.254 0 .424-.145.422-.35-.002-.195-.155-.348-.414-.348h-.3zm-.004-4.699h-.604v-.035c0-.408.295-.844.958-.844.583 0 .96.326.96.756 0 .389-.257.617-.476.848l-.537.572v.03h1.054V9H1.143v-.395l.957-.99c.138-.142.293-.304.293-.508 0-.18-.147-.32-.342-.32a.33.33 0 0 0-.342.338zM2.564 5h-.635V2.924h-.031l-.598.42v-.567l.629-.443h.635z"/>
</svg>`

var tmplSetHeaders = `
	<div class="intro">
	Categories allow for more complex scoring mechanisms. Each ordinary bonus or combo can be marked as belonging to a particular category within each set. Sets of categories can be used to represent entities such as county, country, activity, etc. Such memberships can be used to modify basic bonus scoring and/or apply a second level of scoring using complex calculation rules.
	</div>

<article class="setheads">
	<button id="addset" class="plus" autofocus title="Add new set" onclick="addCatSet(this)">+</button>

	{{range $ix,$el := .Sets}}
		{{if ne $el ""}}
			<fieldset class="sethdr">
				<label for="SetHdr{{inc $ix}}">Set {{inc $ix}} is</label>
				<input type="text" id="SetHdr{{inc $ix}}" name="SetHdr{{inc $ix}}" data-set="{{inc $ix}}" value="{{$el}}" onchange="saveCatSet(this)" onclick="showCatSet(this)">
				<button class="plus" data-set="{{inc $ix}}" onclick="showCatSet(this)" >` + ordered_list_icon + `</button>
			</fieldset>
			{{end}}
	{{end}}


</article>
<hr>
<article class="setcats" id="setcats">
</article>
`

func build_axisLabels() []string {

	sqlx := "SELECT IfNull(Cat1Label,'')"
	for i := 2; i <= NumCategoryAxes; i++ {
		sqlx += ",IfNull(Cat" + strconv.Itoa(i) + "Label,'')"
	}
	sqlx += " FROM rallyparams"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	var res []string
	s := make([]string, NumCategoryAxes)
	for rows.Next() {
		err = rows.Scan(&s[0], &s[1], &s[2], &s[3], &s[4], &s[5], &s[6], &s[7], &s[8])
		checkerr(err)
	}
	res = append(res, s...)
	//log.Printf("AxisLabels = %v\n", res)
	return res

}

func addCatCat(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("addCatCat")
	set := r.FormValue("s")
	if intval(set) < 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad set index"}`)
	}
	sqlx := "SELECT max(Cat) FROM categories WHERE Axis=" + set
	cat := getIntegerFromDB(sqlx, 0) + 1
	sqlx = fmt.Sprintf("INSERT INTO categories(Axis,Cat) VALUES(%v,%v)", set, cat)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprintf(w, `{"ok":true,"msg":"%v"}`, cat)
	//fmt.Printf(" done = %v\n", cat)
}

func delCatCat(w http.ResponseWriter, r *http.Request) {

	fmt.Println("delCatCat")
	set := r.FormValue("s")
	if intval(set) < 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad set index"}`)
	}
	cat := r.FormValue("c")
	if intval(cat) < 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad cat index"}`)
	}
	sqlx := fmt.Sprintf("DELETE FROM categories WHERE Axis=%v AND Cat=%v", set, cat)
	fmt.Println(sqlx)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprintf(w, `{"ok":true,"msg":"%v"}`, cat)
}

func fetchSetCats(set int, alphabetic bool) []CatDefinition {

	res := make([]CatDefinition, 0)
	sqlx := "SELECT Cat,ifnull(BriefDesc,Cat) FROM categories WHERE Axis=" + strconv.Itoa(set)
	if alphabetic {
		sqlx += " ORDER BY BriefDesc"
	} else {
		sqlx += " ORDER BY Cat"
	}
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var cd CatDefinition
		err = rows.Scan(&cd.Cat, &cd.CatName)
		checkerr(err)
		cd.Set = set
		res = append(res, cd)
	}
	return res
}

func showCategoryCats(w http.ResponseWriter, r *http.Request) {

	type catEntry struct {
		Cat     int
		CatDesc string
	}
	var set struct {
		OK      bool   `json:"ok"`
		Msg     string `json:"msg"`
		Set     int
		SetName string
		Cats    []catEntry
	}
	set.Set = intval(r.FormValue("s"))
	if set.Set < 1 || set.Set > NumCategoryAxes {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad set index"}`)
		return
	}
	sqlx := fmt.Sprintf("SELECT Cat%vLabel FROM rallyparams", set.Set)
	set.SetName = getStringFromDB(sqlx, fmt.Sprintf("%v", set.Set))
	sqlx = fmt.Sprintf("SELECT Cat,ifnull(BriefDesc,'') FROM categories WHERE Axis=%v ORDER BY Cat", set.Set)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var ce catEntry
		err = rows.Scan(&ce.Cat, &ce.CatDesc)
		checkerr(err)
		set.Cats = append(set.Cats, ce)
	}
	set.OK = true
	set.Msg = "ok"
	bytes, err := json.Marshal(set)
	checkerr(err)
	fmt.Fprintf(w, "%v", string(bytes))

}
func showCategorySets(w http.ResponseWriter, r *http.Request) {
	type sets struct {
		Sets []string
	}
	var s sets
	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}
	s.Sets = build_axisLabels()

	t, err := template.New("sets").Funcs(funcMap).Parse(tmplSetHeaders)
	checkerr(err)

	startHTML(w, "Categories")
	fmt.Fprint(w, `</header>
	`)
	err = t.Execute(w, s)
	checkerr(err)
}

func updateCatName(w http.ResponseWriter, r *http.Request) {

	set := intval(r.FormValue("s"))
	cat := intval(r.FormValue("c"))
	cd := r.FormValue("setname")
	if set < 1 || set > NumCategoryAxes {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad index"}`)
		return
	}
	sqlx := "UPDATE categories SET BriefDesc=? WHERE Axis=? AND Cat=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(cd, set, cat)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)

}

func updateSetName(w http.ResponseWriter, r *http.Request) {

	ix := intval(r.FormValue("s"))
	nm := r.FormValue("ff")
	xx := r.FormValue(nm)
	if ix < 1 || ix > NumCategoryAxes {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad index"}`)
		return
	}
	if nm != "setname" {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad fieldname"}`)
		return
	}
	nm = fmt.Sprintf("Cat%vLabel", ix)
	sqlx := "UPDATE rallyparams SET " + nm + "=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(xx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)

}
