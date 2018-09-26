package main

import (
	"strings"
	"syscall/js"

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

func (cmp *Component) RenderToVDom() *VNode {
	vRenderer := &VDomRenderer{}

	vdom, err := vRenderer.Render(cmp.Tree)
	must(err)

	return vdom
}

func (cmp *Component) String() string {
	return cmp.RenderToVDom().String()
}

func (cmp *Component) RenderTo(rootID string) {
	changes := make([]Change, 0)
	cmp.RenderToVDom().Diff(nil, &changes)

	for _, change := range changes {
		el := change.NewNode.DomElement()
		js.Global().Get("document").Call("getElementById", rootID).Call("appendChild", el)
	}
}
