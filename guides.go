package main

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed betatest.md
var betatest string

var guideheader = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>chasm</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="/css?file=normalize">
<link rel="stylesheet" href="/css?file=maincss">
<style>
body {
background-color: #F0FFFF;
}
.guide {
width: 98%;
max-width: 200mm;
margin: auto;
padding: 0 0 2em 0;
}
p {
text-align: justify;
margin: .3em 0 0 0;
}
h2,
h3,
h4 {
margin: 1em 0 0 0;
}
h4 {
font-size: smaller;
}
ol,
ul {
margin: .5em 0 0 0;
list-style-position: inside;
}
ol li {
margin: .5em 0 0 0;
}
</style>
</head>
<body>
`

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func showGuides(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	fmt.Fprint(w, guideheader)
	fmt.Fprint(w, `<article class="guide">`)
	md := []byte(betatest)
	html := mdToHTML(md)
	fmt.Fprint(w, string(html))
	fmt.Fprint(w, `</article>`)
}
