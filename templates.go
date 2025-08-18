package main

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const helpicon = "&nbsp;?&nbsp;"
const homeicon = " &#127968; "

//go:embed images/alertalert.b64
var iconalert string

//go:embed images/alertbike.b64
var iconbike string

//go:embed images/alertdaylight.b64
var icondaylight string

//go:embed images/alertface.b64
var iconface string

//go:embed images/alertnight.b64
var iconnight string

//go:embed images/alertrestricted.b64
var iconrestricted string

//go:embed images/alertreceipt.b64
var iconreceipt string

//go:embed chasm.js
var mainscript string

//go:embed normalize.css
var normalize string

//go:embed chasm.css
var maincss string

const TrashcanIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-trash" viewBox="0 0 16 16">
  <path d="M5.5 5.5A.5.5 0 0 1 6 6v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m2.5 0a.5.5 0 0 1 .5.5v6a.5.5 0 0 1-1 0V6a.5.5 0 0 1 .5-.5m3 .5a.5.5 0 0 0-1 0v6a.5.5 0 0 0 1 0z"/>
  <path d="M14.5 3a1 1 0 0 1-1 1H13v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V4h-.5a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1H6a1 1 0 0 1 1-1h2a1 1 0 0 1 1 1h3.5a1 1 0 0 1 1 1zM4.118 4 4 4.059V13a1 1 0 0 0 1 1h6a1 1 0 0 0 1-1V4.059L11.882 4zM2.5 3h11V2h-11z"/>
</svg>`

const FloppyDiskIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-floppy" viewBox="0 0 16 16">
  <path d="M11 2H9v3h2z"/>
  <path d="M1.5 0h11.586a1.5 1.5 0 0 1 1.06.44l1.415 1.414A1.5 1.5 0 0 1 16 2.914V14.5a1.5 1.5 0 0 1-1.5 1.5h-13A1.5 1.5 0 0 1 0 14.5v-13A1.5 1.5 0 0 1 1.5 0M1 1.5v13a.5.5 0 0 0 .5.5H2v-4.5A1.5 1.5 0 0 1 3.5 9h9a1.5 1.5 0 0 1 1.5 1.5V15h.5a.5.5 0 0 0 .5-.5V2.914a.5.5 0 0 0-.146-.353l-1.415-1.415A.5.5 0 0 0 13.086 1H13v4.5A1.5 1.5 0 0 1 11.5 7h-7A1.5 1.5 0 0 1 3 5.5V1H1.5a.5.5 0 0 0-.5.5m3 4a.5.5 0 0 0 .5.5h7a.5.5 0 0 0 .5-.5V1H4zM3 15h10v-4.5a.5.5 0 0 0-.5-.5h-9a.5.5 0 0 0-.5.5z"/>
</svg>`

var htmlheader = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>chasm</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="/css?file=normalize">
<link rel="stylesheet" href="/css?file=maincss">
<script src="/js?file=mainscript"></script>
</head>
<body>
`

var topbar = `
<nav class="topbar">
	<span class="flexitem">
	<!--
	<button id="main_help_button" class="link noprint" onclick="showHelp('')" title="Help">` + helpicon + `</button>
	-->
	<button id="main_home_button" class="link noprint" onclick="goHome(this)" title="Main menu">` + homeicon + `</button>
	<span id="main_rally_title" class="link" onclick="loadPage('%s')">%s</span>
	</span>
	<span class="flexitem">
	<span id="main_current_task">%s</span>
	<button id="about_chasm" class="link noprint" onclick="showAboutChasm(this)" title="About ScoreMaster">&copy;</button>
	</span>
</nav>
`
var reloadticker = `
<div class="reloadticker noprint">
	<progress id="reloadticker" data-active="1" max="30" value="30" title="refreshing soon"></progress>
	<script>
	setInterval(function() {
		let p = document.getElementById('reloadticker')
		let s = p.getAttribute('value')
		let a = p.getAttribute('data-active')
		if (a=='1') {
			s--
			if (s < 1) {window.location.href=#url#}
			p.setAttribute('value',s)
		}
	},1000)
	function killReload() {
		let p = document.getElementById('reloadticker')
		p.setAttribute('data-active','0')
		p.classList.add('hide')
	}
	</script>
</div>
`

func niy(w http.ResponseWriter, r *http.Request) {

	startHTML(w, "NIY")

	fmt.Fprintf(w, `<p class="error">NOT IMPLEMENTED YET</p><p>%v</p>`, r)
}

func printNZ(i int) string {
	if i == 0 {
		return ""
	}
	return strconv.Itoa(i)
}
func jsonBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func builtin_images(w http.ResponseWriter, r *http.Request) {

	img := r.FormValue("i")
	var imgdata string
	switch img {
	case "alert":
		imgdata = iconalert
	case "bike":
		imgdata = iconbike
	case "daylight":
		imgdata = icondaylight
	case "face":
		imgdata = iconface
	case "night":
		imgdata = iconnight
	case "restricted":
		imgdata = iconrestricted
	case "receipt":
		imgdata = iconreceipt
	default:
		fmt.Fprint(w, `<p class="error">no such img</p>`)
		return
	}
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgdata))
	w.Header().Set("Content-Type", "image/png;")
	_, err := io.Copy(w, dec)
	checkerr(err)

	//fmt.Printf("Img %v sent %v bytes (%v)\n", img, n, len(imgdata))
}

func fmtDecimal(ptn string, n float64) string {

	x := fmt.Sprintf(ptn, n)
	if CS.Basics.RallyPointIsComma {
		x = strings.Replace(x, ".", ",", 1)
	}
	return x
}

func send_css(w http.ResponseWriter, r *http.Request) {

	file := r.FormValue("file")
	if file == "" {
		file = "maincss"
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	switch file {
	case "normalize":
		fmt.Fprint(w, normalize)
	case "maincss":
		fmt.Fprint(w, maincss)
	case "certscss":
		fmt.Fprint(w, certscss)
	case "joditcss":
		fmt.Fprint(w, joditcss)
	}
}

func send_js(w http.ResponseWriter, r *http.Request) {

	file := r.FormValue("file")
	if file == "" {
		file = "mainscript"
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	switch file {
	case "mainscript":
		fmt.Fprint(w, mainscript)
	case "certeditjs":
		fmt.Fprint(w, certeditjs)
	case "joditjs":
		fmt.Fprint(w, joditjs)
	case "odosjs":
		fmt.Fprint(w, odosjs)
	}
}

func showReloadTicker(w http.ResponseWriter, url string) {

	fmt.Fprint(w, strings.ReplaceAll(reloadticker, "#url#", fmt.Sprintf("'%v'", url)))
}

func showTopbar(w http.ResponseWriter, currentTask string) {

	showTopbarBL(w, currentTask, "")

}
func showTopbarBL(w http.ResponseWriter, currentTask string, backLink string) {

	itm := "/"
	if backLink != "" {
		itm = backLink
	}
	fmt.Fprintf(w, topbar, itm, CS.Basics.RallyTitle, currentTask)

}

// splitDatetime splits a timestamp as stored in the database
// into separate date and time strings
func splitDatetime(dt string) (string, string) {

	if !strings.Contains(dt, "T") {
		return dt, dt
	}
	res := strings.Split(dt, "T")
	return res[0], res[1]

}
func startHTML(w http.ResponseWriter, currentTask string) {

	startHTMLBL(w, currentTask, "")

}
func startHTMLBL(w http.ResponseWriter, currentTask string, backLink string) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Fprint(w, htmlheader)
	fmt.Fprint(w, `<header>`)
	showTopbarBL(w, currentTask, backLink)

}
