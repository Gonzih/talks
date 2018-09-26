package main

import (
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

type Component struct {
	Template string
	ID       string
	Tree     *El
}

func NewComponent(templateID string) *Component {
	cmp := &Component{}

	template, err := dom.GetInnerHTMLByID(templateID)
	must(err)

	cmp.Template = template
	cmp.ID = uuid.New().String()

	r := strings.NewReader(cmp.Template)
	z := html.NewTokenizer(r)

	el := ParseHTML(z)

	el.Attr = append(el.Attr, &HTMLAttr{
		Key: "go-id",
		Val: cmp.ID,
	})

	cmp.Tree = el

	return cmp
}

func (cmp *Component) Debug() {
	renderer := &DebugRenderer{}
	err := renderer.Render(cmp.Tree)
	must(err)
}

func (cmp *Component) String() string {
	vRenderer := &VDomRenderer{}

	vdom, err := vRenderer.Render(cmp.Tree)
	must(err)

	return vdom.String()
}

func (cmp *Component) RenderTo(rootID string) {
	must(dom.SetInnerHTMLByID(rootID, cmp.String()))
}
