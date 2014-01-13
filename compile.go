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

type htmlTag struct {
	name     string
	value    string
	attrs    map[string]string
	children []*htmlTag
}

func newHtmlTag(value string) *htmlTag {
	return &htmlTag{
		value:    value,
		attrs:    make(map[string]string),
		children: make([]*htmlTag, 0),
	}
}

func (t *htmlTag) string() string {
	if t.value != "" {
		return sanitize(t.value)
	}
	var attrString string
	for key, value := range t.attrs {
		attrString = fmt.Sprintf(`%s %s="%s"`, attrString, key, escapeQuotes(sanitize(value)))
	}
	if len(t.children) > 0 {
		var childrenString string
		for _, child := range t.children {
			childrenString = fmt.Sprint(childrenString, child.string())
		}
		return fmt.Sprintf(`<%s%s>%s</%s>`, t.name, attrString, childrenString, t.name)
	} else {
		return fmt.Sprintf(`<%s%s/>`, t.name, attrString)
	}
}

func (t *htmlTag) appendChild(child *htmlTag) {
	t.children = append(t.children, child)
}

// compile transforms a tag and subexpression into an HTML string.
// It is only used by the generated parser code.
func compile(in bbTag, expr *htmlTag) *htmlTag {
	var out = newHtmlTag("")

	switch {
	case in.key == "url":
		out.name = "a"
		if in.value == "" {
			out.attrs["href"] = safeURL(expr.value)
		} else {
			out.attrs["href"] = safeURL(in.value)
		}
		out.appendChild(expr)
	case in.key == "img":
		out.name = "img"
		if in.value == "" {
			out.attrs["src"] = safeURL(expr.value)
		} else {
			out.attrs["src"] = safeURL(in.value)
			out.attrs["alt"] = expr.value
		}
	case in.key == "i" || in.key == "b":
		out.name = in.key
		out.appendChild(expr)
	}
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