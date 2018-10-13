package main

import (
	"fmt"
	"log"
	"regexp"
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
			content := vn.Data
			for _, attr := range vn.Attr {
				k := attr.Key()
				v := attr.Val()
				re := regexp.MustCompile(fmt.Sprintf(`\{\{\s*%s\s*\}\}`, k))
				content = re.ReplaceAllString(content, v)
			}

			element = js.Global().Get("document").Call("createTextNode", content)
		default:
			element = js.Global().Get("document").Call("createElement", vn.Tag)

			for _, attr := range vn.Attr {
				element.Set(attr.Key(), attr.Val())
			}

			for key, callback := range vn.Callbacks {
				element.Call("addEventListener", key, callback)
			}
		}

		vn.domNode = &element
	}

	return *vn.domNode
}

func (newVD *VNode) Diff(oldVD *VNode, changeset *[]Change, rootID string, parentNode *js.Value) {
	if oldVD == nil {
		el := newVD.DomElement()
		var root js.Value
		if parentNode != nil {
			root = *parentNode
		} else {
			root = js.Global().Get("document").Call("getElementById", rootID)
		}

		*changeset = append(*changeset, Change{
			Type:       "CREATE",
			domNode:    el,
			parentNode: root,
			NewNode:    newVD,
		})

		newVD.DiffChildren(oldVD, changeset, rootID, el)
		return
	}

	if oldVD != nil {
		newVD.domNode = oldVD.domNode
	}

	// if oldVD.Tag != newVD.Tag {
	// 	*changeset = append(*changeset, Change{
	// 		Type:    "REPLACE",
	// 		NewNode: newVD,
	// 	})
	// 	return
	// }

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

	attributesToUpdate := make([]HTMLAttribute, 0)

OUTER:
	for _, newAttr := range newVD.Attr {
		for _, oldAttr := range oldVD.Attr {
			if newAttr.Key() == oldAttr.Key() {
				log.Printf("Comparing %s=%s and %s=%s", newAttr.Key(), newAttr.Val(), oldAttr.Key(), oldAttr.Val())
				if newAttr.Val() != oldAttr.Val() {
					log.Printf("Found changed attribute %#v", newAttr)
					attributesToUpdate = append(attributesToUpdate, newAttr)
				}
				continue OUTER
			}
		}
		log.Printf("end of the loop")
		attributesToUpdate = append(attributesToUpdate, newAttr)
	}

	if len(attributesToUpdate) > 0 || len(attributesToDelete) > 0 {
		*changeset = append(*changeset, Change{
			Type:               "UPDATE",
			domNode:            *oldVD.domNode,
			attributesToDelete: attributesToDelete,
			attributesToUpdate: attributesToUpdate,
			NewNode:            newVD,
		})
	}

	newVD.DiffChildren(oldVD, changeset, rootID, *oldVD.domNode)
}

func (newVD *VNode) DiffChildren(oldVD *VNode, changeset *[]Change, rootID string, parentNode js.Value) {
	for i, newC := range newVD.Children {
		if oldVD == nil {
			newC.Diff(nil, changeset, rootID, &parentNode)
			continue
		}

		if i < len(oldVD.Children) {
			oldC := oldVD.Children[i]
			newC.Diff(oldC, changeset, rootID, &parentNode)
		}
	}
}
