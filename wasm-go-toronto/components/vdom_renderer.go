package main

import (
	"fmt"
	"syscall/js"
)

type VNode struct {
	Tag       string
	ID        string
	Attr      []HTMLAttribute
	Callbacks map[string]js.Callback
	Data      string
	Children  []*VNode
	domNode   *js.Value
}

func (vn *VNode) String() string {
	var out string

	switch vn.Tag {
	case TEXT_TYPE:
		out = vn.Data
	default:
		var attrs string
		for _, attr := range vn.Attr {
			attrs += fmt.Sprintf(` %s="%s"`, attr.Key, attr.Val)
		}
		out += fmt.Sprintf("<%s%s>", vn.Tag, attrs)
		for _, child := range vn.Children {
			out += child.String()
		}
		out += fmt.Sprintf("</%s>", vn.Tag)
	}

	return out
}

func (vn *VNode) DomElement() js.Value {
	if vn.domNode == nil {
		var element js.Value

		switch vn.Tag {
		case TEXT_TYPE:
			element = js.Global().Get("document").Call("createTextNode", vn.Data)
		default:
			element = js.Global().Get("document").Call("createElement", vn.Tag)

			for _, attr := range vn.Attr {
				element.Set(attr.Key(), attr.Val())
			}

			for key, callback := range vn.Callbacks {
				element.Call("addEventListener", key, callback)
			}

			for _, child := range vn.Children {
				childEl := child.DomElement()
				element.Call("appendChild", childEl)
			}
		}

		vn.domNode = &element
	}

	return *vn.domNode
}

func (newVD *VNode) Diff(oldVD *VNode, changeset *[]Change) {
	newVD.domNode = oldVD.domNode

	if oldVD.Tag != newVD.Tag {
		*changeset = append(*changeset, Change{
			Type:    "REPLACE",
			NewNode: newVD,
		})
		return
	}

	attributesToDeleteMap := make(map[string]bool, len(oldVD.Attr))

	for _, attr := range oldVD.Attr {
		attributesToDeleteMap[attr.Key()] = true
	}

	for _, attr := range newVD.Attr {
		attributesToDeleteMap[attr.Key()] = false
	}

	attributesToDelete := make([]string, 0)

	for k, v := range attributesToDeleteMap {
		if v {
			attributesToDelete = append(attributesToDelete, k)
		}
	}

	*changeset = append(*changeset, Change{
		Type:               "UPDATE",
		domNode:            oldVD.domNode,
		attributesToDelete: attributesToDelete,
		NewNode:            newVD,
	})

	// go over children
	// figure out which ones are new
	// figure out which ones need to be deleted
	// figure out which ones are old
}

type VDomRenderer struct {
}

func (rr *VDomRenderer) Render(el *El) (*VNode, error) {
	node := &VNode{}

	node.Attr = el.Attr
	node.Callbacks = el.Callbacks

	switch el.Type {
	case TEXT_TYPE:
		node.Data = el.NodeValue
		node.Tag = TEXT_TYPE
	default:
		node.Tag = el.Type

		for _, child := range el.Children {
			child, err := rr.Render(child)
			if err != nil {
				return node, err
			}
			node.Children = append(node.Children, child)
		}
	}

	return node, nil
}
