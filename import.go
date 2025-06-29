package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

// var SpecialEntrantFields = map[string]string(`[
//
//	{"RiderName": "ifnull(RiderFirst,'') || ' ' || ifnull(RiderLast,'')"},
//
// ]
var SpecialEntrantFields = map[string]string{
	"RiderName":   "ifnull(RiderFirst,'') || ' ' || ifnull(RiderLast,'')",
	"PillionName": "ifnull(PillionFirst,'') || ' '  || ifnull(PillionLast,'')",
	"Bike":        "",
}

const uploadSpace = (10 << 20) // 10MB
const tempfilename = "uploadeddata.csv"

var import1form = `
<script>
function uploadFile(obj) {
	console.log('Uploading file '+obj.value);
	document.getElementById('csvname').value=obj.value;
}
</script>
<article class="import">
<h1>Import data from CSV</h1>
<form action="/upload" method="post" enctype="multipart/form-data">
<fieldset>
	<label for="filetype">What are we importing? </label>
	<select id="filetype" name="filetype">
		<option selected value="entrants"> Entrants </option>
		<option value="bonuses"> Bonuses </option>
		<option value="combos"> Combos </option>
	</select>
</fieldset>
<fieldset>
	<label for="csvfile">Please choose the input file </label>
	<input type="file" id="csvfile" name="csvfile" onchange="uploadFile(this)">
	<input type="hidden" id="csvname" name="csvname">
</fieldset>
<input type="submit" value=" Upload the file ">
</form>
</article>
`

var import2form = `
<article class="import">
<h1>Import data from CSV</h1>
<form action="/upload" method="post" enctype="multipart/form-data">
<fieldset>
	<span>What are we importing? </span>
	<span><strong>%v</strong></span>
	<input type="hidden" name="filetype" value="%v">
</fieldset>
<fieldset>
	<span>Loading from file </span>
	<span><strong>%v</strong></span>
	<input type="hidden" name="csvname" value="%v">
</fieldset>
<fieldset>
<input type="hidden" name="fieldmap" value="%v">

<input type="submit" value="Go ahead, import me">
</fieldset>
<hr>
</article>
`

func tableFields(tab string) []string {

	type table_info struct {
		cid     int
		cname   string
		ctype   string
		notnull int
		defval  any
		pk      int
	}

	rows, err := DBH.Query("pragma table_info(" + tab + ")")
	checkerr(err)
	defer rows.Close()
	cols := make([]string, 0)
	for rows.Next() {
		var ti table_info
		err = rows.Scan(&ti.cid, &ti.cname, &ti.ctype, &ti.notnull, &ti.defval, &ti.pk)
		checkerr(err)
		cols = append(cols, ti.cname)
	}
	rows.Close()
	return cols

}

func importMappedData(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Importing data")
	fmt.Fprint(w, `</header>`)
	hdr := strings.Split(r.FormValue("fieldmap"), ",")
	fmt.Fprint(w, `<ul>`)
	fldx := ""
	valx := ""
	colx := make([]int, 0)
	k := 0
	for i := 0; i < len(hdr); i++ {
		fld := r.FormValue(fmt.Sprintf("mapcol%v", i))
		if fld != "" {
			fmt.Fprintf(w, `<li>%v == %v</li>`, fld, i)
			if fldx != "" {
				fldx += ","
				valx += ","
			}
			fldx += fld
			valx += "?"
			colx = append(colx, i)
			k++

		}
	}
	sqlx := fmt.Sprintf(`INSERT OR IGNORE INTO %v (%v) VALUES(%v)`, r.FormValue("filetype"), fldx, valx)
	fmt.Fprint(w, `</ul>`)
	fmt.Fprintf(w, `<p>%v</p>`, sqlx)
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	csvfile := filepath.Join(CS.UploadsFolder, tempfilename)

	f, err := os.Open(csvfile)
	checkerr(err)
	defer f.Close()

	csvReader := csv.NewReader(f)
	skiphdr := true
	_, err = DBH.Exec("BEGIN TRANSACTION")
	checkerr(err)
	defer DBH.Exec("ROLLBACK")
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		checkerr(err)
		if skiphdr {
			skiphdr = false
			continue
		}
		args := make([]any, 0)
		for i := 0; i < len(colx); i++ {

			args = append(args, rec[colx[i]])
		}
		_, err = stmt.Exec(args...)
		checkerr(err)

	}
	_, err = DBH.Exec("COMMIT")
	checkerr(err)

}

func showImport(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Data import")

	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, import1form)

}

func uploadImport(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(uploadSpace)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	if r.FormValue("fieldmap") != "" {
		importMappedData(w, r)
		return
	}

	file, _, err := r.FormFile("csvfile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	csvfile := filepath.Join(CS.UploadsFolder, tempfilename)
	tempFile, err := os.Create(csvfile)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	tempFile.Close()

	table := r.FormValue("filetype")
	cols := tableFields(table)

	var ignore = []string{""}
	cols = append(cols, ignore...)
	slices.Sort(cols)
	csvname := r.FormValue("csvname")

	f, err := os.Open(csvfile)
	checkerr(err)
	defer f.Close()

	csvReader := csv.NewReader(f)

	hdr, err := csvReader.Read()
	checkerr(err)
	dta, err := csvReader.Read()
	checkerr(err)

	startHTML(w, "Data uploaded")

	fmt.Fprint(w, `</header>`)

	fmt.Fprintf(w, import2form, strings.ToUpper(table), table, csvname, csvname, strings.Join(hdr, ","))

	fmt.Fprint(w, `<article class="importmap">`)

	for i := 0; i < len(hdr); i++ {
		fmt.Fprint(w, `<div class="row">`)
		fmt.Fprintf(w, `<span class="col">%v</span>`, hdr[i])
		fmt.Fprintf(w, `<span class="col">%v</span>`, dta[i])
		fmt.Fprintf(w, `<select name="mapcol%v">`, i)
		for j := 0; j < len(cols); j++ {
			fmt.Fprintf(w, `<option value="%v"`, cols[j])
			if strings.EqualFold(hdr[i], cols[j]) {
				fmt.Fprint(w, ` selected`)
			}
			fmt.Fprintf(w, `>%v</option>`, cols[j])
		}
		fmt.Fprint(w, `</select>`)
		fmt.Fprint(w, `</div>`)
	}
	fmt.Fprint(w, `</article>`)
	fmt.Fprint(w, `</form>`)
}
