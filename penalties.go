package main

import (
	"fmt"
	"net/http"
	"text/template"
)

const tphelp = `
<div id="tphelp" class="popover" popover>
<h1>TIME PENALTIES</h1>
<p>Rally time runs from the rally start time to the rally finish time as specified in [Rally Parameters]</p>
<p>Individual entrants may have less time available. This can happen where variable start times are available (popular during the pandemic or where riders might start anywhere over a large area). The rally time window is set wide enough to cater for the variable starts and each entrant's time window is limited by the <em>Max Rideable hours</em> setting.</p>
<p>Penalties other than DNF can be specified for particular periods within the overall or individual entrant's rally time. Periods are specified either as date/time ranges or as minutes before DNF ranges.</p>
<p>Time penalties are triggered by entrant check-in (or last claim) time.</p>
<p>The start time specified identifies the first minute when a penalty is applied. So for a rally ending at 17:00 and a penalty start of 30 minutes before DNF, the penalty is incurred by anyone checking-in on or after 16:30.</p>
<p>Example:</p>
<ul>
<li> <em>Rally finishes</em> in [Rally Parameters] is "2026-04-12 21:00"</li>
<li> <em>Max Rideable hours</em> in [Rally Parameters] is 8</li>
<li> George Bone starts his rally at "2026-04-12 09:00"</li>
<li> <em>Rally DNF</em> is "2026-04-12 21:00"</li>
<li> <em>Entrant DNF</em> is "2026-04-12 17:00"</li>
</ul>
<p>Typically, penalties result in either a fixed number of points deducted or points deducted per mile/km (depending on the rally's unit of distance). It is also possible to affect the number of multipliers applied to the final score if they're being used in the rally.</p>
</div>
`

const tpheader = `
<div class="intro">
<p>Time penalties are typically imposed to discourage "last minute" finishes. <input type="button" class="popover" popovertarget="tphelp" value="[click here for more info]"></p>
</div>
<article class="timepenalties">
<button class="plus" autofocus title="Add new penalty" onclick="window.location.href='/timep/0?back=timep'">+</button>

<div class="row tp">
<span>Start</span><span>Finish</span>
<span>Time spec</span>
<span>Penalty method</span>
</span>Number</span>
</div>
</article>
<hr>
`

var tmpltListTimepenalties = `
<article class="timepenalties">
	{{range $el := .}}
		<div class="row tp" onclick="window.location.href='/timep/{{.Tpid}}?back=/timep'">
		<span>{{.PenaltyStartX}}</span>{{.PenaltyFinishX}}</span><span>{{if ne $el.TimeSpec 0}} mins < DNF{{end}}</span>
		<span>
		{{if eq .PenaltyMethod 0}}Deduct points{{end}}
		{{if eq .PenaltyMethod 1}}Deduct multipliers{{end}}
		{{if eq .PenaltyMethod 2}}Deduct points per minute{{end}}
		{{if eq .PenaltyMethod 3}}Deduct mults per minute{{end}}
		</span>
		<span>{{.PenaltyFactor}}</span>
		</div>
	{{end}}
</article>
`

func show_timepenalties(w http.ResponseWriter, r *http.Request) {

	Leg := intval(r.FormValue("leg"))
	if Leg < 1 {
		Leg = 1
	}
	TimePenalties = build_timePenaltyArray()

	startHTML(w, "Time penalties")

	fmt.Fprint(w, tphelp)
	fmt.Fprint(w, tpheader)
	fmt.Fprint(w, `</header>`)

	t, err := template.New("tplist").Parse(tmpltListTimepenalties)
	checkerr(err)
	err = t.Execute(w, TimePenalties)
	checkerr(err)
}

var tmplTimepenalty = `
{{$cls := ""}}

<article class="timepenalty">
<input id="tpid" type="hidden" value="{{.Tpid}}">
<fieldset>
	<label for="TimeSpec">Time specification</label>
	<select id="TimeSpec" name="TimeSpec" onchange="saveTimep(this)">
		<option value="0" {{if eq .TimeSpec 0}}selected{{end}}>Date/time</option>
		<option value="1" {{if eq .TimeSpec 1}}selected{{end}}>Minutes before Rally DNF</option>
		<option value="2" {{if eq .TimeSpec 2}}selected{{end}}>Minutes before Entrant DNF</option>
	</select>
</fieldset>
<fieldset>
	<label for="PenaltyStart">Starts from</label>
	{{if ne .TimeSpec 0}}{{$cls = "hide"}}{{end}}
		<input type="date" id="PenaltyStartDate" class="{{$cls}}" name="PenaltyStart" onchange="saveTimep(this)" value="{{.PenaltyStartDate}}">
		<input type="time" id="PenaltyStartTime" class="{{$cls}}" name="PenaltyStart" onchange="saveTimep(this)" value="{{.PenaltyStartTime}}">
	
		{{$cls = ""}}
	{{if eq .TimeSpec 0}}{{$cls = "hide"}}{{end}}

		<input type="number" id="PenaltyStartMins" title="Minutes before DNF" class="{{$cls}}" name="PenaltyStart" oninput="oi(this)" data=save="saveTimep" onchange="saveTimep(this)" value="{{.PenaltyStartMins}}">
	
</fieldset>

<fieldset>
	<label for="PenaltyFinish">Continues until</label>
	{{$cls = ""}}{{if ne .TimeSpec 0}}{{$cls = "hide"}}{{end}}
		<input type="date" id="PenaltyFinishDate" class="{{$cls}}" name="PenaltyFinish" onchange="saveTimep(this)" value="{{.PenaltyFinishDate}}">
		<input type="time" id="PenaltyFinishTime" class="{{$cls}}" name="PenaltyFinish" onchange="saveTimep(this)" value="{{.PenaltyFinishTime}}">
	{{$cls = ""}}
	{{if eq .TimeSpec 0}}{{$cls = "hide"}}{{end}}
		<input type="number" id="PenaltyFinishMins" class="{{$cls}}" title="Minutes before DNF" name="PenaltyFinish" oninput="oi(this)" data=save="saveTimep" onchange="saveTimep(this)" value="{{.PenaltyFinishMins}}">
	

</fieldset>

<fieldset>
	<label for="Penaltymethod">Penalty method</label>
	<select id="PenaltyMethod" name="PenaltyMethod" onchange="saveTimep(this)">
	<option value="0" {{if eq .PenaltyMethod 0}}selected{{end}}>Deduct points</option>
	<option value="1" {{if eq .PenaltyMethod 1}}selected{{end}}>Deduct multipliers</option>
	<option value="2" {{if eq .PenaltyMethod 2}}selected{{end}}>Deduct points per minute</option>
	<option value="3" {{if eq .PenaltyMethod 3}}selected{{end}}>Deduct multipliers per minute</option>
	</select>
</fieldset>

<fieldset>
</fieldset>
	<label for="PenaltyFactor">Points or multipliers</label>
	<input type="number" id="PenaltyFactor" name="PenaltyFactor" onchange="saveTimep(this)" value="{{.PenaltyFactor}}">
</article>
`

const tptopline = `	
<div class="topline">
		
		<fieldset>
			<button title="Delete this penalty?" onclick="enableDelete(!document.getElementById('enableDelete').checked)">   ` + TrashcanIcon + `</button>
			<input type="checkbox" style="display:none;" id="enableDelete" onchange="enableSave(this.checked)">
		</fieldset>
		
		<fieldset>
			<button id="updatedb" class="hideuntil" title="Delete Penalty" disabled onclick="deleteTimep(this)"></button>
		</fieldset>

		<fieldset>
			<button title="back to list" onclick="loadPage('/timep')">↥☰↥</button>
		</fieldset>

	</div>
`
const tpnewline = `	
<div class="topline">
		

		<fieldset>
			<button title="back to list" onclick="loadPage('/timep')">↥☰↥</button>
		</fieldset>

	</div>
`

const tpdetail = `

<div class="intro">
<p>Time penalties are typically imposed to discourage "last minute" finishes. <input type="button" class="popover" popovertarget="tphelp" value="[click here for more info]"></p>
</div>
`

func deleteTimePenalty(w http.ResponseWriter, r *http.Request) {

	tpid := r.PathValue("tpid")
	if tpid == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"incomplete request"}`)
		return
	}
	_, err := DBH.Exec("DELETE FROM timepenalties WHERE rowid=" + tpid)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":ok"}`)
}
func saveTimePenalty(w http.ResponseWriter, r *http.Request) {

	fld := r.FormValue("ff")
	tpid := r.FormValue("tpid")
	if fld == "" || tpid == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"incomplete request"}`)
		return
	}
	val := r.FormValue(fld)
	sqlx := "UPDATE timepenalties SET " + fld + "='" + val + "' WHERE rowid=" + tpid
	fmt.Println(sqlx)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)

}

func show_timepenalty(w http.ResponseWriter, r *http.Request) {

	var tp TimePenalty
	var sqlx string

	tpid := r.PathValue("rec")
	if tpid == "" {
		fmt.Fprint(w, `{"ok":false,"msg":"incomplete request"}`)
		return
	}
	if tpid == "0" { // Create new record
		sqlx = fmt.Sprintf("INSERT INTO timepenalties(TimeSpec,PenaltyStart,PenaltyFinish)VALUES(%v,%v,%v)", TimeSpecEntrantDNF, "30", "0")
		res, err := DBH.Exec(sqlx)
		checkerr(err)
		tp.Tpid, err = res.LastInsertId()
		checkerr(err)
		tp.TimeSpec = TimeSpecEntrantDNF
		tp.PenaltyStartMins = 30
		tp.PenaltyFinishMins = 0
	} else {
		TimePenalties = build_timePenaltyArray()
		for i := range TimePenalties {
			if TimePenalties[i].Tpid == int64(intval(tpid)) {
				tp = TimePenalties[i]
				break
			}
		}
	}
	// Got a good record now
	t, err := template.New("timep").Parse(tmplTimepenalty)
	checkerr(err)
	if r.FormValue("back") != "" {
		startHTMLBL(w, "Timepenalty", r.FormValue("back"))
	} else {
		startHTML(w, "Time penalty")
	}
	if tpid != "0" {
		fmt.Fprint(w, tptopline)
	} else {
		fmt.Fprint(w, tpnewline)
	}
	fmt.Fprint(w, tpdetail)
	fmt.Fprint(w, tphelp)
	fmt.Fprint(w, `</header>`)
	err = t.Execute(w, tp)
	checkerr(err)
}
