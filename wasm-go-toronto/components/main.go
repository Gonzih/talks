package main

import (
	"fmt"
	"log"
	"syscall/js"
)

func must(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func main() {
	store := NewStore()
	store.Set("counter", 0)

	methods := map[string]func(js.Value){
		"ClickHandler": func(event js.Value) {
			c := store.Get("counter").(int)
			store.Set("counter", c+1)
			log.Println(c)
		},
	}

	computed := map[string]func() string{
		"Counter": func() string {
			return fmt.Sprintf("Clicked me %d times", store.Get("counter"))
		},
	}

	cmp := NewComponent("rootTemplate", methods, computed)
	cmp.MountTo("root")

	// Render in to VDOM

	lock := make(chan bool)
	<-lock
}
