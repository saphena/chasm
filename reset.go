package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//go:embed basicdemo.sql
var basicdemosql string

//go:embed basicdemo.json
var basicdemojson string

const authcode = "UtterGobbledygook and then some more"

var resetDatabaseForm = `
<script>
    function doit(obj) {
        obj.disabled = true;
        let choice;
        console.log(obj.id);
        if (obj.id == "firstchoiceplease") choice = document.getElementById('firstchoice');
        if (obj.id == "choice1please") choice = document.getElementById('choice1');
        if (obj.id == "choice2please") choice = document.getElementById('choice2');
        if (obj.id == "choice3please") choice = document.getElementById('choice3');
        if (obj.id == "choice4please") choice = document.getElementById('choice4');
        if (choice) 
            choice = choice.value;
        else
            choice = 0;
        if (choice == 0) {
            window.location.href = "/";
            return;
        }
        if (obj.id == "firstchoiceplease") {
            let c1 = document.getElementById('firstchoice'+choice);
            if (c1) c1.classList.remove('hide');
            return;
        } else {
            let frm = document.getElementById('zapper');
            if (!frm) return;
            let lvl = document.getElementById('zaplevel');
            if (!lvl) return;
            lvl.value = choice;
            frm.submit();
            return false;
        }
        window.location.href = "/";
    }
</script>

    <form id="zapper" action="/reset" method="post">
    <input type="hidden" name="cmd" value="zap">
    <input type="hidden" name="zaplevel" id="zaplevel" value="1">
    <input type="hidden" name="authcode" value="` + authcode + `">
    </form>

    <article class="resetdb">
    <h1>RESET THE DATABASE</h1>
    <p>This procedure will <strong>RESET THE DATABASE</strong> back to an initial state depending on the settings below.</p>
    <p>Once triggered, this procedure cannot be stopped and it <strong>CANNOT BE REVERSED</strong>.</p>
    <p>I offer four levels of reset:</p>
    <ol>
    <li>Remove all scoring info including claims. Rally is ready for live running.</li>
    <li>Remove all claims and entrants. Rally is ready for entrant loading before rally.</li>
    <li>Remove claims, entrants, bonuses, combos and other config data. Need to full configure rally.</li>
	<li>Reload the demo database and bring it up to date.</li>
    </ol>
    <fieldset><label for="firstchoice">What is your desire at this stage?</label>
    <select id="firstchoice">
    <option value="0">Get me back to safety please</option>
    <option value="1">1 - Just clear out my testing claims, etc</option>
    <option value="2">2 - Clear all scoring and entrants</option>
    <option value="3">3 - Rebuild everything from scratch</option>
    <option value="4">4 - Rally testing - reset demo database</option>
    </select> <button id="firstchoiceplease" onclick="doit(this)">Do it now!</button></fieldset>`

func doTheReset(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "Database reset")

	//fmt.Printf("doTheReset called with level %v\n", r.FormValue("zaplevel"))

	//_, err := DBH.Exec("BEGIN TRANSACTION")
	//checkerr(err)
	//defer DBH.Exec("ROLLBACK")

	zl := intval(r.FormValue("zaplevel"))

	switch zl {
	case 1:
		zapAllClaims(true)
		resetScorecardReviews()
		recalc_all()
	case 2:
		zapAllClaims(true)
		zapEntrants()
	case 3:
		zapAllClaims(true)
		zapEntrants()
		sqlx := "UPDATE config SET Settings='{}'"
		_, err := DBH.Exec(sqlx)
		checkerr(err)
		zapRallyConfig()

	case 4:
		zapAllClaims(true)
		zapRallyConfig()
		reloadDemoRally()
		// Rally testing demo reset
		// reset dateranges
		// zap claims but unprocess ebclaims
		// rebuild scorecards
		// reset entrant status/odos
	}
	//_, err = DBH.Exec("COMMIT")
	//checkerr(err)
	fmt.Fprint(w, `</header><p class="thatsall">Reset complete</p>`)
	//fmt.Println("Reset complete")
}

func reloadDemoRally() {

	const OrigDate = "2025-06-21"

	todaysDate := time.Now().Format(time.DateOnly)
	fmt.Printf("Resetting demo from %v to %v\n", OrigDate, todaysDate)

	sqlx := strings.ReplaceAll(basicdemosql, OrigDate, todaysDate)
	//fmt.Println(sqlx)
	_, err := DBH.Exec(sqlx)
	checkerr(err)

	cfg := strings.ReplaceAll(basicdemojson, OrigDate, todaysDate)
	sqlx = "UPDATE config SET Settings=?"
	stmt, err := DBH.Prepare(sqlx)
	checkerr(err)
	defer stmt.Close()
	_, err = stmt.Exec(cfg)
	checkerr(err)
	loadJsonConfigs()
}

func resetScorecardReviews() {

	sqlx := fmt.Sprintf("UPDATE entrants SET ReviewStatus=%v", rs_notreviewed)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
}

func showResetChoiceConfirmation(w http.ResponseWriter, lvl int, txt string) {

	fmt.Fprintf(w, `<div id="firstchoice%v" class="hide">
    <hr><p>You have chosen to %v.</p>
    <p class="yellow">THERE IS NO UNDO FACILITY IF YOU GO AHEAD WITH THIS!</p>
    <fieldset><label for="choice%v">Are you really sure you want to do this?</label>
    <select id="choice%v">
    <option value="0" selected>No! Get me back to safety please</option>
    <option value="%v">I know what I'm doing, just get on with it</option>
    </select> <button id="choice%vplease" onclick="doit(this)">Do it now!</button></fieldset>
    </div>`, lvl, txt, lvl, lvl, lvl, lvl)

}
func showResetOptions(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	//fmt.Printf("showResetOptions ac='%v' == '%v'\n", authcode, r.FormValue("authcode"))

	if r.FormValue("cmd") != "" && r.FormValue("zaplevel") != "" && r.FormValue("authcode") == authcode {
		doTheReset(w, r)
		return
	}

	startHTML(w, "Reset Database")

	fmt.Fprint(w, `</header>`)

	fmt.Fprint(w, resetDatabaseForm)

	showResetChoiceConfirmation(w, 1, "clear out all bonus claims, clear the scorecards, reset start times and make the rally ready for a live start")
	showResetChoiceConfirmation(w, 2, "clear out all bonus claims and DELETE THE ENTRANTS, leaving the rally ready to load the entrants")
	showResetChoiceConfirmation(w, 3, "clear out EVERYTHING and build the rally from scratch")
	showResetChoiceConfirmation(w, 4, "reset the testing database")
	fmt.Fprint(w, `</article>`)

	fmt.Fprint(w, `</body></html>`)

}

func zapAllClaims(zapEBC bool) {

	sqlx := "DELETE FROM claims"
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	if zapEBC {
		sqlx = "DELETE FROM ebclaims"
		_, err = DBH.Exec(sqlx)
		checkerr(err)
		//zapEBPhotoImages()
		sqlx = "DELETE FROM ebcphotos"
		_, err = DBH.Exec(sqlx)
		checkerr(err)

	} else {
		sqlx = "UPDATE ebclaims SET Decision=-1,Processed=0"
		_, err = DBH.Exec(sqlx)
		checkerr(err)
	}

}

func zapRallyConfig() {

	for _, x := range []string{"bonuses", "combos", "categories", "catcompound"} {
		sqlx := "DELETE FROM " + x
		_, err := DBH.Exec(sqlx)
		checkerr(err)
	}
	sqlx := "UPDATE config SET settings='{}'"
	_, err := DBH.Exec(sqlx)
	checkerr(err)

}

func zapEntrants() {

	sqlx := "DELETE FROM entrants"
	_, err := DBH.Exec(sqlx)
	checkerr(err)
}
