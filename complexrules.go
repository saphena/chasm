package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const crintro = `<p>Complex rules and categories are used to enable scoring mechanisms beyond simple bonuses and combos.
	Scores generated here can apply to individual bonuses or to collections or sequences of bonuses.
	</p>`

var tmpltOption = `<option value="%s" %s>%s</option>`

var crtopline = `<div class="topline">
		<fieldset>
			<button title="Delete this Bonus?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" style="display:none;" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
		<fieldset>
			<button id="updatedb" class="hideuntil" title="Delete Rule" disabled onclick="deleteRule(this)"></button>
		</fieldset>

	</div>`

var tmpltSingleRule = `
<div id="singlerule" class="singlerule">
  <form action="updtcrule" method="post">
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="RuleType">Rule type</label>
    <select id="RuleType" name="RuleType" onchange="chgRuleType(this);">
    %v
    </select>
  </fieldset>
  <fieldset class="field help rule3">
  <p>Placeholders are used in some circumstances to improve analysis of a set of rules.</p> 
  <p>For example: A rule of type "DNF unless triggered" will appear on score explanations with a tick if it is triggered but will not appear on score explanations at all if it's not triggered, even though that condition will result in DNF. Mostly it would be better to use a placeholder and a "DNF if triggered" pair as that will give more satisfying results on score explanations.</p>
  </fieldset>
  </fieldset class="field help rule4">
  <p>The sequence refers to the order in which claims are submitted, not the sequence within the rally book or geographical location.</p>
  </fieldset>
  <fieldset class="field rule0">
    <label for="ModBonus">This rule affects the value of</label>
    <select id="ModBonus" name="ModBonus" onchange="saveRule(this)">
    %v
    </select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3">
    <label for="NMethod">Calculate <var>n</var> using </label>
    <select id="NMethod" name="NMethod" onchange="saveRule(this)">
    %v
    </select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="NMin">Triggered when <var>n</var> &ge; </label>
	<input id="NMin" name="NMin" type="number" value="%v" onchange="saveRule(this)">
  </fieldset>
  <fieldset class="field rule0 rule4">
    <label for="PointsMults">This rule results in</label>
	<fieldset class="field rule0 rule4">
	<input id="NPower" name="NPower" type="number" value="%v" onchange="saveRule(this)">
    <select id="PointsMults" name="PointsMults" onchange="saveRule(this)">
    %v
    </select>
	</fieldset>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="Axis">Category set</label>
	<select id="Axis" name="Axis" onchange="chgAxis(this);">
	%v
	</select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="Cat">%v</label>
    <select id="Cat" name="Cat" onchange="saveRule(this)">
    %v
    </select>
  </fieldset>
  <input type="hidden" id="ruleid" name="ruleid" value="%v">
  </form>
</div>
<script>
showCurrentRule()
</script>`

func deleteRule(w http.ResponseWriter, r *http.Request) {

	rule := r.PathValue("rule")
	if rule == "" {
		fmt.Fprint(w, `{"ok":false,"msg":incomplete request"}`)
		return
	}
	sqlx := "DELETE FROM catcompound WHERE rowid=" + rule
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}

func optsSingleAxisCats(axis int, selcat int) []string {

	//fmt.Printf("SAC %v %v\n", axis, selcat)
	sqlx := fmt.Sprintf("SELECT Cat,BriefDesc FROM categories WHERE Axis=%d", axis)
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	res := make([]string, 0)
	sel := ""
	if selcat == 0 {
		sel = "selected"
	}
	res = append(res, fmt.Sprintf(tmpltOption, "0", sel, "any"))
	for rows.Next() {
		var cat int
		var desc string
		err = rows.Scan(&cat, &desc)
		checkerr(err)
		sel = ""
		if cat == selcat {
			sel = "selected"
		}
		x := fmt.Sprintf(tmpltOption, strconv.Itoa(cat), sel, desc)
		res = append(res, x)
	}
	return res
}

func selectOptionArray(vals []int, lbls []string, sel int) []string {

	res := make([]string, 0)
	for i, v := range vals {
		var selx string
		if v == sel {
			selx = "selected"
		}
		x := fmt.Sprintf(tmpltOption, strconv.Itoa(v), selx, lbls[i])
		res = append(res, x)
	}
	return res
}

func show_rule(w http.ResponseWriter, r *http.Request) {

	const Leg = 0

	n, err := strconv.Atoi(r.FormValue("r"))
	if err != nil {
		n = 1
	}
	CompoundRules = build_compoundRuleArray(Leg)
	for _, cr := range CompoundRules {
		if cr.Ruleid == n {
			showSingleRule(w, r, cr)
			return
		}
	}
	fmt.Fprint(w, `OMG`)
}

func showSingleRule(w http.ResponseWriter, r *http.Request, cr CompoundRule) {

	//fmt.Printf("SSR = %v\n", r)
	startHTMLBL(w, "Complex rules", r.FormValue("back"))
	fmt.Fprint(w, crtopline)
	fmt.Fprintf(w, `<div class="intro">%v</div>`, crintro)
	fmt.Fprint(w, `</header>`)
	AxisLabels := build_axisLabels()
	axisopts := make([]string, 0)
	for ix, axis := range AxisLabels {
		if axis == "" {
			continue
		}
		sel := ""
		if cr.Axis == ix+1 {
			sel = "selected"
		}
		xx := strconv.Itoa(ix + 1)
		x := fmt.Sprintf(tmpltOption, xx, sel, axis)
		log.Printf("x=%v, ix=%d, xx=%v\n", x, ix, xx)
		axisopts = append(axisopts, x)
	}

	rtvals := []int{CAT_OrdinaryScoringRule, CAT_DNF_Unless_Triggered, CAT_DNF_If_Triggered, CAT_PlaceholderRule, CAT_OrdinaryScoringSequence}
	rtlabs := []string{"ordinary scoring rule", "DNF unless triggered", "DNF if triggered", "placeholder only", "uninterrupted sequence"}

	page := fmt.Sprintf(tmpltSingleRule,
		selectOptionArray(rtvals, rtlabs, cr.Ruletype),
		selectOptionArray([]int{CAT_ModifyBonusScore, CAT_ModifyAxisScore}, []string{"individual bonuses", "group-based awards"}, cr.Target),
		selectOptionArray([]int{CAT_NumBonusesPerCatMethod, CAT_NumNZCatsPerAxisMethod}, []string{"Bonuses per category", "Categories scored"}, cr.Method),
		cr.Min,
		cr.Power,
		selectOptionArray([]int{CAT_ResultPoints, CAT_ResultMults}, []string{"points", "multipliers"}, cr.PointsMults),
		strings.Join(axisopts, ""),
		AxisLabels[cr.Axis-1],
		strings.Join(optsSingleAxisCats(cr.Axis, cr.Cat), ""),
		cr.Ruleid)

	w.Write([]byte(page))
}

func show_rules(w http.ResponseWriter, r *http.Request) {

	const leg = 0
	var rt = map[int]string{0: "ordinary", 1: "DNF unless", 2: "DNF if", 3: "dummy", 4: "sequence"}
	rules := build_compoundRuleArray(leg)
	axes := build_axisLabels()
	startHTML(w, "Complex rules")
	fmt.Fprintf(w, `<div class="intro">%v</div>`, crintro)
	fmt.Fprint(w, `<div class="ruleset">`)
	fmt.Fprint(w, `<fieldset class="row hdr">`)
	fmt.Fprint(w, `<span class="col">Set</span><span class="col">Category</span>`)
	fmt.Fprint(w, `<span class="col">Type</span><span class="col">Threshold</span>`)
	fmt.Fprint(w, `<span class="col">Target</span><span class="col">Score</span>`)
	fmt.Fprint(w, `</fieldset>`)
	fmt.Fprint(w, `</div><hr></header>`)
	fmt.Fprint(w, `<div class="ruleset">`)
	for _, cr := range rules {
		fmt.Fprintf(w, `<fieldset class="row target" data-rowid="%d" title="%v" onclick="showRule(this);">`, cr.Ruleid, cr.Ruleid)
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, axes[cr.Axis-1])
		sqlx := fmt.Sprintf("SELECT BriefDesc FROM categories WHERE Axis=%d AND Cat=%d", cr.Axis, cr.Cat)
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, getStringFromDB(sqlx, "any"))
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, rt[cr.Ruletype])
		fmt.Fprintf(w, `<fieldset class="col">&ge; %d</fieldset>`, cr.Min)
		target := ""
		pm := "pts"
		if cr.PointsMults == CAT_ResultMults {
			pm = "x"
		}
		pts := strconv.Itoa(cr.Power)
		switch cr.Ruletype {
		case CAT_OrdinaryScoringRule:
			if cr.Target == CAT_ModifyAxisScore {
				target = "group"
			} else {
				target = "bonus"
			}
		case CAT_OrdinaryScoringSequence:
			target = "bonus"
		default:
			pm = ""
			pts = ""
		}
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, target)
		fmt.Fprintf(w, `<fieldset class="col">%s %s</fieldset>`, pts, pm)
		fmt.Fprint(w, `</fieldset>`)
	}
	fmt.Fprint(w, `</div>`)
}

func update_rule(w http.ResponseWriter, r *http.Request) {

	log.Println("updating a rule")
	var cr CompoundRule

	r.ParseForm()

	log.Printf("%v\n", r.Form)
	ruleid, err := strconv.Atoi(r.FormValue("ruleid"))
	if err != nil || ruleid == 0 {
		fmt.Fprint(w, "No ruleid supplied")
		return
	}
	if r.FormValue("delete") != "" {
		log.Printf("Deleting rule %v\n", ruleid)
	}
	cr.Axis, _ = strconv.Atoi(r.FormValue("Axis"))
	cr.Cat, _ = strconv.Atoi(r.FormValue("Cat"))
	cr.Method, _ = strconv.Atoi(r.FormValue("NMethod"))
	cr.Min, _ = strconv.Atoi(r.FormValue("NMin"))
	cr.PointsMults, _ = strconv.Atoi("PointsMults")
	cr.Power, _ = strconv.Atoi(r.FormValue("Power"))
	cr.Ruletype, _ = strconv.Atoi(r.FormValue("RuleType"))

	sqlx := "UPDATE catcompound SET Ruletype=%d,Axis=%d,Cat=%d,NMethod=%d,NMin=%d,PointsMults=%d,NPower=%d WHERE rowid=%d"
	sqlx = fmt.Sprintf(sqlx, cr.Ruletype, cr.Axis, cr.Cat, cr.Method, cr.Min, cr.PointsMults, cr.Power, ruleid)
	log.Println(sqlx)

	_, err = DBH.Exec(sqlx)
	checkerr(err)

	fmt.Fprint(w, `<script>window.location.href="/rules"</script>`)

}

func saveRule(w http.ResponseWriter, r *http.Request) {

	ruleid := r.FormValue("ruleid")
	if ruleid == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no ruleid supplied"}`)
		return
	}
	fld := r.FormValue("ff")
	if fld == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"no fieldname supplied"}`)
		return
	}
	val := r.FormValue(fld)
	sqlx := "UPDATE catcompound SET " + fld + "=? WHERE rowid=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(val, ruleid)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)
}
