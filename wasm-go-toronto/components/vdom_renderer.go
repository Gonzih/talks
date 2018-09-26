package main

import (
	"fmt"
	"syscall/js"
)

type VNode struct {
	Tag      string
	Attr     []*HTMLAttr
	Data     string
	Children []*VNode
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
	var element js.Value

	switch vn.Tag {
	case TEXT_TYPE:
		element = js.Global().Get("document").Call("createTextNode", vn.Data)
	default:
		element = js.Global().Get("document").Call("createElement", vn.Tag)

		for _, attr := range vn.Attr {
			element.Set(attr.Key, attr.Val)
		}

		for _, child := range vn.Children {
			childEl := child.DomElement()
			element.Call("appendChild", childEl)
		}
	}

	return element
}

func (newVD *VNode) Diff(oldVD *VNode, changeset *[]Change) {
	if oldVD == nil {
		*changeset = append(*changeset, Change{
			Type:    "CREATE",
			NewNode: newVD,
		})
	}
}

type VDomRenderer struct {
}

func (rr *VDomRenderer) Render(el *El) (*VNode, error) {
	node := &VNode{}

	for _, attr := range el.Attr {
		node.Attr = append(node.Attr, &HTMLAttr{
			Key: attr.Key,
			Val: attr.Val,
		})
	}

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
