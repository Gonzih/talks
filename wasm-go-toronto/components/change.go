package main

import "syscall/js"

type Change struct {
	Type               string
	NewNode            *VNode
	Children           []Change
	domNode            *js.Value
	attributesToDelete []string
}

func (ch *Change) Apply() {
	for _, attrName := range ch.attributesToDelete {
		ch.domNode.Call("removeAttribute", attrName)
	}

	for _, attr := range ch.NewNode.Attr {
		ch.domNode.Call("setAttribute", attr.Key(), attr.Val())
	}
}
