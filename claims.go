package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed images/alertteam.b64
var flag_team string

//go:embed images/alertalert.b64
var flag_alert string

//go:embed images/alertbike.b64
var flag_bike string

//go:embed images/alertdaylight.b64
var flag_daylight string

//go:embed images/alertface.b64
var flag_face string

//go:embed images/alertnight.b64
var flag_night string

//go:embed images/alertrestricted.b64
var flag_restricted string

//go:embed images/alertreceipt.b64
var flag_receipt string

type ElectronicBonusClaim struct {
	Claimid        int
	EntrantID      int
	RiderName      string
	PillionName    string
	Bonusid        string
	BriefDesc      string
	OdoReading     int
	ClaimTime      string
	Decision       int
	Subject        string
	ExtraField     string
	AttachmentTime string
	DateTime       string
	FirstTime      string
	FinalTime      string
	EmailID        int
}

type BonusClaimVars struct {
	BriefDesc      string
	Points         int
	Notes          string
	Flags          string
	AskPoints      bool
	PointsAreMults bool
	RestMins       int
	AskMins        bool
	Image          string
	Question       string
	Answer         string
	Leg            int
}

type ClaimRecord struct {
	LoggedAt         string
	ClaimTime        string
	EntrantID        int
	BonusID          string
	OdoReading       int
	Decision         int
	Photo            string
	Points           int
	RestMinutes      int
	AskPoints        bool
	AskMinutes       bool
	QuestionAsked    bool
	QuestionAnswered bool
	AnswerSupplied   string
	JudgesNotes      string
	PercentPenalty   bool
	Evidence         string
	Leg              int
}

type EntrantDetails struct {
	EntrantID         int
	RiderName         string
	PillionName       string
	TeamID            int
	RiderFirst        string
	RiderLast         string
	PillionFirst      string
	PillionLast       string
	ReviewedByTeam    int
	AcceptedByEntrant int
	LastReviewed      string
}

const maximg = 3

//var EntrantSelector map[int]string

func emitImage(img string, alt string, title string) string {

	res := fmt.Sprintf(`<img alt="%v", title="%v" class="flagicon" src="data:image/png;base64,`, alt, title)
	for _, xl := range strings.Split(img, "\n") {
		res += xl
	}
	res += `">`
	return res

}

func fetchBonusVars(b string) BonusClaimVars {

	var res BonusClaimVars
	var ap int
	var am int

	sqlx := "SELECT ifnull(BriefDesc,BonusID),Points,ifnull(Notes,''),ifnull(Flags,''),AskPoints,RestMinutes,AskMinutes,ifnull(Image,''),ifnull(Question,''),ifnull(Answer,'')"
	sqlx += " FROM bonuses WHERE BonusID='" + b + "'"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if !rows.Next() {
		res.BriefDesc = b
		return res
	}
	err = rows.Scan(&res.BriefDesc, &res.Points, &res.Notes, &res.Flags, &ap, &res.RestMins, &am, &res.Image, &res.Question, &res.Answer)
	checkerr(err)
	res.AskMins = am != 0
	res.AskPoints = ap == 1
	res.PointsAreMults = ap == 2
	return res
}

func fetchClaimDetails(claimid int) ClaimRecord {

	var cr ClaimRecord
	sqlx := "SELECT ifnull(LoggedAt,ClaimTime),ClaimTime,EntrantID,BonusID,OdoReading,Decision,ifnull(Photo,'')"
	sqlx += ",Points,RestMinutes,AskPoints,AskMinutes,QuestionAsked,QuestionAnswered,ifnull(AnswerSupplied,'')"
	sqlx += ",ifnull(JudgesNotes,''),PercentPenalty,ifnull(Evidence,''),Leg"
	sqlx += " FROM claims WHERE rowid=" + strconv.Itoa(claimid)

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if !rows.Next() {
		return cr
	}

	ap := 0
	am := 0
	qak := 0
	qaw := 0
	pp := 0

	err = rows.Scan(&cr.LoggedAt, &cr.ClaimTime, &cr.EntrantID, &cr.BonusID, &cr.OdoReading, &cr.Decision, &cr.Photo, &cr.Points, &cr.RestMinutes, &ap, &am, &qak, &qaw, &cr.AnswerSupplied, &cr.JudgesNotes, &pp, &cr.Evidence, &cr.Leg)
	checkerr(err)
	cr.AskPoints = ap != 0
	cr.AskMinutes = am != 0
	cr.QuestionAsked = qak != 0
	cr.QuestionAnswered = qaw != 0
	cr.PercentPenalty = pp != 0

	return cr

}

func fetchEntrantDetails(entrant int) EntrantDetails {

	var ed EntrantDetails

	ed.EntrantID = entrant
	if entrant < 1 {
		return ed
	}

	e := strconv.Itoa(entrant)

	ed.RiderName = getStringFromDB("SELECT "+RiderNameSQL+" FROM entrants WHERE EntrantID="+e, e)
	ed.PillionName = getStringFromDB("SELECT ifnull(PillionName,'') FROM entrants WHERE EntrantID="+e, "")
	ed.TeamID = getIntegerFromDB("SELECT TeamID FROM entrants WHERE EntrantID="+e, 0)
	return ed

}

func list_claims(w http.ResponseWriter, r *http.Request) {

	const addnew_icon = "&nbsp;+&nbsp;"

	const tick_icon = "&#10004;"
	const cross_icon = "&#10006;"
	const undecided_icon = "?"

	const filter_icon = "" //"&#65509;"

	r.ParseForm()

	EntrantSelector := loadEntrantsList()

	esel := intval(r.FormValue("esel"))

	startHTML(w, "Claims log")

	fmt.Fprint(w, `<div class="claimslog">`)

	showReloadTicker(w, r.URL.String())
	fmt.Fprint(w, `<h4>Claims log</h4>`)

	fmt.Fprint(w, `<form id="claimslogfrm">`)
	sel := ""
	if esel == 0 {
		sel = "selected"
	}
	fmt.Fprintf(w, `<div class="select"">%v `, filter_icon)
	fmt.Fprintf(w, `<button autofocus title="Add new claim" class="plus" onclick="window.location.href='/claim?c=0';return false">%v</button> <span id="fcc"></span>`, addnew_icon)
	fmt.Fprintf(w, ` <select name="esel" value="%v" onchange="reloadClaimslog()">`, esel)
	fmt.Fprintf(w, `<option value="0" %v>all claims</option>`, sel)

	keys := make([]int, 0, len(EntrantSelector))

	for key := range EntrantSelector {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return EntrantSelector[keys[i]] < EntrantSelector[keys[j]]
	})

	for _, k := range keys {
		sel = ""
		if k == esel {
			sel = "selected"
		}
		fmt.Fprintf(w, `<option value="%v" %v>%v</option>`, k, sel, EntrantSelector[k])
	}
	fmt.Fprint(w, `</select>`)

	bsel := r.FormValue("bsel")
	sel = ""
	if bsel == "" {
		sel = "selected"
	}
	fmt.Fprintf(w, ` <select name="bsel" value="%v" onchange="reloadClaimslog()">`, bsel)
	fmt.Fprintf(w, `<option value="" %v>all claims</option>`, sel)
	sqlx := "SELECT BonusID FROM bonuses ORDER BY BonusID"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var b string
		sel = ""
		err = rows.Scan(&b)
		checkerr(err)
		if b == bsel {
			sel = "selected"
		}
		fmt.Fprintf(w, `<option value="%v" %v>%v</option>`, b, sel, b)
	}
	rows.Close()
	fmt.Fprint(w, `</select>`)

	dsel := r.FormValue("dsel")
	sel = ""
	if dsel == "" {
		sel = "selected"
	}
	fmt.Fprintf(w, ` <select name="dsel" value="%v" onchange="reloadClaimslog()">`, dsel)
	fmt.Fprintf(w, `<option value="" %v>all claims</option>`, sel)

	sel = ""
	if dsel == "g" {
		sel = "selected"
	}
	fmt.Fprintf(w, `<option value="g" %v>Good claims only</option>`, sel)
	sel = ""
	if dsel == "r" {
		sel = "selected"
	}
	fmt.Fprintf(w, `<option value="r" %v>Rejected claims only</option>`, sel)
	sel = ""
	if dsel == "u" {
		sel = "selected"
	}
	fmt.Fprintf(w, `<option value="u" %v>Undecided claims only</option>`, sel)
	fmt.Fprint(w, `</select>`)

	fmt.Fprint(w, `</form>`)

	fmt.Fprint(w, `</div>`)
	sqlx = `SELECT claims.rowid,ifnull(LoggedAt,''),ClaimTime,claims.EntrantID,BonusID,OdoReading,Decision,ifnull(JudgesNotes,'') FROM claims`
	sqlx += " LEFT JOIN entrants ON claims.EntrantID=entrants.EntrantID"
	where := ""
	if esel > 0 {
		where += "  claims.EntrantID=" + strconv.Itoa(esel)
	}
	if bsel != "" {
		if where != "" {
			where += " AND "
		}
		where += "BonusID='" + bsel + "'"
	}
	if dsel != "" {
		if where != "" {
			where += " AND "
		}
		where += " Decision "
		switch dsel {
		case "u":
			where += "< 0"
		case "g":
			where += "= 0"
		default:
			where += "> 0"
		}
	}
	if where != "" {
		sqlx += " WHERE " + where
	}
	sqlx += " ORDER BY ClaimTime DESC"

	rows, err = DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	fmt.Fprint(w, `<fieldset class="row claims hdr">`)
	fmt.Fprint(w, `<fieldset class="col claims hdr">Entrant</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col claims hdr">Bonus</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col claims hdr">Odo</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col claims hdr">Time</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col claims hdr mid">Good?</fieldset>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `</div><hr>`)

	fmt.Fprint(w, `</header><div class="claimslog">`)
	for rows.Next() {
		var cr ClaimRecord
		var claimid int
		err = rows.Scan(&claimid, &cr.LoggedAt, &cr.ClaimTime, &cr.EntrantID, &cr.BonusID, &cr.OdoReading, &cr.Decision, &cr.JudgesNotes)
		checkerr(err)

		rname, ok := EntrantSelector[cr.EntrantID]
		if !ok {
			rname = strconv.Itoa(cr.EntrantID)
		}
		fmt.Fprintf(w, `<fieldset class="row claims" onclick="window.location.href='/claim?c=%v'">`, claimid)
		fmt.Fprintf(w, `<fieldset class="col claims" title="%v">%v</fieldset>`, cr.EntrantID, rname)
		fmt.Fprintf(w, `<fieldset class="col claims">%v</fieldset>`, cr.BonusID)
		fmt.Fprintf(w, `<fieldset class="col claims">%v</fieldset>`, cr.OdoReading)
		fmt.Fprintf(w, `<fieldset class="col claims">%v</fieldset>`, logtime(cr.ClaimTime))
		decision := tick_icon
		if cr.Decision > 0 {
			decision = cross_icon
		} else if cr.Decision < 0 {
			decision = undecided_icon
		}
		fmt.Fprintf(w, `<fieldset class="col claims mid">%v</fieldset>`, decision)
		fmt.Fprint(w, `</fieldset>`)
	}

	fmt.Fprint(w, `</div>`)

}

// Show judgeable claims submitted electronically
func list_EBC_claims(w http.ResponseWriter, r *http.Request) {

	const sorry = "Sorry, no claims need judging at the moment &#128543;"

	sqlx := `SELECT ebclaims.rowid,ebclaims.EntrantID,` + RiderNameSQL + `,ifnull(entrants.PillionName,''),ebclaims.BonusID,xbonus.BriefDesc,ebclaims.OdoReading,ebclaims.ClaimTime
	 		FROM ebclaims LEFT JOIN entrants ON ebclaims.EntrantID=entrants.EntrantID
			LEFT JOIN (SELECT BonusID,BriefDesc FROM bonuses) AS xbonus ON ebclaims.BonusID=xbonus.BonusID
			 WHERE Processed=0 ORDER BY Decision DESC,FinalTime;`

	rows, err := DBH.Query(sqlx)
	checkerr(err)

	startHTML(w, "Process EBC claims")

	fmt.Fprint(w, `<div class="ebclist">`)

	showReloadTicker(w, r.URL.String())
	fmt.Fprint(w, `<h4>Emailed claims ready to be judged</h4>`)

	fmt.Fprint(w, `<div id="judgefc">`)
	fmt.Fprintf(w, `<button autofocus onclick="showFirstClaim()">Judge first claim</button> <span id="fcc"></span>`)

	fmt.Fprint(w, `<fieldset class="row ebc hdr">`)
	fmt.Fprint(w, `<fieldset class="col ebc hdr">Entrant</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col ebc hdr">Bonus</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col ebc hdr">Odo</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col ebc hdr">Time</fieldset>`)
	fmt.Fprint(w, `</fieldset>`)
	fmt.Fprint(w, `</div>`)

	fmt.Fprint(w, `</div><hr></header>`)
	fmt.Fprint(w, `<div class="ebclist">`)
	n := 0
	for rows.Next() {
		var ebc ElectronicBonusClaim
		err := rows.Scan(&ebc.Claimid, &ebc.EntrantID, &ebc.RiderName, &ebc.PillionName, &ebc.Bonusid, &ebc.BriefDesc, &ebc.OdoReading, &ebc.ClaimTime)
		checkerr(err)
		n++
		fmt.Fprintf(w, `<fieldset class="row ebc" data-claimid="%v" onclick="showEBC(this)">`, ebc.Claimid)
		team := ebc.RiderName
		if ebc.PillionName != "" {
			team += " &amp; " + ebc.PillionName
		}
		fmt.Fprintf(w, `<fieldset class="col ebc" title="%v">%v</fieldset>`, ebc.EntrantID, team)
		fmt.Fprintf(w, `<fieldset class="col ebc" title="%v">%v</fieldset>`, ebc.BriefDesc, ebc.Bonusid)
		fmt.Fprintf(w, `<fieldset class="col ebc">%v</fieldset>`, ebc.OdoReading)
		fmt.Fprintf(w, `<fieldset class="col ebc">%v</fieldset>`, logtime(ebc.ClaimTime))
		fmt.Fprint(w, `</fieldset>`)
	}
	fmt.Fprint(w, `</div>`)
	fmt.Fprintf(w, `<script>if (%v>0){let x = document.getElementById('fcc');x.innerHTML='1/%v';}else{let x=document.getElementById('judgefc');x.innerHTML='%v';}</script>`, n, n, sorry)

}

func loadEntrantsList() map[int]string {

	res := make(map[int]string)
	sqlx := "SELECT EntrantID," + RiderNameSQL + ",ifnull(PillionName,'') FROM entrants"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var e int
		var rn string
		var pn string
		err = rows.Scan(&e, &rn, &pn)
		checkerr(err)
		if pn != "" {
			rn += " &amp; " + pn
		}
		res[e] = rn
	}
	return res
}

func logtime(stamp string) string {
	/* We're really only interested in the time of day and which of a few days it's on */

	const showformat = "Mon 15:04"
	ts := parseStoredDate(stamp)
	return fmt.Sprintf(`<span title="%v">%v</span>`, stamp, ts.Format(showformat))
}

func ImgFromURL(url string) string {

	if len(url) == 0 {
		return url
	}
	res := strings.Split(url, "/")
	ix := len(res) - 1
	if ix < 0 {
		ix = 0
	}
	return res[ix]
}

func insertNewClaim(r *http.Request) {

	const Leg = 1

	sqlx := "INSERT INTO claims (LoggedAt, ClaimTime, EntrantID, BonusID, OdoReading, Decision, Photo, Points, RestMinutes, Leg"
	sqlx += ", AnswerSupplied, QuestionAnswered, JudgesNotes, PercentPenalty) "
	sqlx += "VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	//r.FormValue("AnswerSupplied"), r.FormValue("QuestionAnswered"),
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().Format(time.RFC3339), r.FormValue("ClaimTime"), r.FormValue("EntrantID"), strings.ToUpper(r.FormValue("BonusID")), intval(r.FormValue("OdoReading")), intval(r.FormValue("Decision")), filepath.Base(r.FormValue("Photo")), intval(r.FormValue("Points")), intval(r.FormValue(("RestMinutes"))), Leg, r.FormValue("AnswerSupplied"), r.FormValue("QuestionAnswered"), r.FormValue("JudgesNotes"), intval(r.FormValue("PercentPenalty")))
	checkerr(err)

}
func saveClaim(r *http.Request) {

	//fmt.Printf("saveclaim: %v\n", r)
	claimid := intval(r.FormValue("claimid"))

	if claimid < 1 {
		insertNewClaim(r)
		recalc_scorecard(intval(r.FormValue("EntrantID")))
		return
	}
	sqlx := "UPDATE claims SET ClaimTime=?,EntrantID=?,BonusID=?,OdoReading=?,AnswerSupplied=?,QuestionAnswered=?,Points=?,RestMinutes=?,Decision=?,JudgesNotes=?"
	sqlx += ",PercentPenalty=?"
	sqlx += " WHERE rowid=?"

	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(r.FormValue("ClaimTime"), r.FormValue("EntrantID"), strings.ToUpper(r.FormValue("BonusID")), intval(r.FormValue("OdoReading")), r.FormValue("AnswerSupplied"), r.FormValue("QuestionAnswered"), intval(r.FormValue("Points")), intval(r.FormValue(("RestMinutes"))), intval(r.FormValue("Decision")), r.FormValue("JudgesNotes"), intval(r.FormValue("PercentPenalty")), claimid)
	checkerr(err)
	recalc_scorecard(intval(r.FormValue("EntrantID")))

}

func showClaim(w http.ResponseWriter, r *http.Request) {

	var cr ClaimRecord

	var claimdate string
	var claimtime string

	var ed EntrantDetails
	var bd BonusClaimVars

	r.ParseForm()

	startHTMLBL(w, "Individual claim", "claims")

	claimid := intval(r.FormValue("c"))

	fmt.Fprint(w, `</header><div class="claim">`)
	fmt.Fprint(w, `<form id="iclaim">`)

	fmt.Fprintf(w, `<input type="hidden" name="claimid" value="%v">`, claimid)
	if claimid < 1 {
		fmt.Fprint(w, `<h4>Filing new claim</h4>`)
		fmt.Fprint(w, `<input type="text" autofocus tabindex="1" class="subject" placeholder="Paste email Subject line here" oninput="pasteNewClaim(this)">`)
		cr.Decision = 0 // Good claim
		claimdate = time.Now().Format("2006-01-02")
		claimtime = time.Now().Format("15:04")
		cr.ClaimTime = claimdate + "T" + claimtime
	} else {
		cr = fetchClaimDetails((claimid))
		if len(cr.ClaimTime) > 12 {
			claimdate = cr.ClaimTime[0:10]
			claimtime = cr.ClaimTime[11:16]
		}
		ed = fetchEntrantDetails(cr.EntrantID)
		bd = fetchBonusVars(cr.BonusID)
		fmt.Fprint(w, `<h4>Updating claim details</h4>`)
	}
	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="EntrantID">Entrant</label>`)
	fmt.Fprint(w, `<input type="number" tabindex="2" id="EntrantID" name="EntrantID" class="EntrantID" oninput="fetchEntrantDetails(this)"`)
	if claimid > 0 {
		fmt.Fprintf(w, ` value="%v"`, cr.EntrantID)
	}
	fmt.Fprint(w, `>`)
	fmt.Fprint(w, `<span>`)
	fmt.Fprintf(w, ` <span id="entrantDetails">%v</span>`, ed.RiderName)
	hide := ""
	//fmt.Printf("Pillion='%v', Team=%v\n", ed.PillionName, ed.TeamID)
	if ed.PillionName == "" && ed.TeamID < 1 {
		hide = "hide"
	}
	fmt.Fprintf(w, ` <span id="edflag" class="%v">`, hide)
	fmt.Fprint(w, emitImage(flag_team, "TR", CS.FlagTeamTitle))
	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="BonusID">Bonus code</label>`)
	fmt.Fprintf(w, `<input type="text" tabindex="3" id="BonusID" name="BonusID" class="BonusID" oninput="fetchBonusDetails(this)" value="%v">`, cr.BonusID)
	fmt.Fprint(w, `<span>`)
	fmt.Fprintf(w, ` <span id="bonusDetails">%v</span>`, bd.BriefDesc)

	allflags := "ABDFNRT"

	//fmt.Printf("%v has [%v]\n", cr.BonusID, bd.Flags)

	for _, c := range allflags {
		hide := "hide"
		if strings.Contains(bd.Flags, string(c)) {
			hide = ""
		}
		switch c {
		case 'A':
			fmt.Fprintf(w, `<span id="bflagA" class="%v">%v</span>`, hide, emitImage(flag_alert, string(c), CS.FlagAlertTitle))
		case 'B':
			fmt.Fprintf(w, `<span id="bflagB" class="%v">%v</span>`, hide, emitImage(flag_bike, string(c), CS.FlagBikeTitle))
		case 'D':
			fmt.Fprintf(w, `<span id="bflagD" class="%v">%v</span>`, hide, emitImage(flag_daylight, string(c), CS.FlagDaylightTitle))
		case 'F':
			fmt.Fprintf(w, `<span id="bflagF" class="%v">%v</span>`, hide, emitImage(flag_face, string(c), CS.FlagFaceTitle))
		case 'N':
			fmt.Fprintf(w, `<span id="bflagN" class="%v">%v</span>`, hide, emitImage(flag_night, string(c), CS.FlagNightTitle))
		case 'R':
			fmt.Fprintf(w, `<span id="bflagR" class="%v">%v</span>`, hide, emitImage(flag_restricted, string(c), CS.FlagRestrictedTitle))
		case 'T':
			fmt.Fprintf(w, `<span id="bflagT" class="%v">%v</span>`, hide, emitImage(flag_receipt, string(c), CS.FlagReceiptTitle))
		}
	}

	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</fieldset>`)

	ebcimg := strings.Split(cr.Photo, ",")
	for i := 0; i < len(ebcimg); i++ {
		if ebcimg[i] != "" {
			ebcimg[i] = strings.ReplaceAll(filepath.Join(CS.ImgEbcFolder, filepath.Base(ebcimg[i])), "\\", "/")
		}
	}
	hide = "hide"
	if claimid > 0 {
		hide = ""
	}
	fmt.Fprintf(w, `<fieldset class="claimphotos %v">`, hide)

	showPhotoFrame(w, ebcimg, cr.BonusID)

	fmt.Fprint(w, `</fieldset><!-- below photo frame -->`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="OdoReading">Odo reading</label>`)
	fmt.Fprintf(w, `<input type="number" tabindex="4" id="OdoReading" name="OdoReading" class="odo" value="%v">`, cr.OdoReading)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="ClaimDate">Claim time</label>`)
	fmt.Fprintf(w, `<input type="hidden" id="ClaimTimeISO" name="ClaimTime" value="%v">`, cr.ClaimTime)
	fmt.Fprint(w, `<span>`)
	fmt.Fprintf(w, `<input type="date" tabindex="13" id="ClaimDate" value="%v" onchange="fixClaimTimeISO()">`, claimdate)
	fmt.Fprintf(w, ` <input type="time" tabindex="5" id="ClaimTime" value="%v" onchange="fixClaimTimeISO()">`, claimtime)
	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</fieldset>`)

	hide = "hide"
	if CS.RallyUseQA {
		hide = ""
	}
	fmt.Fprintf(w, `<fieldset class="claimfield %v">`, hide)

	const GoodResult = "&#9745;"
	const BadResult = "&#9746;"

	fmt.Fprint(w, `<label for="AnswerSupplied">Answer</label>`)
	fmt.Fprintf(w, `<input id="AnswerSupplied" tabindex="6" name="AnswerSupplied" title="Answer supplied" class="AnswerSupplied" value="%v">`, cr.AnswerSupplied)
	checked := ""
	if cr.QuestionAnswered {
		checked = "checked"
	}
	fmt.Fprintf(w, ` <span>%v <input type="radio" tabindex="7" name="QuestionAnswered" data-pts="%v" onchange="applyCorrectAnswerBonus(this.checked)" id="QuestionAnswered" value="1" %v> `, GoodResult, CS.RallyQAPoints, checked)
	checked = ""
	if !cr.QuestionAnswered {
		checked = "checked"
	}
	fmt.Fprintf(w, ` %v <input class="" type="radio" tabindex="7" name="QuestionAnswered" onchange="applyCorrectAnswerBonus(!this.checked)" id="QuestionAnsweredN" value="0" %v> `, BadResult, checked)

	fmt.Fprintf(w, ` <span id="CorrectAnswer" title="Correct answer" class="correctanswer">%v</span></span>`, bd.Answer)
	fmt.Fprint(w, `</fieldset>`)

	hide = "hide"
	//fmt.Printf("bd=%v\n", bd)
	if bd.AskPoints {
		hide = ""
	}
	pm := "p"
	if bd.PointsAreMults {
		pm = "m"
	}
	fmt.Fprintf(w, `<fieldset id="askpoints" class="claimfield %v">`, hide)
	fmt.Fprint(w, `<label for="Points">Points</label>`)
	fmt.Fprintf(w, `<input type="number" tabindex="8" id="Points" name="Points" class="Points" data-pm="%v" value="%v">`, pm, cr.Points)
	fmt.Fprint(w, `</fieldset>`)

	hide = "hide"
	if bd.AskMins {
		hide = ""
	}
	fmt.Fprintf(w, `<fieldset class="claimfield %v">`, hide)
	fmt.Fprint(w, `<label for="RestMinutes">Rest minutes</label>`)
	fmt.Fprintf(w, `<input type="number" tabindex="9" id="RestMinutes" name="RestMinutes" class="RestMinutes" value="%v">`, cr.RestMinutes)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="DecisionSelect">Decision</label>`)
	fmt.Fprintf(w, `<input type="hidden" id="chosenDecision" name="Decision" value="%v">`, cr.Decision)
	fmt.Fprint(w, `<select id="DecisionSelect" tabindex="10" onchange="updateClaimDecision(this)">`)
	sel := ""
	if cr.Decision < 0 {
		sel = "selected"
	}
	fmt.Fprintf(w, `<option value="-1" %v>undecided</option>`, sel)
	for i, v := range CS.CloseEBC {
		sel = ""
		if i == cr.Decision {
			sel = "selected"
		}
		fmt.Fprintf(w, `<option value="%v" %v>%v</option>`, i, sel, v)
	}

	fmt.Fprint(w, `</select>`)

	hide = "hide"
	if CS.RallyUsePctPen {
		hide = ""
	}
	fmt.Fprintf(w, ` <span class="%v">&nbsp;&nbsp;&nbsp;`, hide)
	fmt.Fprintf(w, ` <label for="PercentPenalty">%v%% Penalty</label>`, CS.RallyPctPenVal)
	fmt.Fprintf(w, `<input type="hidden" id="valPercentPenalty" value="%v">`, CS.RallyPctPenVal)
	checked = ""
	if cr.PercentPenalty {
		checked = "checked"
	}
	fmt.Fprintf(w, `<input type="checkbox" tabindex="14" id="PercentPenalty" onchange="applyPercentPenalty(this.checked)" data-unchecked="0" name="PercentPenalty" value="1" %v>`, checked)
	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="JudgesNotes">Notes</label>`)
	fmt.Fprintf(w, `<textarea tabindex="11" id="JudgesNotes" name="JudgesNotes" class="judgesnotes">%v</textarea>`, cr.JudgesNotes)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<button class="closebutton" tabindex="12" onclick="saveUpdatedClaim(this);return false">Save updated claim</botton>`)
	fmt.Fprint(w, `</form>`)

	fmt.Fprint(w, `</div>`)
}

func showEBC(w http.ResponseWriter, r *http.Request) {

	const email_icon = "&#9993;"

	const answerGood = "&#10003;"
	const answerBad = "&#10007;"
	const answerTest = "&#8773;"

	claimid := r.FormValue("c")
	if claimid == "" {
		return
	}
	sqlx := `SELECT ebclaims.EntrantID,ebclaims.BonusID,ebclaims.OdoReading,ebclaims.ClaimTime,ifnull(ebclaims.Subject,''),ifnull(ebclaims.ExtraField,'')
	,ifnull(AttachmentTime,ebclaims.ClaimTime),ifnull(DateTime,ebclaims.ClaimTime),ifnull(FirstTime,ebclaims.ClaimTime),ifnull(FinalTime,ebclaims.ClaimTime)
	,EmailID
	 FROM ebclaims WHERE Processed=0 AND rowid=` + claimid

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if !rows.Next() {
		return
	}

	startHTML(w, "EBC claim judging")

	var ebc ElectronicBonusClaim
	err = rows.Scan(&ebc.EntrantID, &ebc.Bonusid, &ebc.OdoReading, &ebc.ClaimTime, &ebc.Subject, &ebc.ExtraField, &ebc.AttachmentTime, &ebc.DateTime, &ebc.FirstTime, &ebc.FinalTime, &ebc.EmailID)
	checkerr(err)

	team := getStringFromDB("SELECT "+RiderNameSQL+" FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "***")
	x := getStringFromDB("SELECT ifnull(PillionName,'') FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "")
	if x != "" {
		team += " &amp; " + x
	}
	teamneeded := x != ""
	if !teamneeded {
		teamneeded = getIntegerFromDB("SELECT TeamID FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), 0) > 0
	}

	bcv := fetchBonusVars(ebc.Bonusid)

	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, `<article class="showebc">`)
	//showReloadTicker(w, r.URL.String())
	fmt.Fprint(w, `<h4>Judge this bonus claim or leave it undecided</h4>`)

	fmt.Fprint(w, `<form id="ebcform" action="saveebc" method="post">`)
	fmt.Fprint(w, `<div>`)
	fmt.Fprintf(w, `<input type="hidden" name="EntrantID" value="%v">`, ebc.EntrantID)
	fmt.Fprintf(w, `<input type="hidden" name="BonusID" value="%v">`, ebc.Bonusid)
	//fmt.Fprintf(w, `<input type="hidden" name="ClaimTime" value="%v">`, ebc.ClaimTime)
	//fmt.Fprintf(w, `<input type="hidden" name="OdoReading" value="%v">`, ebc.OdoReading)
	fmt.Fprintf(w, `<input type="hidden" name="claimid" value="%v">`, claimid)
	fmt.Fprint(w, `<input type="hidden" id="chosenDecision" name="Decision" value="-1">`)
	//fmt.Fprintf(w, `<input type="hidden" name="Points" value="%v">`, bcv.Points)
	fmt.Fprintf(w, `<input type="hidden" name="NextURL" value="%v">`, r.URL.String())

	fmt.Fprintf(w, `Entrant <span class="bold">%v %v</span>`, ebc.EntrantID, team)
	x = getStringFromDB("SELECT BriefDesc FROM bonuses WHERE BonusID='"+ebc.Bonusid+"'", ebc.Bonusid)
	fmt.Fprintf(w, ` Bonus <span class="bold">%v %v</span>`, ebc.Bonusid, x)
	fmt.Fprint(w, ` <span id="claimstats" class="link">Claimed @ `)
	evidence := "Photo: " + ebc.AttachmentTime + "\n"
	evidence += "Claim: " + ebc.ClaimTime + "\n"
	evidence += "Email: " + ebc.DateTime + "\n"
	evidence += "Recvd: " + ebc.FinalTime + "\n"
	fmt.Fprintf(w, `<span class="bold" title="%v" onclick="showEvidence(this)">%v, %v %v</span></span>`, evidence, ebc.OdoReading, logtime(ebc.ClaimTime), email_icon)
	fmt.Fprint(w, `</div>`) // row
	x = getStringFromDB("SELECT ifnull(Notes,'') FROM bonuses WHERE BonusID='"+ebc.Bonusid+"'", "")
	fmt.Fprint(w, `<div>`)
	fmt.Fprintf(w, `<span class="bonusnotes">%v</span>`, x)

	x = getStringFromDB("SELECT ifnull(Flags,'') FROM bonuses WHERE BonusID='"+ebc.Bonusid+"'", "")

	if teamneeded && !strings.ContainsRune(x, '2') {
		x += "2"
	}
	for _, c := range x {
		switch c {
		case '2':
			fmt.Fprint(w, emitImage(flag_team, string(c), CS.FlagTeamTitle))
		case 'A':
			fmt.Fprint(w, emitImage(flag_alert, string(c), CS.FlagAlertTitle))
		case 'B':
			fmt.Fprint(w, emitImage(flag_bike, string(c), CS.FlagBikeTitle))
		case 'D':
			fmt.Fprint(w, emitImage(flag_daylight, string(c), CS.FlagDaylightTitle))
		case 'F':
			fmt.Fprint(w, emitImage(flag_face, string(c), CS.FlagFaceTitle))
		case 'N':
			fmt.Fprint(w, emitImage(flag_night, string(c), CS.FlagNightTitle))
		case 'R':
			fmt.Fprint(w, emitImage(flag_restricted, string(c), CS.FlagRestrictedTitle))
		case 'T':
			fmt.Fprint(w, emitImage(flag_receipt, string(c), CS.FlagReceiptTitle))
		}
	}

	fmt.Fprintf(w, `<div id="finetune" class="hide">`)
	fmt.Fprint(w, `<label for="OdoReading">Odo reading</label> `)
	fmt.Fprintf(w, `<input type="number" id="OdoReading" name="OdoReading" class="odo" value="%v"> `, ebc.OdoReading)
	fmt.Fprint(w, ` <label for="ClaimTime">ClaimTime</label> `)

	// format is 'datetime', not 'datetime-local'. It's a simple string, not a local date and time
	fmt.Fprintf(w, `<input type="datetime" id="ClaimTime" name="ClaimTime" class="ClaimTime" value="%v"> `, ebc.ClaimTime)

	fmt.Fprintf(w, `<span class="evidence">&nbsp;&nbsp;Photo: <strong>%v</strong>  &nbsp;Email: <strong>%v</strong>  &nbsp;Recvd: <strong>%v</strong></span>`, ebc.AttachmentTime, ebc.DateTime, ebc.FinalTime)

	fmt.Fprint(w, `</div>`)

	fmt.Fprint(w, `</div>`) // row

	hide := "hide"
	if CS.RallyUseQA && bcv.Answer != "" {
		hide = ""
	}
	fmt.Fprintf(w, `<div class="qa %v">`, hide)
	fmt.Fprint(w, `<label for="AnswerSupplied">QQQ?</label> `)
	fmt.Fprintf(w, `<input type="text" id="AnswerSupplied" name="AnswerSupplied" title="Answer supplied by entrant" class="AnswerSupplied" readonly value="%v">`, ebc.ExtraField)
	answerok := strings.EqualFold(ebc.ExtraField, bcv.Answer) && bcv.Answer != ""
	checked := ""
	if answerok {
		checked = "checked"
	}
	fmt.Fprintf(w, ` <input type="radio" id="QuestionAnswered" data-pts="%v" name="QuestionAnswered" value="1" %v> %v &nbsp; `, CS.RallyQAPoints, checked, answerGood)
	checked = ""
	if !answerok {
		checked = "checked"
	}
	fmt.Fprintf(w, ` <input type="radio" name="QuestionAnswered" id="QuestionAnsweredN"value="0" %v> %v &nbsp; `, checked, answerBad)
	fmt.Fprintf(w, ` %v<span class="correctanswer" title="Correct answer">%v</span>`, answerTest, bcv.Answer)
	fmt.Fprint(w, `</div>`)

	fmt.Fprint(w, `<div>`)

	fmt.Fprintf(w, `<button id="leavebutton" data-result="-1" name="Decision" onclick="closeEBC(this)" class="closebutton">%v</button>`, CS.CloseEBCUndecided)

	hide = "hide"
	//fmt.Printf("bd=%v\n", bd)
	if bcv.AskPoints {
		hide = ""
	}
	pm := "p"
	if bcv.PointsAreMults {
		pm = "m"
	}
	fmt.Fprintf(w, `<span id="askpoints" class=" %v">`, hide)
	fmt.Fprint(w, `<label for="Points">Points</label> `)
	fmt.Fprintf(w, `<input type="number"id="Points" name="Points" class="Points" data-pm="%v" value="%v"> `, pm, bcv.Points)
	fmt.Fprint(w, `</span>`)

	hide = "hide"
	if bcv.AskMins {
		hide = ""
	}
	fmt.Fprintf(w, `<span class=" %v">`, hide)
	fmt.Fprint(w, `<label for="RestMinutes">Rest minutes</label> `)
	fmt.Fprintf(w, `<input type="number" id="RestMinutes" name="RestMinutes" class="RestMinutes" value="%v"> `, bcv.RestMins)
	fmt.Fprint(w, `</span>`)

	fmt.Fprintf(w, `<button data-result="0"  name="Decision" autofocus onclick="closeEBC(this)" class="closebutton">%v</button>`, CS.CloseEBC[0])

	hide = "hide"
	if CS.RallyUsePctPen {
		hide = ""
	}
	fmt.Fprintf(w, ` <span class="%v">`, hide)
	fmt.Fprintf(w, `<input type="hidden" id="valPercentPenalty" value="%v">`, CS.RallyPctPenVal)
	fmt.Fprintf(w, `<button class="closebutton" data-result="0" id="PercentPenalty" onclick="closeEBC(this)">%v%% Penalty</button>`, CS.RallyPctPenVal)
	fmt.Fprint(w, `</span>`)

	x = ""
	fmt.Fprintf(w, `<input type="text" id="judgesnotes" name="JudgesNotes" oninput="killReload(this)" class="judgesnotes" placeholder="Notes"  value="%v">`, x)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div>`)

	for i := 1; i < 10; i++ {
		fmt.Fprintf(w, `<button data-result="%v"  name="Decision" onclick="closeEBC(this)" class="closebutton">%v</button>`, i, CS.CloseEBC[i])
	}
	fmt.Fprint(w, `</div>`)
	showPhotosEBC(w, ebc.EmailID, ebc.Bonusid)

	fmt.Fprint(w, `</form>`)

	fmt.Fprint(w, `</article>`)

	const trapkeys = `
    document.onkeydown = function(evt) {
    evt = evt || window.event;
    var isEscape = false;
    if ("key" in evt) {
        isEscape = (evt.key === "Escape" || evt.key === "Esc");
    } else {
        isEscape = (evt.keyCode === 27);
    }
    if (isEscape) {
        leaveUndecided();
	}}
`
	fmt.Fprintf(w, `<script>%v</script>`, trapkeys)

}

func showPhotosEBC(w http.ResponseWriter, emailid int, BonusID string) {

	sqlx := "SELECT Image FROM ebcphotos WHERE EmailID=" + strconv.Itoa(emailid)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	fmt.Fprint(w, `<div class="imgcomparediv">`)

	// "if(this.width=='50%')this.width='100%';else this.width='50%'"
	showimg := make([]string, maximg)

	ix := 0
	for rows.Next() {
		var img string
		err := rows.Scan(&img)
		checkerr(err)
		if img != "" {
			showimg[ix] = strings.ReplaceAll(filepath.Join(CS.ImgEbcFolder, filepath.Base(img)), `\`, `/`)
			ix++
		}
		if ix >= maximg {
			break
		}
	}

	showPhotoFrame(w, showimg, BonusID)

}

func showPhotoFrame(w http.ResponseWriter, photos []string, BonusID string) {

	maximg := len(photos)
	photo0 := ""
	if maximg > 0 {
		photo0 = photos[0]
	}
	fmt.Fprint(w, `<div class="ebcimgdiv" id="ebcimgdiv" onclick="cycleImgSize(this)">`)

	fmt.Fprintf(w, `<img id="imgdivimg" alt="*" src="%v" title="%v">`, photo0, CS.EBCImgTitle)
	fmt.Fprintf(w, `<input type="hidden" id="chosenPhoto" name="Photo" value="%v">`, photos[0])

	fmt.Fprint(w, `<div id="imgdivs">`)

	for ix := 1; ix < maximg; ix++ {
		if photos[ix] != "" {
			fmt.Fprintf(w, `<img src="%v" alt="*" onclick="swapimg(this)" title="%v">`, photos[ix], CS.EBCImgSwapTitle)
		}
	}
	fmt.Fprint(w, `</div>`) // imgdivs
	fmt.Fprint(w, `</div>`) // ebcimgdiv

	fmt.Fprint(w, `<div class="bonusimgdiv" id="bonusimgdiv">`)
	bimg := strings.ReplaceAll(filepath.Join(CS.ImgBonusFolder, filepath.Base(getStringFromDB("SELECT ifnull(Image,'') FROM bonuses WHERE BonusID='"+BonusID+"'", ""))), `\`, `/`)
	fmt.Fprintf(w, `<img src="%v" alt="*" title="%v">`, bimg, CS.RallyBookImgTitle)
	fmt.Fprint(w, `</div>`)

}

func intval(x string) int {

	res, err := strconv.Atoi(x)
	if err != nil {
		return 0
	}
	return res
}

func saveEBC(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println(r.FormValue("Points"))

	decision := intval(r.FormValue("Decision"))
	processed := 0
	if decision >= 0 {
		processed = 1
	}
	claimid := intval(r.FormValue("claimid"))

	sqlx := fmt.Sprintf("UPDATE ebclaims SET Processed=%v, Decision=%v WHERE Processed=0 AND rowid=%v", processed, decision, claimid)
	//fmt.Println(sqlx)
	res, err := DBH.Exec(sqlx)
	checkerr(err)
	n, err := res.RowsAffected()
	checkerr(err)
	if n == 0 || decision < 0 {
		return
	}

	sqlx = "INSERT INTO claims (LoggedAt, ClaimTime, EntrantID, BonusID, OdoReading, Decision, Photo, Points, RestMinutes, AskPoints, AskMinutes, Leg"
	sqlx += ",Evidence,QuestionAsked, AnswerSupplied, QuestionAnswered, JudgesNotes, PercentPenalty) "
	sqlx += "VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()

	tn := time.Now()
	LoggedAt := tn.Format(time.RFC3339)
	points := intval(r.FormValue("Points"))
	restmins := intval(r.FormValue("RestMinutes"))
	askpoints := intval(r.FormValue("AskPoints"))
	askmins := intval(r.FormValue("AskMinutes"))
	qasked := intval(r.FormValue("QuestionAsked"))
	qanswered := intval(r.FormValue("QuestionAnswered"))
	percent := intval(r.FormValue("PercentPenalty"))
	_, err = stmt.Exec(LoggedAt, r.FormValue("ClaimTime"), r.FormValue("EntrantID"), r.FormValue("BonusID"),
		r.FormValue("OdoReading"), decision, ImgFromURL(r.FormValue("Photo")), points, restmins, askpoints, askmins, CS.CurrentLeg,
		r.FormValue("Evidence"), qasked, r.FormValue("AnswerSupplied"), qanswered, r.FormValue("JudgesNotes"), percent)
	checkerr(err)
	recalc_scorecard(intval(r.FormValue("EntrantID")))
}
