package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"syscall/js"

	"golang.org/x/net/html"
)

const (
	TEXT_TYPE = "TEXT_ELEMENT"
)

func ConstructAnElement(tt html.TokenType, z *html.Tokenizer, c *Component) *El {
	token := z.Token()

	parent := NewEl()

	parent.Type = token.Data

	for _, attr := range token.Attr {
		if strings.HasPrefix(attr.Key, "@") {
			method, ok := c.Methods[attr.Val]
			if !ok {
				name := attr.Val
				id := c.ID
				method = func(_ js.Value) {
					log.Printf(`Handler "%s" was not found for a component "%s"`, name, id)
				}
			}

			callback := js.NewEventCallback(js.PreventDefault, method)
			key := strings.Replace(attr.Key, "@", "", 1)
			parent.Callbacks[key] = callback
		} else {
			var at HTMLAttribute

			if strings.HasPrefix(attr.Key, ":") {
				log.Printf("Found dynamic attr %#v", attr)
				handler, ok := c.Computed[attr.Val]

				if !ok {
					name := attr.Val
					id := c.ID
					handler = func() string {
						msg := fmt.Sprintf(`Computed property "%s" in component "%s" was not found`, name, id)
						log.Println(msg)
						return msg
					}
				}

				at = &DynamicAttribute{
					K:  strings.Replace(attr.Key, ":", "", 1),
					Fn: handler,
				}
			} else {
				at = &StaticAttribute{
					K: attr.Key,
					V: attr.Val,
				}
			}

			parent.Attr = append(parent.Attr, at)
		}
	}

	if tt != html.SelfClosingTagToken {
		for {
			tt := z.Next()
			switch {
			case tt == html.ErrorToken:
				err := z.Err()
				if err == io.EOF {
					return parent
				}
				log.Printf("Error: %s", err)
			case tt == html.StartTagToken:
				child := ConstructAnElement(tt, z, c)
				parent.Children = append(parent.Children, child)
			case tt == html.TextToken:
				t := z.Token()
				data := strings.Trim(t.Data, "\r\n ")
				if len(data) > 0 {
					child := NewEl()
					child.Type = TEXT_TYPE
					child.NodeValue = t.Data
					parent.Children = append(parent.Children, child)
				}
			case tt == html.EndTagToken:
				return parent
			case tt == html.SelfClosingTagToken:
				child := ConstructAnElement(tt, z, c)
				parent.Children = append(parent.Children, child)
			case tt == html.CommentToken:
				break
			case tt == html.DoctypeToken:
				break
			}
		}
	} else {
		parent.SelfClosing = true
	}

	return parent
}

func ParseHTML(z *html.Tokenizer, c *Component) *El {
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			err := z.Err()
			if err == io.EOF {
				return nil
			}
			log.Fatal(err)
		case tt == html.StartTagToken:
			return ConstructAnElement(tt, z, c)
		case tt == html.TextToken:
			t := z.Token()
			data := strings.Trim(t.Data, "\r\n ")
			if len(data) > 0 {
				log.Fatalf(`Element can't start with a text node like "%s"`, data)
			}
			continue
		default:
			log.Fatalf(`Wrong token type "%v"`, tt)
		}
		break
	}

	return nil
}
