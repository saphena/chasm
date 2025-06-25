package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type emailSettings struct {
	SMTP struct {
		Host          string
		Port          string
		UseInboxCreds bool
		Userid        string
		Password      string
		CertName      string // May need to override the certificate name used for TLS
	}
	IMAP struct {
		HostPort string
		Userid   string
		Password string
	}
}
type RallyBasics struct {
	RallyTitle      string
	RallyStarttime  string
	RallyFinishtime string
	RallyMaxHours   int
	RallyUnitKms    bool
	RallyTimezone   string
}
type chasmSettings struct {
	ShowExcludedClaims  bool // If a claim is marked 'excluded' and is not superseded, show it on the scoresheet
	CurrentLeg          int
	UseCheckinForOdo    bool // If true, OdoRallyFinish updated only by check-in, not by individual claims
	Basics              RallyBasics
	UnitMilesLit        string
	UnitKmsLit          string
	PenaltyMilesDNF     int
	RallyMinMiles       int
	DebugRules          bool
	AutoLateDNF         bool
	RallyMinPoints      int
	RallyUseQA          bool
	RallyQAPoints       int
	RallyUsePctPen      bool
	RallyPctPenVal      int
	RallyRankEfficiency bool
	RallySplitTies      bool
	RallyTeamMethod     int
	FlagTeamTitle       string
	FlagAlertTitle      string
	FlagBikeTitle       string
	FlagDaylightTitle   string
	FlagFaceTitle       string
	FlagNightTitle      string
	FlagRestrictedTitle string
	FlagReceiptTitle    string
	CloseEBCUndecided   string
	CloseEBC            []string
	ImgBonusFolder      string // Holds rally book bonus photos
	ImgEbcFolder        string // Holds images captured from emails
	RallyBookImgTitle   string
	EBCImgTitle         string
	EBCImgSwapTitle     string
	Email               emailSettings
}

var CS chasmSettings

const defaultCS = `{
	"ShowExcludedClaims": 	false,
	"CurrentLeg": 			1,
	"UseCheckInForOdo": 	true,
	"RallyUnitKms": 		false,
	"UnitMilesLit":			"miles",
	"UnitKmsLit":			"km",
	"PenaltyMilesDNF":		99999,
	"RallyTitle":			"Brit Butt Rally 2025",
	"RallyMinMiles":		0,
	"DebugRules":			false,
	"AutoLateDNF": 			true,
	"RallyMinPoints":		-99999,
	"RallyTimezone":		"Europe/London",
	"RallyUseQA":			false,
	"RallyQAPoints":		50,
	"RallyUsePctPen":		false,
	"RallyPctPenVal":		10,
	"RallyRankEfficiency":	false,
	"RallySplitTies":		true,
	"RallyTeamMethod":		3,
	"FlagTeamTitle":       	"Team rules",
	"FlagAlertTitle":      	"Read the notes!",
	"FlagBikeTitle":       	"Bike in photo",
	"FlagDaylightTitle":   	"Daylight only",
	"FlagFaceTitle":       	"Face in photo",
	"FlagNightTitle":      	"Night only",
	"FlagRestrictedTitle": 	"Restricted access",
	"FlagReceiptTitle":		"Need a receipt (ticket)",
	"CloseEBCUndecided":	"Leave undecided",
	"CloseEBC":				["Accept good claim","No photo","Wrong/unclear photo","Out of hours/disallowed","Face not in photo","Bike not in photo","Flag not in photo","Missing rider/pillion","Missing receipt","Claim excluded" ],
	"ImgBonusFolder":		"images/bonuses/",
	"ImgEbcFolder":			"images/ebcimg/",
	"RallyBookImgTitle":	"Rally book photo",
	"EBCImgTitle":			"Entrant's image - click to resize",
	"EBCImgSwapTitle":		"Click to view this image",
	"Rally":				{"A1":"AAAAAAAAAAAAAA","A2":"22222222222222"}
}`

const debugDefaults = `{
	"RallyUseQA":			false,
	"RallyUsePctPen":		false,
	"RallyPctPenVal":		10,
	"RallyRankEfficiency":	false

}`

var tzlist = []string{
	"Europe/Amsterdam",
	"Europe/Andorra",
	"Europe/Athens",
	"Europe/Belgrade",
	"Europe/Berlin",
	"Europe/Brussels",
	"Europe/Bucharest",
	"Europe/Copenhagen",
	"Europe/Dublin",
	"Europe/Gibraltar",
	"Europe/Helsinki",
	"Europe/Istanbul",
	"Europe/Kyiv",
	"Europe/Lisbon",
	"Europe/London",
	"Europe/Madrid",
	"Europe/Paris",
	"Europe/Prague",
	"Europe/Riga",
	"Europe/Rome",
	"Europe/Sofia",
	"Europe/Stockholm",
	"Europe/Tallinn",
	"Europe/Vienna",
	"Europe/Vilnius",
	"Europe/Warsaw",
	"Europe/Zurich",
}

func ajaxUpdateSettings(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	fmt.Printf("ajaxUS %v\n", r.Form)

	ok := "false"
	msg := fmt.Sprintf("[%v] not implemented yet", r.FormValue("ff"))
	switch r.FormValue("ff") {
	case "RallyTitle":
		CS.Basics.RallyTitle = r.FormValue("v")
		ok = "true"
		msg = "ok"

		stmt, err := DBH.Prepare("UPDATE rallyparams SET RallyTitle=?")
		checkerr(err)
		defer stmt.Close()
		_, err = stmt.Exec(CS.Basics.RallyTitle)
		checkerr(err)
	case "RallyStart":
		CS.Basics.RallyStarttime = r.FormValue("v")
		ok = "true"
		msg = "ok"
		stmt, err := DBH.Prepare("UPDATE rallyparams SET StartTime=?")
		checkerr(err)
		defer stmt.Close()
		_, err = stmt.Exec(CS.Basics.RallyStarttime)
		checkerr(err)
	case "RallyFinish":
		CS.Basics.RallyFinishtime = r.FormValue("v")
		ok = "true"
		msg = "ok"
		stmt, err := DBH.Prepare("UPDATE rallyparams SET FinishTime=?")
		checkerr(err)
		defer stmt.Close()
		_, err = stmt.Exec(CS.Basics.RallyFinishtime)
		checkerr(err)
	case "MaxHours":
		CS.Basics.RallyMaxHours = intval(r.FormValue("v"))
		ok = "true"
		msg = "ok"
		stmt, err := DBH.Prepare("UPDATE rallyparams SET MaxHours=?")
		checkerr(err)
		defer stmt.Close()
		_, err = stmt.Exec(CS.Basics.RallyMaxHours)
		checkerr(err)
	case "StartOption":
		rs := r.FormValue("v")
		ok = "true"
		msg = "ok"
		stmt, err := DBH.Prepare("UPDATE rallyparams SET StartOption=?")
		checkerr(err)
		defer stmt.Close()
		_, err = stmt.Exec(rs)
		checkerr(err)
	case "LocalTZ":
		CS.Basics.RallyTimezone = r.FormValue("v")
		ok = "true"
		msg = "ok"
		stmt, err := DBH.Prepare("UPDATE rallyparams SET LocalTZ=?")
		checkerr(err)
		defer stmt.Close()
		_, err = stmt.Exec(CS.Basics.RallyTimezone)
		checkerr(err)

		RallyTimezone, err = time.LoadLocation(CS.Basics.RallyTimezone)
		checkerr(err)
	case "MilesKms":
		rs := r.FormValue("v")
		CS.Basics.RallyUnitKms = rs == "1"
		ok = "true"
		msg = "ok"
		stmt, err := DBH.Prepare("UPDATE rallyparams SET MilesKms=?")
		checkerr(err)
		defer stmt.Close()
		_, err = stmt.Exec(rs)
		checkerr(err)

	}

	fmt.Fprintf(w, `{"ok":%v,"msg":"%v"}`, ok, msg)

}
func editConfigMain(w http.ResponseWriter, r *http.Request) {

	var selected string

	startHTML(w, "Rally configuration")
	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<article class="config">`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="RallyTitle">Rally title</label>`)
	fmt.Fprintf(w, `<input type="text" class="RallyTitle" name="RallyTitle" id="RallyTitle" oninput="oi(this)" data-save="saveSetupConfig" value="%v">`, CS.Basics.RallyTitle)
	fmt.Fprint(w, `</fieldset>`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="RallyStartDate">Rally starts</label>`)

	dt, tm := splitDateTime(getStringFromDB("SELECT StartTime FROM rallyparams", "2000-01-01T08:00"))
	fmt.Fprintf(w, `<input type="date" name="RallyStartDate" id="RallyStartDate" onchange="saveSetupStart(this)" value="%v">`, dt)
	fmt.Fprintf(w, ` <input type="time" name="RallyStartTime" id="RallyStartTime" onchange="saveSetupStart(this)" value="%v">`, tm)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="RallyFinishDate">Rally finishes</label>`)

	dt, tm = splitDateTime(getStringFromDB("SELECT FinishTime FROM rallyparams", "2000-01-01T08:00"))
	fmt.Fprintf(w, `<input type="date" name="RallyFinishDate" id="RallyFinishDate" onchange="saveSetupFinish(this)" value="%v">`, dt)
	fmt.Fprintf(w, ` <input type="time" name="RallyFinishTime" id="RallyFinishTime" onchange="saveSetupFinish(this)" value="%v">`, tm)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="MaxHours">Max Rideable hours</label>`)
	mh := getIntegerFromDB("SELECT MaxHours FROM rallyparams", 99)
	fmt.Fprintf(w, ` <input type="number" id="MaxHours" name="MaxHours" class="MaxHours" oninput="oi(this)" data-save="saveSetupConfig" value="%v">`, mh)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="StartOption">Rally Start option</label>`)
	so := getIntegerFromDB("SELECT StartOption FROM rallyparams", 0)
	fmt.Fprint(w, ` <select id="StartOption" name="StartOption" onchange="saveSetupConfig(this)">`)
	selected = ""
	if so != 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="0" %v>Fixed start time</option>`, selected)
	selected = ""
	if so == 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="1" %v>Start by first claim</option>`, selected)
	fmt.Fprint(w, `,</select>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="MilesKms">Unit of distance</label>`)
	mk := getIntegerFromDB("SELECT MilesKms FROM rallyparams", 0)
	fmt.Fprint(w, ` <select id="MilesKms" name="MilesKms" onchange="saveSetupConfig(this)">`)
	selected = ""
	if mk != 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="0" %v>Miles</option>`, selected)
	selected = ""
	if mk == 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="1" %v>Kilometres</option>`, selected)
	fmt.Fprint(w, `,</select>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="LocalTZ">Rally Timezone</label>`)
	rtz := getStringFromDB("SELECT LocalTZ FROM rallyparams", "Europe/London")
	fmt.Fprint(w, ` <select id="LocalTZ" name="LocalTZ" onchange="saveSetupConfig(this)">`)
	selected = ""
	for _, itz := range tzlist {
		selected = ""
		if rtz == itz {
			selected = "selected"
		}
		fmt.Fprintf(w, `<option value="%v" %v>%v</option>`, itz, selected, itz)
	}
	fmt.Fprint(w, `,</select>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `</article>`)
}

func loadRallyBasics(rb *RallyBasics) {

	rb.RallyTitle = getStringFromDB("SELECT RallyTitle FROM rallyparams", "Some Rally")
	rb.RallyStarttime = getStringFromDB("SELECT StartTime FROM rallyparams", "2000-01-01T08:00")
	rb.RallyFinishtime = getStringFromDB("SELECT FinishTime FROM rallyparams", "2000-01-01T18:00")
	rb.RallyTimezone = getStringFromDB("SELECT LocalTZ FROM rallyparams", "Europe/London")
	rb.RallyMaxHours = getIntegerFromDB("SELECT MaxHours FROM rallyparams", 99)
	rb.RallyUnitKms = getIntegerFromDB("SELECT MilesKms FROM rallyparams", 0) == 1
}

func splitDateTime(iso string) (string, string) {

	b4, af, ok := strings.Cut(iso, "T")
	if ok {
		return b4, af
	}
	return iso, ""
}
