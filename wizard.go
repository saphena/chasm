package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func showWizardPage(w http.ResponseWriter, r *http.Request, pg int) {

	var sqlx string = ""

	switch pg {
	case 1:

		fmt.Printf("Running wizard page 1\n")
		type region struct {
			Region  string
			Current bool
		}
		var cfg struct {
			Regions   []region
			Eventname string
			Region    string
			Virtual   bool
		}

		sqlx = "SELECT eventname,region,virtualrally FROM config"
		rows, _ := DBH.Query(sqlx)
		if !rows.Next() {
			return
		}
		var v int
		rows.Scan(&cfg.Eventname, &cfg.Region, &v)
		cfg.Virtual = v != 0

		rows.Close()

		sqlx = "SELECT region FROM regions"
		rows, _ = DBH.Query(sqlx)
		cfg.Regions = make([]region, 0)
		for rows.Next() {
			var reg region
			rows.Scan(&reg.Region)
			reg.Current = reg.Region == cfg.Region
			cfg.Regions = append(cfg.Regions, reg)
		}
		rows.Close()

		tmplt := "wizard"
		t := template.Must(template.New(tmplt + ".html").Option("missingkey=error").ParseFiles("files/" + tmplt + ".html"))

		err := t.Execute(w, cfg)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

	}
}
