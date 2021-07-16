package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

const max_wizard_page = 3

func showWizardPage(w http.ResponseWriter, r *http.Request, pg int) {

	var sqlx string = ""

	//fmt.Printf("Showing wizard page %v\n", pg)
	type region struct {
		Region  string
		Current bool
	}
	var cfg struct {
		Regions           []region
		Eventname         string
		Region            string
		MilesKms          bool
		Virtual           bool
		WizShowNext       bool
		WizShowPrev       bool
		Startdate         string
		Starttime         string
		Finishdate        string
		Finishtime        string
		Maxhours          int
		MaxmilesDNF       int
		MaxmilesPenalties int
		MinmilesDNF       int
		ExcessPPM         int
		Localtz           string
		Locale            string
		Hostcountry       string
		Decimalcomma      bool
		Tankrange         int
		MinsPerStop       int
	}
	cfg.WizShowNext = pg < max_wizard_page
	cfg.WizShowPrev = pg > 1

	sqlx = "SELECT Eventname,Region,MilesKms,VirtualRally,Startdate,Starttime,Finishdate,Finishtime,Maxhours,MaxmilesDNF,Localtz,Locale,Hostcountry,Decimalcomma,Tankrange,MinsPerStop FROM config"
	rows, _ := DBH.Query(sqlx)
	defer rows.Close()
	if !rows.Next() {
		return
	}
	var v, d int
	rows.Scan(&cfg.Eventname, &cfg.Region, &cfg.MilesKms, &v, &cfg.Startdate, &cfg.Starttime, &cfg.Finishdate, &cfg.Finishtime,
		&cfg.Maxhours, &cfg.MaxmilesDNF, &cfg.Localtz, &cfg.Locale, &cfg.Hostcountry, &d, &cfg.Tankrange, &cfg.MinsPerStop)
	cfg.Virtual = v != 0
	cfg.Decimalcomma = d != 0

	rows.Close()

	sqlx = "SELECT Region FROM regions"
	rows, _ = DBH.Query(sqlx)
	cfg.Regions = make([]region, 0)
	for rows.Next() {
		var reg region
		rows.Scan(&reg.Region)
		reg.Current = reg.Region == cfg.Region
		cfg.Regions = append(cfg.Regions, reg)
	}
	rows.Close()

	tt, _ := template.ParseGlob(Docroot + "/" + *Language + "/wiz*")
	tmplt := "wizpage" + strconv.Itoa(pg) + ".html"
	err := tt.ExecuteTemplate(w, tmplt, cfg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}

func saveWizardDetail(r *http.Request) {

	type nameval struct {
		Field, Value string
	}
	nv := make([]nameval, 0)

	if r.FormValue("Startdate") != "" && r.FormValue("Finishdate") != "" {
		datetimefmt := "2006-01-02 03:04"
		loc, err := time.LoadLocation(r.FormValue("Localtz"))
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		startdatetime := r.FormValue("Startdate") + " " + r.FormValue("Starttime")
		finishdatetime := r.FormValue("Finishdate") + " " + r.FormValue("Finishtime")

		t1, _ := time.ParseInLocation(datetimefmt, startdatetime, loc)
		t2, _ := time.ParseInLocation(datetimefmt, finishdatetime, loc)
		td := t1.Sub(t2)
		tdhx := r.FormValue("Maxhours")
		tdh := int(td.Hours())
		if tdhx != "" {
			tdh, _ = strconv.Atoi(tdhx)
			if int(td.Hours()) < tdh {
				tdh = int(td.Hours())
			}
		}
		nv = append(nv, nameval{"Maxhours", strconv.Itoa(tdh)})
	} else if r.FormValue("Maxhours") != "" {
		tdh, _ := strconv.Atoi(r.FormValue("Maxhours"))
		nv = append(nv, nameval{"Maxhours", strconv.Itoa(tdh)})
	}

	sql := "UPDATE config SET "
	for _, fld := range []string{"Eventname", "Region", "Startdate", "Starttime", "Finishdate", "Finishtime",
		"Hostcountry", "Localtz", "Locale", "Decimalcomma", "MilesKms", "VirtualRally", "Tankrange", "MinsPerStop", "DBInitialised"} {
		x := r.FormValue(fld)
		if x != "" {
			nv = append(nv, nameval{fld, x})
		}
	}
	if len(nv) > 0 {
		c := ""
		args := make([]interface{}, len(nv))
		for i, v := range nv {
			sql += c + v.Field + "=?"
			c = ","
			args[i] = v.Value
		}
		stmt, _ := DBH.Prepare(sql)
		defer stmt.Close()
		_, err := stmt.Exec(args...)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

	}

}
func showWizard(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		showWizardPage(w, r, 1)
	}
	pg, _ := strconv.Atoi(r.FormValue("wizpage"))
	if pg < 1 {
		pg = 1
	} else {
		saveWizardDetail(r)
	}
	if r.FormValue("nextpage") != "" {
		pg++
	} else if r.FormValue("prevpage") != "" {
		pg--
	}
	showWizardPage(w, r, pg)
}
