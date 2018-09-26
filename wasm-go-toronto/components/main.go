package main

import (
	"log"
)

func must(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func main() {
	cmp := NewComponent("rootTemplate")
	cmp.RenderTo("root")

	// Render in to VDOM
}
