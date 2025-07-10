package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"text/template"
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

const (
	CheckoutStart = iota
	FirstClaimStart
)

type chasmSettings struct {
	StartOption         int
	AutoFinisher        bool
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
	UploadsFolder       string
}

var CS chasmSettings

const defaultCS = `{
	"StartOption": 			0,
	"AutoFinisher":			false,
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
	"Rally":				{"A1":"AAAAAAAAAAAAAA","A2":"22222222222222"},
	"UploadsFolder":		"uploads"
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

// SQL for safely retrieving full names
const RiderNameSQL = "ifnull(entrants.RiderName,ifnull(entrants.RiderFirst,'') || ' ' || ifnull(entrants.RiderLast,'')) AS RiderName"
const PillionNameSQL = "ifnull(entrants.PillionName,ifnull(entrants.PillionFirst,'') || ' ' || ifnull(entrants.PillionLast,'')) AS PillionName"

var configLiteralsTemplate = `

<article class="config literals">
<fieldset>
	<button onclick="swapconfig(this)">LITERALS</button>
	<fieldset class="hide">
		<legend>Distance</legend>
		<input type="text" id="UnitMilesLit" name="UnitMilesLit" value="{{.UnitMilesLit}}" oninput="oi(this)" data-save="saveSetupConfig">
		<input type="text" id="UnitKmsLit" name="UnitKmsLit" value="{{.UnitKmsLit}}" oninput="oi(this)" data-save="saveSetupConfig" onchange="saveSetupConfig(this)" onblur="saveSetupConfig(this)">
	</fieldset>
	<fieldset class="hide">
		<legend>EBC Decisions</legend>
		{{range $ix,$el := .CloseEBC}}
			<input type="text" name="CloseEBC[{{$ix}}]" id="CloseEBC[{{$ix}}]" value="{{$el}}" oninput="oi(this)" data-save="saveSetupConfig" onchange="saveSetupConfig(this)" onblur="saveSetupConfig(this)">
		{{end}}
	</fieldset>
	<fieldset class="hide">
		<legend>Flag titles<legend>
			##FLAGS##
	</fieldset>
</fieldset>
</article>
`

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
		CS.StartOption = intval(rs)
	case "FinishOption":
		rs := r.FormValue("v")
		ok = "true"
		msg = "ok"
		CS.AutoFinisher = rs == "1"
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
	default:
		fn := r.FormValue("ff")
		fv := r.FormValue("v")
		if strings.Contains(fn, "CloseEBC") {
			index := intval(fn[9:10])
			CS.CloseEBC[index] = fv
			ok = "true"
			msg = "ok"
			break
		}
		ps := reflect.ValueOf(&CS)
		s := ps.Elem()
		f := s.FieldByName(fn)
		if !f.IsValid() {
			break
		}
		if !f.CanSet() {
			break
		}
		switch f.Kind() {
		case reflect.String:
			f.SetString(fv)
		case reflect.Int:
			f.SetInt(int64(intval(fv)))
		case reflect.Bool:
			f.SetBool(fv == "1")
		default:
			break
		}
		ok = "true"
		msg = "ok"

	}
	saveSettings()
	fmt.Fprintf(w, `{"ok":%v,"msg":"%v"}`, ok, msg)

}

func buildFlagTitles() string {

	x := emitConfigText("FlagTeamTitle")
	x += emitConfigText("FlagAlertTitle")
	x += emitConfigText("FlagBikeTitle")
	x += emitConfigText("FlagDaylightTitle")
	x += emitConfigText("FlagFaceTitle")
	x += emitConfigText("FlagNightTitle")
	x += emitConfigText("FlagRestrictedTitle")
	x += emitConfigText("FlagReceiptTitle")

	return x
}

func buildRallyVarsSettings() string {

	du := "mile"
	dus := CS.UnitMilesLit
	if CS.Basics.RallyUnitKms {
		du = "km"
		dus = CS.UnitKmsLit
	}

	x := `<article class="config RallyVars">`

	x += `<fieldset><button onclick="swapconfig(this)">OPTIONS</button>`

	x += `<fieldset class="hide"><label for="PenaltyMilesDNF">`
	x += `DNF if ` + dus + ` &gt;</label>`
	x += emitConfigNum("PenaltyMilesDNF") + `</fieldset>`

	x += `<fieldset class="hide"><label for="RallyMinMiles">`
	x += `DNF if ` + dus + ` &lt;</label>`
	x += emitConfigNum("RallyMinMiles") + `</fieldset>`

	x += `<fieldset class="hide"><label for="RallyMinPoints">`
	x += `DNF if points &lt;</label>`
	x += emitConfigNum("RallyMinPoints") + `</fieldset>`

	x += `<fieldset class="hide"><label for="RallyUseQA">`
	x += `Use questions/answers</label>`
	x += emitConfigBool("RallyUseQA", []string{"no", "yes"}, CS.RallyUseQA)

	x += `<label for="RallyQAPoints">`
	x += `QA points value</label>`
	x += emitConfigNum("RallyQAPoints") + `</fieldset>`

	x += `<fieldset class="hide"><label for="RallyUsePctPen">`
	x += `Offer minor points reduction</label>`
	x += emitConfigBool("RallyUsePctPen", []string{"no", "yes"}, CS.RallyUsePctPen)

	x += `<label for="RallyPctPenVal">`
	x += `Points reduction percentage</label>`
	x += emitConfigNum("RallyPctPenVal") + `</fieldset>`

	x += `<fieldset class="hide"><label for="RallyRankEfficiency">`
	x += `Rank finishers by</label>`
	x += emitConfigBool("RallyRankEfficiency", []string{"total points", "points per " + du}, CS.RallyRankEfficiency) + `</fieldset>`

	x += `<fieldset class="hide"><label for="RallySplitTies">`
	x += `Split ties</label>`
	x += emitConfigBool("RallySplitTies", []string{"leave as tied", "prefer shorter distance"}, CS.RallySplitTies) + `</fieldset>`

	x += `<fieldset class="hide"><label for="RallyTeamMethod">`
	x += `Team ranking</label>`
	x += emitConfigSelect("RallyTeamMethod", []string{"individual placing", "highest ranked member", "lowest ranked member", "clone team member scores"}, CS.RallyTeamMethod) + `</fieldset>`

	x += `</fieldset></article>`
	return x

}

func emitConfigNum(varName string) string {

	x := `<input type="number" class="` + varName + `" id="` + varName + `" name="` + varName + `" `
	x += `value="{{.` + varName + `}}" oninput="oi(this);this.setAttribute('data-chg',1)" data-save="saveSetupConfig" onchange="saveSetupConfig(this)" `
	x += `onblur1="s1aveSetupConfig(this)">`
	return x
}
func emitConfigText(varName string) string {

	x := `<input type="text" class="` + varName + `" id="` + varName + `" name="` + varName + `" placeholder="{{.` + varName + `}}" `
	x += `value="{{.` + varName + `}}" oninput="oi(this);this.setAttribute('data-chg',1)" data-save="saveSetupConfig" onchange="saveSetupConfig(this)" `
	x += `onblur="saveSetupConfig(this)">`
	return x
}

func emitConfigBool(varName string, varOptions []string, varBool bool) string {

	x := `<select id="` + varName + `" name="` + varName + `" `
	varIx := 0
	if varBool {
		varIx = 1
	}
	x += `onchange=";this.setAttribute('data-chg',1);saveSetupConfig(this)">`
	for i, o := range varOptions {
		x += `<option `
		if i == varIx {
			x += ` selected `
		}
		x += fmt.Sprintf(` value="%v">%v</option>`, i, o)
	}
	x += `</select>`
	return x
}

func emitConfigSelect(varName string, varOptions []string, varIx int) string {

	x := `<select id="` + varName + `" name="` + varName + `" `
	x += `onchange=";this.setAttribute('data-chg',1);saveSetupConfig(this)">`
	for i, o := range varOptions {
		x += `<option `
		if i == varIx {
			x += ` selected `
		}
		x += fmt.Sprintf(` value="%v">%v</option>`, i, o)
	}
	x += `</select>`
	return x
}

func editConfigMain(w http.ResponseWriter, r *http.Request) {

	var selected string

	startHTML(w, "Rally configuration")
	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<article class="config basic">`)
	fmt.Fprint(w, `<fieldset><button  onclick="swapconfig(this)">BASIC</button>`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="RallyTitle">Rally title</label>`)
	fmt.Fprintf(w, `<input type="text" autofocus class="RallyTitle" name="RallyTitle" id="RallyTitle" oninput="oi(this)" data-save="saveSetupConfig" value="%v">`, CS.Basics.RallyTitle)
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
	so := CS.StartOption //getIntegerFromDB("SELECT StartOption FROM rallyparams", 0)
	fmt.Fprint(w, ` <select id="StartOption" name="StartOption" onchange="saveSetupConfig(this)">`)
	selected = ""
	if so != 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="0" %v>Fixed with check-out</option>`, selected)
	selected = ""
	if so == 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="1" %v>Start by first claim</option>`, selected)
	fmt.Fprint(w, `,</select>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="FinishOption">Rally Finish option</label>`)
	fmt.Fprint(w, `<select id="FinishOption" name="FinishOption" onchange="saveSetupConfig(this)">`)
	selected = ""
	if !CS.AutoFinisher {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="0" %v>Finish at check-in</option>`, selected)
	selected = ""
	if CS.AutoFinisher {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="1" %v>Autofinish with last claim</option>`, selected)

	fmt.Fprint(w, `</select>`)
	fmt.Fprint(w, `</fieldset>`)

	//fmt.Fprint(w, `</fieldset></article>`) // basic

	//fmt.Fprint(w, `<article class="config regional"><fieldset><legend>REGIONAL</legend>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="MilesKms">Unit of distance</label>`)
	mk := getIntegerFromDB("SELECT MilesKms FROM rallyparams", 0)
	fmt.Fprint(w, ` <select id="MilesKms" name="MilesKms" onchange="saveSetupConfig(this)">`)
	selected = ""
	if mk != 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="0" %v>%v</option>`, selected, CS.UnitMilesLit)
	selected = ""
	if mk == 1 {
		selected = "selected"
	}
	fmt.Fprintf(w, `<option value="1" %v>%v</option>`, selected, CS.UnitKmsLit)
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

	fmt.Fprint(w, `</fieldset></article>`) // basic

	t, err := template.New("rally vars").Parse(buildRallyVarsSettings())
	checkerr(err)
	err = t.Execute(w, CS)
	checkerr(err)

	t, err = template.New("literals").Parse(strings.ReplaceAll(configLiteralsTemplate, "##FLAGS##", buildFlagTitles()))
	checkerr(err)
	err = t.Execute(w, CS)
	checkerr(err)

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

func saveSettings() {

	csb, err := json.Marshal(CS)
	checkerr(err)
	stmt, err := DBH.Prepare("UPDATE config SET Settings=?")
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(csb)
	checkerr(err)

}
