package main

type HTMLAttribute interface {
	Key() string
	Val() string
}

type StaticAttribute struct {
	K string
	V string
}

func (at *StaticAttribute) Key() string {
	return at.K
}

func (at *StaticAttribute) Val() string {
	return at.V
}

type DynamicAttribute struct {
	K  string
	Fn func() string
}

func (at *DynamicAttribute) Key() string {
	return at.K
}

func (at *DynamicAttribute) Val() string {
	return at.Fn()
}
