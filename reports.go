package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

var bonus_anal_hdr = `
    <html>
    <script src="https://bossanova.uk/jspreadsheet/v4/jexcel.js"></script>
    <script src="https://jsuites.net/v4/jsuites.js"></script>
    <link rel="stylesheet" href="https://bossanova.uk/jspreadsheet/v4/jexcel.css" type="text/css" />
    <link rel="stylesheet" href="https://jsuites.net/v4/jsuites.css" type="text/css" />
    
    <h1>%v</h1>
    <h2>Bonus analysis</h2>
    <p><button id="exportcsv">Save as CSV</button>  Flags: 
    <strong>A</strong>lert - <strong>B</strong>ike in photo - <strong>D</strong>aylight only - 
    <strong>F</strong>ace in photo - <strong>N</strong>ight only - <strong>R</strong>estricted - 
    <strong>T</strong>icket/receipt</p>
    <table id="bonusdump"><caption>%v</caption><thead>
    <tr><th>Bonus</th><th>Name</th>
    <th>Claims</th><th>Points</th><th>Flags</th>%v<th>Combos</th>
    </tr></thead><tbody>
	
	`
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
	fmt.Fprintf(w, bonus_anal_hdr, CS.Basics.RallyTitle, CS.Basics.RallyTitle, cathdrs)
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
}
