package main

import (
	"strconv"
)

type RanktabRecord struct {
	EntrantID      int
	TeamID         int
	TotalPoints    int
	CorrectedMiles int
	PPM            float64
	Rank           int
}

type ByPoints []RanktabRecord
type ByPPM []RanktabRecord
type ByMiles []RanktabRecord

func (a ByPoints) Len() int           { return len(a) }
func (a ByPoints) Less(i, j int) bool { return a[i].TotalPoints < a[j].TotalPoints }
func (a ByPoints) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPPM) Len() int              { return len(a) }
func (a ByPPM) Less(i, j int) bool    { return a[i].PPM < a[j].PPM }
func (a ByPPM) Swap(i, j int)         { a[i], a[j] = a[j], a[i] }
func (a ByMiles) Len() int            { return len(a) }
func (a ByMiles) Less(i, j int) bool  { return a[i].CorrectedMiles < a[j].CorrectedMiles }
func (a ByMiles) Swap(i, j int)       { a[i], a[j] = a[j], a[i] }

func loadRankTable(sqlx string) []RanktabRecord {

	//fmt.Println(sqlx)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	ranktab := make([]RanktabRecord, 0)
	for rows.Next() {
		var rt RanktabRecord
		err = rows.Scan(&rt.EntrantID, &rt.TeamID, &rt.TotalPoints, &rt.CorrectedMiles, &rt.PPM, &rt.Rank)
		checkerr(err)
		ranktab = append(ranktab, rt)
	}
	return ranktab

}
func rankEntrants(intransaction bool) {

	var sqlx string
	var err error

	sqlx = "SELECT EntrantID,TeamID,TotalPoints,CorrectedMiles,IfNull((TotalPoints*1.0)/CorrectedMiles,0) AS PPM,0 AS Rank FROM entrants WHERE EntrantStatus = "
	sqlx += strconv.Itoa(EntrantFinisher)
	sqlx += " AND TeamID > 0"
	sqlx += " ORDER BY TeamID"
	if CS.RallyRankEfficiency {
		sqlx += ",PPM"
	} else {
		sqlx += ",TotalPoints"
	}
	if CS.RallyTeamMethod == RankTeamsHighest {
		sqlx += " DESC"
	}

	ranktab := loadRankTable(sqlx)

	// Presort Team records
	LastTeam := -1
	LastTeamPoints := 0
	LastTeamMiles := 0
	stmt, err := DBH.Prepare("UPDATE entrants SET CorrectedMiles=?,TotalPoints=? WHERE TeamID=?")
	checkerr(err)
	defer stmt.Close()
	for _, rtr := range ranktab {

		if rtr.TeamID > 0 {
			if LastTeam != rtr.TeamID {
				LastTeam = rtr.TeamID
				LastTeamPoints = rtr.TotalPoints
				LastTeamMiles = rtr.CorrectedMiles

				_, err = stmt.Exec(LastTeamMiles, LastTeamPoints, LastTeam)
				checkerr(err)

			}
		}
	}
	stmt.Close()

	sqlx = "SELECT EntrantID,TeamID,TotalPoints,CorrectedMiles,IfNull((TotalPoints*1.0)/CorrectedMiles,0) AS PPM,0 AS Rank FROM entrants WHERE EntrantStatus = "
	sqlx += strconv.Itoa(EntrantFinisher)
	sqlx += " ORDER BY TeamID"
	if CS.RallyRankEfficiency {
		sqlx += ",PPM"
	} else {
		sqlx += ",TotalPoints"
	}
	sqlx += " DESC"

	ranktab = loadRankTable(sqlx)

	if !intransaction {
		_, err = DBH.Exec("BEGIN IMMEDIATE TRANSACTION")
		checkerr(err)
		defer DBH.Exec("COMMIT")

	}

	_, err = DBH.Exec("UPDATE entrants SET FinishPosition=0")
	checkerr(err)

	finishPos := 0
	lastTotalPoints := -1
	lastPPM := -1.0
	lastCorrectedMiles := -1
	incN := 1
	LastTeam = -1

	sqlx = "UPDATE entrants SET FinishPosition=? WHERE EntrantID=?"

	stmt, err = DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()

	for _, rr := range ranktab {

		SplitNeeded := CS.RallySplitTies && rr.CorrectedMiles != lastCorrectedMiles

		//fmt.Printf("%v [%v] = %v (%v)\n", rr.EntrantID, rr.TeamID, rr.TotalPoints, rr.Rank)
		if SplitNeeded {
			if CS.RallyRankEfficiency {
				SplitNeeded = rr.PPM == lastPPM
			} else {
				SplitNeeded = rr.TotalPoints == lastTotalPoints
			}
		}

		if SplitNeeded { // This is the same as the last

			if rr.TeamID != LastTeam { // TeamID = 0 or Team not last team
				finishPos += incN
				incN = 1
			} else {
				incN++
			}

		} else {
			if rr.TeamID == LastTeam && CS.RallyTeamMethod != RankTeamsAsIndividuals {

				// All team members assigned the same rank

			} else {
				finishPos += incN
				incN = 1
			}
		}
		if rr.TeamID > 0 {
			LastTeam = rr.TeamID
		}

		_, err = stmt.Exec(finishPos, rr.EntrantID)
		checkerr(err)
	}

}
