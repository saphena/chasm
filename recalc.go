package main

import (
	"fmt"
	"slices"
	"strconv"
)

func checkerr(err error) {
	panic(err)
}

func recalc_all() {

	_, err := DBH.Exec("BEGIN TRANSACTION")
	checkerr(err)
	defer DBH.Exec("COMMIT")
	sqlx := "SELECT EntrantID FROM entrants ORDER BY EntrantID"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var entrant int
		rows.Scan(&entrant)
		recalc_scorecard(entrant)
	}

}

type ScorecardBonusDetail struct {
	BonusID     string
	BriefDesc   string
	Compulsory  bool
	Scored      bool
	Points      int
	RestMinutes int
	CatValue    [NumCategoryAxes]int
}
type ClaimedBonus struct {
	BonusID        string
	Decision       int
	Points         int
	RestMinutes    int
	QuestionScored bool
	MultiplyLast   bool // The points value is used to multiply the value of last bonus claimed
}

const ClaimDecision_ClaimExcluded = 9
const ClaimDecision_GoodClaim = 0
const PointsCalcMethod_MultiplyLast = 2
const NumCategoryAxes = 9

func build_scorecardBonusArray(CurrentLeg int) []ScorecardBonusDetail {

	// Build array of all bonuses for use with this scorecard

	var B []ScorecardBonusDetail

	sqlx := "SELECT BonusID, BriefDesc, Points, RestMinutes"
	for i := 1; i <= NumCategoryAxes; i++ {
		sqlx += ", Cat" + strconv.Itoa(i)
	}
	sqlx += " FROM bonuses"
	sqlx += " WHERE Leg=0 OR <=" + strconv.Itoa(CurrentLeg)
	sqlx += " ORDER BY BonusID"
}
func build_bonusclaim_array(entrant int) []ClaimedBonus {

	// Build list of bonuses claimed

	var B []ClaimedBonus

	sqlx := "SELECT claims.BonusID, Decision, claims.Points, claims.RestMinutes, QuestionScored, ifnull(AskPoints,0)"
	sqlx += " FROM claims"
	sqlx += " LEFT JOIN bonuses ON claims.BonusID=bonuses.BonusID"
	sqlx += " WHERE EntrantID=" + strconv.Itoa(entrant)
	sqlx += " AND Decision >= " + strconv.Itoa(ClaimDecision_GoodClaim) // Decided claim
	sqlx += " AND Decision != " + strconv.Itoa(ClaimDecision_ClaimExcluded)
	sqlx += " ORDER BY ClaimTime, OdoReading"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var bonus ClaimedBonus
		var qs, ap int
		rows.Scan(&bonus.BonusID, &bonus.Decision, &bonus.Points, &bonus.RestMinutes, &qs, &ap)
		bonus.QuestionScored = qs != 0
		bonus.MultiplyLast = ap == PointsCalcMethod_MultiplyLast
		B = append(B, bonus)

	}
	return B
}

type catFields struct {
	catCounts     map[int]int
	sameCatCount  int
	sameCatPoints int
	lastCatScored int
}

const NumberOfAxes = 9

func build_catCountsArray() []catFields {

	var cf [NumberOfAxes + 1]catFields // entry 0 will be ignored. Axes are numbered 1..9
	return []catFields{}
}

// This recalculates the value of the specified scorecard using as
// input the relevant claims records. The results are updated totals
// and score explanation.
func recalc_scorecard(entrant int) {

	const Leg = 1
	ScorecardBonuses := build_scorecardBonusArray(Leg)
	BonusesClaimed := build_bonusclaim_array(entrant)
	fmt.Printf("BonusesClaimed == %v\n", BonusesClaimed)
	for _, BC := range BonusesClaimed {

		// Need to flag the bonus as having been scored
		bix := slices.IndexFunc(ScorecardBonuses, func(c ScorecardBonusDetail) bool { return c.BonusID == BC.BonusID })
		ScorecardBonuses[bix].Scored = BC.Decision == ClaimDecision_GoodClaim
	}

}
