package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

var queryresult string

func execRawSQL(sqlx string) string {

	re := regexp.MustCompile(`(?i)^select `)
	if re.MatchString(sqlx) {
		return queryRawSQL(sqlx)
	}
	res, err := DBH.Exec(sqlx)
	if err != nil {
		return err.Error()
	}
	n, err := res.RowsAffected()
	checkerr(err)
	if n == 1 {
		return "A single record was affected"
	}
	return fmt.Sprintf("%v records affected", n)
}

func queryRawSQL(sqlx string) string {

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	cols, err := rows.Columns()
	checkerr(err)
	data := make([]sql.NullString, len(cols))
	var res string
	res = `<table><thead><tr>`
	for i := range cols {
		res += `<th>` + cols[i] + `</th>`
	}
	res += `</tr></thead><tbody>`
	for rows.Next() {
		columnPointers := make([]any, len(cols))
		for i := range columnPointers {
			columnPointers[i] = &data[i]
		}
		err = rows.Scan(columnPointers...)
		checkerr(err)
		res += `<tr>`
		for i := range data {
			res += `<td>`
			if data[i].Valid {
				res += data[i].String
			}
			res += `</td>`
		}
		res += `</tr>`
	}
	res += `</tbody></table>`

	//b, err := json.Marshal(data)
	//checkerr(err)
	//fmt.Printf("%s\n", b)

	return res
}
func runRawSQL(w http.ResponseWriter, r *http.Request) {

	sqlx := r.FormValue("sql")
	queryresult = ""
	fmt.Printf("RawSQL=%s\n", sqlx)
	if sqlx == "" {
		showRawSQL(w, r)
		return
	}
	res := execRawSQL(sqlx)
	//fmt.Printf("RawSqlResult=%s\n", res)
	queryresult = res
	showRawSQL(w, r)

}

func showRawSQL(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Raw SQL!!!")
	fmt.Fprint(w, `</header>`)
	fmt.Fprint(w, `<div class="sqlquery">`)
	fmt.Fprint(w, `<p>This allows execution of raw SQL against the database. I hope you know what you're doing.</p>`)
	fmt.Fprint(w, `<form action="/sql" method="post">`)
	fmt.Fprintf(w, `<input type="text" class="sqlquery" name="sql" value="%s">`, r.FormValue("sql"))
	fmt.Fprint(w, `<button>Run SQL</button>`)
	fmt.Fprint(w, `</form>`)
	res := queryresult
	if res != "" {
		fmt.Fprintf(w, `<p>%s</p>`, res)
	}
	fmt.Fprint(w, `</div>`)

}

func editRawOptions(w http.ResponseWriter, r *http.Request) {

	b, err := json.MarshalIndent(CS, "", "  ")
	checkerr(err)

	startHTML(w, "Edit raw options")
	fmt.Fprint(w, `</header>`)
	fmt.Fprint(w, `<div class="rawoptions">`)
	fmt.Fprint(w, `<button disabled onclick="saveRawOpts(this)">Save changes</button><br>`)
	fmt.Fprint(w, `<textarea id="rawopts" editable oninput="enableRawSave(this)">`)
	fmt.Fprintf(w, "%s", b)
	fmt.Fprint(w, `</textarea>`)
	fmt.Fprint(w, `</div>`)
}

func saveRawOptions(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("saveRawOptions")
	var mycs chasmSettings
	newcs := r.FormValue("v")
	err := json.Unmarshal([]byte(newcs), &mycs)
	if err != nil {
		editRawOptions(w, r)
		return
	}
	sqlx := "UPDATE config SET Settings=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(newcs)
	checkerr(err)
	CS = mycs
	fmt.Fprint(w, `{"ok":true,"msg":"options updated"}`)
}
