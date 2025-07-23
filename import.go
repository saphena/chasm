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

type SpecialConstructors = map[string]string

// TableSpecials are used to complete 'missing' fields when importing
// data from CSV files.
var TableSpecials = map[string]SpecialConstructors{
	"entrants": {
		"RiderName":   "trim(RiderFirst || ' ' || RiderLast)",
		"PillionName": "trim(PillionFirst || ' '  || PillionLast)",
		"RiderFirst":  "trim(substr(RiderName,1,instr(RiderName,' ')))",
		"RiderLast":   "trim(substr(RiderName,instr(RiderName,' ')))",
		"Bike":        "trim(BikeMake || ' ' || BikeModel)",
		"BikeMake":    "trim(substr(Bike,1,instr(Bike,' ')))",
		"BikeModel":   "trim(substr(Bike,instr(Bike,' ')))",
	},
	"bonuses": {
		"BriefDesc": "BonusID",
	},
	"combos": {
		"BriefDesc": "ComboID",
	},
}

const uploadSpace = (10 << 20) // 10MB
const tempfilename = "uploadeddata.csv"

var importImgform = `
<script>
function uploadFile(obj) {
	console.log('Uploading file '+obj.value);
	document.getElementById('imgname').value=obj.value;
	let frm = obj.form;
	let btn = frm.querySelector('button');
	if (btn) btn.disabled = false;
}
</script>
<article class="import">
<h1>Upload bonus image</h1>
<form action="/upload" method="post" enctype="multipart/form-data">
<input type="hidden" name="filetype" value="bonusimage">
<fieldset>
	<label for="imgfile">Please choose the input file </label>
	<input type="file" id="imgfile" name="imgfile" onchange="uploadFile(this)">
	<input type="hidden" id="imgname" name="imgname">
</fieldset>
<button disabled> Upload the file </button>
</form>
</article>
`

var import1form = `
<script>
function uploadFile(obj) {
	console.log('Uploading file '+obj.value);
	document.getElementById('csvname').value=obj.value;
	let frm = obj.form;
	let btn = frm.querySelector('button');
	if (btn) btn.disabled = false;
}
</script>
<article class="import">
<h1>Import data from CSV</h1>
<form action="/upload" method="post" enctype="multipart/form-data">
<fieldset>
	<label for="filetype">What are we importing? </label>
	<select id="filetype" name="filetype">
	##OPTIONS##
	</select>
</fieldset>
<fieldset>
	<label for="csvfile">Please choose the input file </label>
	<input type="file" id="csvfile" name="csvfile" onchange="uploadFile(this)">
	<input type="hidden" id="csvname" name="csvname">
</fieldset>
<button disabled> Upload the file </button>
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
	<select name="overwrite">
		<option selected value="add">Add new records only</option>
		<option value="replace">Replace whole dataset</option>
	</select>
</fieldset>

<input type="hidden" name="fieldmap" value="%v">

<button>Go ahead, import me</button>

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

func handleUploadedImage(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("imgfile")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgname := strings.Split(r.FormValue("imgname"), "\\")
	imgfile := filepath.Join(CS.ImgBonusFolder, imgname[len(imgname)-1])
	tempFile, err := os.Create(imgfile)
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

}

func importMappedData(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Importing data")
	fmt.Fprint(w, `</header>`)

	hdr := strings.Split(r.FormValue("fieldmap"), ",")
	fldx := ""
	valx := ""
	table := r.FormValue("filetype")
	overwrite := r.FormValue("overwrite") == "replace"
	specials := TableSpecials[table]

	colx := make([]int, 0)
	k := 0
	for i := 0; i < len(hdr); i++ {
		fld := r.FormValue(fmt.Sprintf("mapcol%v", i))
		if fld != "" {
			if fldx != "" {
				fldx += ","
				valx += ","
			}
			fldx += fld

			_, ok := specials[fld]
			if ok {
				delete(specials, fld)
			}

			valx += "?"
			colx = append(colx, i)
			k++

		}
	}
	sqlx := fmt.Sprintf(`INSERT OR IGNORE INTO %v (%v) VALUES(%v)`, table, fldx, valx)
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

	if overwrite {
		_, err = DBH.Exec("DELETE FROM " + table)
		checkerr(err)
	}

	fmt.Fprint(w, `<article class="import">`)

	fmt.Fprintf(w, `<h1>Importing %v</h1>`, table)

	fmt.Fprint(w, `<ul class="loadlist">`)
	nrex := 0
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
		res, err := stmt.Exec(args...)
		checkerr(err)
		n, err := res.RowsAffected()
		checkerr(err)
		if n == 1 {
			fmt.Fprintf(w, `<li>%v</li>`, rec)
			nrex++
		}

	}
	fmt.Fprint(w, `</ul>`)
	for k, v := range specials {
		sqlx = "UPDATE " + table + " SET " + k + "=" + v
		_, err = DBH.Exec(sqlx)
		checkerr(err)
	}
	_, err = DBH.Exec("COMMIT")
	checkerr(err)

	fmt.Fprintf(w, `<hr><p>%v records imported</p>`, nrex)
	fmt.Fprint(w, `</article>`)

}

func showImport(w http.ResponseWriter, r *http.Request) {

	var tabs = []string{"entrants", "bonuses", "combos"}

	startHTML(w, "Data import")

	fmt.Fprint(w, `</header>`)

	tab := r.FormValue("type")
	if tab == "" {
		tab = tabs[0]
	}
	if tab == "img" {
		fmt.Fprint(w, importImgform)
		return
	}
	x := ""
	for k := range tabs {
		x += `<option `
		if tab == tabs[k] {
			x += "selected "
		}
		x += `value="` + tabs[k] + `">` + tabs[k] + `</option>`
	}

	fmt.Fprint(w, strings.ReplaceAll(import1form, "##OPTIONS##", x))

}

func uploadImportDatafile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(uploadSpace)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	if r.FormValue("filetype") == "bonusimage" {
		handleUploadedImage(w, r)
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
	cols := tableFields(table) // get all fields declared for this table

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
		fmt.Fprintf(w, `<span class="col src">%v</span>`, hdr[i])
		fmt.Fprintf(w, `<span class="col dta">%v</span>`, dta[i])
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
