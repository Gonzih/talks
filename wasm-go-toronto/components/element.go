package main

type El struct {
	Type        string
	Attr        []*HTMLAttr
	NodeValue   string
	Children    []*El
	SelfClosing bool
}
