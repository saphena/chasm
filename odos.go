package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//go:embed odos.js
var odosjs string

const timefmt = "2006-01-02T15:04"

const timerticker = `var img = document.getElementById('ticker');

var interval = window.setInterval(function(){
    let paused = document.getElementById('timenow');
	if(paused) {paused = paused.getAttribute('data-paused')=='1';}
	if (img) {
    if(!paused && img.style.visibility == 'hidden'){
        img.style.visibility = 'visible';
    }else{
        img.style.visibility = 'hidden';
    }
}
}, 1000);`

// When showing the odo capture list, this sets the time shown in the header
func get_odolist_start_time(ischeckout bool) string {

	res := storeTimeDB((time.Now()))
	if !ischeckout {
		return res
	}
	// Need to show next available start rather than real time
	st := CS.Basics.RallyStarttime //getStringFromDB("SELECT StartTime FROM rallyparams", "")
	if st == "" {
		return res
	}
	fmt.Println("Starting at " + st)
	res = res[0:11] + st[11:]
	return res

}

func show_odo(w http.ResponseWriter, r *http.Request, showstart bool) {

	if r.FormValue("debug") != "" {
		fmt.Println("show_odo called")
	}

	sqlx := "SELECT EntrantID,RiderFirst,RiderLast,ifnull(OdoRallyStart,''),ifnull(StartTime,''),ifnull(OdoRallyFinish,''),ifnull(FinishTime,''),EntrantStatus,OdoKms"
	sqlx += " FROM entrants WHERE "
	st := get_odolist_start_time(showstart)
	sclist := ""
	if showstart {
		sclist = strconv.Itoa(EntrantDNS) + "," + strconv.Itoa(EntrantOK)
	} else {
		sclist = strconv.Itoa(EntrantOK) + "," + strconv.Itoa(EntrantFinisher)
	}
	sqlx += " EntrantStatus IN (" + sclist + ")"
	sqlx += " ORDER BY RiderLast,RiderFirst"
	//fmt.Println(sqlx)
	rows, err := DBH.Query(sqlx)
	checkerr(err)

	x := "Check-IN"
	if showstart {
		x = "Check-OUT"
	}
	startHTML(w, x)

	fmt.Fprint(w, `<script src="/js?file=odosjs"></script>`)
	fmt.Fprint(w, `<div class="odohdr">`)

	odoname := ""
	if showstart {
		fmt.Fprint(w, " START")
		odoname = "s"
	} else {
		fmt.Fprint(w, " FINISH")
		odoname = "f"
	}
	fmt.Fprintf(w, ` <span id="timenow" data-time="%v" data-refresh="1000" data-pause="120000" data-paused="0"`, st)

	fmt.Fprintf(w, ` >%v</span>`, st[11:16])
	if !showstart {
		fmt.Fprint(w, ` <span id="ticker">&diams;</span>`)
	}

	if !showstart {
		const holdlit = ` stop clock `
		const unholdlit = ` restart clock `
		fmt.Fprintf(w, ` <button data-hold="%v" data-unhold="%v" onclick="clickTimeBtn(this);" id="pauseTime">%v</button>`, holdlit, unholdlit, holdlit)
	}
	fmt.Fprint(w, `<script>`+timerticker+`</script>`)

	fmt.Fprint(w, ` <span id="errlog"></span>`) // Diags only
	fmt.Fprint(w, `</div>`)

	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<script>setTimeout(reloadPage,30000);refreshTime(); timertick = setInterval(refreshTime,1000);</script>`)

	fmt.Fprint(w, `<div class="odolist">`)
	oe := true
	itemno := 0
	for rows.Next() {
		var EntrantID int
		var RiderFirst, RiderLast, OdoStart, StartTime, OdoFinish, FinishTime string
		var EntrantStatus int
		var OdoCounts string
		rows.Scan(&EntrantID, &RiderFirst, &RiderLast, &OdoStart, &StartTime, &OdoFinish, &FinishTime, &EntrantStatus, &OdoCounts)
		itemno++
		fmt.Fprint(w, `<div class="odorow `)
		if oe {
			fmt.Fprint(w, "odd")
		} else {
			fmt.Fprint(w, "even")
		}
		oe = !oe
		fmt.Fprint(w, `">`)

		fmt.Fprintf(w, `<label for="%v" class="name"><strong>%v</strong>, %v</label> `, itemno, RiderLast, RiderFirst)
		pch := "finish odo"
		val := OdoFinish
		if showstart {
			pch = "start odo"
			val = OdoStart
		}
		fmt.Fprintf(w, `<span><input id="%v" data-e="%v" data-st="%v" data-save="saveOdo" name="%v" type="number" class="bignumber" oninput="oi(this);" onchange="oc(this);" min="0" placeholder="%v" value="%v"></span>`, itemno, EntrantID, StartTime, odoname, pch, val)
		fmt.Fprint(w, `</div>`)

	}

	fmt.Fprint(w, `</div><footer></footer></body></html>`)
}

func update_odo(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Here we go")
	fmt.Printf("%v\n\n", r)
	if r.FormValue("e") == "" || r.FormValue("ff") == "" || r.FormValue("v") == "" {
		fmt.Fprint(w, `{"err":false,"msg":"ok"}`)
		return
	}

	dt := r.FormValue("t")
	if dt == "" {
		dt = storeTimeDB(time.Now())
	}
	sqlx := ""
	switch r.FormValue("ff") {
	case "f":
		sqlx = "OdoRallyFinish=" + r.FormValue("v")
		sqlx += ",OdoCheckFinish=" + r.FormValue("v")

		sqlx += ",CorrectedMiles=(" + r.FormValue("v") + " - IfNull(OdoRallyStart,0))"

		ns := EntrantFinisher
		n, _ := strconv.Atoi(r.FormValue("v"))
		if n < 1 {
			ns = EntrantDNF
		}
		sqlx += ",FinishTime='" + dt + "'"
		sqlx += ",EntrantStatus=" + strconv.Itoa(ns)
		sqlx += " WHERE EntrantID=" + r.FormValue("e")
		sqlx += " AND FinishTime IS NULL"
		sqlx += " AND EntrantStatus IN (" + strconv.Itoa(EntrantOK) + "," + strconv.Itoa(EntrantDNF) + ")"
	case "s":
		sqlx = "OdoRallyStart=" + r.FormValue("v")
		sqlx += ",OdoCheckStart=" + r.FormValue("v")

		sqlx += ",StartTime='" + dt + "'"
		sqlx += ",EntrantStatus=" + strconv.Itoa(EntrantOK)
		sqlx += " WHERE EntrantID=" + r.FormValue("e")
		sqlx += " AND EntrantStatus IN (" + strconv.Itoa(EntrantDNS) + "," + strconv.Itoa(EntrantOK) + ")"
	}
	fmt.Println(sqlx)
	_, err := DBH.Exec("UPDATE entrants SET " + sqlx)
	checkerr(err)

	fmt.Fprint(w, `{"err":false,"msg":"ok"}`)
	recalc_scorecard(intval(r.FormValue("e")))
	rankEntrants(false)
}

func show_odo_checks(w http.ResponseWriter, r *http.Request) {

	checkout := true
	if r.FormValue("check") == "in" {
		checkout = false
	}
	show_odo(w, r, checkout)

}

func storeTimeDB(t time.Time) string {

	res := t.Local().Format(timefmt)
	return res
}
