package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

var tmplTeamHeaders = `
	<div class="intro">
	<p>Teams consist of two or more bikes planning and riding together as a team. One bike with a rider and a passenger is not a team, that's a crew.</p>
	<p>Team rules generally require that at least one team member is present in each bonus photo and bonus claims will normally be accepted from one team member throughout the rally.</p>
	<p>The team name can be shown on the certificate as well as the names of the team members.</p>
	</div>

<article id="teamnames" class="teamnames">
	<button id="addteam" class="plus" autofocus title="Add new team" onclick="addNewTeam(this)">+</button>

	{{range $ix,$el := .Teams}}
		{{if ne $el.TeamName ""}}
			<div class="teamname">
				<label for="TeamName{{$el.TeamID}}">Team {{$el.TeamID}} is</label>
				<input type="text" id="TeamName{{$el.TeamID}}" name="BriefDesc" data-team="{{$el.TeamID}}" value="{{$el.TeamName}}" onchange="saveTeamName(this)"  {{if eq $el.TeamID 0}} readonly {{else}}onclick="showTeamMembers(this)"{{end}}>
				{{if ne $el.TeamID 0}}
				<button class="plus" data-team="{{$el.TeamID}}" onclick="showTeamMembers(this)" >` + ordered_list_icon + `</button>
				{{end}}
			</div>
			{{end}}
	{{end}}


</article>
<hr>
<article class="teamMembers" id="teamMembers">
</article>
`

type teamrec struct {
	TeamID   int
	TeamName string
}

func fetchTeams() []teamrec {

	res := make([]teamrec, 0)
	sqlx := "SELECT TeamID,ifnull(BriefDesc,'Team ' || TeamID) FROM teams WHERE TeamID > 0 ORDER BY TeamID"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var x teamrec
		err = rows.Scan(&x.TeamID, &x.TeamName)
		checkerr(err)
		res = append(res, x)
	}
	return res
}

func list_teams(w http.ResponseWriter, r *http.Request) {

	type teams struct {
		Teams []teamrec
	}
	var tm teams
	tm.Teams = fetchTeams()

	t, err := template.New("teams").Parse(tmplTeamHeaders)
	checkerr(err)

	startHTML(w, "Teams")
	fmt.Fprint(w, `</header>`)
	err = t.Execute(w, tm)
	checkerr(err)

}

func setTeam(w http.ResponseWriter, r *http.Request) {

	e := intval(r.FormValue("e"))
	if e < 1 {
		fmt.Fprint(w, `{"ok": false,"msg":"incomplete request"}`)
		return
	}
	t := intval(r.FormValue("t"))
	sqlx := fmt.Sprintf("UPDATE entrants SET TeamID=%v WHERE EntrantID=%v", t, e)
	_, err := DBH.Exec(sqlx)
	checkerr(err)
	fmt.Fprint(w, `{"ok":true,"msg":"ok"}`)

}
func showTeamMembers(w http.ResponseWriter, r *http.Request) {

	type teamEntry struct {
		EntrantID                                        int
		RiderFirst, RiderLast, PillionFirst, PillionLast string
	}
	var team struct {
		OK       bool   `json:"ok"`
		Msg      string `json:"msg"`
		Team     int
		TeamName string
		Members  []teamEntry
	}
	TeamX := r.FormValue("t")
	team.Team = intval(TeamX)
	if team.Team < 1 {
		fmt.Fprint(w, `{"ok":false,"msg":"Bad team index"}`)
		return
	}
	team.TeamName = getStringFromDB("SELECT BriefDesc FROM teams WHERE TeamID="+TeamX, "")
	sqlx := "SELECT Entrantid,ifnull(RiderFirst,''),ifnull(RiderLast,''),ifnull(PillionFirst,''),ifnull(PillionLast,'') FROM entrants"
	sqlx += " WHERE TeamID=" + TeamX + " ORDER BY RiderLast,RiderFirst"
	rows, err := DBH.Query(sqlx)
	checkerr(err)
	defer rows.Close()
	for rows.Next() {
		var te teamEntry
		err = rows.Scan(&te.EntrantID, &te.RiderFirst, &te.RiderLast, &te.PillionFirst, &te.PillionLast)
		checkerr(err)
		team.Members = append(team.Members, te)
	}
	team.OK = true
	team.Msg = "ok"
	bytes, err := json.Marshal(team)
	checkerr(err)
	fmt.Fprintf(w, "%v", string(bytes))

}
