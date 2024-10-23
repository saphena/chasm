package main

import (
	"fmt"
	"net/http"
	"strconv"
)

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
	
	 FROM ebclaims WHERE Processed=0 AND rowid=` + claimid

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	if !rows.Next() {
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<style>%s</style>`, css)
	fmt.Fprintf(w, `<script>%s</script>`, script)

	var ebc ElectronicBonusClaim
	err = rows.Scan(&ebc.EntrantID, &ebc.Bonusid, &ebc.OdoReading, &ebc.ClaimTime, &ebc.Subject, &ebc.ExtraField, &ebc.AttachmentTime, &ebc.DateTime, &ebc.FirstTime, &ebc.FinalTime)
	checkerr(err)

	team := getStringFromDB("SELECT RiderName FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "***")
	x := getStringFromDB("SELECT ifnull(PillionName,'') FROM entrants WHERE EntrantID="+strconv.Itoa(ebc.EntrantID), "")
	if x != "" {
		team += " &amp; " + x
	}

	fmt.Fprintf(w, `<fieldset>%v</fieldset>`, team)
}
