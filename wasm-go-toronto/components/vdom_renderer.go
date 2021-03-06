package main

type VDomRenderer struct {
}

func (rr *VDomRenderer) Render(el *El) (*VNode, error) {
	node := &VNode{}

	for _, attr := range el.Attr {
		node.Attr = append(node.Attr, &StaticAttribute{
			K: attr.Key(),
			V: attr.Val(),
		})
	}

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
