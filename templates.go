package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"
)

const helpicon = "&nbsp;?&nbsp;"
const homeicon = " &#127968; "

//go:embed chasm.js
var mainscript string

//go:embed chasm.css
var maincss string

var htmlheader = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>chasm</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>` + maincss + `
</style>
<script>` + mainscript + `
</script>
</head>
<body>
`

var topbar = `
<nav class="topbar">
	<span class="flexitem">
	<button id="main_help_button" class="link" onclick="showMainHelp(this)" title="Help">` + helpicon + `</button>
	<button id="main_home_button" class="link" onclick="goHome(this)" title="Main menu">` + homeicon + `</button>
	<span id="main_rally_title" class="link" onclick="goHome(this)">%s</span>
	</span>
	<span class="flexitem">
	<span id="main_current_task">%s</span>
	<button id="about_chasm" class="link" onclick="showAboutChasm(this)" title="About ScoreMaster">&copy;</button>
	</span>
</nav>
`
var reloadticker = `
<div class="reloadticker">
	<progress id="reloadticker" max="30" value="30" title="refreshing soon"></progress>
	<script>
	setInterval(function() {
		let p = document.getElementById('reloadticker')
		let s = p.getAttribute('value')
		s--
		if (s < 1) {window.location.href=#url#}
		p.setAttribute('value',s)
	},1000)
	</script>
</div>
`

func showReloadTicker(w http.ResponseWriter, url string) {

	fmt.Fprint(w, strings.ReplaceAll(reloadticker, "#url#", fmt.Sprintf("'%v'", url)))
}

func showTopbar(w http.ResponseWriter, currentTask string) {

	fmt.Fprintf(w, topbar, CS.RallyTitle, currentTask)

}

func startHTML(w http.ResponseWriter, currentTask string) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Fprint(w, htmlheader)
	showTopbar(w, currentTask)

}
