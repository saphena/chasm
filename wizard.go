package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type wizardHostParams struct {
	Host     string
	Port     string
	Userid   string
	Password string
}

type wizardParams struct {
	RallyTitle        string
	RallyStartDate    string
	RallyStartTime    string
	RallyFinishDate   string
	RallyFinishTime   string
	RallyMaxHours     int
	RallyTimezone     string
	RallyTimezones    []string
	RallyUnitKms      bool
	RallyStartOption  int
	RallyFinishOption bool
	RallyMinPoints    int
	RallyUseQA        bool
	RallyQAPoints     int
	RallyMinMiles     int
	UseEBC            bool
	UseSMTP           bool
	SMTP              wizardHostParams
	IMAP              wizardHostParams
	UnitMilesLit      string
	UnitKmsLit        string
}

const wizStartPage = `
<article class="wizard">
	<h1>RALLY SETUP</h1>
	<form action="/wiz" method="post">
	<input type="hidden" name="wizsave" value="setup1">
	<fieldset>
		<label for="RallyTitle">Rally title</label>
		<input type="text" id="RallyTitle" name="RallyTitle" class="RallyTitle" value="{{.RallyTitle}}">
	</fieldset>
	<fieldset title="Earliest possible start time">
		<label for="RallyStartDate">Rally starts</label>
		<input type="date" id="RallyStartDate" name="RallyStartDate" class="RallyStartDate" value="{{.RallyStartDate}}">
		<input type="time" id="RallyStartTime" name="RallyStartTime" class="RallyStartTime" value="{{.RallyStartTime}}">
	</fieldset>
	<fieldset title="Latest finish time - later = DNF">
		<label for="RallyFinishDate">Rally Finishes</label>
		<input type="date" id="RallyFinishDate" name="RallyFinishDate" class="RallyFinishDate" value="{{.RallyFinishDate}}">
		<input type="time" id="RallyFinishTime" name="RallyFinishTime" class="RallyFinishTime" value="{{.RallyFinishTime}}">
	</fieldset>
	<!--
	-->
	<fieldset>
		<label for="RallyTimezone">Rally timezone</label>
		<select id="RallyTimezone" name="RallyTimezone" class="RallyTimezone">
		{{range .RallyTimezones}}
			<option value="{{ . }}" {{if eq . $.RallyTimezone}}selected{{end}}>{{ . }}</option>
		{{end}}
		</select>
	</fieldset>
	<fieldset title="Report distances in miles or kilometres">
		<label for="RallyDistanceUnit">Unit of distance</label>
		<select id="RallyDistanceUnit" name="RallyDistanceUnit" class="RallyDistanceUnit">
			<option value="M" {{if .RallyUnitKms}}{{else}}selected{{end}}>Mile</option>
			<option value="K" {{if .RallyUnitKms}}selected{{end}}</option>
		</select>
	</fieldset>
	<fieldset class="wiznav">
		<input type="hidden" id="wiznext" name="wizpage" value="2">
		<button>Next</button>
	</fieldset>
	</form>
</article>
`

const wizPage2 = `
<article class="wizard">
<script>
function wizp2next(btn) {
let nxt = document.getElementById('wiznext');
if (!nxt) return true;
let ebc = document.getElementById('UseEBC');
let smtp = document.getElementById('UseSMTP');
if ((ebc && ebc.value=='1') || (smtp && smtp.value=='1'))
	nxt.value='email';
else
	nxt.value='score1';
return true;
}
</script>
	<h1>RALLY SETUP</h1>
	<form action="/wiz">
	<input type="hidden" name="wizsave" value="setup2">
	<fieldset class="hide">
		<label for="RallyMaxHours">Max rally hours</label>
		<input type="number" id="RallyMaxHours" name="RallyMaxHours" class="RallyMaxHours" value="{{.RallyMaxHours}}">
		<span>This is the maximum number of hours, including rest stops, available to an entrant. This number might be smaller than the number of hours between <em>Rally starts</em> and <em>Rally finishes</em>.</span>
	</fieldset>

	<fieldset>
		<label for="RallyStartOption">How will riders START their rally</label>
		<select id="RallyStartOption" name="RallyStartOption" class="RallyStartOption" autofocus>
			<option value="0" {{if eq .RallyStartOption 0}}selected{{end}}>during a formal check-out process</option>
			<option value="1" {{if eq .RallyStartOption 1}}selected{{end}}>with their first bonus claim</option>
		</select>
	</fieldset>
	<fieldset>
		<label for="RallyFinishOption">How will riders FINISH their rally</label>
		<select id="RallyFinishOption" name="RallyFinishOption" class="RallyFinishOption">
			<option value="0" {{if .RallyFinishOption}}{{else}}selected{{end}}>with a formal check-in process</option>
			<option value="1" {{if .RallyFinishOption}}selected{{end}}>with their final bonus claim</option>
		</select>
	</fieldset>
	<fieldset>
		<label for="UseEBC">How will bonus claims be entered</label>
		<select id="UseEBC" name="UseEBC" class="UseEBC">
			<option value="0" {{if .UseEBC}}{{else}}selected{{end}}>manually, the old-fashioned way</option>
			<option value="1" {{if .UseEBC}}selected{{end}}>via EBC (email) throughout the rally</option>
		</select>
	</fieldset>
	<fieldset>
		<label for="UseSMTP">Do you want to send emails from the system</label>
		<select id="UseSMTP" name="UseSMTP" class="UseSMTP">
			<option value="0" {{if .UseSMTP}}{{else}}selected{{end}}>no, not even in test mode</option>
			<option value="1" {{if .UseSMTP}}selected{{end}}>yes - EBC test mode, sending files, etc</option>
		</select>
	</fieldset>
	<fieldset class="wiznav">
		<input type="hidden" id="wiznext" name="wizpage" value="email">
		<button onclick="document.getElementById('wiznext').value='1'">Previous</button>
		<button onclick="wizp2next(this)">Next</button>
	</fieldset>

	</form>
</article>
`

const wizPageEmail = `
<article class="wizard">
	<h1>EMAIL SETUP</h1>
	<form action="/wiz" method="post">
	<input type="hidden" name="wizsave" value="email">

	<h2>SMTP</h2>
	<fieldset>
		<label for="HostSMTP">Host</label>
		<input type="text" id="HostSMTP" name="HostSMTP" class="HostSMTP" value="{{.SMTP.Host}}">
		<input type="text" name="PortSMTP" class="PortSMTP" placeholder="port" value="{{.SMTP.Port}}" title="port">
	</fieldset>
	<fieldset>
		<label for="UseridSMTP">Userid</label>
		<input type="text" id="UseridSMTP" name="UseridSMTP" class="UseridSMTP" value="{{.SMTP.Userid}}">
	</fieldset>
	<fieldset>
		<label for="PasswordSMTP">Password</label>
		<input type="text" id="PasswordSMTP" name="PasswordSMTP" class="PasswordSMTP" value="{{.SMTP.Password}}">
	</fieldset>
	<h2>IMAP</h2>
	<fieldset>
		<label for="HostIMAP">Host</label>
		<input type="text" id="HostIMAP" name="HostIMAP" class="HostIMAP" value="{{.IMAP.Host}}">
		<input type="text" name="PortIMAP" class="PortIMAP" placeholder="port" value="{{.IMAP.Port}}" title="port">
	</fieldset>
	<fieldset>
		<label for="UseridIMAP">Userid</label>
		<input type="text" id="UseridIMAP" name="UseridIMAP" class="UseridIMAP" value="{{.IMAP.Userid}}">
	</fieldset>
	<fieldset>
		<label for="PasswordSMTP">Password</label>
		<input type="text" id="PasswordIMAP" name="PasswordIMAP" class="PasswordIMAP" value="{{.IMAP.Password}}">
	</fieldset>
	<div>
		<p>These details can be completed or updated later if you wish.</p>
	</div>
	<fieldset class="wiznav">
		<input type="hidden" id="wiznext" name="wizpage" value="email">
		<button onclick="document.getElementById('wiznext').value='2'">Previous</button>
		<button onclick="document.getElementById('wiznext').value='score1'">Next</button>
	</fieldset>

	</form>
</article>
`

const wizScoring1 = `
<article class="wizard">
	<h1>SCORING OPTIONS</h1>
	<form action="/wiz" method="post">
	<input type="hidden" name="wizsave" value="score1">

	<fieldset>
		<label for="RallyUseQA">Will you use any question/answer pairs</label>
		<select id="RallyUseQA" name="RallyUseQA" class="RallyUseQA" onchange="document.getElementById('RallyQAPointsFS').classList.toggle('hide')">
			<option value="0" {{if .RallyUseQA}}{{else}}selected{{end}}>no, just basic bonuses</option>
			<option value="1" {{if .RallyUseQA}}selected{{end}}>yes, at least one bonus includes a Q/A</option>
		</select>
	</fieldset>
	<fieldset id="RallyQAPointsFS" class="{{if .RallyUseQA}}{{else}}hide{{end}}">
		<label for="RallyQAPoints">Points value of valid answers</label>
		<input type="number" id="RallyQAPoints" name="RallyQAPoints" class="RallyQAPoints" value="{{.RallyQAPoints}}">
	</fieldset>
	<fieldset>
		<label for="RallyMinPoints">Minimum points needed to qualify as Finisher</label>
		<input type="number" id="RallyMinPoints" name="RallyMinPoints" class="RallyMinPoints" value="{{.RallyMinPoints}}">
	</fieldset>
	<fieldset>
		<label for="RallyMinMiles">Minimum {{if .RallyUnitKms}}{{.UnitKmsLit}}{{else}}{{.UnitMilesLit}}{{end}} needed to qualify as Finisher</label>
		<input type="number" id="RallyMinMiles" name="RallyMinMiles" class="RallyMinMiles" value="{{.RallyMinMiles}}">
	</fieldset>
	<fieldset class="wiznav">
		<input type="hidden" id="wiznext" name="wizpage" value="email">
		<button onclick="document.getElementById('wiznext').value='{{if .UseEBC}}email{{else if .UseSMTP}}email{{else}}2{{end}}'">Previous</button>
		<button onclick="document.getElementById('wiznext').value='finish'">Next</button>
	</fieldset>

</article>
`

const wizFinish = `
<article class="wizard">
	<h1>WIZARD COMPLETE</h1>
	<form action="/wiz" method="post">
	<input type="hidden" name="wizsave" value="finish">
	<div>
		<p>		This initial wizard is now complete.		</p>
		<p>You should now continue with the rally setup process using the facilities under "Rally setup &amp; config".</p>
		<p>If you want to use anything more complex than simple bonuses and combos, you will need to establish one or more sets of bonus categories before you can create rules affecting groups or sequences of bonuses.</p>
		<p>You might want to award extra points based on group membership, for example, or make scoring in particular categories mandatory. </p>
		<p>It is also possible to impose time or distance based penalties, different certificates for different classes of Finisher, etc.</p>
	</div>
		<fieldset class="wiznav">
		<input type="hidden" id="wiznext" name="wizpage" value="">
		<button onclick="document.getElementById('wiznext').value='score1'">Previous</button>
		<button onclick="document.getElementById('wiznext').value='end'">Finish</button>
	</fieldset>

	</form>
</article>
`

func buildRallyWizVars() wizardParams {

	var wp wizardParams
	wp.RallyTitle = CS.Basics.RallyTitle
	wp.RallyStartDate = CS.Basics.RallyStarttime[0:10]
	wp.RallyStartTime = CS.Basics.RallyStarttime[11:]
	wp.RallyFinishDate = CS.Basics.RallyFinishtime[0:10]
	wp.RallyFinishTime = CS.Basics.RallyFinishtime[11:]
	wp.RallyMaxHours = CS.Basics.RallyMaxHours
	wp.RallyTimezone = CS.Basics.RallyTimezone
	wp.RallyTimezones = tzlist
	wp.RallyUnitKms = CS.Basics.RallyUnitKms
	wp.RallyStartOption = CS.StartOption
	wp.RallyFinishOption = CS.AutoFinisher
	wp.RallyMinPoints = CS.RallyMinPoints
	wp.RallyUseQA = CS.RallyUseQA
	wp.RallyQAPoints = CS.RallyQAPoints
	wp.RallyMinMiles = CS.RallyMinMiles
	wp.UseEBC = CS.UseEBC
	wp.UseSMTP = CS.Email.UseSMTP
	wp.SMTP.Host = CS.Email.SMTP.Host
	wp.SMTP.Port = CS.Email.SMTP.Port
	wp.SMTP.Userid = CS.Email.SMTP.Userid
	wp.SMTP.Password = CS.Email.SMTP.Password
	wp.IMAP.Host = CS.Email.IMAP.Host
	wp.IMAP.Port = CS.Email.IMAP.Port
	wp.IMAP.Userid = CS.Email.IMAP.Userid
	wp.IMAP.Password = CS.Email.IMAP.Password
	wp.UnitMilesLit = CS.UnitMilesLit
	wp.UnitKmsLit = CS.UnitKmsLit

	return wp
}

func saveWizard(r *http.Request) {

	savepg := r.FormValue("wizsave")
	switch savepg {
	case "setup1":
		CS.Basics.RallyTitle = r.FormValue("RallyTitle")
		CS.Basics.RallyStarttime = r.FormValue("RallyStartDate") + "T" + r.FormValue("RallyStartTime")
		CS.Basics.RallyFinishtime = r.FormValue("RallyFinishDate") + "T" + r.FormValue("RallyFinishTime")
		CS.Basics.RallyTimezone = r.FormValue("RallyTimezone")
		CS.Basics.RallyUnitKms = r.FormValue("RallyUseKms") == "K"
		CS.Basics.RallyMaxHours = calcMaxHours(CS.Basics.RallyStarttime, CS.Basics.RallyFinishtime)
	case "setup2":
		CS.Basics.RallyMaxHours = intval(r.FormValue("RallyMaxHours"))
		CS.StartOption = intval(r.FormValue("RallyStartOption"))
		CS.AutoFinisher = r.FormValue("RallyFinishOption") == "1"
		CS.UseEBC = r.FormValue("UseEBC") == "1"
		CS.Email.UseSMTP = r.FormValue("UseSMTP") == "1"
	case "email":
		CS.Email.SMTP.Host = r.FormValue("HostSMTP")
		CS.Email.SMTP.Port = r.FormValue("PortSMTP")
		CS.Email.SMTP.Userid = r.FormValue("UseridSMTP")
		CS.Email.SMTP.Password = r.FormValue("PasswordSMTP")
		CS.Email.IMAP.Host = r.FormValue("HostIMAP")
		CS.Email.IMAP.Port = r.FormValue("PortIMAP")
		CS.Email.IMAP.Userid = r.FormValue("UseridIMAP")
		CS.Email.IMAP.Password = r.FormValue("PasswordIMAP")
	case "score1":
		CS.RallyUseQA = r.FormValue("RallyUseQA") == "1"
		CS.RallyQAPoints = intval(r.FormValue("RallyQAPoints"))
		CS.RallyMinPoints = intval(r.FormValue("RallyMinPoints"))
		CS.RallyMinMiles = intval(r.FormValue("RallyMinMiles"))
	case "finish":
		CS.ShowSetupWizard = false
	}
	saveSettings()

}

func showWizard(w http.ResponseWriter, r *http.Request) {

	var err error
	var t *template.Template
	savepg := r.FormValue("wizsave")
	if savepg != "" {
		saveWizard(r)
	}
	rv := buildRallyWizVars()
	pg := r.PathValue("page")
	if pg == "" {
		pg = r.FormValue("wizpage")
	}
	switch pg {
	case "2":
		t, err = template.New("wizpage2").Parse(wizPage2)
		checkerr(err)
	case "email":
		t, err = template.New("wizpageemail").Parse(wizPageEmail)
		checkerr(err)
	case "score1":
		t, err = template.New("wizscoring").Parse(wizScoring1)
		checkerr(err)
	case "finish":
		t, err = template.New("wizfinish").Parse(wizFinish)
		checkerr(err)
		CS.ShowSetupWizard = false

	case "end":
		central_dispatch(w, r)
		return
	default:
		t, err = template.New("wizpage").Parse(wizStartPage)
		checkerr(err)
	}
	startHTML(w, "Rally setup wizard")
	fmt.Fprint(w, `</header>`)

	err = t.Execute(w, rv)
	checkerr(err)
}
