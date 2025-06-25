package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var tmpltOption = `<option value="%s" %s>%s</option>`

var tmpltSingleRule = `
<style>` + maincss + `
</style>
<script>` + mainscript + `
</script>
<div class="singlerule">
  <form action="updtcrule" method="post">
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="RuleType">This rule is</label>
    <select id="RuleType" name="RuleType" onchange="chgRuleType(this);">
    %s
    </select>
  </fieldset>
  <fieldset class="field rule0">
    <label for="ModBonus">This rule affects the value of</label>
    <select id="ModBonus" name="ModBonus">
    %s
    </select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3">
    <label for="NMethod">Calculate <var>n</var> using </label>
    <select id="NMethod" name="NMethod">
    %s
    </select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="NMin">Triggered when <var>n</var> &ge; </label>
	<input id="NMin" name="NMin" type="number" value="%d">
  </fieldset>
  <fieldset class="field rule0 rule4">
    <label for="PointsMults">This rule results in</label>
	<fieldset class="field rule0 rule4">
	<input id="NPower" name="NPower" type="number" value="%d">
    <select id="PointsMults" name="PointsMults">
    %s
    </select>
	</fieldset>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="Axis">Bonus category</label>
	<select id="Axis" name="Axis" onchange="chgAxis(this);">
	%s
	</select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="Cat">Category</label>
    <select id="Cat" name="Cat">
    %s
    </select>
  </fieldset>
  <input type="hidden" name="ruleid" value="%d">
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
	<fieldset class="rule0 rule1 rule2 rule3 rule4"></fieldset>
	<fieldset class="rule0 rule1 rule2 rule3 rule4 flexspread">
    <input type="submit" name="save" value=" update database "> <input type="button" name="delete" value=" &#10006; ">
	</fieldset>
  </fieldset>
  </form>
</div>` + `
<script>
setupForm()
</script>`

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

func showSingleRule(w http.ResponseWriter, r CompoundRule) {

	//fmt.Printf("SSR = %v\n", r)
	axisopts := make([]string, 0)
	for ix, axis := range AxisLabels {
		if axis == "" {
			continue
		}
		sel := ""
		if r.Axis == ix+1 {
			sel = "selected"
		}
		xx := strconv.Itoa(ix + 1)
		x := fmt.Sprintf(tmpltOption, xx, sel, axis)
		log.Printf("x=%v, ix=%d, xx=%v\n", x, ix, xx)
		axisopts = append(axisopts, x)
	}

	rtvals := []int{CAT_OrdinaryScoringRule, CAT_DNF_Unless_Triggered, CAT_DNF_If_Triggered, CAT_PlaceholderRule, CAT_OrdinaryScoringSequence}
	rtlabs := []string{"ordinary scoring rule", "DNF unless triggered", "DNF if triggered", "placeholder only", "ordinary scoring sequence"}

	page := fmt.Sprintf(tmpltSingleRule,
		selectOptionArray(rtvals, rtlabs, r.Ruletype),
		selectOptionArray([]int{CAT_ModifyBonusScore, CAT_ModifyAxisScore}, []string{"individual bonuses", "additional group-based awards"}, r.Target),
		selectOptionArray([]int{CAT_NumBonusesPerCatMethod, CAT_NumNZCatsPerAxisMethod}, []string{"Bonuses per category", "Categories scored"}, r.Method),
		r.Min,
		r.Power,
		selectOptionArray([]int{CAT_ResultPoints, CAT_ResultMults}, []string{"points", "multipliers"}, r.PointsMults),
		strings.Join(axisopts, ""),
		strings.Join(optsSingleAxisCats(r.Axis, r.Cat), ""),
		r.Ruleid)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(page))
}

func show_rules(w http.ResponseWriter, r *http.Request) {

	const leg = 0
	var rt = map[int]string{0: "ordinary", 1: "DNF unless", 2: "DNF if", 3: "dummy", 4: "sequence"}
	rules := build_compoundRuleArray(leg)
	axes := build_axisLabels()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<style>%s</style>`, maincss)
	fmt.Fprintf(w, `<script>%s</script>`, mainscript)
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
				target = "axis"
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
