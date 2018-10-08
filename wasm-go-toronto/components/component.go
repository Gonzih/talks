package main

import (
	"strings"
	"syscall/js"

	"github.com/google/uuid"
	"golang.org/x/net/html"
)

type Component struct {
	Template         string
	ID               string
	Tree             *El
	OldVDom          *VNode
	Computed         map[string]func() string
	Methods          map[string]func(js.Value)
	notificationChan chan bool
}

func NewComponent(templateID string, methods map[string]func(js.Value), computed map[string]func() string) *Component {
	cmp := &Component{}
	cmp.Methods = methods
	cmp.Computed = computed
	cmp.notificationChan = make(chan bool)

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
	vdom.Diff(cmp.OldVDom, &changes, rootID, nil)

	for _, change := range changes {
		change.Apply()
	}

	cmp.OldVDom = vdom
}

func (cmp *Component) MountTo(rootID string) {
	cmp.RenderTo(rootID)

	for range cmp.notificationChan {
		var cb js.Callback

		callback := js.NewCallback(func(_ []js.Value) {
			cmp.RenderTo(rootID)
			cb.Release()
		})

		js.Global().Get("window").Call("requestAnimationFrame", callback)
	}
}
