package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

func emitImage(img string, alt string, title string) string {

	res := fmt.Sprintf(`<img alt="%v", title="%v" class="flagicon" src="data:image/png;base64,`, alt, title)
	for _, xl := range strings.Split(img, "\n") {
		res += xl
	}
	res += `">`
	return res

}

// Show judgeable claims submitted electronically
func list_EBC_claims(w http.ResponseWriter, r *http.Request) {

	sqlx := `SELECT ebclaims.rowid,ebclaims.EntrantID,entrants.RiderName,ifnull(entrants.PillionName,''),ebclaims.BonusID,xbonus.BriefDesc,ebclaims.OdoReading,ebclaims.ClaimTime
	 		FROM ebclaims LEFT JOIN entrants ON ebclaims.EntrantID=entrants.EntrantID
			LEFT JOIN (SELECT BonusID,BriefDesc FROM bonuses) AS xbonus ON ebclaims.BonusID=xbonus.BonusID
			 WHERE Processed=0 ORDER BY Decision DESC,FinalTime;`

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<style>%s</style>`, css)
	fmt.Fprintf(w, `<script>%s</script>`, script)
	fmt.Fprint(w, `<div class="ebclist">`)
	for rows.Next() {
		var ebc ElectronicBonusClaim
		err := rows.Scan(&ebc.Claimid, &ebc.EntrantID, &ebc.RiderName, &ebc.PillionName, &ebc.Bonusid, &ebc.BriefDesc, &ebc.OdoReading, &ebc.ClaimTime)
		checkerr(err)
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

}

func logtime(stamp string) string {
	/* We're really only interested in the time of day and which of a few days it's on */

	const showformat = "Mon 15:04"
	ts := parseStoredDate(stamp)
	return fmt.Sprintf(`<span title="%v">%v</span>`, stamp, ts.Format(showformat))
}

func showEBC(w http.ResponseWriter, r *http.Request) {

	const email_icon = "&#9993;"

	/*
		sqlx := `SELECT ebclaims.rowid,ebclaims.EntrantID,RiderName,PillionName,ebclaims.BonusID,xbonus.BriefDesc
			    ,OdoReading,ClaimTime,ExtraField,StrictOK,xphoto.Image,Notes,Flags,TeamID
			    ,ebclaims.AttachmentTime As PhotoTS, ebclaims.DateTime As EmailTS,ebclaims.LoggedAt,ebclaims.Subject
			    ,xbonus.Points,xbonus.AskPoints,xbonus.RestMinutes,xbonus.AskMinutes,xbonus.Image as BImage,Question,Answer
			     FROM ebclaims LEFT JOIN entrants ON ebclaims.EntrantID=entrants.EntrantID
			     LEFT JOIN (SELECT BonusID,BriefDesc,Notes,IfNull(Flags,'') AS Flags,Points,AskPoints,RestMinutes,AskMinutes,
			    IfNull(Image,'') AS Image,IfNull(Question,'') AS Question,IfNull(Answer,'') AS Answer FROM bonuses
			     ) AS xbonus
			     ON ebclaims.BonusID=xbonus.BonusID  LEFT JOIN "
			     (SELECT EmailID,Group_concat(Image) As Image from ebcphotos GROUP BY EmailID) AS xphoto
				 ON ebclaims.EmailID=xphoto.EmailID WHERE Processed=0 ORDER BY Decision DESC,FinalTime;`

	*/

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>`)
	fmt.Fprintf(w, `<style>%s</style>`, css)
	fmt.Fprintf(w, `<script>%s</script>`, script)

	var ebc ElectronicBonusClaim
	err = rows.Scan(&ebc.EntrantID, &ebc.Bonusid, &ebc.OdoReading, &ebc.ClaimTime, &ebc.Subject, &ebc.ExtraField, &ebc.AttachmentTime, &ebc.DateTime, &ebc.FirstTime, &ebc.FinalTime, &ebc.EmailID)
	checkerr(err)

	team := getStringFromDB("SELECT RiderName FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "***")
	x := getStringFromDB("SELECT ifnull(PillionName,'') FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "")
	if x != "" {
		team += " &amp; " + x
	}
	teamneeded := x != ""
	if !teamneeded {
		teamneeded = getIntegerFromDB("SELECT TeamID FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), 0) > 0
	}
	fmt.Fprint(w, `<div>`)
	fmt.Fprintf(w, `Entrant <span class="bold">%v %v</span>`, ebc.EntrantID, team)
	x = getStringFromDB("SELECT BriefDesc FROM bonuses WHERE BonusID='"+ebc.Bonusid+"'", ebc.Bonusid)
	fmt.Fprintf(w, ` Bonus <span class="bold">%v" %v</span>`, ebc.Bonusid, x)
	fmt.Fprint(w, ` Claimed @ `)
	evidence := "Photo: " + ebc.AttachmentTime + "\n"
	evidence += "Claim: " + ebc.ClaimTime + "\n"
	evidence += "Email: " + ebc.DateTime + "\n"
	evidence += "Recvd: " + ebc.FinalTime + "\n"
	fmt.Fprintf(w, `<span class="bold" title="%v" onclick="alert(this.getAttribute('title'))">%v, %v %v</span>`, evidence, ebc.OdoReading, logtime(ebc.ClaimTime), email_icon)
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
	fmt.Fprint(w, `</div>`) // row

	fmt.Fprint(w, `<div>`)

	fmt.Fprintf(w, `<input type="button" data-result="-1" onclick="closeEBC(this)" class="closebutton" value="%v">`, CS.CloseEBCUndecided)
	fmt.Fprintf(w, `<input type="button" data-result="0" onclick="closeEBC(this)" class="closebutton" value="%v">`, CS.CloseEBC[0])
	x = "***"
	fmt.Fprintf(w, `<input type="text" id="judgesnotes" name="judgesnotes" class="judgesnotes" value="%v">`, x)
	fmt.Fprint(w, `</div>`)
	fmt.Fprint(w, `<div>`)
	for i := 1; i < 10; i++ {
		fmt.Fprintf(w, `<input type="button" data-result="%v" onclick="closeEBC(this)" class="closebutton" value="%v">`, i, CS.CloseEBC[i])
	}
	fmt.Fprint(w, `</div>`)
	showPhotos(w, ebc.EmailID, ebc.Bonusid)

}

func showPhotos(w http.ResponseWriter, emailid int, BonusID string) {

	const maximg = 3
	bimg := strings.ReplaceAll(filepath.Join(CS.ImgBonusFolder, filepath.Base(getStringFromDB("SELECT ifnull(Image,'') FROM bonuses WHERE BonusID='"+BonusID+"'", ""))), `\`, `/`)

	sqlx := "SELECT Image FROM ebcphotos WHERE EmailID=" + strconv.Itoa(emailid)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	fmt.Fprint(w, `<div class="imgcomparediv">`)

	// "if(this.width=='50%')this.width='100%';else this.width='50%'"
	fmt.Fprint(w, `<div class="ebcimgdiv" id="ebcimgdiv" onclick="cycleImgSize(this)">`)
	var showimg [maximg]string

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
	fmt.Fprintf(w, `<img id="imgdivimg" alt="*" src="%v"/>`, showimg[0])

	fmt.Fprint(w, `<div id="imgdivs">`)

	for ix = 1; ix < maximg; ix++ {
		if showimg[ix] != "" {
			fmt.Fprintf(w, `<img src="%v" alt="*" onclick="swapimg(this)">`, showimg[ix])
		}
	}
	fmt.Fprint(w, `</div>`) // imgdivs
	fmt.Fprint(w, `</div>`) // ebcimgdiv

	fmt.Fprint(w, `<div class="bonusimgdiv" id="bonusimgdiv">`)
	fmt.Fprintf(w, `<img src="%v">`, bimg)
	fmt.Fprint(w, `</div>`)

	fmt.Fprint(w, `</div>`)
}
