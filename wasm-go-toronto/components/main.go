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

func newCmp() *Component {
	store := NewStore()
	store.Set("counter", 0)

	methods := map[string]func(js.Value){
		"ClickHandler": func(event js.Value) {
			c := store.Get("counter").(int)
			store.Set("counter", c+1)
		},
	}

	computed := map[string]func() string{
		"Count": func() string {
			return fmt.Sprintf("%d", store.Get("counter"))
		},
		"LabelClass": func() string {
			c := store.Get("counter").(int)
			if c < 10 {
				return "small"
			}

			return "large"
		},
	}

	cmp := NewComponent("rootTemplate", methods, computed)

	store.Subscribe(cmp.notificationChan)

	return cmp
}

func main() {
	for i := 0; i < 3; i++ {
		cmp := newCmp()
		go cmp.MountTo(fmt.Sprintf("root%d", i))
	}

	// Render in to VDOM

	lock := make(chan bool)
	<-lock
}
