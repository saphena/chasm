package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ricochet2200/go-disk-usage/du"
)

const PROGRAMTITLE = "ScoreMaster"
const PROGRAMVERSION = "4.0"

const CopyriteYear = "2024"
const ChasmVersion = "0.1"

var EBCFetchVersion string = "0.0" // Loaded at runtime

const PROGDESC = "An application designed to make scoring &amp; administration of IBA style motorcycle rallies easy"

const Author = "Bob Stammers (IBA #51220)"
const CopyriteHolder = "Bob Stammers"

var InspiredBy = []string{
	"Chris Kilner #40058",
	"Steve Eversfield #169",
	"Lee Edwards #59974",
	"Robert Koeber #552",
	"Graeme Dawson #40020",
	"Peter Ihlo #576",
	"Steve Westall #40092",
	"Philip Weston #432",
}

const LicenceMIT = `Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions: The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software. THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`

func showAboutChasm(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	startHTML(w, "About "+PROGRAMTITLE)

	fmt.Fprint(w, `<article class="about">`)
	fmt.Fprintf(w, `<h1>%v v%v <span class="h2" title="CHasm Ain't ScoreMaster">Powered by CHASM</span><span class="h2"> &amp; EBCFetch</span></h1>`, PROGRAMTITLE, PROGRAMVERSION)
	fmt.Fprintf(w, `<p>%v</p>`, PROGDESC)
	fmt.Fprint(w, `<hr>`)
	fmt.Fprint(w, `<dl>`)
	fmt.Fprint(w, `<dt>Host</dt>`)
	host, _ := os.Hostname()
	fmt.Fprintf(w, `<dd>%v [%v]</dd>`, host, runtime.GOOS)
	fmt.Fprint(w, `<dt>Database file</dt>`)
	path, _ := filepath.Abs(*DBNAME)
	fmt.Fprintf(w, `<dd>%v</dd>`, path)
	fmt.Fprintf(w, `<dt>Freespace</dt><dd>%v</dd>`, showFreeSpace(filepath.Dir(path)))
	fmt.Fprint(w, `</dl>`)
	fmt.Fprint(w, `<hr>`)
	fmt.Fprint(w, `<dl class="legal">`)

	fmt.Fprint(w, `<dt>Author &amp; maintainer</dt>`)
	fmt.Fprintf(w, `<dd>%v</dd>`, Author)
	ibs := strings.ReplaceAll(strings.ReplaceAll(strings.Join(InspiredBy, ","), " ", "&nbsp;"), ",", ", ")
	fmt.Fprintf(w, `<dt>Inspired by</dt><dd>%v</dd>`, ibs)
	fmt.Fprintf(w, `<dt>Chasm [v%v]</dt><dd>github.com/ibauk/chasm</dd>`, ChasmVersion)
	fmt.Fprintf(w, `<dt>EBCFetch [v%v]</dt><dd>github.com/ibauk/ebcfetch</dd>`, EBCFetchVersion)
	fmt.Fprint(w, `<dt>Licence</dt>`)
	fmt.Fprintf(w, `<dd class="link" onclick="toggleLicenceMIT()">MIT - Copyright &copy; %v %v</dd>`, CopyriteYear, CopyriteHolder)
	fmt.Fprintf(w, `</dl>`)
	fmt.Fprintf(w, `<p class="legal hide" id="LicenceMIT">%v</p>`, LicenceMIT)
	fmt.Fprint(w, `</article>`)
}

func showFreeSpace(path string) string {

	const KB = uint64(1024)
	const MB = KB * KB
	const GB = MB * KB

	var res string

	usage := du.NewDiskUsage(path)
	freebytes := usage.Free()
	freemb := freebytes / MB
	freegb := freebytes / GB
	if freegb > 0 {
		res = fmt.Sprintf("%vGB", freegb)
	} else if freemb > 0 {
		res = fmt.Sprintf("%vMB", freemb)
	} else {
		res = fmt.Sprintf("%v bytes", freebytes)
	}

	freep := 100.0 - (usage.Usage() * 100.0)

	res += fmt.Sprintf(" (%.2f%%)", freep)
	return res
}
