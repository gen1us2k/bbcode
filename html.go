// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"html"
	"net/url"
	"strings"
)

// HTMLTag represents a DOM node.
type HTMLTag struct {
	Name     string
	Value    string
	Attrs    map[string]string
	Children []*HTMLTag
}

// NewHTMLTag creates a new HTMLTag with string contents specified by value.
func NewHTMLTag(value string) *HTMLTag {
	return &HTMLTag{
		Value:    value,
		Attrs:    make(map[string]string),
		Children: make([]*HTMLTag, 0),
	}
}

func (t *HTMLTag) String() string {
	var value string
	if t.Value != "" {
		value = sanitize(t.Value)
	}
	var attrString string
	for key, value := range t.Attrs {
		attrString = fmt.Sprintf(`%s %s="%s"`, attrString, key, escapeQuotes(sanitize(value)))
	}
	if len(t.Children) > 0 {
		var childrenString string
		for _, child := range t.Children {
			childrenString = fmt.Sprint(childrenString, child.String())
		}
		if t.Name != "" {
			return fmt.Sprintf(`%s<%s%s>%s</%s>`, value, t.Name, attrString, childrenString, t.Name)
		} else {
			return fmt.Sprint(value, childrenString)
		}
	} else if t.Name != "" {
		return fmt.Sprintf(`%s<%s%s>`, value, t.Name, attrString)
	} else {
		return value
	}
}

func (t *HTMLTag) AppendChild(child *HTMLTag) *HTMLTag {
	if child == nil {
		t.Children = append(t.Children, NewHTMLTag(""))
	} else {
		t.Children = append(t.Children, child)
	}
	return t
}

func insertNewlines(out *HTMLTag) {
	if strings.ContainsRune(out.Value, '\n') {
		parts := strings.Split(out.Value, "\n")
		for i, part := range parts {
			if i == 0 {
				out.Value = parts[i]
			} else {
				out.AppendChild(NewlineTag()).AppendChild(NewHTMLTag(part))
			}
		}
	}
}

// Returns a new HTMLTag representing a line break
func NewlineTag() *HTMLTag {
	var out = NewHTMLTag("")
	out.Name = "br"
	return out
}

func escapeQuotes(raw string) string {
	return strings.Replace(strings.Replace(raw, `"`, `\"`, -1), `\`, `\\`, -1)
}

func safeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return strings.Replace(u.String(), `\`, "%5C", -1)
}

func sanitize(raw string) string {
	return html.EscapeString(raw)
}