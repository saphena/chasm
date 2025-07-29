package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

type EntrantDBRecord struct {
	EntrantID            int
	Bike                 string
	BikeReg              string
	RiderFirst           string
	RiderLast            string
	RiderCountry         string
	RiderIBA             string
	RiderPhone           string
	RiderEmail           string
	PillionFirst         string
	PillionLast          string
	PillionIBA           string
	PillionPhone         string
	PillionEmail         string
	OdoKms               string
	OdoStart             int
	OdoFinish            string
	CorrectedMiles       string
	FinishTime           string
	StartTime            string
	EntrantStatus        int
	NokName              string
	NokPhone             string
	NokRelation          string
	EntryDonation        string
	SquiresCheque        string
	SquiresCash          string
	RBLRAccount          string
	JustGivingAmt        string
	JustGivingURL        string
	Route                string
	RiderRBL             string
	PillionRBL           string
	Tshirt1              string
	Tshirt2              string
	Patches              int
	FreeCamping          string
	CertificateDelivered string
	CertificateAvailable string
	FinishPosition       int
	TotalPoints          int
	TeamID               int
}

func ajaxFetchEntrantDetails(w http.ResponseWriter, r *http.Request) {

	e := intval(r.FormValue("e"))
	if e < 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"no e specified"}`)
		return
	}
	ed := fetchEntrantDetails(e)
	if ed.PillionName != "" {
		ed.RiderName += " &amp; " + ed.PillionName
	}
	tr := jsonBool(ed.PillionName != "" || ed.TeamID > 0)

	fmt.Fprintf(w, `{"ok":true,"msg":"ok","name":"%v","team":%v}`, ed.RiderName, tr)

}

func createEntrant(w http.ResponseWriter, r *http.Request) {

	var sqlx string
	entrant := intval(r.FormValue("e"))
	if entrant < 1 {
		entrant = getIntegerFromDB("SELECT max(EntrantID) FROM entrants", 0) + 1
	}
	sqlx = "INSERT INTO entrants(EntrantID) VALUES(?)"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	res, err := stmt.Exec(entrant)
	checkerr(err)
	n, err := res.RowsAffected()
	checkerr(err)
	if n != 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"insert failed"}`)
		return
	}
	fmt.Fprintf(w, `{"ok":true,"msg":"%v"}`, entrant)

}

// deleteEntrant deletes the entrant record only. Associated
// claims and ebclaims will be orphaned
func deleteEntrant(w http.ResponseWriter, r *http.Request) {

	entrant := r.PathValue("e")
	if entrant == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"incomplete request"}`)
		return
	}
	sqlx := "DELETE FROM entrants WHERE EntrantID=" + entrant
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}

func list_entrants(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Entrants")
	mk := CS.UnitMilesLit
	if CS.Basics.RallyUnitKms {
		mk = CS.UnitKmsLit
	}

	sqlx := `SELECT EntrantID,ifnull(RiderLast,''),ifnull(RiderFirst,''),ifnull(PillionLast,''),ifnull(PillionFirst,'')
			,ifnull(Bike,''),EntrantStatus,ifnull(FinishPosition,0),ifnull(TotalPoints,0),ifnull(CorrectedMiles,0)
			FROM entrants`

	seqs := []string{"EntrantID", "RiderLast,RiderFirst", "PillionLast,PillionFirst", "Bike", "EntrantStatus", "Finishposition", "TotalPoints", "CorrectedMiles"}
	hdrs := []string{"Flag", "Name", "Pillion", "Bike", "Status", "Rank", "Points", mk}
	ord := r.FormValue("o")
	if ord == "" {
		ord = seqs[0]
	}
	thisorder := ord
	desc := r.FormValue("d") == "d"
	if desc {
		xx := strings.Split(thisorder, ",")
		for i := 0; i < len(xx); i++ {
			xx[i] += " DESC"
		}
		thisorder = strings.Join(xx, ",")
	}

	sqlx += " ORDER BY " + thisorder
	//fmt.Println(sqlx)
	fmt.Fprint(w, `<div class="entrantlist hdr">`)
	fmt.Fprint(w, `<button class="plus" onclick="window.location.href='/entrant/0'" title="Add new entrant">+</button>`)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div class="entrantlist hdr"><div class="row hdr">`)

	for ix, x := range seqs {
		urlx := "/entrants?o=" + x
		if ord == x && !desc {
			urlx += "&d=d"
		}

		xclass := ""
		if ix+1 >= len(hdrs) {
			xclass = "num"
		}
		fmt.Fprintf(w, `<span class="col hdr %v" title="Click to reorder"><a href="%v">%v</a></span>`, xclass, urlx, hdrs[ix])

	}
	fmt.Fprint(w, `</div></div><hr>`)
	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<article class="entrantlist">`)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var e EntrantDBRecord
		err = rows.Scan(&e.EntrantID, &e.RiderLast, &e.RiderFirst, &e.PillionLast, &e.PillionFirst, &e.Bike, &e.EntrantStatus, &e.FinishPosition, &e.TotalPoints, &e.CorrectedMiles)
		checkerr(err)
		fmt.Fprintf(w, `<div class="row" onclick="window.location.href='/entrant/%v?back=/entrants'">`, e.EntrantID)

		fmt.Fprintf(w, `<span class="col mid">%v</span>`, e.EntrantID)
		fmt.Fprintf(w, `<span class="col"><strong>%v</strong>, %v</span>`, e.RiderLast, e.RiderFirst)
		pillion := ""
		if e.PillionLast != "" {
			pillion = e.PillionLast + ", " + e.PillionFirst
		}
		fmt.Fprintf(w, `<span class="col">%v</span>`, pillion)
		fmt.Fprintf(w, `<span class="col">%v</span>`, e.Bike)
		fmt.Fprintf(w, `<span class="col">%v</span>`, EntrantStatusLits[e.EntrantStatus])
		rank := ""
		if e.EntrantStatus == EntrantFinisher {
			rank = fmt.Sprintf("%v", e.FinishPosition)
		}
		fmt.Fprintf(w, `<span class="col mid">%v</span>`, rank)
		fmt.Fprintf(w, `<span class="col num">%v</span>`, e.TotalPoints)
		fmt.Fprintf(w, `<span class="col num">%v</span>`, e.CorrectedMiles)

		fmt.Fprint(w, `</div>`)
	}
	fmt.Fprint(w, `</article>`)
}

var tmplEntrantBasic = `
{{if ne .EntrantID 0}}
	<div class="topline">
		<fieldset>
			<button title="Delete this Entrant?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" style="display:none;" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
		
		<fieldset>
			<button id="updatedb" class="hideuntil" data-e="{{.EntrantID}}" title="Delete Entrant" disabled onclick="deleteEntrant(this)"></button>
		</fieldset>
	</div>
{{end}}

<article class="entrant basic">
	<fieldset>
		<label for="EntrantID">Entrant</label>
		<input id="EntrantID" name="EntrantID" {{if ne .EntrantID 0}}type="text" readonly{{else}}type="text" autofocus{{end}} class="EntrantID" value="{{if ne .EntrantID 0}}{{.EntrantID}}{{end}}" onblur="addEntrant(this)" title="Flag">
		<input type="text" id="RiderFirst" name="RiderFirst" class="RiderFirst" placeholder="first" value="{{.RiderFirst}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
		<input type="text" id="RiderLast" name="RiderLast" class="RiderLast" {{if ne .EntrantID 0}}autofocus{{end}} placeholder="last" value="{{.RiderLast}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="RiderPhone">Phone</label>
		<input type="text" id="RiderPhone" class="RiderPhone" name="Phone" value="{{.RiderPhone}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="RiderEmail">Email</label>
		<input type="text" id="RiderEmail" class="RiderEmail" name="Email" value="{{.RiderEmail}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="PillionFirst">Pillion</label>
		<input type="text" id="PillionFirst" class="PillionFirst" name="PillionFirst" placeholder="first" value="{{.PillionFirst}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
		<input type="text" id="PillionLast" class="PillionLast" name="PillionLast" placeholder="last" value="{{.PillionLast}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<legend>Emergency contact</legend>
		<input type="text" id="NokName" class="NokName" name="NokName" placeholder="name" value="{{.NokName}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
		<input type="text" id="NokPhone" class="NokPhone" name="NokPhone" placeholder="phone" value="{{.NokPhone}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
		<input type="text" id="NokRelation" class="NokRelation" name="NokRelation" place="relation" value="{{.NokRelation}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="Bike">Bike</label>
		<input type="text" id="Bike" name="Bike" class="Bike" value="{{.Bike}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
		<input type="text" id="BikeReg" name="BikeReg" class="BikeReg" value="{{.BikeReg}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="OdoKms">Odo counts</label>
		<select id="OdoKms" name="OdoKms" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
			<option value="M" {{if ne .OdoKms "K"}}selected{{end}}>miles</option>
			<option value="K" {{if eq .OdoKms "K"}}selected{{end}}>km</option>
		</select>
	</fieldset>
	<fieldset>
		<label for="OdoStart">Readings</label>
		<input type="number" id="OdoStart" class="odo" name="OdoRallyStart" title="start" value="{{.OdoStart}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
		<input type="number" id="OdoFinish" class="odo" name="OdoRallyFinish" title="finish" value="{{.OdoFinish}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="CorrectedMiles">Rally distance</label>
		<input type="number" class="CorrectedMiles" id="CorrectedMiles" name="CorrectedMiles" value="{{.CorrectedMiles}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="TotalPoints">Rally points</label>
		<input type="number" class="TotalPoints" id="TotalPoints" name="TotalPoints" value="{{.TotalPoints}}" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
	</fieldset>
	<fieldset>
		<label for="EntrantStatus">Status</label>
		<select id="EntrantStatus" name="EntrantStatus" data-save="saveEntrant" oninput="oi(this)" onchange="saveEntrant(this)">
			<option value="0" {{if eq .EntrantStatus 0}}selected{{end}}>DNS</option>
			<option value="1" {{if eq .EntrantStatus 1}}selected{{end}}>ok</option>
			<option value="8" {{if eq .EntrantStatus 8}}selected{{end}}>Finisher</option>
			<option value="3" {{if eq .EntrantStatus 3}}selected{{end}}>DNF</option>
		</select>
	</fieldset>
	<fieldset>
		<label for="Team">Team</label>
		<select id="Team" name="TeamID" onchange="saveEntrant(this)">
		##teams##
		</select>
	</fieldset>


</article>
`

func fetchEntrantRecord(entrant int) EntrantDBRecord {

	var er EntrantDBRecord

	if entrant < 1 {
		return er
	}
	sqlx := `SELECT EntrantID,ifnull(Bike,''),ifnull(BikeReg,''),ifnull(RiderFirst,''),ifnull(RiderLast,''),ifnull(Country,'')
		,ifnull(RiderIBA,''),ifnull(Phone,''),ifnull(Email,'')
		,ifnull(PillionFirst,''),ifnull(PillionLast,''),ifnull(PillionIBA,'')
		,ifnull(OdoKms,'M'),ifnull(OdoRallyStart,0),ifnull(OdoRallyFinish,0),ifnull(CorrectedMiles,0)
		,ifnull(FinishTime,''),ifnull(StartTime,''),EntrantStatus,ifnull(NokName,''),ifnull(NokPhone,''),ifnull(NokRelation,'')
		,FinishPosition,TotalPoints,TeamID		FROM entrants`
	sqlx += fmt.Sprintf(" WHERE EntrantID=%v", entrant)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&er.EntrantID, &er.Bike, &er.BikeReg, &er.RiderFirst, &er.RiderLast, &er.RiderCountry, &er.RiderIBA, &er.RiderPhone, &er.RiderEmail, &er.PillionFirst, &er.PillionLast, &er.PillionIBA, &er.OdoKms, &er.OdoStart, &er.OdoFinish, &er.CorrectedMiles, &er.FinishTime, &er.StartTime, &er.EntrantStatus, &er.NokName, &er.NokPhone, &er.NokRelation, &er.FinishPosition, &er.TotalPoints, &er.TeamID)
		checkerr(err)
	} else {
		er.EntrantID = entrant
	}

	return er

}

func saveEntrant(w http.ResponseWriter, r *http.Request) {

	e := r.FormValue("e")
	fld := r.FormValue("ff")
	if e == "" || fld == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"incomplete request"}`)
		return
	}
	val := r.FormValue(fld)

	sqlx := "UPDATE entrants SET " + fld + "=? WHERE EntrantID=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(val, e)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)

}
func showEntrant(w http.ResponseWriter, r *http.Request) {

	entrant := intval(r.PathValue("e"))
	er := fetchEntrantRecord(entrant)

	if r.FormValue("back") != "" {
		startHTMLBL(w, "Entrant detail", r.FormValue("back"))
	} else {
		startHTML(w, "Entrant detail")
	}
	fmt.Fprint(w, `</header>`)

	teamrecs := fetchTeams(true)
	teamopts := ""
	for i := 0; i < len(teamrecs); i++ {
		sel := ""
		if er.TeamID == teamrecs[i].TeamID {
			sel = "selected"
		}
		teamopts += fmt.Sprintf(`<option value="%v" %v>%v</option>`, teamrecs[i].TeamID, sel, teamrecs[i].TeamName)
	}
	t, err := template.New("EntrantDetail").Parse(strings.ReplaceAll(tmplEntrantBasic, "##teams##", teamopts))
	checkerr(err)
	err = t.Execute(w, er)
	checkerr(err)

}
