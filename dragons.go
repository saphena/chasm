package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

	fmt.Println("saveRawOptions")
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
