package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"path/filepath"
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
	Exists         bool
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

const tick_icon = "&#10004;"
const cross_icon = "&#10006;"
const undecided_icon = "?"

const maximg = 3

const judginghelp = `
<article id="judginghelp" class="popover" popover>
<h1>CLAIM JUDGING</h1>
<p>You are asked to assess the bonus claim by comparing it against the specific requirements. Is it a good photo? Is the face in the photo? etc</p>
<p>Specific bonus requirements are shown including icons indicating "Alert", "Bike in photo", "Daylight only", "Face in photo", "Night only", "Restricted hours/access" and "Receipt/ticket required"</p>
<p><strong>Accept good claim</strong> awards the points and other benefits of the claim</p>
<p><strong>Leave undecided</strong> does not judge the claim but returns it to the end of the queue</p>
<p>Other responses apart from "Exclude claim" <strong>reject the claim</strong> for the stated reason. The claim and reason for rejection will appear on the scorecard. This is the normal method of rejecting claims and should be used in preference to excluding the claim.</p>
<p><strong>Exclude claim</strong> excludes the claim from scoring altogether. It should only rarely be used as nothing will appear on the scorecard. It is intended for use with claims which are not judgeable as opposed to those which can be accepted or rejected.</p>
<p>Clicking the info line after '@' will make odo reading and claim time editable</p>
</article>
`

func countFixedClaims(w http.ResponseWriter, r *http.Request) {

	sqlx := "SELECT count(*) FROM claims WHERE BonusID='" + r.FormValue("b") + "' AND AskPoints=0 AND Decision >= 0"
	//fmt.Println(sqlx)
	n := getIntegerFromDB(sqlx, 0)
	fmt.Fprintf(w, `{"ok":true,"msg":"%v"}`, n)
}

func deleteClaim(w http.ResponseWriter, r *http.Request) {

	claimid := r.PathValue("claimid")
	if claimid == "" {
		fmt.Fprint(w, `{"ok":false,"msg:"incomplete request}`)
		return
	}
	sqlx := "SELECT EntrantID FROM claims WHERE rowid=" + claimid
	entrant := getIntegerFromDB(sqlx, 0)
	sqlx = "DELETE FROM claims WHERE rowid=" + claimid
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"claim deleted"}`)
	if entrant > 0 {
		recalc_scorecard(entrant)
		rankEntrants(false)
	}

}

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

	sqlx := "SELECT ifnull(BriefDesc,'" + CS.NoSuchBonus + "'),ifnull(Points,0),ifnull(Notes,''),ifnull(Flags,''),ifnull(AskPoints,0),ifnull(RestMinutes,0),ifnull(AskMinutes,0),ifnull(Image,''),ifnull(Question,''),ifnull(Answer,'')"
	sqlx += " FROM bonuses WHERE BonusID='" + b + "'"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if !rows.Next() {
		res.BriefDesc = CS.NoSuchBonus
		return res
	}
	err = rows.Scan(&res.BriefDesc, &res.Points, &res.Notes, &res.Flags, &ap, &res.RestMins, &am, &res.Image, &res.Question, &res.Answer)
	checkerr(err)
	res.Exists = true
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

func fixFixedClaims(w http.ResponseWriter, r *http.Request) {

	points := intval(r.FormValue("p"))
	sqlx := fmt.Sprintf("UPDATE claims SET Points=%v WHERE BonusID='%v' ", points, r.FormValue("b"))
	if r.FormValue("d") == "new" {
		sqlx += " AND Decision < 0"
	}
	//fmt.Println(sqlx)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)

}

func judge_new_claims(w http.ResponseWriter, r *http.Request) {

	if CS.UseEBC {
		list_EBC_claims(w, r)
		return
	}
	showClaim(w, r)
}

const claimstopline = `	
<div class="topline">
		

		<fieldset>
			<button title="show Scorecard" onclick="loadPage('/score?e=%v')">☷</button>
		</fieldset>
		<fieldset>
			<button title="why not combo?" onclick="loadPage('/ynot?e=%v')">☹??</button>
		</fieldset>

	</div>
`

func list_claims(w http.ResponseWriter, r *http.Request) {

	const addnew_icon = "&nbsp;+&nbsp;"

	const filterclear = `-`

	r.ParseForm()

	EntrantSelector := loadEntrantsList()

	sseq := r.FormValue("sseq")
	dseq := r.FormValue("dseq")

	if sseq == "" {
		sseq = "ClaimTime"
		dseq = "desc"
	}
	esel := intval(r.FormValue("esel"))

	startHTML(w, "CLAIMS LOG")
	if esel > 0 {
		fmt.Fprintf(w, claimstopline, esel, esel)
	}

	fmt.Fprint(w, `<div class="claimslog">`)

	showReloadTicker(w, r.URL.String())

	fmt.Fprintf(w, `<button autofocus title="Add new claim" class="plus" onclick="window.location.href='/claim/0';return false">%v</button> `, addnew_icon)

	fmt.Fprint(w, `<form id="claimslogfrm">`)

	fmt.Fprintf(w, `<input type="hidden" id="sseq" name="sseq" value="%v">`, sseq)
	fmt.Fprintf(w, `<input type="hidden" id="dseq" name="dseq" value="%v">`, dseq)

	fmt.Fprint(w, `<fieldset class="inline filter"><legend>Filters</legend>`)

	fmt.Fprintf(w, `<span class="button" title="Clear filters, show everything" onclick="resetClaimslogFilter();"> %v </span> &nbsp;`, filterclear)

	fmt.Fprintf(w, `<input name="esel" type="number" onchange="reloadClaimslog()" placeholder="flag" class="EntrantID" value="%v">`, r.FormValue("esel"))

	bsel := r.FormValue("bsel")
	fmt.Fprintf(w, `<input name="bsel" onchange="reloadClaimslog()" placeholder="bonus" class="BonusID" value="%v">`, bsel)

	sel := ""
	if bsel == "" {
		sel = "selected"
	}

	dsel := r.FormValue("dsel")
	sel = ""
	if dsel == "" {
		sel = "selected"
	}
	fmt.Fprintf(w, ` <select name="dsel" value="%v" title="Filter by Decision" onchange="reloadClaimslog()">`, dsel)
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

	for ix, msg := range CS.CloseEBC {
		if ix < 1 {
			continue
		}
		sel = ""
		if dsel == strconv.Itoa(ix) {
			sel = "selected"
		}
		fmt.Fprintf(w, `<option value="%v" %v>%v</option>`, ix, sel, msg)
	}
	fmt.Fprint(w, `</select>`)

	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `</form>`)

	sqlx := `SELECT claims.rowid,ifnull(LoggedAt,''),ClaimTime,claims.EntrantID,BonusID,OdoReading,Decision,ifnull(JudgesNotes,'') FROM claims`
	sqlx += " LEFT JOIN entrants ON claims.EntrantID=entrants.EntrantID"
	where := ""
	if esel > 0 {
		where += "  claims.EntrantID=" + strconv.Itoa(esel)
	}
	if bsel != "" {
		if where != "" {
			where += " AND "
		}
		where += "BonusID='" + strings.ToUpper(bsel) + "'"
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
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			where += "= " + dsel
		default:
			where += "> 0"
		}
	}
	if where != "" {
		sqlx += " WHERE " + where
	}
	sqlx += " ORDER BY " + sseq + " " + dseq

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	fmt.Fprint(w, `<fieldset class="row claims hdr">`)
	fmt.Fprint(w, `<fieldset class="col claims hdr sort" onclick="reseqClaimslog('RiderLast')" >Entrant</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col claims hdr sort" onclick="reseqClaimslog('BonusID')" >Bonus</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col claims hdr sort" onclick="reseqClaimslog('OdoReading')" >Odo</fieldset>`)
	fmt.Fprint(w, `<fieldset class="col claims hdr sort" onclick="reseqClaimslog('ClaimTime')" >Time</fieldset>`)
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
		fmt.Fprintf(w, `<fieldset class="row claims" onclick="window.location.href='/claim/%v?back=/claims'">`, claimid)
		fmt.Fprintf(w, `<fieldset class="col claims" >%v #%v</fieldset>`, rname, cr.EntrantID)
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

	sqlx := `SELECT ebclaims.rowid,ebclaims.EntrantID,` + RiderNameSQL + `,` + PillionNameSQL + `,ebclaims.BonusID,ifnull(xbonus.BriefDesc,'` + CS.NoSuchBonus + `'),ebclaims.OdoReading,ebclaims.ClaimTime
			 		FROM ebclaims LEFT JOIN entrants ON ebclaims.EntrantID=entrants.EntrantID
					LEFT JOIN (SELECT BonusID,BriefDesc FROM bonuses) AS xbonus ON ebclaims.BonusID=xbonus.BonusID
					 WHERE Processed=0 ORDER BY Decision DESC,FinalTime;`

	rows, err := DBH.Query(sqlx)
	checkerr(err)

	startHTML(w, "JUDGING NEW CLAIMS")

	fmt.Fprint(w, `<div class="ebclist">`)

	showReloadTicker(w, r.URL.String())

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
		//fmt.Printf("%v == %v\n", ebc.Bonusid, ebc.BriefDesc)
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
	sqlx := "SELECT EntrantID, ifnull(RiderFirst,''),ifnull(RiderLast,'')," + PillionNameSQL + " FROM entrants"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var e int
		var rl string
		var rf string
		var rn string
		var pn string
		err = rows.Scan(&e, &rf, &rl, &pn)
		checkerr(err)
		rn = "<strong>" + rl + "</strong>" + ", " + rf
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

	sqlx := "SELECT ifnull(OdoRallyStart,0) FROM entrants WHERE EntrantID=" + r.FormValue("EntrantID")
	checkoutodo := getIntegerFromDB(sqlx, 0)
	sqlx = fmt.Sprintf("SELECT ifnull(OdoReading,%v) FROM claims WHERE EntrantID=%v", checkoutodo, r.FormValue("EntrantID"))
	sqlx += " ORDER BY EntrantID,ClaimTime DESC,OdoReading DESC"

	lastOdo := getIntegerFromDB(sqlx, checkoutodo)
	thisOdo := intval(r.FormValue("OdoReading"))
	if thisOdo == 0 {
		thisOdo = lastOdo + 1
	}

	sqlx = "INSERT INTO claims (LoggedAt, ClaimTime, EntrantID, BonusID, OdoReading, Decision, Photo, Points, RestMinutes, Leg"
	sqlx += ", AnswerSupplied, QuestionAnswered, JudgesNotes, PercentPenalty) "
	sqlx += "VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	//r.FormValue("AnswerSupplied"), r.FormValue("QuestionAnswered"),
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().Format(time.RFC3339), r.FormValue("ClaimTime"), r.FormValue("EntrantID"), strings.ToUpper(r.FormValue("BonusID")), thisOdo, intval(r.FormValue("Decision")), filepath.Base(r.FormValue("Photo")), intval(r.FormValue("Points")), intval(r.FormValue(("RestMinutes"))), Leg, r.FormValue("AnswerSupplied"), r.FormValue("QuestionAnswered"), r.FormValue("JudgesNotes"), intval(r.FormValue("PercentPenalty")))
	checkerr(err)

}
func saveClaim(r *http.Request) {

	//fmt.Printf("saveclaim: %v\n", r)
	claimid := intval(r.FormValue("claimid"))

	if claimid < 1 {
		insertNewClaim(r)
		recalc_scorecard(intval(r.FormValue("EntrantID")))
		rankEntrants(false)
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
	rankEntrants(false)

}

const claimTopline = `
	<div class="topline">
		<fieldset>
			<button title="Delete this Claim?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" style="display:none;" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
		<fieldset>
			<button id="updatedb" class="hideuntil" title="Delete Claim" disabled onclick="deleteClaim(this)"></button>
		</fieldset>
		<fieldset>
			<button title="back to list" onclick="window.location.href='/claims'">↥☰↥</button>
		</fieldset>
	</div>
`

const lastClaimTopline = `
	<div class="topline">
		<fieldset>
			Last claim: %v [%v] - %v %v - %v
		</fieldset>
		<fieldset>
			<button title="back to list" onclick="window.location.href='/claims'">↥☰↥</button>
		</fieldset>
	</div>
`
const unloadTrapper = `

<script>
window.addEventListener("beforeunload",function(e) {
	let ifrm = document.getElementById('iclaim');
	if (ifrm.getAttribute('data-unloadok')!='0') return true;
	var confirmationMessage = 'OMG';
	(e || window.event).returnValue = confirmationMessage;
	return confirmationMessage;
});
</script>
->
`

func emitLastClaimTopline(w http.ResponseWriter) {

	sqlx := "SELECT claims.EntrantID," + RiderNameSQL + ",claims.BonusID,BriefDesc,Decision"
	sqlx += " FROM claims LEFT JOIN entrants ON claims.EntrantID=entrants.EntrantID"
	sqlx += " LEFT JOIN bonuses ON claims.BonusID=bonuses.BonusID"
	sqlx += " ORDER BY LoggedAt DESC LIMIT 1"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if rows.Next() {
		var entrant int
		var rname string
		var bid, briefdesc string
		var decision int
		err = rows.Scan(&entrant, &rname, &bid, &briefdesc, &decision)
		checkerr(err)
		xd := undecided_icon
		if decision == 0 {
			xd = tick_icon
		} else {
			xd = cross_icon
		}
		fmt.Fprintf(w, lastClaimTopline, rname, entrant, bid, briefdesc, xd)
	}
}

func showClaim(w http.ResponseWriter, r *http.Request) {

	var cr ClaimRecord

	var claimdate string
	var claimtime string

	var ed EntrantDetails
	var bd BonusClaimVars

	r.ParseForm()

	claimhdr := "Claims log"
	claimid := intval(r.PathValue("claim"))
	if claimid < 1 {
		claimhdr = "Making new claim"
	}
	startHTMLBL(w, claimhdr, "/claims")

	if claimid > 0 {
		fmt.Fprint(w, claimTopline)
	} else if !CS.UseEBC {
		emitLastClaimTopline(w)
	}

	fmt.Fprint(w, `</header><div class="claim">`)
	fmt.Fprint(w, `<form id="iclaim" data-unloadok="1">`)

	datasave := "claims"
	if claimid < 1 && !CS.UseEBC {
		datasave = "judgenew"
	}
	fmt.Fprintf(w, `<input type="hidden" id="claimid" name="claimid" data-save="%v" value="%v">`, datasave, claimid)
	if claimid < 1 {
		//fmt.Fprint(w, `<input type="text" autofocus tabindex="1" class="subject" placeholder="Paste email Subject line here" oninput="pasteNewClaim(this)">`)
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
	}
	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="EntrantID">Entrant</label>`)
	fmt.Fprint(w, `<input type="number" tabindex="2" autofocus id="EntrantID" name="EntrantID" class="EntrantID" oninput="setdirty(this);fetchEntrantDetails(this);"`)
	if claimid > 0 {
		fmt.Fprintf(w, ` readonly value="%v"`, cr.EntrantID)
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
	fmt.Fprint(w, ` <span id="edwarn" class="warn hide"></span>`)
	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="BonusID">Bonus code</label>`)
	af := ""
	if claimid > 0 {
		af = "autofocus"
	}
	fmt.Fprintf(w, `<input type="text" tabindex="3" %v id="BonusID" name="BonusID" class="BonusID" oninput="setdirty(this);fetchBonusDetails(this)" value="%v">`, af, cr.BonusID)
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

	manclm := "hide"
	if CS.UseEBC {
		manclm = ""
	}

	fmt.Fprintf(w, `<fieldset class="claimfield %v" title="This defaults to last reading +1">`, manclm)
	fmt.Fprint(w, `<label for="OdoReading">Odo reading</label>`)
	fmt.Fprintf(w, `<input type="number" tabindex="4" id="OdoReading" name="OdoReading" onchange="setdirty(this)" class="odo" value="%v">`, cr.OdoReading)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprintf(w, `<fieldset class="claimfield %v">`, manclm)
	fmt.Fprint(w, `<label for="ClaimDate">Claim time</label>`)
	fmt.Fprintf(w, `<input type="hidden" id="ClaimTimeISO" name="ClaimTime" value="%v">`, cr.ClaimTime)
	fmt.Fprint(w, `<span>`)
	fmt.Fprintf(w, `<input type="date" tabindex="13" id="ClaimDate" value="%v" onchange="setdirty(this);fixClaimTimeISO()">`, claimdate)
	fmt.Fprintf(w, ` <input type="time" tabindex="5" id="ClaimTime" value="%v" onchange="setdirty(this);fixClaimTimeISO()">`, claimtime)
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
	fmt.Fprintf(w, `<input id="AnswerSupplied" tabindex="6" name="AnswerSupplied"  title="Answer supplied" class="AnswerSupplied" value="%v">`, cr.AnswerSupplied)
	checked := ""
	if cr.QuestionAnswered {
		checked = "checked"
	}
	fmt.Fprintf(w, ` <span>%v <input type="radio" tabindex="7" name="QuestionAnswered" data-pts="%v" onchange="setdirty(this);applyCorrectAnswerBonus(this.checked)" id="QuestionAnswered" value="1" %v> `, GoodResult, CS.RallyQAPoints, checked)
	checked = ""
	if !cr.QuestionAnswered {
		checked = "checked"
	}
	fmt.Fprintf(w, ` %v <input class="" type="radio" tabindex="7" name="QuestionAnswered" onchange="setdirty(this);applyCorrectAnswerBonus(!this.checked)" id="QuestionAnsweredN" value="0" %v> `, BadResult, checked)

	fmt.Fprintf(w, ` <span id="CorrectAnswer" title="Correct answer" class="correctanswer">%v</span></span>`, bd.Answer)
	fmt.Fprint(w, `</fieldset>`)

	hide = "hide"
	//fmt.Printf("bd=%v\n", bd)
	if bd.AskPoints || false {
		hide = ""
	}
	pm := "p"
	if bd.PointsAreMults {
		pm = "m"
	}
	fmt.Fprintf(w, `<fieldset id="askpoints" class="claimfield %v">`, hide)
	fmt.Fprint(w, `<label for="Points">Points</label>`)
	fmt.Fprintf(w, `<input type="number" tabindex="8" id="Points" name="Points" class="Points" onchange="setdirty(this)" data-pm="%v" value="%v">`, pm, cr.Points)
	fmt.Fprint(w, `</fieldset>`)

	hide = "hide"
	if bd.AskMins {
		hide = ""
	}
	fmt.Fprintf(w, `<fieldset class="claimfield %v">`, hide)
	fmt.Fprint(w, `<label for="RestMinutes">Rest minutes</label>`)
	fmt.Fprintf(w, `<input type="number" tabindex="9" id="RestMinutes" name="RestMinutes" class="RestMinutes" onchange="setdirty(this)" value="%v">`, cr.RestMinutes)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="DecisionSelect">Decision</label>`)
	fmt.Fprintf(w, `<input type="hidden" id="chosenDecision" name="Decision" value="%v">`, cr.Decision)
	fmt.Fprint(w, `<select id="DecisionSelect" tabindex="10" onchange="setdirty(this);updateClaimDecision(this)">`)
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
	fmt.Fprintf(w, `<input type="checkbox" tabindex="14" id="PercentPenalty" onchange="setdirty(this);applyPercentPenalty(this.checked)" data-unchecked="0" name="PercentPenalty" value="1" %v>`, checked)
	fmt.Fprint(w, `</span>`)
	fmt.Fprint(w, `</fieldset>`)

	fmt.Fprint(w, `<fieldset class="claimfield">`)
	fmt.Fprint(w, `<label for="JudgesNotes">Notes</label>`)
	fmt.Fprintf(w, `<textarea tabindex="11" id="JudgesNotes" name="JudgesNotes" oninput="setdirty(this)" class="judgesnotes">%v</textarea>`, cr.JudgesNotes)
	fmt.Fprint(w, `</fieldset>`)

	savex := "Save updated claim"
	if claimid < 1 {
		savex = "Save new claim"
	}
	fmt.Fprintf(w, `<button class="closebutton" id="closebutton" tabindex="12" onclick="saveUpdatedClaim(this);return false">%v</button>`, savex)

	ebcimg := strings.Split(cr.Photo, ",")
	for i := 0; i < len(ebcimg); i++ {
		if ebcimg[i] != "" {
			ebcimg[i] = "/" + strings.ReplaceAll(filepath.Join(CS.ImgEbcFolder, filepath.Base(ebcimg[i])), "\\", "/")
		}
	}
	fmt.Fprint(w, `<fieldset class="claimphotos ">`)

	showPhotoFrame(w, ebcimg, cr.BonusID)

	fmt.Fprint(w, `</fieldset><!-- below photo frame -->`)

	fmt.Fprint(w, `</form>`)

	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, unloadTrapper)
}

func showEBC(w http.ResponseWriter, r *http.Request) {

	const email_icon = "&#9993;"

	const answerGood = "&#10003;"
	const answerBad = "&#10007;"
	const answerTest = "&#8773;"

	claimid := r.PathValue("claim")
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

	startHTML(w, "Claim judging")

	var ebc ElectronicBonusClaim
	err = rows.Scan(&ebc.EntrantID, &ebc.Bonusid, &ebc.OdoReading, &ebc.ClaimTime, &ebc.Subject, &ebc.ExtraField, &ebc.AttachmentTime, &ebc.DateTime, &ebc.FirstTime, &ebc.FinalTime, &ebc.EmailID)
	checkerr(err)

	team := getStringFromDB("SELECT "+RiderNameSQL+" FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "***")
	x := getStringFromDB("SELECT "+PillionNameSQL+" FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "")
	if x != "" {
		team += " &amp; " + x
	}
	teamneeded := x != ""
	if !teamneeded {
		teamneeded = getIntegerFromDB("SELECT TeamID FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), 0) > 0
	}

	bcv := fetchBonusVars(ebc.Bonusid)

	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, judginghelp)

	fmt.Fprint(w, `<article class="showebc">`)
	//showReloadTicker(w, r.URL.String())
	fmt.Fprint(w, `<p class="h4">Judge this bonus claim or leave it undecided <input type="button" class="popover" style="font-size: .8em;" popovertarget="judginghelp" value="[click for help]"></p>`)

	fmt.Fprint(w, `<form id="ebcform" action="saveebc" onsubmit="event.preventDefault()" method="post">`)
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
	x = bcv.BriefDesc
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
	answerok := CS.RallyUseQA && strings.EqualFold(ebc.ExtraField, bcv.Answer) && bcv.Answer != ""
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
	fmt.Fprintf(w, `<input type="text" id="judgesnotes" name="JudgesNotes" class="judgesnotes" placeholder="Notes"  value="%v">`, x)
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
			showimg[ix] = "/" + strings.ReplaceAll(filepath.Join(CS.ImgEbcFolder, filepath.Base(img)), `\`, `/`)
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
	fmt.Fprint(w, `<div class="ebcimgdiv cangrow" id="ebcimgdiv" onclick="cycleImgSize(this)">`)

	fmt.Fprintf(w, `<img id="imgdivimg" alt=" " src="%v" title="%v">`, photo0, CS.EBCImgTitle)
	fmt.Fprintf(w, `<input type="hidden" id="chosenPhoto" name="Photo" value="%v">`, photos[0])

	fmt.Fprint(w, `<div id="imgdivs">`)

	for ix := 1; ix < maximg; ix++ {
		if photos[ix] != "" {
			fmt.Fprintf(w, `<img src="%v" alt=" " onclick="swapimg(this)" title="%v">`, photos[ix], CS.EBCImgSwapTitle)
		}
	}
	fmt.Fprint(w, `</div>`) // imgdivs
	fmt.Fprint(w, `</div>`) // ebcimgdiv

	fmt.Fprint(w, `<div class="bonusimgdiv" id="bonusimgdiv">`)
	bimg := "/" + strings.ReplaceAll(filepath.Join(CS.ImgBonusFolder, filepath.Base(getStringFromDB("SELECT ifnull(Image,'') FROM bonuses WHERE BonusID='"+BonusID+"'", ""))), `\`, `/`)
	fmt.Fprintf(w, `<img src="%v" id="bonusPhoto" alt=" " title="%v" data-folder="%v">`, bimg, CS.RallyBookImgTitle, CS.ImgBonusFolder)
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
	//	fmt.Println(r.Form)
	//	fmt.Println(r.FormValue("Points"))

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
		fmt.Fprint(w, `{"ok":false,"msg":update failed}`)
		return
	}

	sqlx = "SELECT ifnull(OdoRallyStart,0) FROM entrants WHERE EntrantID=" + r.FormValue("EntrantID")
	checkoutodo := getIntegerFromDB(sqlx, 0)
	sqlx = fmt.Sprintf("SELECT ifnull(OdoReading,%v) FROM claims WHERE EntrantID=%v", checkoutodo, r.FormValue("EntrantID"))
	sqlx += " ORDER BY EntrantID,ClaimTime DESC,OdoReading DESC"

	lastOdo := getIntegerFromDB(sqlx, checkoutodo)
	thisOdo := intval(r.FormValue("OdoReading"))
	if thisOdo == 0 {
		thisOdo = lastOdo + 1
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
		thisOdo, decision, ImgFromURL(r.FormValue("Photo")), points, restmins, askpoints, askmins, CS.CurrentLeg,
		r.FormValue("Evidence"), qasked, r.FormValue("AnswerSupplied"), qanswered, r.FormValue("JudgesNotes"), percent)
	checkerr(err)
	recalc_scorecard(intval(r.FormValue("EntrantID")))
	rankEntrants(false)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}
