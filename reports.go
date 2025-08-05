// T O D O
//
// # Combo analysis
//
// v3 used CombosTicked stored in entrant record
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FinisherExport struct {
	RiderFirst   string
	RiderLast    string
	RiderIBA     string
	PillionFirst string
	PillionLast  string
	PillionIBA   string
	Bike         string
	BikeReg      string
	Rank         int
	Distance     int
	Class        string
	Phone        string
	Email        string
	Country      string
}

var bonus_anal_hdr = `
    <html>
	<head>
	<title>Bonus analysis</title>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="stylesheet" href="/css?file=normalize">
	<link rel="stylesheet" href="/css?file=maincss">
	<script src="/js?file=mainscript"></script>

    <script src="https://bossanova.uk/jspreadsheet/v4/jexcel.js"></script>
    <script src="https://jsuites.net/v4/jsuites.js"></script>
    <link rel="stylesheet" href="https://bossanova.uk/jspreadsheet/v4/jexcel.css" type="text/css" />
    <link rel="stylesheet" href="https://jsuites.net/v4/jsuites.css" type="text/css" />
    </head><body>
	<header>
`

var bonus_anal_top = `
<article class="reportba">
    <p><button id="exportcsv">Save as CSV</button> <br> Flags: 
    <strong>A</strong>lert - <strong>B</strong>ike in photo - <strong>D</strong>aylight only - 
    <strong>F</strong>ace in photo - <strong>N</strong>ight only - <strong>R</strong>estricted - 
    <strong>T</strong>icket/receipt</p>
    <table id="bonusdump"><caption>%v</caption><thead>
    <tr><th>Bonus</th><th>Name</th>
    <th>Claims</th><th>Points</th><th>Flags</th>%v<th>Combos</th>
    </tr></thead><tbody>
`

/*
var combo_anal_top = `
<article class="reportcmb">
<p><button id="exportcsv">Save as CSV</button></p>
    <table id="bonusdump"><caption>%v</caption><thead>
    <tr><th>Combo</th><th>Name</th>
    <th>Points</th><th>Bonuses</th>Needed<th>Scored</th>
    </tr></thead><tbody>
`
*/

var bonus_anal_script = `
    <script>
        var table = jspreadsheet(document.getElementById('bonusdump'), {
            filters: true,
            includeHeadersOnDownload: true,
            editable: false,
            rowDrag: false,
            allowInsertRow: false,
            allowManualInsertRow: false,
            allowInsertColumn: false,
            allowManualInsertColumn: false,
            allowDeleteRow: false,
            allowDeleteColumn: false,
            csvFileName: 'bonuses',
        })
        document.getElementById('exportcsv').onclick = function() {
            table.download();
        }
    </script>
`

const EXPORT_BONUS_SELECT = "SELECT bonuses.BonusID AS Bonus,bonuses.BriefDesc AS Name,IFNULL(bonusclaims.Claims,0) AS Claims,Points,IFNULL(Flags,'')"
const EXPORT_BONUS_FILES = " FROM bonuses LEFT JOIN (SELECT BonusID,COUNT(DISTINCT EntrantID) AS Claims FROM claims GROUP BY BonusID) AS bonusclaims ON bonuses.BonusID=bonusclaims.BonusID;"

type BonusAnalRec struct {
	Bonus     string
	BriefDesc string
	Claims    int
	Points    int
	Flags     string
	Cat       [NumCategoryAxes]int
}

// combo_bonus_array shows which bonuses are used in which combos
// Key is bonus, array lists combos with that bonus
type bonus_combos []string

var combo_bonus_array = make(map[string]bonus_combos, 0)

func buildCBA() {

	sqlx := "SELECT ComboID,ifnull(Bonuses,'') FROM combos"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var cmb string
		var bonuses string
		err = rows.Scan(&cmb, &bonuses)
		checkerr(err)
		bs := strings.Split(bonuses, ",")
		for i := 0; i < len(bs); i++ {
			_, ok := combo_bonus_array[bs[i]]
			if !ok {
				combo_bonus_array[bs[i]] = make(bonus_combos, 0)
			}
			combo_bonus_array[bs[i]] = append(combo_bonus_array[bs[i]], cmb)
		}
	}
}

func exportBonusesReport(w http.ResponseWriter, r *http.Request) {

	sets := build_axisLabels()
	cathdrs := ""
	cats := ""
	for i := range sets {
		cats += fmt.Sprintf(`,Cat%v`, i+1)
		if sets[i] == "" {
			continue
		}
		cathdrs += `<th>` + sets[i] + `</th>`
	}
	buildCBA()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Fprint(w, bonus_anal_hdr)
	showTopbar(w, "Bonus Analysis")
	fmt.Fprint(w, `</header>`)
	fmt.Fprintf(w, bonus_anal_top, CS.Basics.RallyTitle, cathdrs)
	sqlx := EXPORT_BONUS_SELECT + cats + EXPORT_BONUS_FILES
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var br BonusAnalRec

		s := reflect.ValueOf(&br).Elem()
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

		fmt.Fprint(w, `<tr>`)
		fmt.Fprintf(w, `<td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td>`, br.Bonus, br.BriefDesc, br.Claims, br.Points, br.Flags)
		for i := range sets {
			if sets[i] == "" {
				continue
			}
			catx := ""
			if br.Cat[i] > 0 {
				catx = getStringFromDB(fmt.Sprintf("SELECT ifnull(BriefDesc,'') FROM categories WHERE Axis=%v AND Cat=%v", i+1, br.Cat[i]), fmt.Sprintf("%v", br.Cat[i]))
			}
			fmt.Fprintf(w, `<td>%v</td>`, catx)
		}
		fmt.Fprint(w, `<td>`)
		cmb, ok := combo_bonus_array[br.Bonus]
		if ok {
			fmt.Fprint(w, cmb)
		}
		fmt.Fprint(w, `</td>`)

		fmt.Fprint(w, `</tr>`)
	}
	fmt.Fprint(w, `</tbody></table>`)
	fmt.Fprint(w, bonus_anal_script)
	fmt.Fprint(w, `</article>`)
}

const exFinisherSQL = `
SELECT ifnull(RiderFirst,''),ifnull(RiderLast,''),ifnull(RiderIBA,''),ifnull(PillionFirst,''),ifnull(PillionLast,''),ifnull(PillionIBA,'')
,ifnull(Bike,''),ifnull(BikeReg,''),FinishPosition AS Rank,CorrectedMiles,ifnull(BriefDesc,''),ifnull(Phone,''),ifnull(Email,''),ifnull(Country,'UK')
 FROM entrants LEFT JOIN classes ON entrants.Class=classes.Class WHERE Rank > 0 ORDER BY Rank
`
const exFinisherCols = `RiderFirst,RiderLast,RiderIBA,PillionFirst,PillionLast,PillionIBA,Bike,BikeReg,Rank,Distance,Class,Phone,Email,Country`

func exportFinisherCSV(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=finishers.csv")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	rows, err := DBH.Query(exFinisherSQL)
	checkerr(err)
	defer rows.Close()
	fmt.Fprintln(w, exFinisherCols)
	x := csv.NewWriter(w)
	for rows.Next() {
		var fx FinisherExport
		err = rows.Scan(&fx.RiderFirst, &fx.RiderLast, &fx.RiderIBA, &fx.PillionFirst, &fx.PillionLast, &fx.PillionIBA, &fx.Bike, &fx.BikeReg, &fx.Rank, &fx.Distance, &fx.Class, &fx.Phone, &fx.Email, &fx.Country)
		checkerr(err)
		fxx := make([]string, 0)
		fxx = append(fxx, fx.RiderFirst)
		fxx = append(fxx, fx.RiderLast)
		fxx = append(fxx, fx.RiderIBA)
		fxx = append(fxx, fx.PillionFirst)
		fxx = append(fxx, fx.PillionLast)
		fxx = append(fxx, fx.PillionIBA)
		fxx = append(fxx, fx.Bike)
		fxx = append(fxx, fx.BikeReg)
		fxx = append(fxx, strconv.Itoa(fx.Rank))
		fxx = append(fxx, strconv.Itoa(fx.Distance))
		fxx = append(fxx, fx.Class)
		fxx = append(fxx, fx.Phone)
		fxx = append(fxx, fx.Email)
		fxx = append(fxx, fx.Country)
		err = x.Write(fxx)
		checkerr(err)
	}
	x.Flush()
}

func exportFinisherJSON(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=finishers.json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	mk := "M"
	if CS.Basics.RallyUnitKms {
		mk = "K"
	}

	fmt.Fprintf(w, `{"filetype":"ibachasm","fileversion": 1,"rally":"%v","mk":"%v","asat":"%v","finishers":[`, CS.Basics.RallyTitle, mk, time.Now().Format(timefmt))

	rows, err := DBH.Query(exFinisherSQL)
	checkerr(err)
	defer rows.Close()
	comma := false
	for rows.Next() {
		var fx FinisherExport
		err = rows.Scan(&fx.RiderFirst, &fx.RiderLast, &fx.RiderIBA, &fx.PillionFirst, &fx.PillionLast, &fx.PillionIBA, &fx.Bike, &fx.BikeReg, &fx.Rank, &fx.Distance, &fx.Class, &fx.Phone, &fx.Email, &fx.Country)
		checkerr(err)
		b, err := json.Marshal(fx)
		checkerr(err)
		if comma {
			fmt.Fprint(w, `,`)
		}
		fmt.Fprintln(w, string(b))
		comma = true

	}
	fmt.Fprint(w, `]}`)

}
