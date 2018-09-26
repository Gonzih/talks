package main

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

const (
	TEXT_TYPE = "TEXT_ELEMENT"
)

func ConstructAnElement(tt html.TokenType, z *html.Tokenizer) *El {
	token := z.Token()

	parent := &El{}

	parent.Type = token.Data

	for _, attr := range token.Attr {
		parent.Attr = append(parent.Attr, &HTMLAttr{
			Key: attr.Key,
			Val: attr.Val,
		})
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
				child := ConstructAnElement(tt, z)
				parent.Children = append(parent.Children, child)
			case tt == html.TextToken:
				t := z.Token()
				data := strings.Trim(t.Data, "\r\n ")
				if len(data) > 0 {
					child := &El{Type: TEXT_TYPE, NodeValue: t.Data}
					parent.Children = append(parent.Children, child)
				}
			case tt == html.EndTagToken:
				return parent
			case tt == html.SelfClosingTagToken:
				child := ConstructAnElement(tt, z)
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

func ParseHTML(z *html.Tokenizer) *El {
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
			return ConstructAnElement(tt, z)
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
