package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var tmpltOption = `<option value="%d" %s>%s</option>`

var tmpltSingleRule = `
<div class="singlerule">
  <fieldset class="field">
    <label for="Axis">Axis</label>
	<select id="Axis" name="Axis">
	%s
	</select>
  </fieldset>
  <fieldset class="field">
    <label for="Cat">Category</label>
    <select id="Cat" name="Cat">
    %s
    </select>
  </fieldset>
  <fieldset class="field">
    <label for="NMethod">Calculate 'n' =</label>
    <select id="NMethod" name="NMethod">
    %s
    </select>
  </fieldset>
  <fieldset class="field">
    <label for="ModBonus">This rule affects</label>
    <select id="ModBonus" name="ModBonus">
    %s
    </select>
  </fieldset>
  <fieldset class="field">
    <label for="NMin">Minimum value of 'n'</label>
	<input id="NMin" name="NMin" type="number" value="%d">
  </fieldset>
  <fieldset class="field">
    <label for="PointsMults">This rule results in</label>
    <select id="PointsMults" name="PointsMults">
    %s
    </select>
  </fieldset>
  <fieldset class="field">
    <label for="NPower">Number of points or multipliers</label>
	<input id="NPower" name="NPower" type="number" value="%d">
  </fieldset>
  <fieldset class="field">
    <label for="RuleType">This rule is</label>
    <select id="RuleType" name="RuleType">
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
	for rows.Next() {
		var cat int
		var desc string
		err = rows.Scan(&cat, &desc)
		checkerr(err)
		sel := ""
		if cat == selcat {
			sel = "selected"
		}
		x := fmt.Sprintf(tmpltOption, cat, sel, desc)
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
		x := fmt.Sprintf(tmpltOption, v, selx, lbls[i])
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
		if r.Axis == ix {
			sel = "selected"
		}
		x := fmt.Sprintf(tmpltOption, strconv.Itoa(ix), sel, axis)
		axisopts = append(axisopts, x)
	}

	opts := make([]string, 2)

	if r.Method == CAT_NumBonusesPerCatMethod {
		opts[0] = fmt.Sprintf(tmpltOption, CAT_NumBonusesPerCatMethod, "selected", "Bonuses per category")
		opts[1] = fmt.Sprintf(tmpltOption, CAT_NumNZCatsPerAxisMethod, "", "Nonzero categories")
	} else {
		opts[0] = fmt.Sprintf(tmpltOption, CAT_NumBonusesPerCatMethod, "", "Bonuses per category")
		opts[1] = fmt.Sprintf(tmpltOption, CAT_NumNZCatsPerAxisMethod, "selected", "Nonzero categories")
	}
	targets := make([]string, 2)

	if r.Method == CAT_ModifyBonusScore {
		targets[0] = fmt.Sprintf(tmpltOption, CAT_ModifyBonusScore, "selected", "Individual bonuses")
		targets[1] = fmt.Sprintf(tmpltOption, CAT_ModifyAxisScore, "", "Groups of bonuses")
	} else {
		targets[0] = fmt.Sprintf(tmpltOption, CAT_ModifyBonusScore, "", "Individual bonuses")
		targets[1] = fmt.Sprintf(tmpltOption, CAT_ModifyAxisScore, "selected", "Groups of bonuses")
	}

	pms := make([]string, 2)

	if r.Method == CAT_ResultPoints {
		pms[0] = fmt.Sprintf(tmpltOption, CAT_ResultPoints, "selected", "points")
		pms[1] = fmt.Sprintf(tmpltOption, CAT_ResultMults, "", "multipliers")
	} else {
		pms[0] = fmt.Sprintf(tmpltOption, CAT_ResultPoints, "", "points")
		pms[1] = fmt.Sprintf(tmpltOption, CAT_ResultMults, "selected", "multipliers")
	}

	if r.Method == CAT_ResultPoints {
		pms[0] = fmt.Sprintf(tmpltOption, CAT_ResultPoints, "selected", "points")
		pms[1] = fmt.Sprintf(tmpltOption, CAT_ResultMults, "", "multipliers")
	} else {
		pms[0] = fmt.Sprintf(tmpltOption, CAT_ResultPoints, "", "points")
		pms[1] = fmt.Sprintf(tmpltOption, CAT_ResultMults, "selected", "multipliers")
	}

	rtvals := []int{CAT_OrdinaryScoringRule, CAT_DNF_Unless_Triggered, CAT_DNF_If_Triggered, CAT_PlaceholderRule, CAT_OrdinaryScoringSequence}
	rtlabs := []string{"ordinary scoring rule", "DNF unless triggered", "DNF if triggered", "placeholder only", "Ordinary scoring sequence"}

	page := fmt.Sprintf(tmpltSingleRule,
		strings.Join(axisopts, ""),
		strings.Join(optsSingleAxisCats(r.Axis, r.Cat), ""),
		strings.Join(opts, ""),
		strings.Join(targets, ""),
		r.Min,
		strings.Join(pms, ""),
		r.Power,
		selectOptionArray(rtvals, rtlabs, r.Ruletype))
	w.Write([]byte(page))
}
