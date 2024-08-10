package main

import (
	"fmt"
	"html"
	"log"
	"reflect"
	"slices"
	"strconv"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func recalc_all() {

	_, err := DBH.Exec("BEGIN TRANSACTION")
	checkerr(err)
	defer DBH.Exec("COMMIT")
	sqlx := "SELECT EntrantID FROM entrants ORDER BY EntrantID"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	n := 2
	for rows.Next() {
		var entrant int
		rows.Scan(&entrant)
		recalc_scorecard(entrant)
		n--
		if n < 1 {
			break
		}
	}

}

type ScorecardBonusDetail struct {
	Bonusid     string
	BriefDesc   string
	Compulsory  bool
	Points      int
	RestMinutes int
	CatValue    [NumCategoryAxes]int
	Scored      bool // Don't ry to Scan this
}

var ScorecardBonuses []ScorecardBonusDetail

type ClaimedBonus struct {
	Bonusid          string
	BonusScorecardIX int
	Decision         int
	Points           int
	RestMinutes      int
	QuestionScored   bool
	MultiplyLast     bool // The points value is used to multiply the value of last bonus claimed
}

type ClaimedBonusMap = []ClaimedBonus

var BonusesClaimed ClaimedBonusMap

type CompoundRule struct {
	Ruleid      int
	Axis        int
	Cat         int
	Method      int
	Target      int
	Value       int
	Min         int
	PointsMults int
	Power       int
	Ruletype    int
	Triggered   bool
}

var CompoundRules []CompoundRule

type ScorexLine struct {
	IsValidLine bool
	Code        string
	Desc        string
	PointsDesc  string
	Points      int
}

const ClaimDecision_ClaimExcluded = 9
const ClaimDecision_GoodClaim = 0
const PointsCalcMethod_MultiplyLast = 2
const NumCategoryAxes = 9
const CompoundRuleType_SimpleSequence = 4
const CatPointsMults_Points = 0

func build_compoundRuleArray(CurrentLeg int) []CompoundRule {

	var res []CompoundRule
	sqlx := "SELECT rowid AS id,Axis,Cat,NMethod,ModBonus,NMin,PointsMults,NPower,Ruletype"
	sqlx += " FROM catcompound WHERE Leg=0 OR Leg=" + strconv.Itoa(CurrentLeg)
	sqlx += " ORDER BY Axis,NMin DESC"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var cr CompoundRule
		rows.Scan(&cr.Ruleid, &cr.Axis, &cr.Cat, &cr.Method, &cr.Target, &cr.Min, &cr.PointsMults, &cr.Power, &cr.Ruletype)
		res = append(res, cr)
	}
	return res
}

func build_scorecardBonusArray(CurrentLeg int) []ScorecardBonusDetail {

	// Build array of all bonuses for use with this scorecard

	var res []ScorecardBonusDetail
	var b ScorecardBonusDetail

	s := reflect.ValueOf(&b).Elem()
	numCols := s.NumField() - 1 + NumCategoryAxes - 1
	//	log.Printf("numCols is %v\n", numCols)
	columns := make([]interface{}, numCols)
	for i := 0; i < s.NumField()-1; i++ { // -1 limit to avoid Scored
		//		log.Printf("Loading %v\n", i)
		field := s.Field(i)

		//log.Printf("Field(%v).Kind is %v\n", i, field.Kind())
		if field.Kind() == reflect.Array {
			//			k := field.Len()
			//			L := len(columns)
			for j := 0; j < field.Len(); j++ {
				//				log.Printf("array index %v  < Len(field)=%v  i+j=%v < Len(columns)=%v\n", j, k, i+j, L)
				columns[i+j] = field.Index(j).Addr().Interface()
				//				log.Println("ok")
			}
		} else {
			//			log.Printf("index %v\n", i)
			columns[i] = field.Addr().Interface()
		}
	}

	//	log.Println("Got here")
	sqlx := "SELECT Bonusid, BriefDesc, Compulsory, Points, RestMinutes"
	for i := 1; i <= NumCategoryAxes; i++ {
		sqlx += ", Cat" + strconv.Itoa(i)
	}
	sqlx += " FROM bonuses"
	sqlx += " WHERE Leg=0 OR Leg<=" + strconv.Itoa(CurrentLeg)
	sqlx += " ORDER BY Bonusid"
	//	log.Println(sqlx)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(columns...)
		checkerr(err)
		res = append(res, b)
	}
	//log.Printf("BonusArray is %v\n", res)
	return res

}
func build_bonusclaim_array(entrant int) ClaimedBonusMap {

	// Build list of bonuses claimed

	Bid := make(map[string]int)
	B := make(ClaimedBonusMap, 0)

	sqlx := "SELECT claims.Bonusid, Decision, claims.Points, claims.RestMinutes, claims.QuestionAnswered, ifnull(claims.AskPoints,0)"
	sqlx += " FROM claims"
	sqlx += " LEFT JOIN bonuses ON claims.Bonusid=bonuses.Bonusid"
	sqlx += " WHERE EntrantID=" + strconv.Itoa(entrant)
	sqlx += " AND Decision >= " + strconv.Itoa(ClaimDecision_GoodClaim) // Decided claim
	sqlx += " AND Decision != " + strconv.Itoa(ClaimDecision_ClaimExcluded)
	sqlx += " ORDER BY ClaimTime, OdoReading"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	bix := 0
	for rows.Next() {
		var bonus ClaimedBonus
		var qs, ap int
		rows.Scan(&bonus.Bonusid, &bonus.Decision, &bonus.Points, &bonus.RestMinutes, &qs, &ap)
		bonus.QuestionScored = qs != 0
		bonus.MultiplyLast = ap == PointsCalcMethod_MultiplyLast
		if bonus.Decision == ClaimDecision_ClaimExcluded {
			log.Printf("Excluding %v\n", bonus.Bonusid)
			continue
		}
		log.Printf("Including %v\n", bonus.Bonusid)
		ix, ok := Bid[bonus.Bonusid]
		if ok { // Supersede the earlier claim
			B[ix] = bonus
		} else {
			B = append(B, bonus)
			Bid[bonus.Bonusid] = bix
			bix++
		}

	}
	return B
}

type catFields struct {
	catCounts     map[int]int
	sameCatCount  int
	sameCatPoints int
	lastCatScored int
}

var CatCounts []catFields

type catLabel struct {
	Axis      int
	Cat       int
	BriefDesc string
}

const NumberOfAxes = 9

func build_emptyCatCountsArray() []catFields {

	var res []catFields

	res = make([]catFields, 0)
	for i := 0; i <= NumberOfAxes; i++ {
		var cf catFields
		cf.catCounts = make(map[int]int, 0)
		res = append(res, cf)
	}
	return res
}

func build_axisLabels() []string {

	sqlx := "SELECT Cat1Label"
	for i := 2; i <= NumberOfAxes; i++ {
		sqlx += ",Cat" + strconv.Itoa(i) + "Label"
	}
	sqlx += " FROM rallyparams"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	var res []string
	var s string
	for rows.Next() {
		err = rows.Scan(&s)
		checkerr(err)
		res = append(res, s)
	}
	return res

}

func checkApplySequences(BC ClaimedBonus, LastBonusClaimed ClaimedBonus) ScorexLine {

	var sx ScorexLine

	// Check for sequence bonus
	for _, CR := range CompoundRules {
		if CR.Ruletype != CompoundRuleType_SimpleSequence {
			continue
		}

		if LastBonusClaimed.Bonusid == "" {
			continue
		}

		Cat := ScorecardBonuses[BC.BonusScorecardIX].CatValue[CR.Axis-1]
		LastCat := ScorecardBonuses[LastBonusClaimed.BonusScorecardIX].CatValue[CR.Axis-1]
		// If this bonus is in the same category as the last one
		if Cat == LastCat {
			// Still building the sequence so

			continue
		}

		if CatCounts[CR.Axis].sameCatCount < CR.Min {
			continue
		}

		// Trigger sequential bonus

		log.Println("Triggering sequential bonus")
		const sequential_bonus_symbol = "&#8752;"
		//const atleast_symbol = "&ge;"
		//const checkmark_symbol = "&#x2713;"
		sqlx := fmt.Sprintf("SELECT BriefDesc FROM categories WHERE Axis=%d AND Cat=%d", CR.Axis, LastCat)
		defaultValue := fmt.Sprintf("%d/%d", CR.Axis, LastCat)
		BonusDesc := fmt.Sprintf("%s %d X %s", sequential_bonus_symbol, CatCounts[CR.Axis].sameCatCount, getStringFromDB(sqlx, defaultValue))
		//BonusDesc += fmt.Sprintf(" (%d %s %d)", , atleast_symbol, CR.Min)

		PointsDesc := ""
		ExtraBonusPoints := 0
		if CR.PointsMults == CatPointsMults_Points { // Result is specified number of points
			ExtraBonusPoints = CR.Power
		} else { // Result is sequence length * multiplier
			ExtraBonusPoints = CatCounts[CR.Axis].sameCatPoints * CR.Power
			if CR.Power != 1 && CR.Power != 0 {
				PointsDesc = fmt.Sprintf(" (+ %dx%d)", CatCounts[CR.Axis].sameCatPoints, CR.Power)
			}
		}
		sx.Desc = BonusDesc
		sx.PointsDesc = PointsDesc
		sx.Points = ExtraBonusPoints
		sx.IsValidLine = true
		break // Only apply the first matching rule

	}

	return sx
}

// This recalculates the value of the specified scorecard using as
// input the relevant claims records. The results are updated totals
// and score explanation.
func recalc_scorecard(entrant int) {

	const Leg = 1

	log.Printf("recalc for %v\n", entrant)
	ScorecardBonuses = build_scorecardBonusArray(Leg)

	BonusesClaimed = build_bonusclaim_array(entrant)

	CompoundRules = build_compoundRuleArray(Leg)

	CatCounts = build_emptyCatCountsArray()

	var sx ScorexLine
	TotalPoints := 0
	var Scorex []ScorexLine

	log.Printf("BonusesClaimed == %v\n", BonusesClaimed)

	var LastBonusClaimed ClaimedBonus
	for _, BC := range BonusesClaimed {

		// ClaimExcluded means ignore it, treat it as if it didn't exist
		// This might need to be a switchable response
		if BC.Decision == ClaimDecision_ClaimExcluded {
			continue
		}

		// Need to flag the bonus as having been scored
		BC.BonusScorecardIX = slices.IndexFunc(ScorecardBonuses, func(c ScorecardBonusDetail) bool { return c.Bonusid == BC.Bonusid })

		SB := ScorecardBonuses[BC.BonusScorecardIX] // Convenient shorthand

		//log.Printf("ScorecardIX = %v\n", BC.BonusScorecardIX)
		SB.Scored = BC.Decision == ClaimDecision_GoodClaim // Only good claims count against "must score" flag

		sx = checkApplySequences(BC, LastBonusClaimed)
		if sx.IsValidLine {
			TotalPoints += sx.Points
			Scorex = append(Scorex, sx)
		}

		if BC.Decision != ClaimDecision_GoodClaim {
			var sx ScorexLine

			// Firstly, let's zap any sequence in progress
			for i := 1; i <= NumCategoryAxes; i++ {
				CatCounts[i].sameCatCount = 0
				CatCounts[i].sameCatPoints = 0
				CatCounts[i].lastCatScored = -1
			}
			sx.IsValidLine = true
			sx.Code = SB.Bonusid
			sx.Desc = fmt.Sprintf("%v<br>CLAIM REJECTED - %v", SB.BriefDesc, "JUST BECAUSE")

			Scorex = append(Scorex, sx)
			continue
		}
		updateCatCounts(SB, BC.Points)
		LastBonusClaimed = BC

		BasicPoints := BC.Points

		TotalPoints += BasicPoints
		var sx ScorexLine

		sx.Code = SB.Bonusid
		sx.Desc = SB.BriefDesc
		sx.Points = BasicPoints
		Scorex = append(Scorex, sx)
		//log.Printf("counts=%v\n", CatCounts)
	}

	// Final check for a sequence
	var BC ClaimedBonus
	sx = checkApplySequences(BC, LastBonusClaimed)
	if sx.IsValidLine {
		TotalPoints += sx.Points
		Scorex = append(Scorex, sx)
	}
	//log.Printf("Scorex == %v\n", Scorex)
	for x := range Scorex {
		log.Printf("%-3s %-20s %-10s %v\n", Scorex[x].Code, html.UnescapeString(Scorex[x].Desc), html.UnescapeString(Scorex[x].PointsDesc), Scorex[x].Points)
	}
	log.Printf("Total points is %d\n", TotalPoints)
}

func updateCatCounts(BS ScorecardBonusDetail, Points int) {

	for i := 1; i <= NumCategoryAxes; i++ {

		cat := BS.CatValue[i-1]

		if cat <= 0 {
			CatCounts[i].sameCatCount = 0
			CatCounts[i].sameCatPoints = 0
			CatCounts[i].lastCatScored = cat
			continue
		} else if cat == CatCounts[i].lastCatScored {
			CatCounts[i].sameCatCount++
			CatCounts[i].sameCatPoints += Points
		} else {
			CatCounts[i].sameCatCount = 1
			CatCounts[i].sameCatPoints = Points
			CatCounts[i].lastCatScored = cat
		}
		_, ok := CatCounts[i].catCounts[cat]

		if ok {
			CatCounts[i].catCounts[cat]++
		} else {
			CatCounts[i].catCounts[cat] = 1
		}

		// Now accrue overall axis totals
		_, ok = CatCounts[i].catCounts[0]
		if ok {
			CatCounts[i].catCounts[0]++
		} else {
			CatCounts[i].catCounts[0] = 1
		}
	}
}
