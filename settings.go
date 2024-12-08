package main

import (
	"fmt"
	"net/http"
	"strings"
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
type chasmSettings struct {
	ShowExcludedClaims  bool // If a claim is marked 'excluded' and is not superseded, show it on the scoresheet
	CurrentLeg          int
	UseCheckinForOdo    bool // If true, OdoRallyFinish updated only by check-in, not by individual claims
	RallyUnitKms        bool // Report in Kms(true) or Miles(false)
	UnitMilesLit        string
	UnitKmsLit          string
	PenaltyMilesDNF     int
	RallyMinMiles       int
	DebugRules          bool
	AutoLateDNF         bool
	RallyTitle          string
	RallyMinPoints      int
	RallyTimezone       string
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
	"ShowExcludedClaims": 	true,
	"CurrentLeg": 			5,
	"UseCheckInForOdo": 	false,
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
	"RallyUseQA":			true,
	"RallyUsePctPen":		true,
	"RallyPctPenVal":		10,
	"RallyRankEfficiency":	false

}`

func ajaxUpdateSettings(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	fmt.Printf("ajaxUS %v\n", r.Form)

	if r.FormValue("tab") == "1" {
		if r.FormValue("RallyTitle") != "" {
			CS.RallyTitle = r.FormValue("RallyTitle")
		}
		stmt, err := DBH.Prepare("UPDATE rallyparams SET RallyTitle=?")
		checkerr(err)
		defer stmt.Close()
	}
}
func editConfigMain(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Rally configuration")
	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<article class="config">`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="RallyTitle">Rally title</label>`)
	fmt.Fprintf(w, `<input type="text" class="RallyTitle" name="RallyTitle" id="RallyTitle" value="%v">`, CS.RallyTitle)
	fmt.Fprint(w, `</fieldset>`)
	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="RallyStartDate">Rally starts</label>`)

	dt, tm := splitDateTime(getStringFromDB("SELECT StartTime FROM rallyparams", "2000-01-01T08:00"))
	fmt.Fprintf(w, `<input type="date" name="RallyStartDate" id="RallyStartDate" value="%v">`, dt)
	fmt.Fprintf(w, ` <input type="time" name="RallyStartTime" value="%v">`, tm)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset>`)
	fmt.Fprint(w, `<label for="RallyFinishDate">Rally finishes</label>`)

	dt, tm = splitDateTime(getStringFromDB("SELECT FinishTime FROM rallyparams", "2000-01-01T08:00"))
	fmt.Fprintf(w, `<input type="date" name="RallyFinishDate" id="RallyFinishDate" value="%v">`, dt)
	fmt.Fprintf(w, ` <input type="time" name="RallyFinishTime" value="%v">`, tm)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `</article>`)
}

func splitDateTime(iso string) (string, string) {

	b4, af, ok := strings.Cut(iso, "T")
	if ok {
		return b4, af
	}
	return iso, ""
}
