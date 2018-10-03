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
	OldVDom  *VNode
	Computed map[string]func() string
	Methods  map[string]func(js.Value)
}

func NewComponent(templateID string, methods map[string]func(js.Value), computed map[string]func() string) *Component {
	cmp := &Component{}
	cmp.Methods = methods
	cmp.Computed = computed

	template, err := dom.GetInnerHTMLByID(templateID)
	must(err)

	cmp.Template = template
	cmp.ID = uuid.New().String()

	r := strings.NewReader(cmp.Template)
	z := html.NewTokenizer(r)

	el := ParseHTML(z, cmp)

	el.Attr = append(el.Attr, &StaticAttribute{
		K: "go-id",
		V: cmp.ID,
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
	vdom := cmp.RenderToVDom()
	vdom.Diff(cmp.OldVDom, &changes)

	for _, change := range changes {
		el := change.NewNode.DomElement()
		root := js.Global().Get("document").Call("getElementById", rootID)
		root.Set("innerHTML", "")
		root.Call("appendChild", el)
	}

	// cmp.OldVDom = vdom
}

func (cmp *Component) MountTo(rootID string) {
	cmp.RenderTo(rootID)

	for range notificationChan {
		cmp.RenderTo(rootID)
	}
}
