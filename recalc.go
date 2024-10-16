package main

/*

Bonus is multiplier is implemented using the flag held on the bonus record, not the flag held on the claim record

Axis is in the range 1..NumCategoryAxes. Axis 0 is used for aggregate totals.const
Cat is a non-zero integer but -1 when used to access arrays
*/

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"math"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Compound category rules

const CAT_ModifyAxisScore = 0
const CAT_ModifyBonusScore = 1
const CAT_ResultMults = 1
const CAT_ResultPoints = 0
const CAT_NumBonusesPerCatMethod = 0
const CAT_NumNZCatsPerAxisMethod = 1
const CAT_OrdinaryScoringRule = 0
const CAT_DNF_Unless_Triggered = 1
const CAT_DNF_If_Triggered = 2
const CAT_PlaceholderRule = 3
const CAT_OrdinaryScoringSequence = 4

const checkmark_symbol = "&#x2713;"
const sequential_bonus_symbol = "&#8752;"
const compoundrule_symbol = "&#9783;"
const penalty_symbol = "&#9785;"

const ClaimDecision_ClaimExcluded = 9
const ClaimDecision_GoodClaim = 0

// Bonus points calculation method
const PointsCalcMethod_MultiplyLast = 2

// Combo ScoreMethod values
const ScoreMethodPoints = 0
const ScoreMethodMults = 1

const NumCategoryAxes = 9

const NumberOfAxes = 9

const EntrantDNS = 0
const EntrantOK = 1
const EntrantFinisher = 8
const EntrantDNF = 3

var EntrantStatusLits = map[int]string{EntrantDNS: "DNS", EntrantOK: "ok", EntrantFinisher: "Finisher", EntrantDNF: "DNF"}

const KmsPerMile = float64(1.60934)

type ScorecardBonusDetail struct {
	Bonusid      string
	BriefDesc    string
	Compulsory   bool
	Points       int
	AskPoints    int // A value of 2 here indicates MultiplyLast
	RestMinutes  int
	CatValue     [NumCategoryAxes]int
	Scored       bool // Don't ry to Scan
	MultiplyLast bool // The points value is used to multiply the value of last bonus claimed
}

var ScorecardBonuses []ScorecardBonusDetail

type ClaimedBonus struct {
	Bonusid          string
	BonusScorecardIX int
	Decision         int
	Points           int
	RestMinutes      int
	QuestionScored   bool
	OdoReading       int
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

type ComboBonus struct {
	Comboid     string
	BriefDesc   string
	ScoreMethod int
	MinTicks    int
	PointsList  string
	BonusList   string
	Compulsory  bool
	Cat         [NumCategoryAxes]int
	Points      []int
	Bonuses     []string
	Scored      bool
}

var ComboBonuses []ComboBonus

type ScorexLine struct {
	IsValidLine bool
	Code        string
	Desc        string
	PointsDesc  string
	Points      int
}

type catFields struct {
	CatCounts     map[int]int
	SameCatCount  int
	SameCatPoints int
	LastCatScored int
}

type axisCounts map[int]catFields

var AxisCounts axisCounts
var AxisLabels []string

type ScorexParams struct {
	RiderName      string
	PillionName    string
	DistanceRidden int
}

type RejectReason struct {
	Code      int
	BriefDesc string
	Action    int
	Param     string
}

type RejectReasonsMap map[int]RejectReason

var RejectReasons RejectReasonsMap

const myTimestamp = "2006-01-02T15:04"  // Stored format
const myTimestampX = "2006-01-02 15:04" // Display format

// Time penalties
const (
	TimeSpecDatetime   = 0
	TimeSpecRallyDNF   = 1
	TimeSpecEntrantDNF = 2
)

// TPM = Time Penalty Method
const (
	TPM_FixedPoints  = 0
	TPM_FixedMult    = 1
	TPM_PointsPerMin = 2
	TPM_MultPerMin   = 3
)

type TimePenalty struct {
	TimeSpec       int
	PenaltyStartX  string
	PenaltyFinishX string
	PenaltyMethod  int
	PenaltyFactor  int
	PenaltyStart   time.Time
	PenaltyFinish  time.Time
}

var TimePenalties []TimePenalty

const (
	MMM_FixedPoints   = 0
	MMM_Multipliers   = 1
	MMM_PointsPerMile = 2
)

var RallyParametersLoaded bool

func build_axisLabels() []string {

	sqlx := "SELECT IfNull(Cat1Label,'')"
	for i := 2; i <= NumberOfAxes; i++ {
		sqlx += ",IfNull(Cat" + strconv.Itoa(i) + "Label,'')"
	}
	sqlx += " FROM rallyparams"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	var res []string
	s := make([]string, NumberOfAxes)
	for rows.Next() {
		err = rows.Scan(&s[0], &s[1], &s[2], &s[3], &s[4], &s[5], &s[6], &s[7], &s[8])
		checkerr(err)
	}
	res = append(res, s...)
	//log.Printf("AxisLabels = %v\n", res)
	return res

}

// Build list of bonuses claimed
func build_bonusclaim_array(entrant int) ClaimedBonusMap {

	Bid := make(map[string]int)
	B := make(ClaimedBonusMap, 0)

	sqlx := "SELECT claims.Bonusid, Decision, claims.Points, claims.RestMinutes, claims.QuestionAnswered, claims.OdoReading"
	sqlx += " FROM claims"
	sqlx += " LEFT JOIN bonuses ON claims.Bonusid=bonuses.Bonusid"
	sqlx += " WHERE EntrantID=" + strconv.Itoa(entrant)
	sqlx += " AND Decision >= " + strconv.Itoa(ClaimDecision_GoodClaim) // Decided claim
	//sqlx += " AND Decision != " + strconv.Itoa(ClaimDecision_ClaimExcluded)
	sqlx += " ORDER BY ClaimTime, OdoReading"

	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	bix := 0
	for rows.Next() {
		var bonus ClaimedBonus
		var qs int
		rows.Scan(&bonus.Bonusid, &bonus.Decision, &bonus.Points, &bonus.RestMinutes, &qs, &bonus.OdoReading)
		bonus.QuestionScored = qs != 0
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

func build_compoundRuleArray(CurrentLeg int) []CompoundRule {

	var res []CompoundRule
	sqlx := "SELECT rowid AS id,IfNull(Axis,1),IfNull(Cat,0),IfNull(NMethod,0),IfNull(ModBonus,0),IfNull(NMin,1),IfNull(PointsMults,0),IfNull(NPower,0),IfNull(Ruletype,0)"
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

func build_emptyAxisCountsArray() axisCounts {

	res := make(axisCounts, 0)
	for i := 1; i <= NumberOfAxes; i++ {
		var cf catFields
		cf.CatCounts = make(map[int]int, 0)
		res[i] = cf
	}
	return res
}

// Build array of all bonuses for use with this scorecard
func build_scorecardBonusArray(CurrentLeg int) []ScorecardBonusDetail {

	var res []ScorecardBonusDetail
	var b ScorecardBonusDetail

	s := reflect.ValueOf(&b).Elem()
	numCols := s.NumField() - 2 + NumCategoryAxes - 1
	columns := make([]interface{}, numCols)
	for i := 0; i < s.NumField()-2; i++ { // -2 limit to avoid Scored, MultiplyLast
		field := s.Field(i)

		if field.Kind() == reflect.Array {
			for j := 0; j < field.Len(); j++ {
				columns[i+j] = field.Index(j).Addr().Interface()
			}
		} else {
			columns[i] = field.Addr().Interface()
		}
	}

	//	log.Println("Got here")
	sqlx := "SELECT Bonusid, BriefDesc, Compulsory, Points, AskPoints, RestMinutes"
	for i := 1; i <= NumCategoryAxes; i++ {
		sqlx += ", Cat" + strconv.Itoa(i)
	}
	sqlx += " FROM bonuses"
	sqlx += " WHERE Leg=0 OR Leg<=" + strconv.Itoa(CurrentLeg)
	sqlx += " ORDER BY Bonusid"
	//log.Println(sqlx)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(columns...)
		checkerr(err)
		b.MultiplyLast = b.AskPoints == PointsCalcMethod_MultiplyLast
		res = append(res, b)
	}
	//log.Printf("BonusArray is %v\n", res)
	return res

}

func build_timePenaltyArray() []TimePenalty {

	//	const currentLeg = "0"

	rallyFinishtimex := getStringFromDB("SELECT FinishTime FROM rallyparams", "")

	rft, err := time.Parse(myTimestamp, rallyFinishtimex)
	checkerr(err)

	sqlx := "SELECT TimeSpec,PenaltyStart,PenaltyFinish,PenaltyMethod,PenaltyFactor FROM timepenalties"
	//sqlx += " WHERE Leg=0 OR Leg=" + currentLeg
	sqlx += " ORDER BY PenaltyStart,PenaltyFinish"

	res := make([]TimePenalty, 0)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var tp TimePenalty
		err = rows.Scan(&tp.TimeSpec, &tp.PenaltyStartX, &tp.PenaltyFinishX, &tp.PenaltyMethod, &tp.PenaltyFactor)
		checkerr(err)

		switch tp.TimeSpec {

		case TimeSpecRallyDNF:
			ps, _ := strconv.Atoi(tp.PenaltyStartX)
			tp.PenaltyStart = rft.Add(time.Minute * time.Duration(0-ps))
			pf, _ := strconv.Atoi(tp.PenaltyFinishX)
			tp.PenaltyFinish = rft.Add(time.Minute * time.Duration(0-pf))

		case TimeSpecDatetime:
			tp.PenaltyStart, _ = time.Parse(myTimestamp, tp.PenaltyStartX)
			tp.PenaltyFinish, _ = time.Parse(myTimestamp, tp.PenaltyFinishX)
		}

		res = append(res, tp)
	}
	return res
}

func BuildRallyParameters(Leg int) {

	if RallyParametersLoaded {
		return
	}

	ScorecardBonuses = build_scorecardBonusArray(Leg)
	CompoundRules = build_compoundRuleArray(Leg)
	AxisLabels = build_axisLabels()
	ComboBonuses = loadCombos("")
	RejectReasons = loadRejectReasons()
	TimePenalties = build_timePenaltyArray()

	RallyParametersLoaded = true

}

func calcEntrantStatus(Miles int) (ScorexLine, int) {

	const DNF_icon = "&#9760;"

	var sx ScorexLine

	mklit := CS.UnitMilesLit
	if CS.RallyUnitKms {
		mklit = CS.UnitKmsLit
	}
	if Miles < CS.RallyMinMiles {
		sx.IsValidLine = true
		sx.Code = DNF_icon
		sx.Desc = fmt.Sprintf("%v < %v", mklit, CS.RallyMinMiles)
		return sx, EntrantDNF
	}
	if Miles > CS.PenaltyMilesDNF {
		sx.IsValidLine = true
		sx.Code = DNF_icon
		sx.Desc = fmt.Sprintf("%v > %v", mklit, CS.PenaltyMilesDNF)
		return sx, EntrantDNF
	}

	return sx, EntrantFinisher
}

// Miles is either miles or kilometres depending on rally setting, ridden by the current entrant
func calcMileagePenalty(Miles int) ([]ScorexLine, int) {

	const MileagePenaltyIcon = "&#10238;"

	res := make([]ScorexLine, 0)
	var numx int

	usingKms := getStringFromDB("SELECT MilesKms FROM rallyparams", "0") != "0"
	mklit := "miles"
	if usingKms {
		mklit = "km"
	}
	penaltyMaxMiles := getIntegerFromDB("SELECT PenaltyMaxMiles FROM rallyparams", 99999)
	if penaltyMaxMiles < 1 {
		return res, numx
	}
	penaltyMiles := Miles - penaltyMaxMiles
	if penaltyMiles < 1 {
		return res, numx
	}
	penaltyMethod := getIntegerFromDB("SELECT MaxMilesMethod FROM rallyparams", MMM_FixedPoints)
	penaltyPoints := getIntegerFromDB("SELECT MaxMilesPoints FROM rallyparams", 0)
	if penaltyPoints < 1 {
		return res, numx
	}

	var sx ScorexLine
	sx.IsValidLine = true
	sx.Code = penalty_symbol
	switch penaltyMethod {
	case MMM_Multipliers:
		numx = 0 - penaltyPoints
		sx.Desc = fmt.Sprintf("%v %v %v > %v ", MileagePenaltyIcon, Miles, mklit, penaltyMaxMiles)
		sx.PointsDesc = fmt.Sprintf("-&sum;&times;%v", numx)
	case MMM_PointsPerMile:
		sx.Points = 0 - penaltyPoints*penaltyMiles
		sx.Desc = fmt.Sprintf("%v %v %v > %v ", MileagePenaltyIcon, Miles, mklit, penaltyMaxMiles)
		sx.PointsDesc = fmt.Sprintf("%v x %v", penaltyMiles, penaltyPoints)
	default:
		sx.Points = 0 - penaltyPoints
		sx.Desc = fmt.Sprintf("%v %v %v > %v ", MileagePenaltyIcon, Miles, mklit, penaltyMaxMiles)
	}
	res = append(res, sx)
	return res, numx
}

func calcRallyDistance(last int, this int, isKm bool) int {

	delta := this - last
	if isKm == CS.RallyUnitKms {
		return delta
	}
	if isKm {
		delta = int(float64(delta) / KmsPerMile)
	} else {
		delta = int(float64(delta) * KmsPerMile)
	}
	return delta
}

// 2nd return is number of multipliers
func calcTimePenalty(entrant int) ([]ScorexLine, int) {

	const TimePenaltyIcon = "&#x23F0;"
	const FinishTimeIcon = "" // Blank for now

	res := make([]ScorexLine, 0)
	var numx int

	rallyFinishtimex := getStringFromDB("SELECT FinishTime FROM rallyparams", "")

	starttimex := getStringFromDB(fmt.Sprintf("SELECT StartTime FROM entrants WHERE EntrantID=%v", entrant), "")
	if starttimex == "" {
		starttimex = getStringFromDB("SELECT StartTime FROM rallyparams", "")
	}
	if starttimex == "" {
		return res, numx
	}
	finishtimex := getStringFromDB(fmt.Sprintf("SELECT FinishTime FROM entrants WHERE EntrantID=%v", entrant), "")
	if finishtimex == "" {
		finishtimex = rallyFinishtimex
	}
	if finishtimex == "" {
		return res, numx
	}

	rft, err := time.Parse(myTimestamp, rallyFinishtimex)
	if err != nil {
		return res, numx
	}
	st, err := time.Parse(myTimestamp, starttimex)
	if err != nil {
		return res, numx
	}
	ft, err := time.Parse(myTimestamp, finishtimex)
	if err != nil {
		return res, numx
	}
	if ft.Compare(rft) > 0 {
		ft = rft
	}

	maxhrs, _ := strconv.Atoi(getStringFromDB("SELECT MaxHours FROM rallyparams", "999"))
	myDNF := st.Add(time.Hour * time.Duration(maxhrs))
	//fmt.Printf("%v: %v << %v [%v]\n", entrant, starttimex, finishtimex, myDNF.Format(myTimestamp))

	for _, tp := range TimePenalties {
		switch tp.TimeSpec {
		case TimeSpecEntrantDNF:
			ps, _ := strconv.Atoi(tp.PenaltyStartX)
			tp.PenaltyStart = myDNF.Add(time.Minute * time.Duration(0-ps))
			pf, _ := strconv.Atoi(tp.PenaltyFinishX)
			tp.PenaltyFinish = myDNF.Add(time.Minute * time.Duration(0-pf))
		case TimeSpecRallyDNF:
			ps, _ := strconv.Atoi(tp.PenaltyStartX)
			tp.PenaltyStart = rft.Add(time.Minute * time.Duration(0-ps))
			pf, _ := strconv.Atoi(tp.PenaltyFinishX)
			tp.PenaltyFinish = rft.Add(time.Minute * time.Duration(0-pf))
		default:
			// Was specified as actual date/time stamps
			// so already done
		}
		//fmt.Printf("pstart=%v pfinish=%v\n", tp.PenaltyStart.Format(myTimestamp), tp.PenaltyFinish.Format(myTimestamp))
		if ft.Before(tp.PenaltyStart) {
			continue
		}
		if ft.After(tp.PenaltyFinish) {
			ft = tp.PenaltyFinish
		}
		//fmt.Println("Hello sailor")
		var sx ScorexLine
		penaltyTime := ft.Sub(tp.PenaltyStart)
		penaltyMinutes := penaltyTime.Minutes()
		sx.IsValidLine = true
		sx.Code = penalty_symbol
		switch tp.PenaltyMethod {
		case TPM_PointsPerMin:
			sx.Points = 0 - tp.PenaltyFactor*int(penaltyMinutes)
			sx.Desc = fmt.Sprintf("%v %v &ge; %v", TimePenaltyIcon, FinishTimeIcon, tp.PenaltyStart.Format(myTimestampX))
			sx.PointsDesc = fmt.Sprintf("%vm @ %v/m", int(penaltyMinutes), tp.PenaltyFactor)
		case TPM_FixedPoints:
			sx.Points = 0 - tp.PenaltyFactor
			sx.Desc = fmt.Sprintf("%v %v &ge; %v", TimePenaltyIcon, FinishTimeIcon, tp.PenaltyStart.Format(myTimestampX))
		case TPM_FixedMult:
			numx -= tp.PenaltyFactor
			sx.Desc = fmt.Sprintf("%v %v &ge; %v", TimePenaltyIcon, FinishTimeIcon, tp.PenaltyStart.Format(myTimestampX))
			sx.PointsDesc = fmt.Sprintf(" -&sum;&times;%v", tp.PenaltyFactor)
		case TPM_MultPerMin:
			numx -= tp.PenaltyFactor * int(penaltyMinutes)
			sx.Desc = fmt.Sprintf("%v %v &ge; %v", TimePenaltyIcon, FinishTimeIcon, tp.PenaltyStart.Format(myTimestampX))
			sx.PointsDesc = fmt.Sprintf("%vm @ -&sum;&times;%v/m", int(penaltyMinutes), tp.PenaltyFactor)
		}
		res = append(res, sx)
	}
	return res, numx

}

func checkApplySequences(BC ClaimedBonus, LastBonusClaimed ClaimedBonus) ScorexLine {

	var sx ScorexLine

	// Check for sequence bonus
	for _, CR := range CompoundRules {
		if CR.Ruletype != CAT_OrdinaryScoringSequence {
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

		if AxisCounts[CR.Axis].SameCatCount < CR.Min {
			continue
		}

		// Trigger sequential bonus

		sqlx := fmt.Sprintf("SELECT BriefDesc FROM categories WHERE Axis=%d AND Cat=%d", CR.Axis, LastCat)
		defaultValue := fmt.Sprintf("%d/%d", CR.Axis, LastCat)
		BonusDesc := fmt.Sprintf("%s %d x %s", sequential_bonus_symbol, AxisCounts[CR.Axis].SameCatCount, getStringFromDB(sqlx, defaultValue))

		PointsDesc := ""
		ExtraBonusPoints := 0
		if CR.PointsMults == CAT_ResultPoints { // Result is specified number of points
			ExtraBonusPoints = CR.Power
		} else { // Result is sequence length * multiplier
			ExtraBonusPoints = AxisCounts[CR.Axis].SameCatPoints * CR.Power
			if CR.Power != 1 && CR.Power != 0 {
				PointsDesc = fmt.Sprintf("= %d&times;%d", AxisCounts[CR.Axis].SameCatPoints, CR.Power)
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

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func excludeClaim(BC ClaimedBonus) ScorexLine {

	var res ScorexLine

	fmt.Println("excludeClaim called")
	if !CS.ShowExcludedClaims {
		return res
	}
	res.IsValidLine = true
	res.Code = BC.Bonusid
	SB := ScorecardBonuses[BC.BonusScorecardIX] // Convenient shorthand

	res.Desc = fmt.Sprintf("%v<br>*** CLAIM EXCLUDED ***", SB.BriefDesc)
	return res

}

func fetchCatDesc(axis int, cat int) string {

	sqlx := fmt.Sprintf("SELECT BriefDesc FROM categories WHERE Axis=%d AND Cat=%d", axis, cat)
	return getStringFromDB(sqlx, strconv.Itoa(cat))
}

func htmlScorex(sx []ScorexLine, e int, es int, tp int) string {

	const NoScoreIcon = "&#10007;"

	var sp ScorexParams

	KmsRally := getStringFromDB("SELECT MilesKms FROM rallyparams", "0") == "1"
	mk := "miles"
	if KmsRally {
		mk = "km"
	}
	sqlx := "SELECT RiderName,IfNull(PillionName,''),IfNull(CorrectedMiles,0) FROM entrants WHERE EntrantID=" + strconv.Itoa(e)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&sp.RiderName, &sp.PillionName, &sp.DistanceRidden)
		checkerr(err)
	}

	xx := sp.RiderName
	if sp.PillionName != "" {
		xx += " &amp; " + sp.PillionName
	}
	esx, ok := EntrantStatusLits[es]
	if !ok {
		esx = strconv.Itoa(es)
	}
	res := fmt.Sprintf(`<table class="sxtable"><caption>#%d %s [ <span class="sxsfs">%s</span> ]<br><span class="explain">%d %s</span></caption>`, e, xx, esx, sp.DistanceRidden, mk)
	for _, sl := range sx {
		if !sl.IsValidLine {
			continue
		}
		pv := strconv.Itoa(sl.Points)
		if sl.Points == 0 {
			if sl.PointsDesc == "" {
				pv = NoScoreIcon
			} else {
				pv = ""
			}
		}
		res += fmt.Sprintf(`<tr><td class="sxcode">%s</td><td class="sxdesc">%s<span class="sxdescx">%s</span></td><td class="sxitempoints">%s</td></tr>`, sl.Code, sl.Desc, sl.PointsDesc, pv)
	}

	res += fmt.Sprintf(`<tr><td class="sxcode"></td><td class="sxdesc">TOTAL<span class="sxdescx"></span></td><td class="sxitempoints">%d</td></tr>`, tp)
	res += `</table>`

	return res
}

func loadCombos(comboid string) []ComboBonus {

	const cbFieldsB4Cats = 8

	var cb ComboBonus
	res := make([]ComboBonus, 0)

	sqlx := "SELECT ComboID,BriefDesc,ScoreMethod,MinimumTicks,ScorePoints,Bonuses,Compulsory,Cat1"
	for i := 2; i <= NumCategoryAxes; i++ {
		sqlx += fmt.Sprintf(",Cat%d", i)
	}
	sqlx += " FROM combinations"
	if comboid != "" {
		sqlx += " WHERE ComboID='" + comboid + "'"
	}
	sqlx += " ORDER BY ComboID"
	log.Println(sqlx)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()

	s := reflect.ValueOf(&cb).Elem()
	numCols := cbFieldsB4Cats + NumCategoryAxes - 1
	columns := make([]interface{}, numCols)
	for i := 0; i < cbFieldsB4Cats; i++ {
		field := s.Field(i)

		if field.Kind() == reflect.Array {
			for j := 0; j < field.Len(); j++ {
				columns[i+j] = field.Index(j).Addr().Interface()
			}
		} else {
			columns[i] = field.Addr().Interface()
		}
	}
	for rows.Next() {
		err := rows.Scan(columns...)
		checkerr(err)
		x := strings.Split(cb.BonusList, ",")
		cb.Bonuses = make([]string, len(x))
		k := len(x)
		for i := 0; i < k; i++ {
			cb.Bonuses[i] = x[i]
		}
		cb.Points = make([]int, k)
		x = strings.Split(cb.PointsList, ",")
		if cb.MinTicks > k {
			cb.MinTicks = k
		}
		if cb.MinTicks < 1 {
			cb.MinTicks = k
		}
		j := 0
		n := 0
		for i := cb.MinTicks - 1; i < k; i++ {
			if j < len(x) {
				n, _ = strconv.Atoi(x[j])
			}
			cb.Points[i] = n
			j++
		}
		res = append(res, cb)
	}
	return res

}

func loadRejectReasons() RejectReasonsMap {

	const useScoreMasterDB = true

	res := make(RejectReasonsMap, 0)

	if useScoreMasterDB {
		sqlx := "SELECT RejectReasons FROM rallyparams"
		rr := strings.Split(strings.ReplaceAll(getStringFromDB(sqlx, ""), "\r", ""), "\n")
		for _, rx := range rr {
			ct := strings.Split(rx, "=")
			if len(ct) > 1 {
				var r RejectReason
				r.Code, _ = strconv.Atoi(ct[0])
				r.BriefDesc = ct[1]
				res[r.Code] = r
			}
		}
	} else {
		sqlx := "SELECT Code,BriefDesc,Action,IfNull(Param,'') FROM reasons"
		rows, err := DBH.Query(sqlx)
		checkerr(err)
		defer rows.Close()
		for rows.Next() {
			var r RejectReason
			err = rows.Scan(&r.Code, &r.BriefDesc, &r.Action, &r.Param)
			checkerr(err)
			res[r.Code] = r
		}
	}
	return res
}

func processCombos() ([]ScorexLine, int) {

	mults := 0
	res := make([]ScorexLine, 0)

	for cix, cb := range ComboBonuses {
		scoredbonuses := 0
		for _, b := range cb.Bonuses {
			for _, sb := range ScorecardBonuses {
				if sb.Bonusid == b {
					if sb.Scored {
						scoredbonuses++
					}
					break
				}
			}
			for _, sb := range ComboBonuses {
				if sb.Comboid == b {
					if sb.Scored {
						scoredbonuses++
					}
					break
				}
			}
		}
		if scoredbonuses >= cb.MinTicks {
			var sx ScorexLine
			ComboBonuses[cix].Scored = true
			sx.Code = "[" + cb.Comboid + "]"
			sx.Desc = cb.BriefDesc
			sx.IsValidLine = true
			sx.PointsDesc = fmt.Sprintf("(%v/%v)", scoredbonuses, len(cb.Bonuses))
			if cb.ScoreMethod == ScoreMethodMults {
				mults += cb.Points[scoredbonuses-1]
				sx.PointsDesc += fmt.Sprintf(" = &sum;&times;%d", cb.Points[scoredbonuses-1])
			} else {
				sx.Points = cb.Points[scoredbonuses-1]
			}

			res = append(res, sx)
			updateComboCatCounts(cb)
			updateComboCatPoints(cb, sx.Points)
		}
	}
	return res, mults

}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

// This deals with compound rules involving the number of categories scored per axis
func processCompoundCats() ([]ScorexLine, int) {

	mults := 0
	res := make([]ScorexLine, 0)

	lastAxis := -1
	lastCat := -1
	lastMin := 0
	for cix, cr := range CompoundRules {
		if cr.Ruletype == CAT_OrdinaryScoringSequence {
			continue
		}
		if cr.Method != CAT_NumBonusesPerCatMethod {
			continue
		}
		if cr.Target == CAT_ModifyBonusScore {
			continue
		}

		if cr.Axis <= lastAxis && cr.Cat <= lastCat { // Process each axis only once
			continue
		}

		myCount := AxisCounts[cr.Axis].CatCounts[cr.Cat]

		if myCount < cr.Min {
			lastMin = cr.Min
			continue
		}

		// Apply this rule
		lastAxis = cr.Axis
		CompoundRules[cix].Triggered = true

		if cr.Ruletype == CAT_PlaceholderRule {
			continue
		}
		Points := 0
		Pointsdesc := ""
		if cr.Power > 0 {
			Points = cr.Power
		} else {
			Points = myCount
		}
		if cr.Ruletype == CAT_DNF_Unless_Triggered {
			Pointsdesc = checkmark_symbol
			Points = 0
		} else if cr.Ruletype == CAT_DNF_If_Triggered {
			Pointsdesc = "DNF"
			Points = 0
		} else if cr.PointsMults == CAT_ResultMults {
			mults = Points
			Points = 0
			Pointsdesc = fmt.Sprintf("= &sum;&times;%d", mults)
		}
		var sx ScorexLine
		sx.IsValidLine = true
		sx.Code = compoundrule_symbol
		sx.Desc = fmt.Sprintf("%s: <em>n</em>=%d", AxisLabels[cr.Axis-1], myCount)

		sx.Desc += fmt.Sprintf(" rc%v ", cr.Ruleid)

		if cr.Cat > 0 {
			catx := fetchCatDesc(cr.Axis, cr.Cat)
			sx.Desc += fmt.Sprintf(" [%s]", catx)
		}
		if Points < 1 && lastMin > Points {
			sx.Desc += fmt.Sprintf(" &lt; %d", lastMin)
		}
		sx.Points = Points
		sx.PointsDesc = Pointsdesc
		res = append(res, sx)

	}

	return res, mults

}

// This deals with compound rules involving the number of categories scored per axis
func processCompoundNZ() ([]ScorexLine, int) {

	res := make([]ScorexLine, 0)
	mults := 0
	/*
		nzAxisCounts := make([]int, NumCategoryAxes)
		for i := 1; i <= NumCategoryAxes; i++ {
			for j, n := range AxisCounts[i].catCounts {
				if j > 0 && n > 0 { // Skip j=0 as that is sum of other cols
					nzAxisCounts[i]++
				}
			}
		}
	*/

	lastAxis := -1
	lastMin := 0
	for cix, cr := range CompoundRules {
		if cr.Ruletype == CAT_OrdinaryScoringSequence {
			continue
		}
		if cr.Method != CAT_NumNZCatsPerAxisMethod {
			continue
		}
		if cr.Target == CAT_ModifyBonusScore {
			continue
		}

		if cr.Axis <= lastAxis { // Process each axis only once
			continue
		}

		nzCount := 0
		if cr.Axis > 0 {
			nzCount = AxisCounts[cr.Axis].CatCounts[0]
		} else {
			for i := 1; i <= NumCategoryAxes; i++ {
				nzCount += AxisCounts[i].CatCounts[0]
			}
		}

		dbg, _ := json.Marshal(AxisCounts)
		fmt.Printf("A=%v Cnt=%v, Counts=%v\n", cr.Axis, nzCount, string(dbg))
		if nzCount < cr.Min {
			lastMin = cr.Min
			continue
		}

		// Apply this rule
		lastAxis = cr.Axis
		CompoundRules[cix].Triggered = true

		if cr.Ruletype == CAT_PlaceholderRule {
			continue
		}
		Points := 0
		Pointsdesc := ""
		if cr.Power > 0 {
			Points = cr.Power
		} else {
			Points = nzCount
		}
		if cr.Ruletype == CAT_DNF_Unless_Triggered {
			Pointsdesc = checkmark_symbol
			Points = 0
		} else if cr.Ruletype == CAT_DNF_If_Triggered {
			Pointsdesc = "DNF"
			Points = 0
		} else if cr.PointsMults == CAT_ResultMults {
			mults = Points
			Points = 0
			Pointsdesc = fmt.Sprintf("= &sum;&times%d", mults)
		}
		var sx ScorexLine
		sx.IsValidLine = true
		sx.Code = compoundrule_symbol
		sx.Desc = fmt.Sprintf("%s: <em>n</em>=%d", AxisLabels[cr.Axis-1], nzCount)

		sx.Desc += fmt.Sprintf(" rz%v ", cr.Ruleid)
		if cr.Cat > 0 {
			catx := fetchCatDesc(cr.Axis, cr.Cat)
			sx.Desc += fmt.Sprintf(" [%s]", catx)
		}
		if Points < 1 && lastMin > Points {
			sx.Desc += fmt.Sprintf(" &lt; %d", lastMin)
		}
		sx.Points = Points
		sx.PointsDesc = Pointsdesc
		res = append(res, sx)

	}

	return res, mults

}

func recalc_all() {

	/* 	_, err := DBH.Exec("BEGIN TRANSACTION")
	   	checkerr(err)
	   	defer DBH.Exec("COMMIT")
	*/
	sqlx := "SELECT EntrantID FROM entrants ORDER BY EntrantID"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	n := 3
	entrants := make([]int, 0)
	for rows.Next() {
		var entrant int
		rows.Scan(&entrant)
		entrants = append(entrants, entrant)
		//		recalc_scorecard(entrant)
		n--
		if n < 1 {
			break
		}
	}
	rows.Close()
	for _, e := range entrants {
		recalc_scorecard(e)
	}

}

// This recalculates the value of the specified scorecard using as
// input the relevant claims records. The results are updated totals
// and score explanation.
func recalc_scorecard(entrant int) {

	// Multipliers can accrue from combos or compound rules
	// The final points score, after penalties, is multiplied by this value
	Multipliers := 0

	log.Printf("recalc for %v\n", entrant)

	/* 	_, err := DBH.Exec("BEGIN TRANSACTION")
	   	checkerr(err)
	   	defer DBH.Exec("COMMIT")
	*/

	BuildRallyParameters(CS.CurrentLeg)

	AxisCounts = build_emptyAxisCountsArray()

	BonusesClaimed = build_bonusclaim_array(entrant)

	CorrectedMiles := 0
	StartOdo := getIntegerFromDB("SELECT ifnull(OdoRallyStart,0) FROM entrants WHERE EntrantID="+strconv.Itoa(entrant), 0)
	LastOdo := StartOdo
	FinishOdo := getIntegerFromDB("SELECT ifnull(OdoRallyFinish,0) FROM entrants WHERE EntrantID="+strconv.Itoa(entrant), 0)
	OdoIsKm := getIntegerFromDB("SELECT OdoKms FROM entrants WHERE EntrantID="+strconv.Itoa(entrant), 0) != 0

	var sx ScorexLine
	TotalPoints := 0
	var Scorex []ScorexLine

	var LastBonusClaimed ClaimedBonus
	for _, BC := range BonusesClaimed {

		//fmt.Printf("%v d==%v\n", BC.Bonusid, BC.Decision)

		// ClaimExcluded means ignore it, treat it as if it didn't exist
		// This might need to be a switchable response
		if BC.Decision == ClaimDecision_ClaimExcluded {
			sx = excludeClaim(BC)
			if sx.IsValidLine {
				TotalPoints += sx.Points
				Scorex = append(Scorex, sx)
			}
			continue
		}

		// Need to flag the bonus as having been scored
		BC.BonusScorecardIX = slices.IndexFunc(ScorecardBonuses, func(c ScorecardBonusDetail) bool { return c.Bonusid == BC.Bonusid })

		SB := ScorecardBonuses[BC.BonusScorecardIX] // Convenient shorthand

		ScorecardBonuses[BC.BonusScorecardIX].Scored = BC.Decision == ClaimDecision_GoodClaim // Only good claims count against "must score" flag

		LastOdo = BC.OdoReading

		sx = checkApplySequences(BC, LastBonusClaimed)
		if sx.IsValidLine {
			TotalPoints += sx.Points
			Scorex = append(Scorex, sx)
		}

		if BC.Decision != ClaimDecision_GoodClaim {

			// Firstly, let's zap any sequence in progress
			for i := 1; i <= NumCategoryAxes; i++ {
				cc := AxisCounts[i]
				cc.SameCatCount = 0
				cc.SameCatPoints = 0
				cc.LastCatScored = -1
				AxisCounts[i] = cc
			}
			sx.IsValidLine = true
			sx.Code = SB.Bonusid
			reject, ok := RejectReasons[BC.Decision]
			errmsg := "***"
			if ok {
				errmsg = reject.BriefDesc
			}

			sx.Desc = fmt.Sprintf("%v<br>CLAIM REJECTED - %v", SB.BriefDesc, errmsg)
			sx.Points = 0
			sx.PointsDesc = ""
			Scorex = append(Scorex, sx)
			continue
		}

		PointsDesc := ""

		// Handle multipliers

		if SB.MultiplyLast {
			if ScorecardBonuses[LastBonusClaimed.BonusScorecardIX].MultiplyLast {
				BC.Points = 0
			} else {
				PointsDesc = fmt.Sprintf("= %v &times;%v", LastBonusClaimed.Points, BC.Points)
				BC.Points = LastBonusClaimed.Points * BC.Points
			}
		}

		// Question/answer logic. Points awarded at claim time so just flag on scorex
		if BC.QuestionScored {
			const QuestionScoredMark = "?" + checkmark_symbol
			PointsDesc += QuestionScoredMark
		}

		updateBonusCatCounts(SB) // Updating here gets the counts right but not points upgraded below

		fmt.Printf("Feature Cat is %v %v\n", SB.CatValue[0], SB.CatValue)
		LastBonusClaimed = BC

		BasicPoints := BC.Points

		// Compound rules affecting individual bonuses

		for _, cr := range CompoundRules {

			if cr.Ruletype != CAT_OrdinaryScoringRule {
				continue
			}
			if cr.Target != CAT_ModifyBonusScore {
				continue
			}
			if cr.Cat > 0 { // Rule applies only to one category
				if SB.CatValue[cr.Axis-1] != cr.Cat { // not this one
					continue
				}
			}

			// Check how many hits
			catcount := 0
			if cr.Cat == 0 {
				for _, cc := range AxisCounts[cr.Axis].CatCounts {
					catcount += cc
				}
			} else {
				catcount += AxisCounts[cr.Axis].CatCounts[cr.Cat]
			}
			if catcount < cr.Min {
				continue
			}
			dbg, _ := json.Marshal(cr)
			fmt.Printf("%v [[ %v ]] %v\n", BC.Bonusid, string(dbg), catcount)
			if cr.Power == 0 {
				PointsDesc = fmt.Sprintf("%d x %d", BasicPoints, catcount-1)
				BasicPoints = BasicPoints * (catcount - 1)
			} else if cr.PointsMults == CAT_ResultMults {
				PointsDesc = fmt.Sprintf("%d x %d x %d", BasicPoints, cr.Power, catcount-1)
				BasicPoints = BasicPoints * cr.Power * (catcount - 1)
			} else {
				PointsDesc = fmt.Sprintf("%d x %d^%d", BasicPoints, cr.Power, catcount-1)
				BasicPoints = BasicPoints * powInt(cr.Power, catcount-1)
			}
			PointsDesc += fmt.Sprintf("r%v", cr.Ruleid)
			break // Only apply the first matching rule
		}

		updateBonusCatPoints(SB, BasicPoints) // Updating here gets wrong counts but correctly upgraded points

		TotalPoints += BasicPoints
		sx.IsValidLine = true
		sx.Code = SB.Bonusid
		sx.Desc = SB.BriefDesc
		sx.Points = BasicPoints
		sx.PointsDesc = PointsDesc
		Scorex = append(Scorex, sx)
	}

	// Final check for a sequence
	var BC ClaimedBonus
	sx = checkApplySequences(BC, LastBonusClaimed)
	if sx.IsValidLine {
		TotalPoints += sx.Points
		Scorex = append(Scorex, sx)
	}

	combosx, nm := processCombos()
	for _, cx := range combosx {
		TotalPoints += cx.Points
	}

	Multipliers += nm
	Scorex = append(Scorex, combosx...)

	// Now let's calculate the axis scores - starting with non-zero numbers of categories
	nzAxisCounts := make([]int, NumCategoryAxes)
	for i := 0; i < NumCategoryAxes; i++ {
		for _, n := range AxisCounts[i].CatCounts {
			if n > 0 {
				nzAxisCounts[i]++
			}
		}
	}

	nz, nm := processCompoundNZ()
	for _, cx := range nz {
		TotalPoints += cx.Points
	}

	Multipliers += nm
	Scorex = append(Scorex, nz...)

	nc, nm := processCompoundCats()
	for _, cx := range nc {
		TotalPoints += cx.Points
	}

	Multipliers += nm
	Scorex = append(Scorex, nc...)

	ntp, nm := calcTimePenalty(entrant)
	for _, px := range ntp {
		TotalPoints += px.Points
	}
	Multipliers += nm
	Scorex = append(Scorex, ntp...)

	if CS.UseCheckinForOdo && FinishOdo > StartOdo {
		LastOdo = FinishOdo
	}
	if !CS.UseCheckinForOdo {
		FinishOdo = LastOdo
	}

	CorrectedMiles = calcRallyDistance(StartOdo, LastOdo, OdoIsKm)

	ntp, nm = calcMileagePenalty(CorrectedMiles)
	for _, px := range ntp {
		TotalPoints += px.Points
	}
	Multipliers += nm
	Scorex = append(Scorex, ntp...)

	if Multipliers > 1 {
		var sx ScorexLine
		sx.IsValidLine = true
		n := TotalPoints * Multipliers
		sx.Desc = fmt.Sprintf("%d &times; %d", TotalPoints, Multipliers)
		sx.PointsDesc = fmt.Sprintf("= %v", n)
		sx.Points = 0 //n - TotalPoints
		TotalPoints = n

		Scorex = append(Scorex, sx)
	}

	sx, status := calcEntrantStatus(CorrectedMiles)
	Scorex = append(Scorex, sx)

	htmlSX := htmlScorex(Scorex, entrant, status, TotalPoints)

	/* 	_, err = DBH.Exec("COMMIT")
	   	checkerr(err)
	*/
	sqlx := "UPDATE entrants SET ScoreX=?,EntrantStatus=?,TotalPoints=?,CorrectedMiles=?,OdoRallyFinish=? WHERE EntrantID=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()

	_, err = stmt.Exec(htmlSX, status, TotalPoints, CorrectedMiles, FinishOdo, entrant)
	checkerr(err)

	// Debugging code below
	const dbgstyle = `<style>
    .sxitempoints { text-align: right;}
    .sxdescx { padding-left: 1em; font-style: italic;}
    </style>
`
	f, err := os.Create("crap.md")
	checkerr(err)
	f.WriteString(dbgstyle)
	f.WriteString(htmlSX)
	f.Close()
	for x := range Scorex {
		log.Printf("%-3s %-20s %-10s %7d\n", Scorex[x].Code, html.UnescapeString(Scorex[x].Desc), html.UnescapeString(Scorex[x].PointsDesc), Scorex[x].Points)
	}
	log.Printf("Total points is %d\n", TotalPoints)

}

func updateBonusCatCounts(BS ScorecardBonusDetail) {

	updateCatCounts(BS.CatValue[:])

}

func updateBonusCatPoints(BS ScorecardBonusDetail, Points int) {

	updateCatPoints(BS.CatValue[:], Points)

}

func updateCatCounts(CatValue []int) {

	for i := 1; i <= NumCategoryAxes; i++ {

		cat := CatValue[i-1]

		cc := AxisCounts[i]
		if cat <= 0 {
			cc.SameCatCount = 0
			cc.SameCatPoints = 0
			cc.LastCatScored = cat
			AxisCounts[i] = cc
			continue

		} else if cat == cc.LastCatScored {
			cc.SameCatCount++
		} else {
			cc.SameCatCount = 1
			cc.SameCatPoints = 0
			cc.LastCatScored = cat
			// Now accrue overall axis totals
			_, ok := cc.CatCounts[0]
			if ok {
				cc.CatCounts[0]++
			} else {
				cc.CatCounts[0] = 1
			}
		}
		AxisCounts[i] = cc

		_, ok := AxisCounts[i].CatCounts[cat]

		if ok {
			AxisCounts[i].CatCounts[cat]++
		} else {
			AxisCounts[i].CatCounts[cat] = 1
		}

	}

}

func updateCatPoints(CatValue []int, Points int) {

	for i := 1; i <= NumCategoryAxes; i++ {

		cat := CatValue[i-1]

		cc := AxisCounts[i]
		if cat <= 0 {
			cc.SameCatPoints = 0
			AxisCounts[i] = cc
			continue
		} else if cat == cc.LastCatScored {
			cc.SameCatPoints += Points
		} else {
			cc.SameCatPoints = Points
		}
		AxisCounts[i] = cc
	}
}

func updateComboCatCounts(CB ComboBonus) {

	updateCatCounts(CB.Cat[:])

}

func updateComboCatPoints(CB ComboBonus, Points int) {

	updateCatPoints(CB.Cat[:], Points)

}
