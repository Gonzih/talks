package main

import "syscall/js"

type El struct {
	Type        string
	Attr        []HTMLAttribute
	Callbacks   map[string]js.Callback
	NodeValue   string
	Children    []*El
	SelfClosing bool
}

func NewEl() *El {
	el := &El{}
	el.Callbacks = make(map[string]js.Callback, 0)

	return el
}
