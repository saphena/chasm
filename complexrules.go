package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const complexhelp = `
<article id="complexhelp" class="popover" popover>
<h1>Complex rules</h1>

<p>Complex rules are grouped into <em>rulesets</em> each belonging to a single <em>Set</em>, <em>Category</em> and <em>Target</em> and ranked according to the value of the <em>Threshold</em>. The rule matching the highest value of threshold will be applied, others in the ruleset are ignored.</p>
<p>Scores are built up as follows:-</p>
<ol>
	<li>Basic bonus score, maybe multiplier of last bonus</li>
	<li>Uninterrupted sequence bonus</li>
	<li>Complex rules affecting individual bonuses</li>
	<li>Complex rules affecting groups of bonuses</li>
</ol>
<h2>How do I ?</h2>

<p>The following gives examples of how to configure a ruleset to achieve certain specific scoring goals.</p>
<h3>Award extra points for scoring N categories within an set</h3>
<ul>
    <li>One or more rules, each of type Categories Scored, Affects Group-based awards, Ordinary scoring rule.</li>
    <li>Set Triggered when and Results in to the required values.</li>
</ul>
<h3>Deduct points for scoring less than N categories within an set</h3>
<ul>
    <li>Set a placeholder with a Triggered when value = N;</li>
    <li>Set an Ordinary scoring rule with Triggered when = 0 and Results in to the negative value.</li>
</ul>
<h3>Award extra points for scoring N bonuses within a category</h3>
<ul>
    <li>One or more rules, each of type Bonuses per category, Affects Group-based awards, Ordinary scoring rule.</li>
    <li>Set Triggered when and Results in to the required values.</li>
</ul>
<h3>Award DNF if not enough categories scored</h3>
<ul>
    <li>Set a placeholder with a Triggered when value = N;</li>
    <li>Set a DNF if triggered rule with Triggered when = 0</li>
</ul>
<h3>Award DNF if too many categories scored</h3>
<ul>
    <li>Set a DNF is triggered rule with Triggered when set to the limit</li>
</ul>
</article>
`
const crintro = `<p>Complex rules and categories are used to enable scoring mechanisms beyond simple bonuses and combos.
	Scores generated here can apply to individual bonuses or to groups of bonuses.
	<input type="button" class="popover" popovertarget="complexhelp" value="[more details here]">
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
	<article class="popover" id="individualbonushelp" popover>
	<h1>Individual Bonuses</h1>
		<p>Some definitions:- <strong>BV</strong> = points value of current bonus; <strong>RV</strong> = the "results in" value of current rule; <strong>N</strong> is the number of bonuses within the category, <strong>N1</strong> is <strong>N</strong>-1; <strong>SV</strong> is the resulting score.</p>
		<p><strong>N</strong> is calculated as the number of bonuses per category, regardless of the setting of the "Calculate" flag. If the category is set to "any", "then </p>
		<p>If <strong>RV</strong> is 0, <strong>SV</strong> = <strong>BV</strong> * <strong>N1</strong>  simple multiplication.</p>
		<p>If <strong>RV</strong> is set to "multipliers", <strong>SV</strong> = <strong>BV</strong> * <strong>RV</strong> * <strong>N1</strong>  simple multiplication.</p>
		<p>If <strong>RV</strong> is set to "points", <strong>SV</strong> = <strong>BV</strong> * <strong>RV</strong> ^ <strong>N1</strong> exponential score.</p>
	</article>

<div id="singlerule" class="singlerule">
  <form action="updtcrule" method="post">
  <fieldset class="field rule0 rule1 rule2 rule3 rule4 rule5">
    <label for="RuleType">Rule type</label>
    <select id="RuleType" name="RuleType" onchange="chgRuleType(this);">
    %v
    </select>
  </fieldset>

  													<!-- Help texts -->
  <fieldset class="field help rule0">
	<p>GROUP AWARDS: Multipliers are applied to the entrant's total score, not just the bonuses contributing to this rule.</p>
	<p><input type="button" class="popover" popovertarget="individualbonushelp" value="INDIVIDUAL BONUSES: [click here for full explanation]"></p>
  </fieldset>

  <fieldset class="field help rule3">
  	<p>Placeholders are used in some circumstances to improve analysis of a set of rules.</p> 
  	<p>For example: A rule of type "DNF unless triggered" will appear on score explanations with a tick if it is triggered but will not appear on score explanations at all if it's not triggered, even though that condition will result in DNF. Mostly it would be better to use a placeholder and a "DNF if triggered" pair as that will give more satisfying results on score explanations.</p>
  </fieldset>
  
  <fieldset class="field help rule4">
  	<p>The sequence refers to the order in which claims are submitted, not the sequence within the rally book or geographical location. The result is either a fixed number of points or a multiple of the points scored by the bonuses forming the sequence.</p>
  </fieldset>

  <fieldset class="field help rule5">
	<p>Cat ratio rules describe an optional final test on a scorecard involving comparing N, the number of successful claims in 'First category' with M, the number of successful claims in 'Second category'.</p>
	<p>If R, the 'Ratio between cats' value is 1, N must equal M.
	If R is &gt; 1, N must be at least R x M</p>
  </fieldset>
												  <!-- End of help texts -->

  <fieldset class="field rule0 rule1 rule2 rule3 rule4 rule5">
    <label for="Axis">Category set</label>
	<select id="Axis" name="Axis" onchange="chgAxis(this);">
	%v
	</select>
  </fieldset>
  <fieldset class="field rule0 rule1 rule2 rule3 rule4 rule5">
    <label class="rule0 rule1 rule2 rule3 rule4" for="Cat">Which %v</label>
	<label class="rule5" for="Cat">First category</label>
    <select id="Cat" name="Cat" onchange="saveRule(this)">
    %v
    </select>
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
  <fieldset class="field rule0 rule1 rule2 rule3 rule4 rule5">
    <label class="rule0 rule1 rule2 rule3 rule4" for="NMin">Triggered when <var>n</var> &ge; </label>
	<label class="rule5" for="NMin">Ratio between cats</label>
	<input id="NMin" name="NMin" type="number" value="%v" onchange="saveRule(this)">
  </fieldset>

  <fieldset class="field rule0 rule4 ">
    <label for="PointsMults">This rule results in</label>
	<fieldset class="field rule0 rule4 " title="Result value">
		<input id="NPower" name="NPower" type="number" value="%v" onchange="saveRule(this)">
    	<select id="PointsMults" name="PointsMults" onchange="saveRule(this)">
    		%v
    	</select>
	</fieldset>
  </fieldset>

  <fieldset class="field rule5">
	<label for="NPowerSelect">Second category</label>
	<select id="NPowerSelect" name="NPower" onchange="saveRule(this)">
		%v
	</select>
  </fieldset>

  <input type="hidden" id="ruleid" name="ruleid" value="%v">
  </form>
</div>
<script>
showCurrentRule()
</script>`

func createRule(w http.ResponseWriter, r *http.Request) {

	sqlx := "INSERT INTO catcompound(Axis,Cat,NMethod,NMin,PointsMults,NPower,Ruletype,ModBonus) VALUES(1,0,0,0,0,0,0,0)"
	res, err := DBH.Exec(sqlx)
	checkerr(err)
	n, err := res.RowsAffected()
	checkerr(err)
	if n != 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"insert failed"}`)
		return
	}
	ruleid, err := res.LastInsertId()
	checkerr(err)
	fmt.Fprintf(w, `{"ok":true,"msg":"%v"}`, ruleid)

}

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
	fmt.Fprint(w, complexhelp)
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

	rtvals := []int{CAT_OrdinaryScoringRule, CAT_DNF_Unless_Triggered, CAT_DNF_If_Triggered, CAT_PlaceholderRule, CAT_OrdinaryScoringSequence, CAT_RatioRule}
	rtlabs := []string{"ordinary scoring rule", "DNF unless triggered", "DNF if triggered", "placeholder only", "uninterrupted sequence", "Cat ratio DNF"}

	page := fmt.Sprintf(tmpltSingleRule,
		selectOptionArray(rtvals, rtlabs, cr.Ruletype),
		strings.Join(axisopts, ""),
		AxisLabels[cr.Axis-1],
		strings.Join(optsSingleAxisCats(cr.Axis, cr.Cat), ""),

		selectOptionArray([]int{CAT_ModifyBonusScore, CAT_ModifyAxisScore}, []string{"individual bonuses", "group-based awards"}, cr.Target),
		selectOptionArray([]int{CAT_NumBonusesPerCatMethod, CAT_NumNZCatsPerAxisMethod}, []string{"Bonuses per category", "Categories scored"}, cr.Method),
		cr.Min,
		cr.Power,
		selectOptionArray([]int{CAT_ResultPoints, CAT_ResultMults, CAT_ResultCount}, []string{"points", "multipliers", "count"}, cr.PointsMults),

		strings.Join(optsSingleAxisCats(cr.Axis, cr.Power), ""),

		cr.Ruleid)

	w.Write([]byte(page))
}

func show_rules(w http.ResponseWriter, r *http.Request) {

	const leg = 0
	var rt = map[int]string{0: "ordinary", 1: "DNF unless", 2: "DNF if", 3: "dummy", 4: "sequence", 5: "cat ratio"}
	rules := build_compoundRuleArray(leg)
	axes := build_axisLabels()
	startHTML(w, "Complex rules")
	fmt.Fprint(w, complexhelp)
	fmt.Fprintf(w, `<div class="intro">%v</div>`, crintro)

	fmt.Fprint(w, `<div class="ruleset">`)
	fmt.Fprint(w, `<button class="plus" autofocus title="Add new rule" onclick="addRule()">+</button>`)
	fmt.Fprint(w, `<fieldset class="row hdr">`)
	fmt.Fprint(w, `<span class="col">Set</span><span class="col">Category</span>`)
	fmt.Fprint(w, `<span class="col">Type</span><span class="col">Threshold</span>`)
	fmt.Fprint(w, `<span class="col">Target</span><span class="col">Score</span>`)
	fmt.Fprint(w, `</fieldset>`)
	fmt.Fprint(w, `</div><hr></header>`)
	fmt.Fprint(w, `<div class="ruleset">`)
	lastAxis := -1
	lastCat := -1
	LastTarget := -1
	maincolor := true
	for _, cr := range rules {
		if cr.Axis != lastAxis || cr.Cat != lastCat || cr.Target != LastTarget {
			maincolor = !maincolor
			lastAxis = cr.Axis
			lastCat = cr.Cat
			LastTarget = cr.Target
		}
		rowcls := ""
		if !maincolor {
			rowcls = "altrow"
		}
		fmt.Fprintf(w, `<fieldset class="row target %v" data-rowid="%d" title="%v" onclick="showRule(this);">`, rowcls, cr.Ruleid, cr.Ruleid)
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, axes[cr.Axis-1])

		sqlx := fmt.Sprintf("SELECT BriefDesc FROM categories WHERE Axis=%d AND Cat=%d", cr.Axis, cr.Cat)
		cb := getStringFromDB(sqlx, "any")
		fmt.Fprintf(w, `<fieldset class="col">%s</fieldset>`, cb)
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
		case CAT_RatioRule:
			sqlx = fmt.Sprintf("SELECT BriefDesc FROM categories WHERE Axis=%d AND Cat=%d", cr.Axis, cr.Power)
			target = getStringFromDB(sqlx, "any")
			pm = "DNF"
			pts = ""
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
