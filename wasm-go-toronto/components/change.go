package main

import (
	"log"
	"syscall/js"
)

type Change struct {
	Type    string
	NewNode *VNode
	// Children           []Change
	domNode            js.Value
	parentNode         js.Value
	attributesToDelete []string
}

func (ch *Change) Apply() {
	log.Printf("Applying %s change: %#v", ch.Type, ch)

	switch ch.Type {
	case "UPDATE":
		ch.update()
	case "CREATE":
		ch.create()
	}
}

func (ch *Change) update() {
	for _, attrName := range ch.attributesToDelete {
		ch.domNode.Call("removeAttribute", attrName)
	}

	for _, attr := range ch.NewNode.Attr {
		log.Printf("setAttribute %s %s", attr.Key(), attr.Val())
		ch.domNode.Call("setAttribute", attr.Key(), attr.Val())
	}
}

func (ch *Change) create() {
	ch.parentNode.Call("appendChild", ch.domNode)
}
