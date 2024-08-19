package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//go:embed compoundrules.js
var script string

var tmpltOption = `<option value="%s" %s>%s</option>`

var tmpltSingleRule = `
<style>
.singlerule { max-width: 50em; margin: auto; }
.singlerule input[type='number'] {width: 5em;}
.hide {display:none;}
</style>
<script>` + script + `
</script>
<div class="singlerule">
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
    <label for="NMethod">Calculate 'n' using </label>
    <select id="NMethod" name="NMethod">
    %s
    </select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="NMin">Triggered when 'n' &ge; </label>
	<input id="NMin" name="NMin" type="number" value="%d">
  </fieldset>
  <fieldset class="field rule0 rule4">
    <label for="PointsMults">This rule results in</label>
	<input id="NPower" name="NPower" type="number" value="%d">
    <select id="PointsMults" name="PointsMults">
    %s
    </select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4">
    <label for="Axis">Axis</label>
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
</div>`

func optsSingleAxisCats(axis int, selcat int) []string {

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
	rtlabs := []string{"ordinary scoring rule", "DNF unless triggered", "DNF if triggered", "placeholder only", "Ordinary scoring sequence"}

	page := fmt.Sprintf(tmpltSingleRule,
		selectOptionArray(rtvals, rtlabs, r.Ruletype),
		selectOptionArray([]int{CAT_ModifyBonusScore, CAT_ModifyAxisScore}, []string{"Individual bonuses", "Additional group-based awards"}, r.Target),
		selectOptionArray([]int{CAT_NumBonusesPerCatMethod, CAT_NumNZCatsPerAxisMethod}, []string{"Bonuses per category", "Categories scored"}, r.Method),
		r.Min,
		r.Power,
		selectOptionArray([]int{CAT_ResultPoints, CAT_ResultMults}, []string{"points", "multipliers"}, r.PointsMults),
		strings.Join(axisopts, ""),
		strings.Join(optsSingleAxisCats(r.Axis, r.Cat), ""))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	w.Write([]byte(page))
}
