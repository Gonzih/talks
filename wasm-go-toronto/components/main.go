package main

import (
	"log"
	"strings"

	"golang.org/x/net/html"
)

func must(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func main() {
	template, err := dom.GetInnerHTMLByID("rootTemplate")
	must(err)

	log.Println(template)

	r := strings.NewReader(template)
	z := html.NewTokenizer(r)

	el := ParseHTML(z)

	renderer := &DebugRenderer{}

	log.Println(el)

	err = renderer.Render(el)
	must(err)

	// Render in to VDOM
	vRenderer := &VDomRenderer{}

	vdom, err := vRenderer.Render(el)
	must(err)

	log.Println(vdom.String())
}
